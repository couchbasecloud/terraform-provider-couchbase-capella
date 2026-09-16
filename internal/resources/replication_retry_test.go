package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// fastXDCRLockRetry shrinks the package-level retry interval for the duration of a test so
// retryOnXDCRLock resolves in milliseconds.
func fastXDCRLockRetry(t *testing.T) {
	t.Helper()
	orig := xdcrLockRetryInterval
	xdcrLockRetryInterval = time.Millisecond
	t.Cleanup(func() { xdcrLockRetryInterval = orig })
}

func lockedBody(t *testing.T, message, hint string) []byte {
	t.Helper()
	body, err := json.Marshal(apigen.Error{
		Code:           4242,
		Message:        message,
		Hint:           hint,
		HttpStatusCode: http.StatusLocked,
	})
	assert.NilError(t, err)
	return body
}

// lockResponder returns an fn that replies 423 for the first lockedReplies calls and 200 after.
func lockResponder(t *testing.T, lockedReplies int, calls *int) func() (*http.Response, []byte, error) {
	t.Helper()
	return func() (*http.Response, []byte, error) {
		*calls++
		if *calls <= lockedReplies {
			return &http.Response{StatusCode: http.StatusLocked}, lockedBody(t, "XDCROperationInProgress", "wait"), nil
		}
		return &http.Response{StatusCode: http.StatusOK}, []byte(`{"ok":true}`), nil
	}
}

func TestRetryOnXDCRLock_SucceedsImmediately(t *testing.T) {
	fastXDCRLockRetry(t)

	var calls int
	body, err := retryOnXDCRLock(t.Context(), xdcrShortLockRetryBudget, lockResponder(t, 0, &calls))

	assert.NilError(t, err)
	assert.Equal(t, string(body), `{"ok":true}`)
	assert.Equal(t, calls, 1)
}

func TestRetryOnXDCRLock_RetriesUntilLockReleased(t *testing.T) {
	fastXDCRLockRetry(t)

	var calls int
	body, err := retryOnXDCRLock(t.Context(), xdcrShortLockRetryBudget, lockResponder(t, 3, &calls))

	assert.NilError(t, err)
	assert.Equal(t, string(body), `{"ok":true}`)
	assert.Equal(t, calls, 4)
}

func TestRetryOnXDCRLock_ExhaustsBudget(t *testing.T) {
	fastXDCRLockRetry(t)

	var calls int
	_, err := retryOnXDCRLock(t.Context(), 20*time.Millisecond, func() (*http.Response, []byte, error) {
		calls++
		return &http.Response{StatusCode: http.StatusLocked}, lockedBody(t, "XDCROperationInProgress", "wait"), nil
	})

	assert.ErrorContains(t, err, errors.ErrXDCROperationInProgress.Error())
	assert.ErrorContains(t, err, "XDCROperationInProgress")
	assert.Assert(t, calls > 1, "expected more than one attempt, got %d", calls)
}

// TestRetryOnXDCRLock_NeverOverrunsBudget covers the clamp that stops a final wait from running
// past the deadline.
func TestRetryOnXDCRLock_NeverOverrunsBudget(t *testing.T) {
	origInterval := xdcrLockRetryInterval
	xdcrLockRetryInterval = time.Second
	t.Cleanup(func() { xdcrLockRetryInterval = origInterval })

	const budget = 50 * time.Millisecond
	start := time.Now()
	_, err := retryOnXDCRLock(t.Context(), budget, func() (*http.Response, []byte, error) {
		return &http.Response{StatusCode: http.StatusLocked}, lockedBody(t, "XDCROperationInProgress", ""), nil
	})
	elapsed := time.Since(start)

	assert.ErrorContains(t, err, errors.ErrXDCROperationInProgress.Error())
	// One clamped wait of at most the budget, plus scheduling slack.
	assert.Assert(t, elapsed < 500*time.Millisecond, "helper overran its budget: %v", elapsed)
}

// TestRetryOnXDCRLock_PassesThroughOtherStatuses proves the helper only reacts to 423 and leaves
// every other outcome to the caller, including conflict statuses the API does document.
func TestRetryOnXDCRLock_PassesThroughOtherStatuses(t *testing.T) {
	fastXDCRLockRetry(t)

	for _, status := range []int{http.StatusConflict, http.StatusUnprocessableEntity, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var calls int
			body, err := retryOnXDCRLock(t.Context(), xdcrShortLockRetryBudget, func() (*http.Response, []byte, error) {
				calls++
				return &http.Response{StatusCode: status}, []byte(`{"code":1}`), nil
			})

			assert.NilError(t, err)
			assert.Equal(t, string(body), `{"code":1}`)
			assert.Equal(t, calls, 1)
		})
	}
}

func TestRetryOnXDCRLock_ReturnsTransportError(t *testing.T) {
	fastXDCRLockRetry(t)

	var calls int
	_, err := retryOnXDCRLock(t.Context(), xdcrShortLockRetryBudget, func() (*http.Response, []byte, error) {
		calls++
		return nil, nil, context.DeadlineExceeded
	})

	assert.Assert(t, stderrors.Is(err, context.DeadlineExceeded))
	assert.Equal(t, calls, 1)
}

func TestRetryOnXDCRLock_StopsOnContextCancellation(t *testing.T) {
	fastXDCRLockRetry(t)
	xdcrLockRetryInterval = 50 * time.Millisecond

	ctx, cancel := context.WithCancel(t.Context())
	var calls int
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err := retryOnXDCRLock(ctx, xdcrCreateLockRetryBudget, func() (*http.Response, []byte, error) {
		calls++
		return &http.Response{StatusCode: http.StatusLocked}, lockedBody(t, "XDCROperationInProgress", ""), nil
	})

	assert.ErrorContains(t, err, errors.ErrXDCROperationInProgress.Error())
	assert.Assert(t, stderrors.Is(err, context.Canceled))
}

func TestParseXDCRLockMessage(t *testing.T) {
	tests := []struct {
		name string
		body []byte
		want string
	}{
		{"message and hint", lockedBody(t, "operation in progress", "retry later"), "operation in progress (retry later)"},
		{"message only", lockedBody(t, "operation in progress", ""), "operation in progress"},
		{"unparsable body falls back to raw", []byte("not json"), "not json"},
		{"empty message falls back to raw", []byte(`{"code":1}`), `{"code":1}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, parseXDCRLockMessage(test.body), test.want)
		})
	}
}
