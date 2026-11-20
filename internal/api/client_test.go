package api

import (
	"context"
	goer "errors"
	"fmt"
	"testing"
	"time"

	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_Success(t *testing.T) {
	fn := func() (*Response, time.Duration, error) {
		return &Response{}, 0, nil
	}
	resp, err := exec(context.Background(), fn, defaultWaitAttempt)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestExec_NonRetryableError(t *testing.T) {
	fn := func() (*Response, time.Duration, error) {
		return nil, 0, assert.AnError
	}
	_, err := exec(context.Background(), fn, defaultWaitAttempt)
	require.Error(t, err)
	assert.True(t, goer.Is(err, assert.AnError))
}

func TestExec_AttemptCapEnforced(t *testing.T) {
	// Use short Retry-After values so the test runs quickly.
	// The exponential-backoff path is exercised by TestBackoffWithJitter_*.
	shortBackoff := 1 * time.Millisecond

	tests := []struct {
		name           string
		retryableErr   error
		backOff        time.Duration
		expectAttempts int
	}{
		{
			name:           "503 service unavailable",
			retryableErr:   internalerrors.ErrServiceUnavailable,
			backOff:        shortBackoff,
			expectAttempts: maxRetryAttempts + 1,
		},
		{
			name:           "504 gateway timeout",
			retryableErr:   internalerrors.ErrGatewayTimeout,
			backOff:        shortBackoff,
			expectAttempts: maxRetryAttempts + 1,
		},
		{
			name:           "429 rate limit with Retry-After",
			retryableErr:   internalerrors.ErrRatelimit,
			backOff:        shortBackoff,
			expectAttempts: maxRetryAttempts + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var callCount int
			fn := func() (*Response, time.Duration, error) {
				callCount++
				return nil, tt.backOff, tt.retryableErr
			}
			_, err := exec(context.Background(), fn, defaultWaitAttempt)
			require.Error(t, err)

			assert.Equal(t, tt.expectAttempts, callCount, "should make initial call + maxRetryAttempts retries")

			var retryErr *RetriesExhaustedError
			require.True(t, goer.As(err, &retryErr), "error should be *RetriesExhaustedError")
			assert.Equal(t, tt.expectAttempts, retryErr.Attempts)
			assert.True(t, goer.Is(retryErr, tt.retryableErr), "RetriesExhaustedError should wrap the original sentinel")
			assert.Greater(t, retryErr.Elapsed, time.Duration(0), "elapsed time should be recorded")
		})
	}
}

func TestExec_ExhaustionErrorIsTyped(t *testing.T) {
	fn := func() (*Response, time.Duration, error) {
		return nil, 1 * time.Millisecond, internalerrors.ErrServiceUnavailable
	}
	_, err := exec(context.Background(), fn, defaultWaitAttempt)

	var retryErr *RetriesExhaustedError
	require.True(t, goer.As(err, &retryErr))
	assert.Contains(t, err.Error(), "request failed after")
	assert.Contains(t, err.Error(), "attempts over")
	assert.True(t, goer.Is(err, internalerrors.ErrServiceUnavailable),
		"errors.Is should find the wrapped sentinel through Unwrap()")
}

func TestExec_RetryAfterHonored(t *testing.T) {
	retryAfter := 10 * time.Millisecond
	var callCount int
	fn := func() (*Response, time.Duration, error) {
		callCount++
		return nil, retryAfter, internalerrors.ErrRatelimit
	}

	start := time.Now()
	_, err := exec(context.Background(), fn, defaultWaitAttempt)
	elapsed := time.Since(start)

	// maxRetryAttempts retries × retryAfter = total wait time, plus overhead.
	minExpected := time.Duration(maxRetryAttempts) * retryAfter
	assert.GreaterOrEqual(t, elapsed, minExpected,
		"total elapsed should be at least retryAfter × maxRetryAttempts")

	var retryErr *RetriesExhaustedError
	require.True(t, goer.As(err, &retryErr))
	assert.Equal(t, maxRetryAttempts+1, callCount)
}

func TestExec_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var callCount int
	fn := func() (*Response, time.Duration, error) {
		callCount++
		// Return a retryable error so we stay in the loop, but with a long
		// backoff that gives the context time to expire.
		return nil, 200 * time.Millisecond, internalerrors.ErrServiceUnavailable
	}

	_, err := exec(ctx, fn, defaultWaitAttempt)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out executing request against api")
	assert.Equal(t, 1, callCount, "only the initial call should fire before context expires")
}

