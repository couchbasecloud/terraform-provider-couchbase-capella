package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

// roundTripFunc allows stubbing http.RoundTripper in tests.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// trackingBody tracks reads and close calls.
type trackingBody struct {
	io.Reader
	closed    bool
	readBytes int
}

func newTrackingBody(s string) *trackingBody {
	return &trackingBody{Reader: bytes.NewBufferString(s)}
}

func (tb *trackingBody) Read(p []byte) (int, error) {
	n, err := tb.Reader.Read(p)
	if n > 0 {
		tb.readBytes += n
	}
	return n, err
}

func (tb *trackingBody) Close() error {
	tb.closed = true
	return nil
}

func TestRoundTrip_SuccessPassthrough(t *testing.T) {
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ok")),
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want %d", res.StatusCode, http.StatusOK)
	}

	b, _ := io.ReadAll(res.Body)
	if string(b) != "ok" {
		t.Fatalf("body: got %q, want %q", string(b), "ok")
	}
}

func TestRoundTrip_504_Code7001_NoRetryAndBodyPreserved(t *testing.T) {
	var calls int
	body := `{"code":7001}`
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusGatewayTimeout,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	if calls != 1 {
		t.Fatalf("round trip calls: got %d, want 1", calls)
	}

	b, _ := io.ReadAll(res.Body)
	if string(b) != body {
		t.Fatalf("body preserved: got %q, want %q", string(b), body)
	}
}

func TestRoundTrip_504_Other_CancelContext_ReturnsErrorAndBodyPreserved(t *testing.T) {
	body := `{"code":1234}`
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusGatewayTimeout,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		}, nil
	})}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err == nil {
		t.Fatalf("expected error due to context timeout, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context deadline exceeded or canceled, got %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil response even when error occurs")
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	if string(b) != body {
		t.Fatalf("body preserved: got %q, want %q", string(b), body)
	}
}

func TestRoundTrip_429_RetryAfter_DrainsAndClosesAndCancel(t *testing.T) {
	tb := newTrackingBody("some content to drain")
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		res := &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       tb,
			Header:     make(http.Header),
		}
		res.Header.Set("Retry-After", "1")
		return res, nil
	})}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err == nil {
		t.Fatalf("expected error due to context timeout, got nil")
	}
	if res == nil {
		t.Fatalf("expected response even when error occurs")
	}

	if !tb.closed {
		t.Fatalf("expected body to be closed after draining on 429")
	}
	if tb.readBytes == 0 {
		t.Fatalf("expected body to be drained on 429")
	}
}

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		minWant time.Duration
		maxWant time.Duration
	}{
		{"first retry", 0, 750 * time.Millisecond, 1250 * time.Millisecond},   // 1s ± 25%
		{"second retry", 1, 1500 * time.Millisecond, 2500 * time.Millisecond}, // 2s ± 25%
		{"third retry", 2, 3 * time.Second, 5 * time.Second},                   // 4s ± 25%
		{"fourth retry", 3, 6 * time.Second, 10 * time.Second},                 // 8s ± 25%
		{"fifth retry", 4, 12 * time.Second, 20 * time.Second},                 // 16s ± 25%
		{"max delay cap", 10, 22500 * time.Millisecond, 37500 * time.Millisecond}, // Should cap at 30s ± 25%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backoff := calculateBackoff(tt.attempt)
			if backoff < tt.minWant || backoff > tt.maxWant {
				t.Errorf("calculateBackoff(%d) = %v, want between %v and %v", 
					tt.attempt, backoff, tt.minWant, tt.maxWant)
			}
		})
	}
}

func TestRoundTrip_429_NoRetryAfter_UsesExponentialBackoff(t *testing.T) {
	var calls int64
	var callTimes []time.Time
	
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt64(&calls, 1)
		callTimes = append(callTimes, time.Now())
		
		// Always return 429 - after maxRetryAttempts retries, it should stop trying
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(bytes.NewBufferString("rate limited")),
			Header:     make(http.Header), // No Retry-After header
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	start := time.Now()
	res, err := transport.RoundTrip(req)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	finalCalls := atomic.LoadInt64(&calls)
	if finalCalls != maxRetryAttempts+1 {
		t.Fatalf("expected %d calls (max retries + 1), got %d", maxRetryAttempts+1, finalCalls)
	}

	// Should have taken at least some time due to backoff
	minExpectedDuration := calculateBackoff(0) // At least the first backoff
	if duration < minExpectedDuration {
		t.Errorf("expected at least %v total duration, got %v", minExpectedDuration, duration)
	}

	if res.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected final status %d after max retries, got %d", http.StatusTooManyRequests, res.StatusCode)
	}
}

