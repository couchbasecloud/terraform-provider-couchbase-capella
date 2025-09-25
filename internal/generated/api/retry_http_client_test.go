package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

	shouldRetry, err := customRetryPolicy(nil, resp, nil)
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

	shouldRetry, err = customRetryPolicy(nil, resp, nil)
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

	shouldRetry, err = customRetryPolicy(nil, resp, nil)
	if !shouldRetry {
		t.Error("expected retry for 429 case")
	}
	if err != nil {
		t.Errorf("expected no error for 429, got: %v", err)
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

	client := NewRetryHTTPClient(5 * time.Second)
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

	client := NewRetryHTTPClient(30 * time.Second) // Longer timeout for retries

	start := time.Now()
	resp, err := client.Get(server.URL)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have made maxRetryAttempts + 1 calls (initial + retries)
	expectedCalls := int32(maxRetryAttempts + 1)
	if atomic.LoadInt32(&callCount) != expectedCalls {
		t.Fatalf("expected %d calls, got %d", expectedCalls, atomic.LoadInt32(&callCount))
	}

	// Should have taken some time due to backoff
	if duration < 500*time.Millisecond {
		t.Errorf("expected some delay due to retries, got %v", duration)
	}

	if resp.StatusCode != http.StatusGatewayTimeout {
		t.Errorf("expected final status %d, got %d", http.StatusGatewayTimeout, resp.StatusCode)
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

	client := NewRetryHTTPClient(30 * time.Second)

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

	// Should have taken some time due to backoff
	if duration < 500*time.Millisecond {
		t.Errorf("expected some delay due to retries, got %v", duration)
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

	// Create test server that returns 429 with Retry-After, then success
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

	client := NewRetryHTTPClient(30 * time.Second)

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

	// Should have waited approximately 1 second as specified by Retry-After
	if duration < 900*time.Millisecond || duration > 1200*time.Millisecond {
		t.Errorf("expected ~1s delay due to Retry-After, got %v", duration)
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

	client := NewRetryHTTPClient(5 * time.Second)
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

			client := NewRetryHTTPClient(5 * time.Second)
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
	client := NewRetryHTTPClient(timeout)

	if client.Timeout != timeout {
		t.Fatalf("timeout: got %v, want %v", client.Timeout, timeout)
	}

	if client.Transport == nil {
		t.Fatalf("expected non-nil transport")
	}
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

	client := NewRetryHTTPClient(30 * time.Second) // Longer timeout for retries
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have retried maxRetryAttempts times due to invalid JSON (treats as regular 504)
	expectedCalls := int32(maxRetryAttempts + 1)
	if atomic.LoadInt32(&callCount) != expectedCalls {
		t.Fatalf("expected %d calls, got %d", expectedCalls, atomic.LoadInt32(&callCount))
	}
}
