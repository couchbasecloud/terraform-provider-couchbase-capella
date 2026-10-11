package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// replicationJobPollInterval and replicationJobTimeout are overridden in tests to avoid waiting on
// real wall-clock time while exercising waitForReplicationJob's polling loop.
var (
	replicationJobPollInterval = 15 * time.Second
	replicationJobTimeout      = 30 * time.Minute
)

// maxConsecutiveJobNotFound is how many notfound states in a row are allowed before the job was
// seen running. The job record is not always readable at once after the create call returns 202.
const maxConsecutiveJobNotFound = 3

// waitForReplicationJob polls an async replication create job until it stops, and returns the
// completed job. replicationId and reverseReplicationId are set only on a completed job.
//
// Server errors and network errors are logged and polled again, because the job keeps running on
// the server. Client errors (4xx) stop the poll.
func waitForReplicationJob(
	ctx context.Context,
	data *providerschema.Data,
	organizationId, projectId, clusterId, jobId string,
) (*replicationapi.GetReplicationJobResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, replicationJobTimeout)
	defer cancel()

	ticker := time.NewTicker(replicationJobPollInterval)
	defer ticker.Stop()

	var (
		consecutiveNotFound int
		seenAlive           bool
		lastJob             *replicationapi.GetReplicationJobResponse
	)

	for {
		job, err := getReplicationJob(ctx, data, organizationId, projectId, clusterId, jobId)
		switch {
		case err != nil:
			if isClientError(err) {
				return nil, fmt.Errorf("could not get replication job %s: %w", jobId, err)
			}
			tflog.Warn(ctx, "could not get replication job status, retrying", map[string]interface{}{
				"job_id": jobId,
				"error":  err.Error(),
			})

		case job.State == replicationapi.JobStateNotFound:
			// A job that was seen running cannot disappear, so only the first read race is allowed.
			if seenAlive {
				return nil, fmt.Errorf(
					"%w: job %s was running, then reported notfound",
					errors.ErrReplicationJobNotFound, jobId,
				)
			}
			consecutiveNotFound++
			if consecutiveNotFound > maxConsecutiveJobNotFound {
				return nil, fmt.Errorf(
					"%w: job %s reported notfound %d times in a row and never started",
					errors.ErrReplicationJobNotFound, jobId, consecutiveNotFound,
				)
			}
			tflog.Warn(ctx, "replication job is not visible yet, retrying", map[string]interface{}{
				"job_id": jobId,
			})

		case replicationapi.IsFailedJobState(job.State):
			return nil, fmt.Errorf(
				"%w: job %s, state %s, retry %d, last error: %s",
				errors.ErrReplicationJobFailed, jobId, job.State, job.RetryNumber, lastErrorOrNone(job),
			)

		case job.State == replicationapi.JobStateComplete:
			if job.ReplicationId == "" {
				return nil, fmt.Errorf("%w: job %s", errors.ErrReplicationJobMissingId, jobId)
			}
			return job, nil

		default:
			consecutiveNotFound = 0
			seenAlive = true
			lastJob = job
			tflog.Info(ctx, "waiting for replication job to complete", map[string]interface{}{
				"job_id": jobId,
				"state":  job.State,
			})
		}

		select {
		case <-ctx.Done():
			return nil, replicationJobTimeoutError(jobId, lastJob, ctx.Err())
		case <-ticker.C:
		}
	}
}

func getReplicationJob(
	ctx context.Context,
	data *providerschema.Data,
	organizationId, projectId, clusterId, jobId string,
) (*replicationapi.GetReplicationJobResponse, error) {
	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationJobURL(data.HostURL, organizationId, projectId, clusterId, jobId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := data.ClientV1.ExecuteWithRetry(ctx, cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}

	var job replicationapi.GetReplicationJobResponse
	if err := json.Unmarshal(response.Body, &job); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &job, nil
}

// isClientError reports whether err is an API error with a 4xx status. ExecuteWithRetry already
// retries 429, so any 4xx that reaches the caller is final.
func isClientError(err error) bool {
	var apiError *api.Error
	return stderrors.As(err, &apiError) &&
		apiError.HttpStatusCode >= http.StatusBadRequest &&
		apiError.HttpStatusCode < http.StatusInternalServerError
}

// replicationJobTimeoutError includes the last state seen, because the API gives no other
// information about a job that did not finish.
func replicationJobTimeoutError(jobId string, lastJob *replicationapi.GetReplicationJobResponse, cause error) error {
	if lastJob == nil {
		return fmt.Errorf("%w: job %s, no job status was seen: %w", errors.ErrReplicationJobTimeout, jobId, cause)
	}
	return fmt.Errorf(
		"%w: job %s, last state %s, retry %d, last error: %s: %w",
		errors.ErrReplicationJobTimeout, jobId, lastJob.State, lastJob.RetryNumber, lastErrorOrNone(lastJob), cause,
	)
}

func lastErrorOrNone(job *replicationapi.GetReplicationJobResponse) string {
	if job.LastError == "" {
		return "none"
	}
	return job.LastError
}
