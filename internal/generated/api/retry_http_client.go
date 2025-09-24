// Package api provides HTTP client functionality with intelligent retry logic
// for the Couchbase Capella API. This package includes exponential backoff
// retry mechanisms for handling rate limits and temporary server errors.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RetryTransport implements an http.RoundTripper that provides intelligent retry logic
// with exponential backoff for transient HTTP errors. It is designed to handle
// rate limiting (429) and gateway timeout (504) responses gracefully.
//
// Retry Behavior:
//   - HTTP 429 (Too Many Requests): Respects the Retry-After header if present.
//     If no Retry-After header is provided, uses exponential backoff with jitter.
//   - HTTP 504 (Gateway Timeout): Uses exponential backoff with jitter unless
//     the response body contains a JSON object with "code": 7001, which indicates
//     a non-retryable error condition.
//   - All other status codes: Returns immediately without retry.
//
// Backoff Strategy:
// The exponential backoff follows the formula: baseBackoffDelay * 2^attempt + jitter
// where jitter is ±25% of the calculated delay to prevent thundering herd problems.
// The backoff delay is capped at maxBackoffDelay (30 seconds).
//
// Retry Limits:
// Both 429 and 504 responses are retried up to maxRetryAttempts (5) times.
// After the maximum attempts, the last response is returned to the caller.
//
// Request Body Handling:
// The transport buffers request bodies to enable retries of non-idempotent requests.
// This ensures that POST, PUT, and other requests with bodies can be safely retried.
type RetryTransport struct {
	// Base is the underlying RoundTripper to use for making requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

// base returns the underlying RoundTripper to use for making HTTP requests.
// If Base is nil, it returns http.DefaultTransport as a fallback.
func (t *RetryTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// baseBackoffDelay is the initial delay used as the base for exponential backoff calculations.
// The first retry will wait approximately this duration (plus jitter).
// Value: 1 second.
const baseBackoffDelay = time.Second * 1

// maxBackoffDelay is the maximum delay that will be applied between retry attempts.
// This prevents exponentially growing delays from becoming excessively long.
// Even with jitter, delays will not exceed this value by more than 25%.
// Value: 30 seconds.
const maxBackoffDelay = time.Second * 30

// maxRetryAttempts defines the maximum number of retry attempts that will be made
// for retryable HTTP responses (429 and 504). After this many failed attempts,
// the last response will be returned to the caller without further retries.
// Value: 5 attempts (6 total requests including the initial request)
const maxRetryAttempts = 5

// calculateBackoff computes the delay duration for retry attempts using exponential
// backoff with jitter. This helps to distribute retry attempts over time and prevents
// thundering herd problems when multiple clients are retrying simultaneously.
//
// Parameters:
//   - attempt: The retry attempt number (0-based). The first retry is attempt 0,
//     the second retry is attempt 1, etc.
//
// Returns:
//   - A time.Duration representing how long to wait before the next retry attempt.
//
// Algorithm:
//  1. Calculate base delay: baseBackoffDelay * 2^attempt
//  2. Cap the delay at maxBackoffDelay to prevent excessively long waits
//  3. Add random jitter of ±25% to prevent synchronized retries
//  4. Ensure the result is never negative
//
// Examples:
//   - Attempt 0: ~1s (750ms - 1.25s with jitter)
//   - Attempt 1: ~2s (1.5s - 2.5s with jitter)
//   - Attempt 2: ~4s (3s - 5s with jitter)
//   - Attempt 3: ~8s (6s - 10s with jitter)
//   - Attempt 4: ~16s (12s - 20s with jitter)
//   - Attempt 5+: ~30s (22.5s - 37.5s with jitter, capped at maxBackoffDelay)
func calculateBackoff(attempt int) time.Duration {
	// Calculate exponential backoff: baseDelay * 2^attempt.
	delay := baseBackoffDelay * time.Duration(1<<attempt)

	// Cap the delay at maxBackoffDelay.
	if delay > maxBackoffDelay {
		delay = maxBackoffDelay
	}

	// Add jitter: +/-25% of the calculated delay to prevent thundering herd.
	jitter := float64(delay) * 0.25 * (rand.Float64() - 0.5) * 2 // #nosec G404 -- non-cryptographic jitter
	finalDelay := float64(delay) + jitter

	// Ensure delay is not negative.
	if finalDelay < 0 {
		finalDelay = float64(baseBackoffDelay)
	}

	return time.Duration(finalDelay)
}

// RoundTrip executes a single HTTP transaction with intelligent retry logic.
// This method implements the http.RoundTripper interface and provides automatic
// retries for transient failures with exponential backoff.
//
// The method handles the following scenarios:
//
// 1. Request Body Buffering:
//   - Buffers the entire request body in memory to enable retries
//   - Closes the original body and replaces it with a reusable reader
//   - This allows non-idempotent requests (POST, PUT, etc.) to be safely retried
//
// 2. HTTP 429 (Too Many Requests) Handling:
//   - Checks for a Retry-After header and respects its value if present
//   - Falls back to exponential backoff if no Retry-After header is provided
//   - Completely drains and closes the response body before retrying
//   - Retries up to maxRetryAttempts times
//
// 3. HTTP 504 (Gateway Timeout) Handling:
//   - Reads and parses the response body to check for special error codes
//   - If the body contains "code": 7001, treats it as non-retryable
//   - For all other 504 responses, uses exponential backoff for retries
//   - Restores the response body for the caller if returning the response
//   - Retries up to maxRetryAttempts times
//
// 4. All Other HTTP Status Codes:
//   - Returns immediately without retry attempts
//   - Preserves the original response and any errors
//
// Context Cancellation:
// The method respects context cancellation during backoff delays. If the request
// context is cancelled or times out during a delay, the method returns the last
// response along with the context error.
//
// Parameters:
//   - req: The HTTP request to execute. Must not be nil.
//
// Returns:
//   - *http.Response: The HTTP response from the server (may be from original or retry attempt)
//   - error: Any error that occurred during the request or retry process
//
// Thread Safety:
// This method is safe for concurrent use by multiple goroutines.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the body so we can retry non-idempotent requests as well.
	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if closeErr := req.Body.Close(); closeErr != nil {
			return nil, closeErr
		}
		bodyBytes = b
	}

	makeReq := func(ctx context.Context) *http.Request {
		r := req.Clone(ctx)
		if bodyBytes != nil {
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		return r
	}

	retryCount := 0

	for {
		r := makeReq(req.Context())
		res, err := t.base().RoundTrip(r)
		if err != nil {
			return nil, err
		}

		switch res.StatusCode {
		case http.StatusTooManyRequests:
			if retryCount >= maxRetryAttempts {
				return res, nil
			}

			var backoff time.Duration
			// Respect Retry-After header if present, otherwise use exponential backoff.
			if ra := res.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil && secs > 0 {
					backoff = time.Second * time.Duration(secs)
				} else {
					backoff = calculateBackoff(retryCount)
				}
			} else {
				backoff = calculateBackoff(retryCount)
			}

			// Drain and retry.
			if _, err := io.Copy(io.Discard, res.Body); err != nil {
				return res, err
			}
			if err := res.Body.Close(); err != nil {
				return res, err
			}
			if err := sleepOrDone(req.Context(), backoff); err != nil {
				return res, err
			}
			retryCount++
			continue
		case http.StatusGatewayTimeout:
			if retryCount >= maxRetryAttempts {
				return res, nil
			}
			// Peek body to check for special code 7001 (do not retry).
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return res, err
			}
			if err := res.Body.Close(); err != nil {
				return res, err
			}
			// Restore body for caller if we return.
			res.Body = io.NopCloser(bytes.NewReader(b))

			var apiErr struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(b, &apiErr); err != nil {
				// If we can't parse the error, treat it as a regular 504 and retry.
				apiErr.Code = 0
			}
			if apiErr.Code == 7001 {
				return res, nil
			}

			// Use exponential backoff for 504 errors.
			backoff := calculateBackoff(retryCount)
			if err := sleepOrDone(req.Context(), backoff); err != nil {
				return res, err
			}
			retryCount++
			continue
		default:
			return res, nil
		}
	}
}

