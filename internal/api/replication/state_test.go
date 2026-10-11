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
