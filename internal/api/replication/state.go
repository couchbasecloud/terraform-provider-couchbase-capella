package replication

import (
	"slices"
	"strings"
)

// JobState is the state of an async replication create job.
type JobState string

const (
	// JobStatePending and JobStateProcessing communicate that the job is still running.
	JobStatePending    JobState = "pending"
	JobStateProcessing JobState = "processing"

	// JobStateComplete communicates that the replication was created.
	JobStateComplete JobState = "complete"

	// JobStateFailed, JobStateSkipped, JobStateKilled, JobStateNotFound and
	// JobStateUnknown communicate that the job stopped without creating the replication.
	JobStateFailed   JobState = "failed"
	JobStateSkipped  JobState = "skipped"
	JobStateKilled   JobState = "killed"
	JobStateNotFound JobState = "notfound"
	JobStateUnknown  JobState = "unknown"
)

// IsFailedJobState reports whether the job stopped without creating the replication.
func IsFailedJobState(state JobState) bool {
	failedStates := []JobState{
		JobStateFailed,
		JobStateSkipped,
		JobStateKilled,
		JobStateNotFound,
		JobStateUnknown,
	}
	return slices.Contains(failedStates, state)
}

// Capella does not retry a job whose last error has one of these prefixes.
const (
	fatalJobErrorPrefix       = "job-fatal:"
	goxdcrFatalJobErrorPrefix = "goxdcr-fatal-error:"
)

// IsFatalJobError reports whether Capella will not retry a failed job.
func IsFatalJobError(lastError string) bool {
	return strings.Contains(lastError, fatalJobErrorPrefix) || strings.Contains(lastError, goxdcrFatalJobErrorPrefix)
}

// Status is the status of a replication.
type Status string

const (
	// StatusPending and StatusPausing communicate that the replication is
	// moving to running or paused.
	StatusPending Status = "pending"
	StatusPausing Status = "pausing"

	StatusRunning Status = "running"
	StatusPaused  Status = "paused"
	StatusFailed  Status = "failed"
)
