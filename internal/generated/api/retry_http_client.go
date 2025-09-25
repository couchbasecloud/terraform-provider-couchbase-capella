// Package api provides HTTP client functionality with intelligent retry logic
// for the Couchbase Capella API. This package uses the Hashicorp retryablehttp
// library for robust retry mechanisms with exponential backoff and jitter.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// maxRetryAttempts defines the maximum number of retry attempts that will be made
// for retryable HTTP responses (429 and 504). After this many failed attempts,
// the last response will be returned to the caller without further retries.
// Value: 5 attempts (6 total requests including the initial request).
const maxRetryAttempts = 5

// tflogBridge implements retryablehttp.Logger interface and forwards log messages to tflog.
// This ensures consistent logging throughout the provider using Terraform's structured logging system.
type tflogBridge struct {
	ctx context.Context
}

// Printf implements the retryablehttp.Logger interface by forwarding to tflog.Debug.
// This integrates retry logging with Terraform's logging system, allowing users to control
// retry logging via TF_LOG environment variable alongside other provider logs.
func (t *tflogBridge) Printf(format string, args ...interface{}) {
	// Use tflog.Debug for retry messages as they are diagnostic information
	// The [RETRY] prefix helps identify these logs among other provider logs
	tflog.Debug(t.ctx, fmt.Sprintf("[RETRY] "+format, args...))
}

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
//
//nolint:nilerr // returning nil error is expected behavior for retryablehttp.CheckRetry interface
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

// RetryOption configures retry behavior for the HTTP client.
type RetryOption func(*retryablehttp.Client)

// WithFastBackoff configures the client to use fast backoff delays suitable for testing.
// This reduces retry delays to: 50ms, 100ms, 200ms, 400ms, 800ms for quick test feedback.
func WithFastBackoff() RetryOption {
	return func(client *retryablehttp.Client) {
		client.Backoff = func(minDelay, maxDelay time.Duration, attemptNum int, resp *http.Response) time.Duration {
			// Use very short delays for testing: 50ms, 100ms, 200ms, 400ms, 800ms
			delay := 50 * time.Millisecond * (1 << attemptNum)
			if delay > 800*time.Millisecond {
				delay = 800 * time.Millisecond
			}
			return delay
		}
	}
}

// WithMaxRetries configures the maximum number of retry attempts.
func WithMaxRetries(maxRetries int) RetryOption {
	return func(client *retryablehttp.Client) {
		client.RetryMax = maxRetries
	}
}

// WithLogger configures a custom logger for the retry client.
func WithLogger(logger retryablehttp.Logger) RetryOption {
	return func(client *retryablehttp.Client) {
		client.Logger = logger
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
// Default Configuration:
//   - Uses retryablehttp.Client with custom retry policy for Couchbase Capella API
//   - Applies LinearJitterBackoff for optimal retry spacing and thundering herd prevention
//   - Sets the specified timeout for all requests (applies to entire retry sequence)
//   - Respects Retry-After headers automatically when present
//   - Limits retries to maxRetryAttempts (5) per request
//   - Handles special 7001 error codes for index DDL operations
//   - Enables/disables retry logging based on debugLogging parameter
//
// Parameters:
//   - timeout: The maximum duration for each individual HTTP request (including retries).
//     This timeout applies to the entire retry sequence, not individual attempts.
//     Use 0 for no timeout, though this is not recommended for production use.
//   - debugLogging: When true, enables retry attempt logging. Should be controlled by
//     the provider's debug logging settings (e.g., TF_LOG=DEBUG).
//   - opts: Optional configuration functions to customize retry behavior
//
// Returns:
//   - *http.Client: A fully configured HTTP client ready for use
//
// Example Usage:
//
//	// Production client with no retry logging
//	client := NewRetryHTTPClient(ctx, 30 * time.Second, false)
//
//	// Debug client with retry logging enabled  
//	debugClient := NewRetryHTTPClient(ctx, 30 * time.Second, true)
//
//	// Testing client with fast backoff and debug logging
//	testClient := NewRetryHTTPClient(ctx, 30 * time.Second, true, WithFastBackoff())
//
//	// Custom configuration
//	customClient := NewRetryHTTPClient(ctx, 30 * time.Second, false,
//		WithMaxRetries(3),
//		WithFastBackoff())
//
// Thread Safety:
// The returned client is safe for concurrent use by multiple goroutines.
func NewRetryHTTPClient(ctx context.Context, timeout time.Duration, debugLogging bool, opts ...RetryOption) *http.Client {
	retryClient := retryablehttp.NewClient()

	// Configure default retry behavior for Couchbase Capella API
	retryClient.RetryMax = maxRetryAttempts
	retryClient.CheckRetry = customRetryPolicy
	retryClient.Backoff = retryablehttp.LinearJitterBackoff

	// Configure logging based on debug setting
	if debugLogging {
		// Use tflogBridge for consistent structured logging throughout the provider
		// This integrates retry logs with Terraform's logging system via tflog.Debug
		// The provider can control this via TF_LOG=DEBUG or TF_LOG_PROVIDER=DEBUG
		retryClient.Logger = &tflogBridge{ctx: ctx}
	} else {
		retryClient.Logger = nil
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(retryClient)
	}

	// Convert to standard HTTP client and set timeout
	httpClient := retryClient.StandardClient()
	httpClient.Timeout = timeout

	return httpClient
}
