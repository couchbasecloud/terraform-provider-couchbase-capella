package resources

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// transientErrorBody is the transient backend 500 Capella returns while it finishes
// propagating state after a write.
const transientErrorBody = `{"hint":"Something went wrong on our end. We are actively investigating the issue.",` +
	`"message":"An internal server error occurred.","code":10000,"httpStatusCode":500}`

// scriptedResponse is one canned HTTP reply in a queue served by newTransientRetryServer.
type scriptedResponse struct {
	status int
	body   string
}

// fastTransientRetry shrinks the package-level retry timing for the duration of a test
// so the retry loop resolves in milliseconds rather than minutes.
func fastTransientRetry(t *testing.T, window time.Duration) {
	t.Helper()
	origWindow, origInterval := transientRetryWindow, transientRetryInterval
	transientRetryWindow = window
	transientRetryInterval = time.Millisecond
	t.Cleanup(func() {
		transientRetryWindow = origWindow
		transientRetryInterval = origInterval
	})
}

// newTransientRetryServer serves responses in order, repeating the last one once the
// queue is drained, and counts every request it receives.
func newTransientRetryServer(t *testing.T, responses []scriptedResponse) (*httptest.Server, *int32) {
	t.Helper()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		idx := int(atomic.AddInt32(&calls, 1)) - 1
		if idx >= len(responses) {
			idx = len(responses) - 1
		}
		w.WriteHeader(responses[idx].status)
		if responses[idx].body != "" {
			_, _ = w.Write([]byte(responses[idx].body))
		}
	}))
	t.Cleanup(srv.Close)

	return srv, &calls
}

// callTransientRetry runs the helper against a server scripted with responses.
func callTransientRetry(t *testing.T, ctx context.Context, srv *httptest.Server) error {
	t.Helper()
	cfg := api.EndpointCfg{Url: srv.URL, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := executeWithTransientRetry(ctx, &api.Client{Client: srv.Client()}, "token", "test operation", cfg, nil)
	return err
}

func TestExecuteWithTransientRetry(t *testing.T) {
	const (
		notFoundBody         = `{"hint":"","message":"resource not found","code":4025,"httpStatusCode":404}`
		otherServerErrorBody = `{"hint":"","message":"boom","code":9999,"httpStatusCode":500}`
	)

	tests := []struct {
		name           string
		responses      []scriptedResponse
		wantCalls      int32
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "retries the transient 500 until the request succeeds",
			responses: []scriptedResponse{
				{http.StatusInternalServerError, transientErrorBody},
				{http.StatusInternalServerError, transientErrorBody},
				{http.StatusAccepted, ""},
			},
			wantCalls: 3,
		},
		{
			name:      "succeeds on the first attempt without retrying",
			responses: []scriptedResponse{{http.StatusAccepted, ""}},
			wantCalls: 1,
		},
		{
			name:           "surfaces a 500 with a different code immediately",
			responses:      []scriptedResponse{{http.StatusInternalServerError, otherServerErrorBody}},
			wantCalls:      1,
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "surfaces a 404 immediately so callers can treat it as not found",
			responses:      []scriptedResponse{{http.StatusNotFound, notFoundBody}},
			wantCalls:      1,
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fastTransientRetry(t, 2*time.Second)
			srv, calls := newTransientRetryServer(t, tt.responses)

			err := callTransientRetry(t, context.Background(), srv)

			assert.Equal(t, atomic.LoadInt32(calls), tt.wantCalls)

			if !tt.wantErr {
				assert.NilError(t, err)
				return
			}

			assert.Assert(t, err != nil)
			var apiErr *api.Error
			assert.Assert(t, stderrors.As(err, &apiErr))
			assert.Equal(t, apiErr.HttpStatusCode, tt.wantStatusCode)
		})
	}
}

func TestExecuteWithTransientRetryExhaustsWindow(t *testing.T) {
	fastTransientRetry(t, 20*time.Millisecond)
	srv, calls := newTransientRetryServer(t, []scriptedResponse{
		{http.StatusInternalServerError, transientErrorBody},
	})

	err := callTransientRetry(t, context.Background(), srv)

	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "retry window"), "got: %s", err)
	assert.Assert(t, strings.Contains(err.Error(), "test operation"), "got: %s", err)

	// The exhausted error wraps the last transient 500 rather than discarding it.
	var apiErr *api.Error
	assert.Assert(t, stderrors.As(err, &apiErr))
	assert.Equal(t, apiErr.Code, 10000)

	assert.Assert(t, atomic.LoadInt32(calls) > 1, "expected more than one attempt, got %d", atomic.LoadInt32(calls))
}

func TestExecuteWithTransientRetryStopsOnCanceledContext(t *testing.T) {
	fastTransientRetry(t, time.Minute)
	srv, _ := newTransientRetryServer(t, []scriptedResponse{
		{http.StatusInternalServerError, transientErrorBody},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() {
		done <- callTransientRetry(t, ctx, srv)
	}()

	select {
	case err := <-done:
		assert.Assert(t, err != nil)
	case <-time.After(10 * time.Second):
		t.Fatal("executeWithTransientRetry did not return on a canceled context")
	}
}