// sleepOrDone implements a cancellable sleep operation that respects context cancellation.
// This function will wait for the specified duration or until the context is cancelled,
// whichever occurs first. This is essential for allowing retry delays to be interrupted
// when a request context times out or is explicitly cancelled.
//
// Parameters:
//   - ctx: The context to monitor for cancellation signals
//   - d: The duration to sleep (must be >= 0)
//
// Returns:
//   - nil if the full duration elapsed without context cancellation
//   - ctx.Err() if the context was cancelled or timed out before the duration elapsed
//
// Behavior:
//   - Uses a timer to implement the delay
//   - Properly cleans up the timer resources regardless of how the function exits
//   - Returns immediately if the context is already cancelled
//
// This function is crucial for maintaining responsiveness during retry delays and
// preventing resource leaks from abandoned retry attempts.
func sleepOrDone(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// NewRetryHTTPClient creates and returns a new HTTP client configured with intelligent
// retry logic. The client uses RetryTransport to automatically handle transient failures
// such as rate limiting (429) and gateway timeouts (504) with exponential backoff.
//
// The returned client is fully configured and ready for use with any HTTP operations.
// It provides the same interface as a standard http.Client but with enhanced reliability
// for API interactions.
//
// Configuration:
//   - Uses http.DefaultTransport as the base transport
//   - Configures RetryTransport for intelligent retry behavior
//   - Sets the specified timeout for all requests
//   - Applies exponential backoff with jitter for retries
//   - Respects Retry-After headers when present
//   - Limits retries to maxRetryAttempts (5) per request
//
// Parameters:
//   - timeout: The maximum duration for each individual HTTP request (including retries).
//     This timeout applies to the entire retry sequence, not individual attempts.
//     Use 0 for no timeout, though this is not recommended for production use.
//
// Returns:
//   - *http.Client: A fully configured HTTP client ready for use
//
// Example Usage:
//
//	client := NewRetryHTTPClient(30 * time.Second)
//	resp, err := client.Get("https://api.example.com/data")
//	if err != nil {
//	    // Handle error (this will be returned only after all retries are exhausted)
//	}
//	defer resp.Body.Close()
//
// Thread Safety:
// The returned client is safe for concurrent use by multiple goroutines.
func NewRetryHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: &RetryTransport{Base: http.DefaultTransport},
	}
}
