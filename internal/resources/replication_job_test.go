package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// replicationJobStep is one scripted reply from the fake control plane. A step either serves a job
// document or, when status is set, a bare HTTP status.
type replicationJobStep struct {
	state         apigen.GetReplicationJobResponseState
	replicationID string
	reverseID     string
	retryNumber   int
	lastError     string
	status        int
}

// replicationJobBackend is an in-memory control-plane stand-in that answers GET replication job
// polls from a scripted queue of steps (the last entry repeats once the queue is drained) so tests
// can drive waitForReplicationJob's poll loop deterministically.
type replicationJobBackend struct {
	mu        sync.Mutex
	steps     []replicationJobStep
	idx       int
	calls     int
	lastPath  string
	jobID     string
	omitRepID bool
}

func (b *replicationJobBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.calls++
	b.lastPath = r.URL.Path

	idx := b.idx
	if idx >= len(b.steps) {
		idx = len(b.steps) - 1
	}
	b.idx++
	step := b.steps[idx]

	if step.status != 0 {
		w.WriteHeader(step.status)
		return
	}

	job := apigen.GetReplicationJobResponse{
		JobId:                b.jobID,
		State:                step.state,
		LastUpdatedTimestamp: time.Now(),
	}
	if step.replicationID != "" {
		job.ReplicationId = &step.replicationID
	}
	if step.reverseID != "" {
		job.ReverseReplicationId = &step.reverseID
	}
	if step.retryNumber != 0 {
		job.RetryNumber = &step.retryNumber
	}
	if step.lastError != "" {
		job.LastError = &step.lastError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(job)
}

func (b *replicationJobBackend) callCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls
}

func (b *replicationJobBackend) path() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.lastPath
}

// fastReplicationJobPolling shrinks the package-level poll timing for the duration of a test so
// waitForReplicationJob's poll loop resolves in milliseconds.
func fastReplicationJobPolling(t *testing.T) {
	t.Helper()
	origInterval, origTimeout := replicationJobPollInterval, replicationJobTimeout
	replicationJobPollInterval = time.Millisecond
	replicationJobTimeout = 2 * time.Second
	t.Cleanup(func() {
		replicationJobPollInterval = origInterval
		replicationJobTimeout = origTimeout
	})
}

func newTestReplicationJobClient(t *testing.T, handler http.HandlerFunc) *apigen.ClientWithResponses {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client, err := apigen.NewClientWithResponses(srv.URL)
	assert.NilError(t, err)
	return client
}

func TestIsFinalReplicationJobState(t *testing.T) {
	tests := []struct {
		state apigen.GetReplicationJobResponseState
		final bool
	}{
		{apigen.GetReplicationJobResponseStatePending, false},
		{apigen.GetReplicationJobResponseStateProcessing, false},
		{apigen.GetReplicationJobResponseStateComplete, true},
		{apigen.GetReplicationJobResponseStateFailed, true},
		{apigen.GetReplicationJobResponseStateKilled, true},
		{apigen.GetReplicationJobResponseStateSkipped, true},
		{apigen.GetReplicationJobResponseStateNotfound, true},
		{apigen.GetReplicationJobResponseStateUnknown, true},
		{apigen.GetReplicationJobResponseState("somethingElse"), false},
	}

	for _, test := range tests {
		t.Run(string(test.state), func(t *testing.T) {
			assert.Equal(t, isFinalReplicationJobState(test.state), test.final)
		})
	}
}

func TestIsFailureReplicationJobState(t *testing.T) {
	tests := []struct {
		state   apigen.GetReplicationJobResponseState
		failure bool
	}{
		{apigen.GetReplicationJobResponseStatePending, false},
		{apigen.GetReplicationJobResponseStateProcessing, false},
		{apigen.GetReplicationJobResponseStateComplete, false},
		{apigen.GetReplicationJobResponseStateFailed, true},
		{apigen.GetReplicationJobResponseStateKilled, true},
		{apigen.GetReplicationJobResponseStateSkipped, true},
		{apigen.GetReplicationJobResponseStateNotfound, true},
		{apigen.GetReplicationJobResponseStateUnknown, true},
		{apigen.GetReplicationJobResponseState("somethingElse"), false},
	}

	for _, test := range tests {
		t.Run(string(test.state), func(t *testing.T) {
			assert.Equal(t, isFailureReplicationJobState(test.state), test.failure)
		})
	}
}

// TestReplicationJobFailureStatesAreFinal guards the invariant documented on
// isFailureReplicationJobState.
func TestReplicationJobFailureStatesAreFinal(t *testing.T) {
	states := []apigen.GetReplicationJobResponseState{
		apigen.GetReplicationJobResponseStatePending,
		apigen.GetReplicationJobResponseStateProcessing,
		apigen.GetReplicationJobResponseStateComplete,
		apigen.GetReplicationJobResponseStateFailed,
		apigen.GetReplicationJobResponseStateKilled,
		apigen.GetReplicationJobResponseStateSkipped,
		apigen.GetReplicationJobResponseStateNotfound,
		apigen.GetReplicationJobResponseStateUnknown,
	}

	for _, state := range states {
		if isFailureReplicationJobState(state) {
			assert.Assert(t, isFinalReplicationJobState(state), "failure state %s must also be final", state)
		}
	}
}

