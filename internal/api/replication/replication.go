// Package replication provides the API types for replication operations.
package replication

import (
	"time"
)

// AuditData contains the audit data for replication (simpler than full CouchbaseAuditData).
type AuditData struct {
	// CreatedAt The RFC3339 timestamp associated with when the replication was initially created.
	CreatedAt time.Time `json:"createdAt"`

	// CreatedBy The user who created the replication.
	CreatedBy string `json:"createdBy"`
}

// GetReplicationSummaryResponse represents a replication summary item from the list replications endpoint.
type GetReplicationSummaryResponse struct {
	// Audit contains the audit data for the replication.
	Audit AuditData `json:"audit"`

	// Direction specifies the replication flow — whether it's oneWay (source to target only) or twoWay (also from target back to source).
	Direction *string `json:"direction,omitempty"`

	// Id is the ID of the specified replication.
	Id string `json:"id"`

	// SourceCluster is the name of the source cluster.
	SourceCluster string `json:"sourceCluster"`

	// Status is the status of the replication (pending, pausing, failed, paused, running).
	Status string `json:"status"`

	// TargetCluster is the name of the target cluster.
	TargetCluster string `json:"targetCluster"`
}
