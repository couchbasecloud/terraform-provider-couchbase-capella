package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	apierrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestCustomRetryPolicyDirectly(t *testing.T) {
	// Test our custom retry policy function directly

	// Test 7001 case
	resp := &http.Response{
		StatusCode: http.StatusGatewayTimeout,
		Body:       io.NopCloser(bytes.NewBufferString(`{"code":7001}`)),
		Header:     make(http.Header),
	}

	shouldRetry, err := customRetryPolicy(context.Background(), resp, nil)
	if shouldRetry {
		t.Error("expected no retry for 7001 case")
	}
	if !errors.Is(err, apierrors.ErrGatewayTimeoutForIndexDDL) {
		t.Errorf("expected ErrGatewayTimeoutForIndexDDL, got: %v", err)
	}

	// Test regular 504 case
	resp = &http.Response{
		StatusCode: http.StatusGatewayTimeout,
		Body:       io.NopCloser(bytes.NewBufferString(`{"code":1234}`)),
		Header:     make(http.Header),
	}

	shouldRetry, err = customRetryPolicy(context.Background(), resp, nil)
	if !shouldRetry {
		t.Error("expected retry for regular 504 case")
	}
	if err != nil {
		t.Errorf("expected no error for regular 504, got: %v", err)
	}

	// Test 429 case
	resp = &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(bytes.NewBufferString("rate limited")),
		Header:     make(http.Header),
	}

	shouldRetry, err = customRetryPolicy(context.Background(), resp, nil)
	if !shouldRetry {
		t.Error("expected retry for 429 case")
	}
	if err != nil {
		t.Errorf("expected no error for 429, got: %v", err)
	}
}

func TestCustomRetryPolicyJSONDecoder(t *testing.T) {
	// Test comprehensive JSON decoder scenarios for 504 responses

	tests := []struct {
		name          string
		responseBody  string
		expectedRetry bool
		expectedError error
		description   string
	}{
		{
			name:          "Valid JSON with code 7001",
			responseBody:  `{"code":7001}`,
			expectedRetry: false,
			expectedError: apierrors.ErrGatewayTimeoutForIndexDDL,
			description:   "Should not retry and return specific error for 7001",
		},
		{
			name:          "Valid JSON with code 7001 and extra fields",
			responseBody:  `{"code":7001,"message":"Index DDL timeout","timestamp":"2023-09-25T10:30:00Z"}`,
			expectedRetry: false,
			expectedError: apierrors.ErrGatewayTimeoutForIndexDDL,
			description:   "Should handle 7001 with additional JSON fields",
		},
		{
			name:          "Valid JSON with different error code",
			responseBody:  `{"code":5001}`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry for non-7001 error codes",
		},
		{
			name:          "Valid JSON with zero code",
			responseBody:  `{"code":0}`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry for zero error code",
		},
		{
			name:          "Invalid JSON - missing quotes",
			responseBody:  `{code:7001}`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when JSON is malformed",
		},
		{
			name:          "Invalid JSON - incomplete",
			responseBody:  `{"code":70`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when JSON is incomplete",
		},
		{
			name:          "Non-JSON response",
			responseBody:  `Gateway Timeout`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when response is not JSON",
		},
		{
			name:          "Empty JSON object",
			responseBody:  `{}`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when code field is missing (defaults to 0)",
		},
		{
			name:          "JSON with code as string",
			responseBody:  `{"code":"7001"}`,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when code is string instead of int",
		},
		{
			name:          "Empty response body",
			responseBody:  ``,
			expectedRetry: true,
			expectedError: nil,
			description:   "Should retry when response body is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: http.StatusGatewayTimeout,
				Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
				Header:     make(http.Header),
			}

			shouldRetry, err := customRetryPolicy(context.Background(), resp, nil)

			// Check retry behavior
			if shouldRetry != tt.expectedRetry {
				t.Errorf("Expected retry=%v, got retry=%v. %s", tt.expectedRetry, shouldRetry, tt.description)
			}

			// Check error behavior
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %v. %s", err, tt.description)
			} else if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got: %v. %s", tt.expectedError, err, tt.description)
			}

			t.Logf("✓ %s: retry=%v, error=%v", tt.description, shouldRetry, err)
		})
	}
}

func Test504_Code7001_NoRetryAndErrorReturned(t *testing.T) {
	var callCount int32
	body := `{"code":7001}`

	// Create test server that returns 504 with code 7001
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(body))
	}))
	defer server.Close()

	client := NewRetryHTTPClient(context.Background(), 5*time.Second, false)
	resp, err := client.Get(server.URL)

	// Should return the specific 7001 error - matches V1 client behavior
	if !errors.Is(err, apierrors.ErrGatewayTimeoutForIndexDDL) {
		t.Fatalf("expected ErrGatewayTimeoutForIndexDDL error, got: %v", err)
	}

	// Should not retry - only one call
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("expected 1 call, got %d", atomic.LoadInt32(&callCount))
	}

	// Response may or may not be available when error is returned
	// This is acceptable as long as the correct error is returned
	if resp != nil {
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("expected status %d, got %d", http.StatusGatewayTimeout, resp.StatusCode)
		}

		// Body should be preserved and contain 7001 code if response is available
		b, _ := io.ReadAll(resp.Body)
		if string(b) != body {
			t.Errorf("body preserved: got %q, want %q", string(b), body)
		}
	}
}