func waitForTestReplicationJob(t *testing.T, backend *replicationJobBackend) (*apigen.GetReplicationJobResponse, error) {
	t.Helper()
	fastReplicationJobPolling(t)

	jobUUID := uuid.New()
	backend.jobID = jobUUID.String()
	client := newTestReplicationJobClient(t, backend.handler)

	return waitForReplicationJob(t.Context(), client, uuid.New(), uuid.New(), uuid.New(), jobUUID)
}

func TestWaitForReplicationJob_Complete(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStatePending},
		{state: apigen.GetReplicationJobResponseStateProcessing},
		{state: apigen.GetReplicationJobResponseStateComplete, replicationID: "repl-1", reverseID: "repl-1-reverse"},
	}}

	job, err := waitForTestReplicationJob(t, backend)

	assert.NilError(t, err)
	assert.Equal(t, *job.ReplicationId, "repl-1")
	assert.Equal(t, *job.ReverseReplicationId, "repl-1-reverse")
}

func TestWaitForReplicationJob_TerminalFailureStates(t *testing.T) {
	tests := []apigen.GetReplicationJobResponseState{
		apigen.GetReplicationJobResponseStateFailed,
		apigen.GetReplicationJobResponseStateKilled,
		apigen.GetReplicationJobResponseStateSkipped,
		apigen.GetReplicationJobResponseStateUnknown,
	}

	for _, state := range tests {
		t.Run(string(state), func(t *testing.T) {
			backend := &replicationJobBackend{steps: []replicationJobStep{
				{state: state, retryNumber: 2, lastError: "boom"},
			}}

			_, err := waitForTestReplicationJob(t, backend)

			assert.ErrorContains(t, err, errors.ErrReplicationJobFailed.Error())
			assert.ErrorContains(t, err, string(state))
			assert.ErrorContains(t, err, "retry 2")
			assert.ErrorContains(t, err, "boom")
		})
	}
}

// TestWaitForReplicationJob_ToleratesInitialNotFound covers the read race between the create call
// returning 202 and the job record becoming readable.
func TestWaitForReplicationJob_ToleratesInitialNotFound(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateNotfound},
		{state: apigen.GetReplicationJobResponseStateNotfound},
		{state: apigen.GetReplicationJobResponseStateComplete, replicationID: "repl-2"},
	}}

	job, err := waitForTestReplicationJob(t, backend)

	assert.NilError(t, err)
	assert.Equal(t, *job.ReplicationId, "repl-2")
}

func TestWaitForReplicationJob_ToleratesInitialHTTP404(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{status: http.StatusNotFound},
		{state: apigen.GetReplicationJobResponseStateComplete, replicationID: "repl-3"},
	}}

	job, err := waitForTestReplicationJob(t, backend)

	assert.NilError(t, err)
	assert.Equal(t, *job.ReplicationId, "repl-3")
}

func TestWaitForReplicationJob_NotFoundBeyondTolerance(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateNotfound},
	}}

	_, err := waitForTestReplicationJob(t, backend)

	assert.ErrorContains(t, err, errors.ErrReplicationJobNotFound.Error())
	assert.Equal(t, backend.callCount(), maxConsecutiveJobNotFound+1)
}

// TestWaitForReplicationJob_NotFoundAfterSeenLive proves a job that was observed running is not
// given the initial read-race tolerance if it later disappears.
func TestWaitForReplicationJob_NotFoundAfterSeenLive(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateProcessing},
		{state: apigen.GetReplicationJobResponseStateNotfound},
	}}

	_, err := waitForTestReplicationJob(t, backend)

	assert.ErrorContains(t, err, errors.ErrReplicationJobNotFound.Error())
	assert.Equal(t, backend.callCount(), 2)
}

func TestWaitForReplicationJob_RetriesTransientErrors(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{status: http.StatusInternalServerError},
		{status: http.StatusBadGateway},
		{state: apigen.GetReplicationJobResponseStateComplete, replicationID: "repl-4"},
	}}

	job, err := waitForTestReplicationJob(t, backend)

	assert.NilError(t, err)
	assert.Equal(t, *job.ReplicationId, "repl-4")
	assert.Assert(t, backend.callCount() >= 3)
}

func TestWaitForReplicationJob_TimesOutWithoutReachingTerminalState(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateProcessing, retryNumber: 1, lastError: "still going"},
	}}

	_, err := waitForTestReplicationJob(t, backend)

	assert.ErrorContains(t, err, errors.ErrReplicationJobTimeout.Error())
	assert.ErrorContains(t, err, string(apigen.GetReplicationJobResponseStateProcessing))
	assert.ErrorContains(t, err, "still going")
}

// TestWaitForReplicationJob_CompleteWithoutReplicationID guards against handing the caller a blank
// ID when the platform reports completion but omits the field.
func TestWaitForReplicationJob_CompleteWithoutReplicationID(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateComplete},
	}}

	_, err := waitForTestReplicationJob(t, backend)

	assert.ErrorContains(t, err, errors.ErrReplicationJobFailed.Error())
	assert.ErrorContains(t, err, "without returning a replication ID")
}

func TestWaitForReplicationJob_UsesJobEndpoint(t *testing.T) {
	backend := &replicationJobBackend{steps: []replicationJobStep{
		{state: apigen.GetReplicationJobResponseStateComplete, replicationID: "repl-5"},
	}}

	_, err := waitForTestReplicationJob(t, backend)

	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(backend.path(), "/replications/jobs/"), "unexpected path %q", backend.path())
}
