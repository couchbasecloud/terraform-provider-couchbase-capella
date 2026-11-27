package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		header string
		want   time.Duration
	}{
		{name: "missing", header: "", want: 0},
		{name: "delay seconds", header: "3", want: 3 * time.Second},
		{name: "delay seconds with whitespace", header: " 5 ", want: 5 * time.Second},
		{name: "negative seconds", header: "-1", want: 0},
		{name: "http date in future", header: now.Add(7 * time.Second).Format(http.TimeFormat), want: 7 * time.Second},
		{name: "http date in past", header: now.Add(-time.Minute).Format(http.TimeFormat), want: 0},
		{name: "garbage", header: "soon", want: 0},
		{name: "fractional seconds", header: "1.5", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRetryAfter(tt.header, now); got != tt.want {
				t.Fatalf("parseRetryAfter(%q) = %s, want %s", tt.header, got, tt.want)
			}
		})
	}
}

// A 429 whose Retry-After is missing or not an integer must be retried with
// the default back-off rather than failing the request outright.
func TestExecuteWithRetry_429UnparseableRetryAfterIsRetried(t *testing.T) {
	for _, header := range []string{"", "not-a-number"} {
		t.Run("header="+header, func(t *testing.T) {
			var calls atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if calls.Add(1) == 1 {
					if header != "" {
						w.Header().Set("Retry-After", header)
					}
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"id":"x"}`))
			}))
			defer srv.Close()

			c := NewClient(10 * time.Second)
			resp, err := c.ExecuteWithRetry(
				context.Background(),
				EndpointCfg{Url: srv.URL, Method: http.MethodPost, SuccessStatus: http.StatusAccepted},
				map[string]int{"retention": 168},
				"token",
				nil,
			)
			if err != nil {
				t.Fatalf("ExecuteWithRetry returned error: %v", err)
			}
			if resp == nil || resp.Response.StatusCode != http.StatusAccepted {
				t.Fatalf("unexpected response: %+v", resp)
			}
			if got := calls.Load(); got != 2 {
				t.Fatalf("server called %d times, want 2", got)
			}
		})
	}
}