func TestRoundTrip_504_UsesExponentialBackoff(t *testing.T) {
	var calls int64
	
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt64(&calls, 1)
		
		// Always return 504 without special error code
		return &http.Response{
			StatusCode: http.StatusGatewayTimeout,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code":1234}`)), // Not 7001
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	start := time.Now()
	res, err := transport.RoundTrip(req)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	finalCalls := atomic.LoadInt64(&calls)
	if finalCalls != maxRetryAttempts+1 {
		t.Fatalf("expected %d calls (max retries + 1), got %d", maxRetryAttempts+1, finalCalls)
	}

	// Should have taken at least some time due to backoff
	minExpectedDuration := calculateBackoff(0) // At least the first backoff
	if duration < minExpectedDuration {
		t.Errorf("expected at least %v total duration, got %v", minExpectedDuration, duration)
	}

	if res.StatusCode != http.StatusGatewayTimeout {
		t.Errorf("expected final status %d after max retries, got %d", http.StatusGatewayTimeout, res.StatusCode)
	}
}

func TestRoundTrip_RetryLimitEnforced_429(t *testing.T) {
	var calls int64
	
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt64(&calls, 1)
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(bytes.NewBufferString("rate limited")),
			Header:     make(http.Header),
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	finalCalls := atomic.LoadInt64(&calls)
	expectedCalls := int64(maxRetryAttempts + 1) // Initial call + max retries
	if finalCalls != expectedCalls {
		t.Fatalf("expected exactly %d calls, got %d", expectedCalls, finalCalls)
	}
}

func TestRoundTrip_RetryLimitEnforced_504(t *testing.T) {
	var calls int64
	
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt64(&calls, 1)
		return &http.Response{
			StatusCode: http.StatusGatewayTimeout,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code":1234}`)),
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	res, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	finalCalls := atomic.LoadInt64(&calls)
	expectedCalls := int64(maxRetryAttempts + 1) // Initial call + max retries
	if finalCalls != expectedCalls {
		t.Fatalf("expected exactly %d calls, got %d", expectedCalls, finalCalls)
	}
}

func TestRoundTrip_429_RetryAfter_StillRespected(t *testing.T) {
	var calls int64
	var callTimes []time.Time
	
	transport := &RetryTransport{Base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt64(&calls, 1)
		callTimes = append(callTimes, time.Now())
		
		if atomic.LoadInt64(&calls) == 1 {
			res := &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(bytes.NewBufferString("rate limited")),
				Header:     make(http.Header),
			}
			res.Header.Set("Retry-After", "1") // 1 second
			return res, nil
		}
		// Second call succeeds
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ok")),
		}, nil
	})}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	start := time.Now()
	res, err := transport.RoundTrip(req)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("round trip unexpected error: %v", err)
	}
	defer res.Body.Close()

	finalCalls := atomic.LoadInt64(&calls)
	if finalCalls != 2 {
		t.Fatalf("expected 2 calls, got %d", finalCalls)
	}

	// Should have waited approximately 1 second as specified by Retry-After
	if duration < 900*time.Millisecond || duration > 1200*time.Millisecond {
		t.Errorf("expected ~1s delay due to Retry-After, got %v", duration)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected successful status, got %d", res.StatusCode)
	}
}

func TestCalculateBackoff_ExponentialGrowth(t *testing.T) {
	// Test multiple samples to ensure exponential growth pattern overall
	samples := 10
	attempts := []int{0, 1, 2, 3, 4}
	
	for _, attempt := range attempts {
		delays := make([]time.Duration, samples)
		
		// Collect multiple samples for this attempt
		for i := 0; i < samples; i++ {
			delays[i] = calculateBackoff(attempt)
		}
		
		// Calculate average delay for this attempt
		var total time.Duration
		for _, d := range delays {
			total += d
		}
		avgDelay := total / time.Duration(samples)
		
		// Expected base delay without jitter: baseBackoffDelay * 2^attempt
		multiplier := 1 << uint(attempt)
		expectedBase := time.Duration(float64(baseBackoffDelay) * float64(multiplier))
		if expectedBase > maxBackoffDelay {
			expectedBase = maxBackoffDelay
		}
		
		// Average should be close to expected base (within reasonable bounds considering jitter)
		tolerance := expectedBase / 4 // 25% tolerance
		if avgDelay < expectedBase-tolerance || avgDelay > expectedBase+tolerance {
			t.Errorf("attempt %d: average delay %v should be close to expected %v (±%v)", 
				attempt, avgDelay, expectedBase, tolerance)
		}
		
		// All individual delays should be positive and within reasonable bounds
		for i, delay := range delays {
			if delay <= 0 {
				t.Errorf("attempt %d sample %d: delay %v should be positive", attempt, i, delay)
			}
			
			// Should not exceed max + jitter tolerance
			maxWithJitter := maxBackoffDelay + maxBackoffDelay/4
			if delay > maxWithJitter {
				t.Errorf("attempt %d sample %d: delay %v exceeds maximum %v", 
					attempt, i, delay, maxWithJitter)
			}
		}
	}
}

func TestNewRetryHTTPClient_ConfiguresTransportAndTimeout(t *testing.T) {
	tout := 123 * time.Second
	cli := NewRetryHTTPClient(tout)

	if cli.Timeout != tout {
		t.Fatalf("timeout: got %v, want %v", cli.Timeout, tout)
	}

	rt, ok := cli.Transport.(*RetryTransport)
	if !ok {
		t.Fatalf("transport type: got %T, want *RetryTransport", cli.Transport)
	}
	if rt.Base != http.DefaultTransport {
		t.Fatalf("base transport: got %v, want %v", rt.Base, http.DefaultTransport)
	}
}
