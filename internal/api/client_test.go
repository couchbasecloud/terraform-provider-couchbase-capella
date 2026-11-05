package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stalledServer never answers. It returns when the client cancels the request, or
// after a guard so a regression fails the test instead of hanging it.
func stalledServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(10 * time.Second):
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestExecuteWithRetryEndsStalledRequestWithContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     func() (context.Context, context.CancelFunc)
		wantErr error
	}{
		{
			name: "deadline",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(t.Context(), 100*time.Millisecond)
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "cancel",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(t.Context())
				time.AfterFunc(100*time.Millisecond, cancel)
				return ctx, cancel
			},
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := stalledServer(t)
			client := &Client{Client: srv.Client()}

			ctx, cancel := tt.ctx()
			defer cancel()

			start := time.Now()
			_, err := client.ExecuteWithRetry(ctx, EndpointCfg{Url: srv.URL, Method: http.MethodGet, SuccessStatus: http.StatusOK}, nil, "token", nil)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Less(t, time.Since(start), 5*time.Second, "the request must end with its context, not the HTTP client timeout")
		})
	}
}

func TestPollError(t *testing.T) {
	errRequest := errors.New("request failed")
	errTimeout := errors.New("poller timed out")

	live := t.Context()
	ended, cancel := context.WithCancel(t.Context())
	cancel()

	assert.Equal(t, errRequest, PollError(live, errRequest, errTimeout), "a live context keeps the request error")
	assert.Equal(t, errTimeout, PollError(ended, errRequest, errTimeout), "an ended context reports the poller's timeout")
}