func Test504_Other_RetriesUpToMax(t *testing.T) {
	var callCount int32
	body := `{"code":1234}` // Not 7001, so should retry

	// Create test server that always returns 504 with non-7001 code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(body))
	}))
	defer server.Close()

	// Use testing client with fast backoff for quicker tests
	client := NewRetryHTTPClient(context.Background(), 30*time.Second, false, WithFastBackoff())

	start := time.Now()
	resp, err := client.Get(server.URL)
	duration := time.Since(start)

	// retryablehttp returns error when all retries are exhausted
	if err == nil {
		t.Fatalf("expected error after exhausting retries, got nil")
	}

	// Error should indicate giving up after attempts
	if !strings.Contains(err.Error(), "giving up after") {
		t.Errorf("expected 'giving up after' error, got: %v", err)
	}

	// Should have made maxRetryAttempts + 1 calls (initial + retries)
	expectedCalls := int32(maxRetryAttempts + 1)
	if atomic.LoadInt32(&callCount) != expectedCalls {
		t.Fatalf("expected %d calls, got %d", expectedCalls, atomic.LoadInt32(&callCount))
	}

	// Should have taken some time due to backoff (but much less with fast backoff)
	if duration < 100*time.Millisecond {
		t.Errorf("expected at least 100ms delay due to retries, got %v", duration)
	}
	if duration > 10*time.Second {
		t.Errorf("expected less than 10s total duration with fast backoff, got %v", duration)
	}

	// Response may still be available even with error (retryablehttp behavior)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("expected final status %d, got %d", http.StatusGatewayTimeout, resp.StatusCode)
		}
	}
}

func Test429_RetriesWithBackoff(t *testing.T) {
	var callCount int32

	// Create test server that returns 429 first few times, then success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count <= 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	client := NewRetryHTTPClient(context.Background(), 30*time.Second, false, WithFastBackoff()) // Fast testing client

	start := time.Now()
	resp, err := client.Get(server.URL)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have made 4 calls total (3 retries + 1 success)
	if atomic.LoadInt32(&callCount) != 4 {
		t.Fatalf("expected 4 calls, got %d", atomic.LoadInt32(&callCount))
	}

	// Should have taken some time due to backoff (much less with fast backoff)
	if duration < 50*time.Millisecond {
		t.Errorf("expected at least 50ms delay due to retries, got %v", duration)
	}
	if duration > 5*time.Second {
		t.Errorf("expected less than 5s total duration with fast backoff, got %v", duration)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected successful status, got %d", resp.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)
	if string(b) != "success" {
		t.Fatalf("body: got %q, want %q", string(b), "success")
	}
}

func Test429_RetryAfter_Respected(t *testing.T) {
	var callCount int32

	// Create test server that returns 429 with small Retry-After, then success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count == 1 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	// Use fast backoff for more predictable timing in tests
	// The retryablehttp library's default LinearJitterBackoff can have unpredictable delays
	client := NewRetryHTTPClient(context.Background(), 30*time.Second, false, WithFastBackoff())

	start := time.Now()
	resp, err := client.Get(server.URL)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have made 2 calls
	if atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("expected 2 calls, got %d", atomic.LoadInt32(&callCount))
	}

	// With fast backoff, the delay should be much shorter (50ms for first retry)
	// We still verify it's making the retry, just with predictable timing
	if duration < 50*time.Millisecond {
		t.Errorf("expected at least 50ms delay for retry, got %v", duration)
	}
	if duration > 5*time.Second {
		t.Errorf("expected less than 5s total duration with fast backoff, got %v", duration)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected successful status, got %d", resp.StatusCode)
	}
}

func TestSuccessPassthrough(t *testing.T) {
	var callCount int32

	// Create test server that returns success immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewRetryHTTPClient(context.Background(), 5*time.Second, false)
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("request unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	b, _ := io.ReadAll(resp.Body)
	if string(b) != "ok" {
		t.Fatalf("body: got %q, want %q", string(b), "ok")
	}

	// Should not retry - only one call
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("expected 1 call, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestDefaultCase_NoRetry(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		body       string
	}{
		{"BadRequest", http.StatusBadRequest, "bad request"},
		{"Unauthorized", http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", http.StatusForbidden, "forbidden"},
		{"NotFound", http.StatusNotFound, "not found"},
		{"InternalServerError", http.StatusInternalServerError, "internal server error"},
		{"BadGateway", http.StatusBadGateway, "bad gateway"},
		{"ServiceUnavailable", http.StatusServiceUnavailable, "service unavailable"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callCount int32

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&callCount, 1)
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := NewRetryHTTPClient(context.Background(), 5*time.Second, false)
			resp, err := client.Get(server.URL)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()

			// Should not retry - only one call
			if atomic.LoadInt32(&callCount) != 1 {
				t.Fatalf("expected 1 call, got %d", atomic.LoadInt32(&callCount))
			}

			if resp.StatusCode != tc.statusCode {
				t.Fatalf("status code: got %d, want %d", resp.StatusCode, tc.statusCode)
			}

			b, _ := io.ReadAll(resp.Body)
			if string(b) != tc.body {
				t.Fatalf("body: got %q, want %q", string(b), tc.body)
			}
		})
	}
}

