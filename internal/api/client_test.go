package api

import (
	"context"
	goer "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// testPolicy is a retryPolicy whose backoff is short enough that the retry loop
// can be exercised without waiting for the production budget. A maxServiceErrors
// of zero reproduces the default, uncapped policy.
func testPolicy(maxServiceErrors int) retryPolicy {
	return retryPolicy{
		backoff:          func(int) time.Duration { return time.Millisecond },
		maxServiceErrors: maxServiceErrors,
	}
}

// scriptedFn returns a fn for exec that yields each error in errs in turn and
// succeeds once they are exhausted. calls reports how many times it was invoked.
func scriptedFn(errs []error, calls *int) func() (*Response, time.Duration, error) {
	return func() (*Response, time.Duration, error) {
		i := *calls
		*calls++
		if i < len(errs) {
			return nil, 0, errs[i]
		}
		return &Response{}, 0, nil
	}
}

func repeat(err error, n int) []error {
	errs := make([]error, n)
	for i := range errs {
		errs[i] = err
	}
	return errs
}

func TestExecRetryBudget(t *testing.T) {
	const maxServiceErrors = 5

	nonRetryable := &Error{Code: 4000, Message: "bad request", HttpStatusCode: http.StatusBadRequest}

	tests := []struct {
		name string
		// errs is returned by successive calls to fn; fn succeeds afterwards.
		errs []error
		// wantErrIs, when set, must match the returned error via errors.Is.
		wantErrIs error
		// wantCalls is the exact number of requests the loop should issue.
		wantCalls int
		// wantExhausted asserts the returned error is a *RetryExhaustedError.
		wantExhausted bool
	}{
		{
			name:      "success on first attempt",
			errs:      nil,
			wantCalls: 1,
		},
		{
			name:      "service unavailable then success",
			errs:      repeat(errors.ErrServiceUnavailable, 2),
			wantCalls: 3,
		},
		{
			name:      "gateway timeout then success",
			errs:      repeat(errors.ErrGatewayTimeout, 2),
			wantCalls: 3,
		},
		{
			name:          "persistent service unavailable is capped",
			errs:          repeat(errors.ErrServiceUnavailable, 100),
			wantErrIs:     errors.ErrServiceUnavailable,
			wantCalls:     maxServiceErrors + 1,
			wantExhausted: true,
		},
		{
			name:          "persistent gateway timeout is capped",
			errs:          repeat(errors.ErrGatewayTimeout, 100),
			wantErrIs:     errors.ErrGatewayTimeout,
			wantCalls:     maxServiceErrors + 1,
			wantExhausted: true,
		},
		{
			name:          "mixed 5xx share a single budget",
			errs:          append(repeat(errors.ErrServiceUnavailable, 3), repeat(errors.ErrGatewayTimeout, 100)...),
			wantErrIs:     errors.ErrGatewayTimeout,
			wantCalls:     maxServiceErrors + 1,
			wantExhausted: true,
		},
		{
			name:      "non retryable error returns immediately",
			errs:      repeat(nonRetryable, 100),
			wantErrIs: nonRetryable,
			wantCalls: 1,
		},
		{
			name:      "index DDL gateway timeout is not retried",
			errs:      repeat(errors.ErrGatewayTimeoutForIndexDDL, 100),
			wantErrIs: errors.ErrGatewayTimeoutForIndexDDL,
			wantCalls: 1,
		},
		{
			// 429 carries a server-dictated Retry-After, so the attempt cap does
			// not apply to it; only the overall deadline bounds those retries.
			name:      "rate limit is not subject to the attempt cap",
			errs:      repeat(errors.ErrRatelimit, maxServiceErrors*3),
			wantCalls: maxServiceErrors*3 + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			_, err := exec(context.Background(), scriptedFn(tt.errs, &calls), testPolicy(maxServiceErrors))

			if calls != tt.wantCalls {
				t.Errorf("issued %d requests, want %d", calls, tt.wantCalls)
			}

			if tt.wantErrIs == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if !goer.Is(err, tt.wantErrIs) {
				t.Fatalf("got error %v, want it to match %v", err, tt.wantErrIs)
			}

			var exhausted *RetryExhaustedError
			if got := goer.As(err, &exhausted); got != tt.wantExhausted {
				t.Fatalf("errors.As(*RetryExhaustedError) = %v, want %v (err: %v)", got, tt.wantExhausted, err)
			}

			if tt.wantExhausted {
				if exhausted.Attempts != tt.wantCalls {
					t.Errorf("RetryExhaustedError.Attempts = %d, want %d", exhausted.Attempts, tt.wantCalls)
				}
				if exhausted.Elapsed <= 0 {
					t.Errorf("RetryExhaustedError.Elapsed = %v, want a positive duration", exhausted.Elapsed)
				}
			}
		})
	}
}