func TestExec_MixedRetryableAndNonRetryable(t *testing.T) {
	// A non-retryable error on the second call should stop retries immediately.
	var callCount int
	fn := func() (*Response, time.Duration, error) {
		callCount++
		switch callCount {
		case 1:
			return nil, 1 * time.Millisecond, internalerrors.ErrServiceUnavailable
		default:
			return nil, 0, fmt.Errorf("non-retryable failure")
		}
	}

	_, err := exec(context.Background(), fn, defaultWaitAttempt)
	require.Error(t, err)
	assert.Equal(t, 2, callCount, "should stop after encountering a non-retryable error")
	assert.False(t, goer.Is(err, internalerrors.ErrServiceUnavailable))
	var retryErr *RetriesExhaustedError
	assert.False(t, goer.As(err, &retryErr))
}

func TestBackoffWithJitter_Grows(t *testing.T) {
	// Verify that backoff values are strictly increasing between attempts
	// due to non-overlapping jitter ranges.
	var previous time.Duration
	for attempt := range maxRetryAttempts {
		b := backoffWithJitter(attempt)
		if attempt > 0 {
			assert.Greater(t, b, previous,
				"backoff attempt %d (%v) should be > attempt %d (%v)", attempt, b, attempt-1, previous)
		}
		previous = b
	}
}

func TestBackoffWithJitter_Capped(t *testing.T) {
	// At high attempt numbers the raw value is capped at maxBackoff.
	b := backoffWithJitter(10) // 2 * 2^10 = 2048s, capped at 30s
	assert.LessOrEqual(t, b, maxBackoff)
	assert.GreaterOrEqual(t, b, maxBackoff/2)
}

func TestBackoffWithJitter_Jittered(t *testing.T) {
	// Run multiple times for the same attempt and verify the values differ
	// (proof of jitter), but stay within the expected range [raw/2, raw].
	const attempt = 2 // raw = 2 * 2^2 = 8s, range [4s, 8s]
	const iterations = 10

	var values []time.Duration
	for range iterations {
		values = append(values, backoffWithJitter(attempt))
	}

	allSame := true
	for i := 1; i < len(values); i++ {
		if values[i] != values[0] {
			allSame = false
			break
		}
	}
	assert.False(t, allSame, "jitter should produce varying values across calls")

	raw := defaultWaitAttempt * (1 << attempt)
	for _, v := range values {
		assert.GreaterOrEqual(t, v, raw/2, "jittered value should be >= raw/2")
		assert.LessOrEqual(t, v, raw, "jittered value should be <= raw")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{"ErrRatelimit is retryable", internalerrors.ErrRatelimit, true},
		{"ErrServiceUnavailable is retryable", internalerrors.ErrServiceUnavailable, true},
		{"ErrGatewayTimeout is retryable", internalerrors.ErrGatewayTimeout, true},
		{"ErrGatewayTimeoutForIndexDDL is NOT retryable", internalerrors.ErrGatewayTimeoutForIndexDDL, false},
		{"wrapped ErrRatelimit is retryable", fmt.Errorf("wrap: %w", internalerrors.ErrRatelimit), true},
		{"nil is not retryable", nil, false},
		{"arbitrary error is not retryable", assert.AnError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.retryable, isRetryable(tt.err))
		})
	}
}
