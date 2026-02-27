package api

import (
	"context"
	goer "errors"
	"net"
	"testing"
	"time"
)

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "context deadline exceeded",
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			name: "wrapped context deadline",
			err:  goer.Join(goer.New("wrapped"), context.DeadlineExceeded),
			want: true,
		},
		{
			name: "timeout in error message",
			err:  goer.New("context deadline exceeded while awaiting headers"),
			want: true,
		},
		{
			name: "net timeout error",
			err:  &timeoutError{},
			want: true,
		},
		{
			name: "other error",
			err:  goer.New("some other error"),
			want: false,
		},
		{
			name: "connection refused",
			err:  goer.New("connection refused"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTimeoutError(tt.err)
			if got != tt.want {
				t.Errorf("isTimeoutError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateExponentialBackoff(t *testing.T) {
	baseWait := 2 * time.Second

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name:    "first attempt",
			attempt: 1,
			want:    2 * time.Second,
			wantMin: 1600 * time.Millisecond, // 2s - 20%
			wantMax: 2400 * time.Millisecond, // 2s + 20%
		},
		{
			name:    "second attempt",
			attempt: 2,
			want:    4 * time.Second,
			wantMin: 3200 * time.Millisecond, // 4s - 20%
			wantMax: 4800 * time.Millisecond, // 4s + 20%
		},
		{
			name:    "third attempt",
			attempt: 3,
			want:    8 * time.Second,
			wantMin: 6400 * time.Millisecond, // 8s - 20%
			wantMax: 9600 * time.Millisecond, // 8s + 20%
		},
		{
			name:    "fourth attempt",
			attempt: 4,
			want:    16 * time.Second,
			wantMin: 12800 * time.Millisecond, // 16s - 20%
			wantMax: 19200 * time.Millisecond, // 16s + 20%
		},
		{
			name:    "fifth attempt (capped at 32s)",
			attempt: 5,
			want:    32 * time.Second,
			wantMin: 25600 * time.Millisecond, // 32s - 20%
			wantMax: 38400 * time.Millisecond, // 32s + 20%
		},
		{
			name:    "tenth attempt (still capped at 32s)",
			attempt: 10,
			want:    32 * time.Second,
			wantMin: 25600 * time.Millisecond, // 32s - 20%
			wantMax: 38400 * time.Millisecond, // 32s + 20%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to verify jitter is working
			for i := 0; i < 10; i++ {
				got := calculateExponentialBackoff(tt.attempt, baseWait)
				if got < tt.wantMin || got > tt.wantMax {
					t.Errorf("calculateExponentialBackoff(%d) = %v, want between %v and %v",
						tt.attempt, got, tt.wantMin, tt.wantMax)
				}
			}
		})
	}
}

func TestCalculateExponentialBackoff_Jitter(t *testing.T) {
	baseWait := 2 * time.Second
	attempt := 3

	// Run multiple times and ensure we get different values (jitter is working)
	results := make(map[time.Duration]bool)
	for i := 0; i < 50; i++ {
		backoff := calculateExponentialBackoff(attempt, baseWait)
		results[backoff] = true
	}

	// Should have multiple different values due to jitter
	if len(results) < 5 {
		t.Errorf("calculateExponentialBackoff produced only %d unique values in 50 attempts, jitter may not be working",
			len(results))
	}
}

// timeoutError is a mock implementation of net.Error for testing
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout error" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

// Verify timeoutError implements net.Error
var _ net.Error = (*timeoutError)(nil)