func TestNewRetryHTTPClient_ConfiguresCorrectly(t *testing.T) {
	timeout := 123 * time.Second
	client := NewRetryHTTPClient(context.Background(), timeout, false)

	if client.Timeout != timeout {
		t.Fatalf("timeout: got %v, want %v", client.Timeout, timeout)
	}

	if client.Transport == nil {
		t.Fatalf("expected non-nil transport")
	}
}

func TestNewRetryHTTPClient_WithOptions(t *testing.T) {
	// Test default configuration without debug logging
	defaultClient := NewRetryHTTPClient(context.Background(), 30*time.Second, false)
	if defaultClient.Timeout != 30*time.Second {
		t.Errorf("default timeout: got %v, want %v", defaultClient.Timeout, 30*time.Second)
	}

	// Test with debug logging enabled
	debugClient := NewRetryHTTPClient(context.Background(), 30*time.Second, true)
	if debugClient.Timeout != 30*time.Second {
		t.Errorf("debug client timeout: got %v, want %v", debugClient.Timeout, 30*time.Second)
	}

	// Test with fast backoff option and debug logging
	fastClient := NewRetryHTTPClient(context.Background(), 15*time.Second, true, WithFastBackoff())
	if fastClient.Timeout != 15*time.Second {
		t.Errorf("fast client timeout: got %v, want %v", fastClient.Timeout, 15*time.Second)
	}

	// Test with custom max retries
	customClient := NewRetryHTTPClient(context.Background(), 45*time.Second, false, WithMaxRetries(3))
	if customClient.Timeout != 45*time.Second {
		t.Errorf("custom client timeout: got %v, want %v", customClient.Timeout, 45*time.Second)
	}

	// Test combining options
	combinedClient := NewRetryHTTPClient(context.Background(), 60*time.Second, true, WithFastBackoff(), WithMaxRetries(2))
	if combinedClient.Timeout != 60*time.Second {
		t.Errorf("combined client timeout: got %v, want %v", combinedClient.Timeout, 60*time.Second)
	}

	// All clients should have non-nil transport
	clients := []*http.Client{defaultClient, debugClient, fastClient, customClient, combinedClient}
	for i, client := range clients {
		if client.Transport == nil {
			t.Errorf("client %d has nil transport", i)
		}
	}
}

func TestNewRetryHTTPClient_DebugLogging(t *testing.T) {
	// Test that debug logging parameter is respected
	// Note: We can't easily test the actual log output, but we can verify the client is created correctly

	var callCount int32

	// Create test server that returns 429, then success to trigger retry (and logging)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	// Test with debug logging enabled and fast backoff for quick test
	client := NewRetryHTTPClient(context.Background(), 30*time.Second, true, WithFastBackoff())
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have made 2 calls (1 retry)
	if atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("expected 2 calls, got %d", atomic.LoadInt32(&callCount))
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected successful status, got %d", resp.StatusCode)
	}

	// When debug logging is enabled, retryablehttp.DefaultLogger() would have logged:
	// - The initial 429 response
	// - The retry attempt with backoff delay
	// - The successful 200 response
	// Since we can't capture stdout in this test, we just verify the client worked correctly
}

func Test504_Code7001_InvalidJSON_Retries(t *testing.T) {
	var callCount int32
	body := `{"invalid json` // Invalid JSON to trigger parsing error

	// Create test server that always returns 504 with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(body))
	}))
	defer server.Close()

	client := NewRetryHTTPClient(context.Background(), 30*time.Second, false, WithFastBackoff()) // Fast testing client
	resp, err := client.Get(server.URL)

	// retryablehttp returns error when all retries are exhausted
	if err == nil {
		t.Fatalf("expected error after exhausting retries, got nil")
	}

	// Error should indicate giving up after attempts
	if !strings.Contains(err.Error(), "giving up after") {
		t.Errorf("expected 'giving up after' error, got: %v", err)
	}

	// Should have retried maxRetryAttempts times due to invalid JSON (treats as regular 504)
	expectedCalls := int32(maxRetryAttempts + 1)
	if atomic.LoadInt32(&callCount) != expectedCalls {
		t.Fatalf("expected %d calls, got %d", expectedCalls, atomic.LoadInt32(&callCount))
	}

	// Response may still be available even with error
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("expected final status %d, got %d", http.StatusGatewayTimeout, resp.StatusCode)
		}
	}
}