// The default policy must not cap 5xx retries. Every endpoint relies on retrying
// until the request deadline to wait out long transient conditions such as an
// in-flight bucket delete, so a cap only applies when a caller opts in.
func TestExecDefaultPolicyDoesNotCapServiceErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	var calls int
	_, err := exec(ctx, scriptedFn(repeat(errors.ErrServiceUnavailable, 10000), &calls), testPolicy(0))

	if calls <= FastFailServiceErrorRetries+1 {
		t.Errorf("issued %d requests, want more than the opt-in budget of %d: the default must not cap",
			calls, FastFailServiceErrorRetries+1)
	}

	var exhausted *RetryExhaustedError
	if !goer.As(err, &exhausted) {
		t.Fatalf("got error %v, want *RetryExhaustedError once the deadline expires", err)
	}
	if !goer.Is(err, errors.ErrServiceUnavailable) {
		t.Errorf("error %v does not unwrap to ErrServiceUnavailable", err)
	}
	if exhausted.Attempts != calls {
		t.Errorf("RetryExhaustedError.Attempts = %d, want %d", exhausted.Attempts, calls)
	}
}

func TestExecHonoursReturnedBackoff(t *testing.T) {
	const backoff = 120 * time.Millisecond

	var calls int
	fn := func() (*Response, time.Duration, error) {
		calls++
		if calls == 1 {
			return nil, backoff, errors.ErrRatelimit
		}
		return &Response{}, 0, nil
	}

	// The policy backoff is 1ms, so anything close to backoff can only have come
	// from the duration fn returned.
	start := time.Now()
	if _, err := exec(context.Background(), fn, testPolicy(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed := time.Since(start); elapsed < backoff {
		t.Errorf("waited %v before retrying, want at least %v", elapsed, backoff)
	}
}

func TestExecRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var calls int
	_, err := exec(ctx, scriptedFn(repeat(errors.ErrServiceUnavailable, 100), &calls), testPolicy(5))

	if !goer.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want context.Canceled", err)
	}
}

func TestExponentialJitterBackoff(t *testing.T) {
	for attempt := 1; attempt <= 8; attempt++ {
		ceiling := retryWaitMax
		if shifted := retryWaitMin << (attempt - 1); shifted > 0 && shifted < retryWaitMax {
			ceiling = shifted
		}

		// Assert bounds rather than exact values: the delay is jittered, so an
		// equality assertion here would be flaky by construction.
		for i := 0; i < 100; i++ {
			got := exponentialJitterBackoff(attempt)
			if got < ceiling/2 || got > ceiling {
				t.Fatalf("attempt %d returned %v, want within [%v, %v]", attempt, got, ceiling/2, ceiling)
			}
		}
	}
}

func TestExponentialJitterBackoffIsJittered(t *testing.T) {
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 100; i++ {
		seen[exponentialJitterBackoff(5)] = struct{}{}
	}

	if len(seen) < 2 {
		t.Errorf("100 samples produced %d distinct delay(s), want jittered values", len(seen))
	}
}

