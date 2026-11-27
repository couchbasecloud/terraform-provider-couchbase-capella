package replication

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFailedJobState(t *testing.T) {
	tests := []struct {
		name       string
		state      JobState
		wantFailed bool
	}{
		{name: "pending is not a failure", state: JobStatePending, wantFailed: false},
		{name: "processing is not a failure", state: JobStateProcessing, wantFailed: false},
		{name: "complete is not a failure", state: JobStateComplete, wantFailed: false},
		{name: "failed is a failure", state: JobStateFailed, wantFailed: true},
		{name: "skipped is a failure", state: JobStateSkipped, wantFailed: true},
		{name: "killed is a failure", state: JobStateKilled, wantFailed: true},
		{name: "notfound is a failure", state: JobStateNotFound, wantFailed: true},
		{name: "unknown is a failure", state: JobStateUnknown, wantFailed: true},
		{name: "unrecognised state is not a failure", state: JobState("somethingElse"), wantFailed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFailed, IsFailedJobState(tt.state))
		})
	}
}

func TestIsFatalJobError(t *testing.T) {
	tests := []struct {
		name      string
		lastError string
		wantFatal bool
	}{
		{name: "job-fatal prefix is fatal", lastError: "job-fatal:invalid create replication job params", wantFatal: true},
		{name: "goxdcr-fatal-error prefix is fatal", lastError: "goxdcr-fatal-error:invalid filter expression", wantFatal: true},
		{name: "job-error prefix is retried", lastError: "job-error:Replication creation rate-limited due to peering deletion", wantFatal: false},
		{name: "an error without a prefix is retried", lastError: "connection reset", wantFatal: false},
		{name: "no error is retried", lastError: "", wantFatal: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFatal, IsFatalJobError(tt.lastError))
		})
	}
}
