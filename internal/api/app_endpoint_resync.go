package api

import (
	"time"
)

// CreateResyncRequest represents the request body for creating a resync operation.
// It contains a mapping of scopes to their respective collections that need to be resynced.
type CreateResyncRequest struct {
	// Scopes is a map where each key represents a scope name and the value
	// is a slice of collection names within that scope to be resynced.
	// Example:
	//   {
	//     "scope1": ["collection1", "collection2"],
	//     "scope2": ["collection3"]
	//   }
	Scopes map[string][]string `json:"scopes"`
}

// CreateResyncResponse represents the response returned after initiating a resync operation.
// It contains the status and progress information of the resync process.
type CreateResyncResponse struct {
	// CollectionsProcessing contains a map of collections currently being processed,
	// organized by scope. This field is optional and may be nil.
	CollectionsProcessing map[string][]string `json:"collections_processing,omitempty"`

	// DocsChanged represents the number of documents that have been changed
	// during the resync operation.
	DocsChanged int64 `json:"docsChanged"`

	// DocsProcessed represents the total number of documents that have been
	// processed during the resync operation.
	DocsProcessed int64 `json:"docsProcessed"`

	// LastError contains the last error message encountered during the resync
	// operation, if any. Empty string indicates no errors.
	LastError string `json:"lastError"`

	// StartTime represents the timestamp when the resync operation was initiated.
	StartTime time.Time `json:"startTime"`

	// State indicates the current state of the resync operation
	// (e.g., "running", "completed", "error", etc).
	State ResyncStatusState `json:"state"`
}

// ResyncStatusState represents the current state of a resync operation.
type ResyncStatusState string

// Defines values for ResyncStatusState.
const (
	ResyncStatusStateCompleted ResyncStatusState = "completed"
	ResyncStatusStateError     ResyncStatusState = "error"
	ResyncStatusStateRunning   ResyncStatusState = "running"
	ResyncStatusStateStopped   ResyncStatusState = "stopped"
	ResyncStatusStateStopping  ResyncStatusState = "stopping"
)
