// Package api provides HTTP client functionality with intelligent retry logic
// for the Couchbase Capella API. This package uses the Hashicorp retryablehttp
// library for robust retry mechanisms with exponential backoff and jitter.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// maxRetryAttempts defines the maximum number of retry attempts that will be made
// for retryable HTTP responses (429 and 504). After this many failed attempts,
// the last response will be returned to the caller without further retries.
// Value: 5 attempts (6 total requests including the initial request).
const maxRetryAttempts = 5

// customRetryPolicy implements retry logic for the Couchbase Capella API.
// It handles specific cases for 429 (rate limiting) and 504 (gateway timeout) responses.
//
// Retry Behavior:
//   - HTTP 429 (Too Many Requests): Always retries, respecting Retry-After header if present
//   - HTTP 504 (Gateway Timeout): Retries unless the response body contains "code": 7001,
//     which indicates a non-retryable error condition for index DDL operations
//   - All other status codes: Returns immediately without retry
//
// For 504 responses with code 7001, it returns the specific ErrGatewayTimeoutForIndexDDL
// error to provide better error context to callers, matching V1 client behavior.
func customRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Handle connection errors - retry these
	if err != nil {
		return true, nil
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		// Always retry 429 responses - retryablehttp handles Retry-After header automatically
		return true, nil

	case http.StatusGatewayTimeout:
		// Check for special code 7001 (do not retry)
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			// If we can't read the body, treat it as a regular 504 and retry
			return true, nil
		}
		// Restore body for subsequent reads
		resp.Body = io.NopCloser(bytes.NewReader(body))

		var apiErr struct {
			Code int `json:"code"`
		}
		if err := json.Unmarshal(body, &apiErr); err != nil {
			// If we can't parse the error, treat it as a regular 504 and retry
			return true, nil
		}
		if apiErr.Code == 7001 {
			// Don't retry and return the specific error - matches V1 client behavior
			return false, errors.ErrGatewayTimeoutForIndexDDL
		}
		// Retry other 504 responses
		return true, nil

	default:
		// Don't retry other status codes
		return false, nil
	}
}

// NewRetryHTTPClient creates and returns a new HTTP client configured with intelligent
// retry logic using the Hashicorp retryablehttp library. The client automatically handles
// transient failures such as rate limiting (429) and gateway timeouts (504) with
// exponential backoff and jitter.
//
// The returned client is fully configured and ready for use with any HTTP operations.
// It provides the same interface as a standard http.Client but with enhanced reliability
// for API interactions.
//
// Configuration:
//   - Uses retryablehttp.Client with custom retry policy for Couchbase Capella API
//   - Applies LinearJitterBackoff for optimal retry spacing and thundering herd prevention
//   - Sets the specified timeout for all requests (applies to entire retry sequence)
//   - Respects Retry-After headers automatically when present
//   - Limits retries to maxRetryAttempts (5) per request
//   - Handles special 7001 error codes for index DDL operations
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
	retryClient := retryablehttp.NewClient()

	// Configure retry behavior for Couchbase Capella API
	retryClient.RetryMax = maxRetryAttempts
	retryClient.CheckRetry = customRetryPolicy
	retryClient.Backoff = retryablehttp.LinearJitterBackoff

	// Disable default logging to avoid noise - clients can configure their own
	retryClient.Logger = nil

	// Convert to standard HTTP client and set timeout
	httpClient := retryClient.StandardClient()
	httpClient.Timeout = timeout

	return httpClient
}