func TestExponentialJitterBackoffFloorsAttempt(t *testing.T) {
	// serviceErrors is zero when a 429 asks for no wait at all; that must not
	// shift by a negative amount or collapse the delay to zero.
	for _, attempt := range []int{-1, 0} {
		got := exponentialJitterBackoff(attempt)
		if got < retryWaitMin/2 || got > retryWaitMin {
			t.Errorf("attempt %d returned %v, want within [%v, %v]", attempt, got, retryWaitMin/2, retryWaitMin)
		}
	}
}

// retryServer returns an httptest server whose handler is invoked for each
// request, along with a counter of the requests it received.
func retryServer(t *testing.T, handler func(w http.ResponseWriter, requests int32)) (*httptest.Server, *int32) {
	t.Helper()

	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handler(w, atomic.AddInt32(&requests, 1))
	}))
	t.Cleanup(server.Close)

	return server, &requests
}

func alwaysStatus(status int, body string) func(http.ResponseWriter, int32) {
	return func(w http.ResponseWriter, _ int32) {
		w.WriteHeader(status)
		if body != "" {
			_, _ = fmt.Fprint(w, body)
		}
	}
}

func TestExecuteWithRetryServiceUnavailable(t *testing.T) {
	t.Run("retries then succeeds", func(t *testing.T) {
		server, requests := retryServer(t, func(w http.ResponseWriter, requests int32) {
			if requests < 3 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		cfg := EndpointCfg{Url: server.URL, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		if _, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := atomic.LoadInt32(requests); got != 3 {
			t.Errorf("issued %d requests, want 3", got)
		}
	})

	t.Run("default budget is bounded only by the deadline", func(t *testing.T) {
		server, requests := retryServer(t, alwaysStatus(http.StatusServiceUnavailable, ""))

		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		// No MaxServiceErrorRetries: this is what every existing call site gets.
		cfg := EndpointCfg{Url: server.URL, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
		_, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(ctx, cfg, nil, "token", nil)

		var exhausted *RetryExhaustedError
		if !goer.As(err, &exhausted) {
			t.Fatalf("got error %v, want *RetryExhaustedError", err)
		}
		if got := atomic.LoadInt32(requests); got <= FastFailServiceErrorRetries+1 {
			t.Errorf("issued %d requests, want more than %d: the default must not cap retries",
				got, FastFailServiceErrorRetries+1)
		}
	})

	t.Run("opt-in budget is honoured", func(t *testing.T) {
		server, requests := retryServer(t, alwaysStatus(http.StatusServiceUnavailable, ""))

		cfg := EndpointCfg{
			Url:                    server.URL,
			Method:                 http.MethodPut,
			SuccessStatus:          http.StatusNoContent,
			MaxServiceErrorRetries: FastFailServiceErrorRetries,
		}
		_, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil)

		var exhausted *RetryExhaustedError
		if !goer.As(err, &exhausted) {
			t.Fatalf("got error %v, want *RetryExhaustedError", err)
		}
		if !goer.Is(err, errors.ErrServiceUnavailable) {
			t.Errorf("error %v does not unwrap to ErrServiceUnavailable", err)
		}
		if exhausted.Attempts != FastFailServiceErrorRetries+1 {
			t.Errorf("Attempts = %d, want %d", exhausted.Attempts, FastFailServiceErrorRetries+1)
		}
		if got := atomic.LoadInt32(requests); got != FastFailServiceErrorRetries+1 {
			t.Errorf("issued %d requests, want %d", got, FastFailServiceErrorRetries+1)
		}
		if msg := ParseError(err); msg != err.Error() {
			t.Errorf("ParseError rendered %q, want the RetryExhaustedError message %q", msg, err.Error())
		}
	})

	t.Run("a custom budget of one still retries once", func(t *testing.T) {
		server, requests := retryServer(t, alwaysStatus(http.StatusServiceUnavailable, ""))

		cfg := EndpointCfg{
			Url:                    server.URL,
			Method:                 http.MethodGet,
			SuccessStatus:          http.StatusOK,
			MaxServiceErrorRetries: 1,
		}
		if _, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil); err == nil {
			t.Fatal("expected the retry budget to be exhausted")
		}

		if got := atomic.LoadInt32(requests); got != 2 {
			t.Errorf("issued %d requests, want 2", got)
		}
	})
}

func TestExecuteWithRetryGatewayTimeout(t *testing.T) {
	const timeoutBody = `{"code":4000,"hint":"retry later","httpStatusCode":504,"message":"gateway timeout"}`

	// A 504 is retried for every method, including non-idempotent ones. Whether a
	// POST should be replayed after a gateway timeout is a separate question.
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
		t.Run(method+" is retried", func(t *testing.T) {
			server, requests := retryServer(t, alwaysStatus(http.StatusGatewayTimeout, timeoutBody))

			cfg := EndpointCfg{
				Url:                    server.URL,
				Method:                 method,
				SuccessStatus:          http.StatusOK,
				MaxServiceErrorRetries: FastFailServiceErrorRetries,
			}
			_, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil)

			if !goer.Is(err, errors.ErrGatewayTimeout) {
				t.Fatalf("got error %v, want it to match ErrGatewayTimeout", err)
			}
			if got := atomic.LoadInt32(requests); got != FastFailServiceErrorRetries+1 {
				t.Errorf("issued %d requests, want %d", got, FastFailServiceErrorRetries+1)
			}
		})
	}

	t.Run("index DDL code 7001 is not retried", func(t *testing.T) {
		server, requests := retryServer(t, alwaysStatus(http.StatusGatewayTimeout, `{"code":7001,"message":"index DDL timed out"}`))

		cfg := EndpointCfg{Url: server.URL, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		_, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil)

		if !goer.Is(err, errors.ErrGatewayTimeoutForIndexDDL) {
			t.Fatalf("got error %v, want ErrGatewayTimeoutForIndexDDL", err)
		}
		if got := atomic.LoadInt32(requests); got != 1 {
			t.Errorf("issued %d requests, want 1", got)
		}
	})
}

