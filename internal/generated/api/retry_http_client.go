package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RetryTransport implements retry logic with exponential backoff for HTTP 429 and 504 responses.
// - 429: respects Retry-After header (seconds) if present, otherwise uses exponential backoff with jitter
// - 504: uses exponential backoff with jitter unless body contains code 7001 (do not retry)
// - Both retry up to maxRetryAttempts times with exponential backoff: baseBackoffDelay * 2^attempt + jitter
type RetryTransport struct {
	Base http.RoundTripper
}

func (t *RetryTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// baseBackoffDelay is the base delay for exponential backoff.
const baseBackoffDelay = time.Second * 1

// maxBackoffDelay is the maximum delay for exponential backoff.
const maxBackoffDelay = time.Second * 30

// maxRetryAttempts is the maximum number of retry attempts.
const maxRetryAttempts = 5

// calculateBackoff calculates exponential backoff delay with jitter.
// attempt is 0-based (0 for first retry, 1 for second retry, etc.).
func calculateBackoff(attempt int) time.Duration {
	// Calculate exponential backoff: baseDelay * 2^attempt
	delay := float64(baseBackoffDelay) * math.Pow(2, float64(attempt))
	
	// Cap the delay at maxBackoffDelay
	if delay > float64(maxBackoffDelay) {
		delay = float64(maxBackoffDelay)
	}
	
	// Add jitter: +/-25% of the calculated delay to prevent thundering herd
	jitter := delay * 0.25 * (rand.Float64() - 0.5) * 2
	finalDelay := delay + jitter
	
	// Ensure delay is not negative
	if finalDelay < 0 {
		finalDelay = float64(baseBackoffDelay)
	}
	
	return time.Duration(finalDelay)
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the body so we can retry non-idempotent requests as well
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
			// Respect Retry-After header if present, otherwise use exponential backoff
			if ra := res.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil && secs > 0 {
					backoff = time.Second * time.Duration(secs)
				} else {
					backoff = calculateBackoff(retryCount)
				}
			} else {
				backoff = calculateBackoff(retryCount)
			}
			
			// Drain and retry
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
			// Peek body to check for special code 7001 (do not retry)
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return res, err
			}
			if err := res.Body.Close(); err != nil {
				return res, err
			}
			// restore body for caller if we return
			res.Body = io.NopCloser(bytes.NewReader(b))

			var apiErr struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(b, &apiErr); err != nil {
				// If we can't parse the error, treat it as a regular 504 and retry
				apiErr.Code = 0
			}
			if apiErr.Code == 7001 {
				return res, nil
			}
			
			// Use exponential backoff for 504 errors
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

// NewRetryHTTPClient returns an *http.Client with RetryTransport configured.
func NewRetryHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: &RetryTransport{Base: http.DefaultTransport},
	}
}
