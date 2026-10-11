package resources

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

const testReplicationJobPath = "/v4/organizations/org-1/projects/proj-1/clusters/cluster-1/replications/jobs/job-1"

// jobReply is one scripted answer of replicationJobBackend: a job body, or an HTTP error status.
type jobReply struct {
	job        replicationapi.GetReplicationJobResponse
	errStatus  int
	errMessage string
}

// replicationJobBackend answers job polls from a scripted queue. The last entry repeats once the
// queue is drained.
type replicationJobBackend struct {
	t       *testing.T
	mu      sync.Mutex
	replies []jobReply
	calls   int
}

func (b *replicationJobBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	assert.Equal(b.t, http.MethodGet, r.Method)
	assert.Equal(b.t, testReplicationJobPath, r.URL.Path)

	idx := min(b.calls, len(b.replies)-1)
	b.calls++
	reply := b.replies[idx]

	if reply.errStatus != 0 {
		w.WriteHeader(reply.errStatus)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":           2,
			"httpStatusCode": reply.errStatus,
			"message":        reply.errMessage,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reply.job)
}

func (b *replicationJobBackend) callCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls
}

// fastReplicationJobPolling shrinks the package-level poll timing for the duration of a test so
// waitForReplicationJob's poll loop resolves in milliseconds.
func fastReplicationJobPolling(t *testing.T, timeout time.Duration) {
	t.Helper()
	origInterval, origTimeout := replicationJobPollInterval, replicationJobTimeout
	replicationJobPollInterval = time.Millisecond
	replicationJobTimeout = timeout
	t.Cleanup(func() {
		replicationJobPollInterval = origInterval
		replicationJobTimeout = origTimeout
	})
}

func jobInState(state replicationapi.JobState) jobReply {
	return jobReply{job: replicationapi.GetReplicationJobResponse{JobId: "job-1", State: state}}
}

func jobError(status int) jobReply {
	return jobReply{errStatus: status, errMessage: http.StatusText(status)}
}