func TestExecuteWithRetryRateLimit(t *testing.T) {
	t.Run("Retry-After is honoured", func(t *testing.T) {
		server, _ := retryServer(t, func(w http.ResponseWriter, requests int32) {
			if requests == 1 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		cfg := EndpointCfg{Url: server.URL, Method: http.MethodGet, SuccessStatus: http.StatusOK}

		start := time.Now()
		if _, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// WithFastBackoff would have retried in ~1ms, so a wait of a second can
		// only have come from the Retry-After header.
		if elapsed := time.Since(start); elapsed < time.Second {
			t.Errorf("retried after %v, want at least the 1s the server asked for", elapsed)
		}
	})

	t.Run("is not subject to an opt-in 5xx budget", func(t *testing.T) {
		const rateLimited = FastFailServiceErrorRetries * 3

		server, requests := retryServer(t, func(w http.ResponseWriter, requests int32) {
			if requests <= rateLimited {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		cfg := EndpointCfg{
			Url:                    server.URL,
			Method:                 http.MethodGet,
			SuccessStatus:          http.StatusOK,
			MaxServiceErrorRetries: FastFailServiceErrorRetries,
		}
		if _, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := atomic.LoadInt32(requests); got != rateLimited+1 {
			t.Errorf("issued %d requests, want %d", got, rateLimited+1)
		}
	})
}

func TestExecuteWithRetrySuccessPassthrough(t *testing.T) {
	server, requests := retryServer(t, alwaysStatus(http.StatusOK, `{"id":"abc"}`))

	cfg := EndpointCfg{Url: server.URL, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := NewClient(time.Minute, WithFastBackoff()).ExecuteWithRetry(context.Background(), cfg, nil, "token", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(response.Body) != `{"id":"abc"}` {
		t.Errorf("body = %q, want the response body verbatim", response.Body)
	}
	if got := atomic.LoadInt32(requests); got != 1 {
		t.Errorf("issued %d requests, want 1", got)
	}
}
