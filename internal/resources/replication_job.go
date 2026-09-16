package resources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// replicationJobPollInterval and replicationJobTimeout are overridden in tests to avoid waiting on
// real wall-clock time while exercising waitForReplicationJob's polling loop.
//
// The timeout has to outlast the control plane's internal reattempt of a failed replication job,
// which it makes 15 minutes after the first attempt, plus the reattempted job's own runtime.
var (
	replicationJobPollInterval = 15 * time.Second
	replicationJobTimeout      = 30 * time.Minute
)

// maxConsecutiveJobNotFound bounds how long a job may report notfound immediately after the create
// call was accepted. The job record is not always readable the instant the API returns 202, so a
// short burst of notfound is a read race rather than a real failure.
const maxConsecutiveJobNotFound = 3

// isFinalReplicationJobState reports whether the job has stopped moving.
func isFinalReplicationJobState(state apigen.GetReplicationJobResponseState) bool {
	switch state {
	case apigen.GetReplicationJobResponseStateComplete,
		apigen.GetReplicationJobResponseStateFailed,
		apigen.GetReplicationJobResponseStateKilled,
		apigen.GetReplicationJobResponseStateSkipped,
		apigen.GetReplicationJobResponseStateNotfound,
		apigen.GetReplicationJobResponseStateUnknown:
		return true
	default:
		return false
	}
}

// isFailureReplicationJobState reports whether the job finished without creating the replication.
// Every state listed here must also be final in isFinalReplicationJobState.
func isFailureReplicationJobState(state apigen.GetReplicationJobResponseState) bool {
	switch state {
	case apigen.GetReplicationJobResponseStateFailed,
		apigen.GetReplicationJobResponseStateKilled,
		apigen.GetReplicationJobResponseStateSkipped,
		apigen.GetReplicationJobResponseStateNotfound,
		apigen.GetReplicationJobResponseStateUnknown:
		return true
	default:
		return false
	}
}

// waitForReplicationJob polls an asynchronous replication creation job until it reaches a terminal
// state, and returns the completed job so the caller can read both replicationId and
// reverseReplicationId without a second fetch. Those two fields are only populated once the job is
// complete.
//
// Errors from the job endpoint itself are treated as transient and retried, because the job
// outlives any single request; only a terminal job state or the expiry of the polling budget stops
// the loop.
func waitForReplicationJob(
	ctx context.Context,
	clientV2 *apigen.ClientWithResponses,
	orgUUID, projUUID, clusterUUID, jobUUID uuid.UUID,
) (*apigen.GetReplicationJobResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, replicationJobTimeout)
	defer cancel()

	ticker := time.NewTicker(replicationJobPollInterval)
	defer ticker.Stop()

	var (
		consecutiveNotFound int
		seenLive            bool
		lastJob             *apigen.GetReplicationJobResponse
	)

	for {
		response, err := clientV2.GetReplicationJobWithResponse(ctx, orgUUID, projUUID, clusterUUID, jobUUID)
		switch {
		case err != nil:
			tflog.Warn(ctx, "could not get replication job status, retrying", map[string]any{"error": err.Error()})

		case response.StatusCode() == http.StatusNotFound,
			response.JSON200 != nil && response.JSON200.State == apigen.GetReplicationJobResponseStateNotfound:
			// A job that was already observed running cannot legitimately vanish, so only the
			// initial read race is tolerated.
			consecutiveNotFound++
			if seenLive || consecutiveNotFound > maxConsecutiveJobNotFound {
				return nil, fmt.Errorf(
					"%w: job %s reported notfound %d time(s)",
					errors.ErrReplicationJobNotFound, jobUUID, consecutiveNotFound,
				)
			}
			tflog.Warn(ctx, "replication job not visible yet, retrying", map[string]any{"job_id": jobUUID.String()})

		case response.JSON200 == nil:
			tflog.Warn(ctx, "unexpected replication job response, retrying", map[string]any{
				"status": response.StatusCode(),
				"body":   string(response.Body),
			})

		default:
			job := response.JSON200
			lastJob = job
			consecutiveNotFound = 0
			seenLive = true

			switch {
			case !isFinalReplicationJobState(job.State):
				tflog.Info(ctx, "waiting for replication job to complete the execution", map[string]any{
					"job_id": job.JobId,
					"state":  string(job.State),
				})

			case isFailureReplicationJobState(job.State):
				return nil, fmt.Errorf(
					"%w: job %s, state %s, retry %d, last error: %s",
					errors.ErrReplicationJobFailed,
					job.JobId, job.State, replicationJobRetryNumber(job), replicationJobLastError(job),
				)

			// Final and not a failure leaves only complete, where the replication IDs are populated.
			default:
				if job.ReplicationId == nil || *job.ReplicationId == "" {
					return nil, fmt.Errorf(
						"%w: job %s completed without returning a replication ID",
						errors.ErrReplicationJobFailed, job.JobId,
					)
				}
				return job, nil
			}
		}

		select {
		case <-ctx.Done():
			return nil, replicationJobTimeoutError(jobUUID, lastJob, ctx.Err())
		case <-ticker.C:
		}
	}
}

// replicationJobTimeoutError reports the budget expiry along with whatever the last poll saw, which
// is the only diagnostic the API offers for a job that never settled.
func replicationJobTimeoutError(jobUUID uuid.UUID, lastJob *apigen.GetReplicationJobResponse, cause error) error {
	if lastJob == nil {
		return fmt.Errorf("%w: job %s, no job status was observed: %w", errors.ErrReplicationJobTimeout, jobUUID, cause)
	}
	return fmt.Errorf(
		"%w: job %s, last state %s, retry %d, last error: %s: %w",
		errors.ErrReplicationJobTimeout,
		lastJob.JobId, lastJob.State, replicationJobRetryNumber(lastJob), replicationJobLastError(lastJob), cause,
	)
}

func replicationJobRetryNumber(job *apigen.GetReplicationJobResponse) int {
	if job.RetryNumber == nil {
		return 0
	}
	return *job.RetryNumber
}

func replicationJobLastError(job *apigen.GetReplicationJobResponse) string {
	if job.LastError == nil {
		return "none"
	}
	return *job.LastError
}