func TestWaitForReplicationJob(t *testing.T) {
	complete := jobReply{job: replicationapi.GetReplicationJobResponse{
		JobId:                "job-1",
		State:                replicationapi.JobStateComplete,
		ReplicationId:        "cmVwbC0x",
		ReverseReplicationId: "cmV2LTE=",
	}}
	failed := func(state replicationapi.JobState) jobReply {
		return jobReply{job: replicationapi.GetReplicationJobResponse{
			JobId:       "job-1",
			State:       state,
			RetryNumber: 1,
			LastError:   "source bucket not found",
		}}
	}

	tests := []struct {
		name        string
		replies     []jobReply
		timeout     time.Duration
		wantErr     error
		errContains []string
		wantCalls   int
	}{
		{
			name:      "complete on the first poll returns the job with both replication IDs",
			replies:   []jobReply{complete},
			wantCalls: 1,
		},
		{
			name:      "pending then processing then complete returns the job after three polls",
			replies:   []jobReply{jobInState(replicationapi.JobStatePending), jobInState(replicationapi.JobStateProcessing), complete},
			wantCalls: 3,
		},
		{
			name:        "failed after processing returns ErrReplicationJobFailed with state, retry number and last error",
			replies:     []jobReply{jobInState(replicationapi.JobStateProcessing), failed(replicationapi.JobStateFailed)},
			wantErr:     errors.ErrReplicationJobFailed,
			errContains: []string{"state failed", "retry 1", "source bucket not found"},
			wantCalls:   2,
		},
		{
			name:        "skipped state returns ErrReplicationJobFailed",
			replies:     []jobReply{failed(replicationapi.JobStateSkipped)},
			wantErr:     errors.ErrReplicationJobFailed,
			errContains: []string{"state skipped"},
			wantCalls:   1,
		},
		{
			name:        "killed state returns ErrReplicationJobFailed",
			replies:     []jobReply{failed(replicationapi.JobStateKilled)},
			wantErr:     errors.ErrReplicationJobFailed,
			errContains: []string{"state killed"},
			wantCalls:   1,
		},
		{
			name:        "unknown state returns ErrReplicationJobFailed",
			replies:     []jobReply{failed(replicationapi.JobStateUnknown)},
			wantErr:     errors.ErrReplicationJobFailed,
			errContains: []string{"state unknown"},
			wantCalls:   1,
		},
		{
			name:        "failed state with an empty last error reports last error as none",
			replies:     []jobReply{jobInState(replicationapi.JobStateFailed)},
			wantErr:     errors.ErrReplicationJobFailed,
			errContains: []string{"last error: none"},
			wantCalls:   1,
		},
		{
			name: "notfound before the job starts is tolerated and the poll continues to complete",
			replies: []jobReply{
				jobInState(replicationapi.JobStateNotFound),
				jobInState(replicationapi.JobStateNotFound),
				jobInState(replicationapi.JobStateProcessing),
				complete,
			},
			wantCalls: 4,
		},
		{
			name:        "notfound four times in a row before the job starts returns ErrReplicationJobNotFound",
			replies:     []jobReply{jobInState(replicationapi.JobStateNotFound)},
			wantErr:     errors.ErrReplicationJobNotFound,
			errContains: []string{"4 times in a row and never started"},
			wantCalls:   4,
		},
		{
			name:        "notfound after the job was seen running returns ErrReplicationJobNotFound at once",
			replies:     []jobReply{jobInState(replicationapi.JobStateProcessing), jobInState(replicationapi.JobStateNotFound)},
			wantErr:     errors.ErrReplicationJobNotFound,
			errContains: []string{"was running, then reported notfound"},
			wantCalls:   2,
		},
		{
			name:      "server error (5xx) is retried and the poll continues to complete",
			replies:   []jobReply{jobError(http.StatusInternalServerError), complete},
			wantCalls: 2,
		},
		{
			name:        "client error (4xx, 403) stops the poll after one call",
			replies:     []jobReply{jobError(http.StatusForbidden), complete},
			errContains: []string{"could not get replication job job-1", "Forbidden"},
			wantCalls:   1,
		},
		{
			name:        "complete without a replication ID returns ErrReplicationJobMissingId",
			replies:     []jobReply{jobInState(replicationapi.JobStateComplete)},
			wantErr:     errors.ErrReplicationJobMissingId,
			errContains: []string{"job job-1"},
			wantCalls:   1,
		},
		{
			name:        "timeout while processing returns ErrReplicationJobTimeout with the last state seen",
			replies:     []jobReply{failed(replicationapi.JobStateProcessing)},
			timeout:     50 * time.Millisecond,
			wantErr:     errors.ErrReplicationJobTimeout,
			errContains: []string{"last state processing", "retry 1", "source bucket not found"},
		},
		{
			name:        "timeout with only server errors returns ErrReplicationJobTimeout with no job status seen",
			replies:     []jobReply{jobError(http.StatusInternalServerError)},
			timeout:     50 * time.Millisecond,
			wantErr:     errors.ErrReplicationJobTimeout,
			errContains: []string{"no job status was seen"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := tt.timeout
			if timeout == 0 {
				timeout = 5 * time.Second
			}
			fastReplicationJobPolling(t, timeout)

			backend := &replicationJobBackend{t: t, replies: tt.replies}
			data := newTestProviderData(t, backend.handler)

			job, err := waitForReplicationJob(t.Context(), data, "org-1", "proj-1", "cluster-1", "job-1")

			if tt.wantErr != nil || len(tt.errContains) > 0 {
				require.Error(t, err)
				assert.Nil(t, job)
				if tt.wantErr != nil {
					assert.True(t, stderrors.Is(err, tt.wantErr), "error %q is not %q", err, tt.wantErr)
				}
				for _, s := range tt.errContains {
					assert.ErrorContains(t, err, s)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, "cmVwbC0x", job.ReplicationId)
				assert.Equal(t, "cmV2LTE=", job.ReverseReplicationId)
			}

			if tt.wantCalls > 0 {
				assert.Equal(t, tt.wantCalls, backend.callCount())
			}
		})
	}
}
