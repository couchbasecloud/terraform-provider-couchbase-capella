package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
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
