package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

// RetryTransport implements simple retry logic for HTTP 429 and 504 responses.
// - 429: respects Retry-After header (seconds) if present
// - 504: retries with default backoff unless body contains code 7001.
type RetryTransport struct {
	Base http.RoundTripper
}

func (t *RetryTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// defaultWaitAttempt re-attempt http request after 2 seconds.
const defaultWaitAttempt = time.Second * 2

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the body so we can retry non-idempotent requests as well
	var bodyBytes []byte
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		_ = req.Body.Close()
		bodyBytes = b
	}

	makeReq := func(ctx context.Context) *http.Request {
		r := req.Clone(ctx)
		if bodyBytes != nil {
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		return r
	}

	backoff := defaultWaitAttempt

	for {
		r := makeReq(req.Context())
		res, err := t.base().RoundTrip(r)
		if err != nil {
			return nil, err
		}

		switch res.StatusCode {
		case http.StatusTooManyRequests:
			// Respect Retry-After header
			if ra := res.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil && secs > 0 {
					backoff = time.Second * time.Duration(secs)
				}
			}
			// Drain and retry
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()
			if err := sleepOrDone(req.Context(), backoff); err != nil {
				return res, err
			}
			continue
		case http.StatusGatewayTimeout:
			// Peek body to check for special code 7001 (do not retry)
			b, _ := io.ReadAll(res.Body)
			_ = res.Body.Close()
			// restore body for caller if we return
			res.Body = io.NopCloser(bytes.NewReader(b))

			var apiErr struct {
				Code int `json:"code"`
			}
			_ = json.Unmarshal(b, &apiErr)
			if apiErr.Code == 7001 {
				return res, nil
			}
			if err := sleepOrDone(req.Context(), backoff); err != nil {
				return res, err
			}
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
