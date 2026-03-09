package replication

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// GetReplicationResponse represents the response from the GET replication endpoint.
// Fetches the details of the given replication.
type GetReplicationResponse struct {
	// Id is the ID of the specified replication.
	Id string `json:"id"`

	// Status is the status of the replication (pending, pausing, failed, paused, running).
	Status string `json:"status"`

	// ChangesLeft is the number of remaining mutations to the replication.
	ChangesLeft int64 `json:"changesLeft"`

	// Error is the error message if the replication has failed.
	Error string `json:"error,omitempty"`

	// Source contains all the metadata about a replication source.
	Source Source `json:"source"`

	// Target contains all the metadata about a replication target.
	Target Target `json:"target"`

	// Mappings defines mappings from source to target scopes and collections.
	Mappings []Mapping `json:"mappings,omitempty"`

	// Direction specifies the replication flow — whether it's oneWay (source to target only) or twoWay (also from target back to source).
	Direction string `json:"direction"`

	// Priority represents the resource allocation to the replication (low, medium, high).
	Priority string `json:"priority,omitempty"`

	// NetworkUsageLimit is the network usage limit in MiB per second. 0 means unlimited.
	NetworkUsageLimit int64 `json:"networkUsageLimit"`

	// Filter contains the replication settings which are passed to the Couchbase server API while creating a replication.
	Filter *Filter `json:"filter,omitempty"`

	// Audit contains the audit data for the replication.
	Audit ReplicationAuditData `json:"audit"`
}

// ReplicationSummary represents a replication summary item in the list response.
type ReplicationSummary struct {
	// Id is the ID of the specified replication.
	Id string `json:"id"`

	// SourceCluster is the name of the source cluster.
	SourceCluster string `json:"sourceCluster"`

	// TargetCluster is the name of the target cluster.
	TargetCluster string `json:"targetCluster"`

	// Status is the status of the replication (pending, pausing, failed, paused, running).
	Status string `json:"status"`

	// Direction specifies the replication flow — whether it's oneWay or twoWay.
	Direction string `json:"direction,omitempty"`

	// Audit contains the audit data for the replication.
	Audit ReplicationAuditData `json:"audit"`
}

// Source contains all the metadata about a replication source.
type Source struct {
	// Project contains the metadata for the source from which the replication is established.
	Project ProjectReference `json:"project"`

	// Cluster contains the metadata for the source from which the replication is established.
	Cluster ClusterReference `json:"cluster"`

	// Bucket contains the metadata for the source from which the replication is established.
	Bucket BucketReference `json:"bucket"`

	// Scopes contains scope and collection details for the source.
	Scopes []Scope `json:"scopes,omitempty"`

	// Type tells us if the source cluster is capella or external.
	Type string `json:"type,omitempty"`
}

// Target contains all the metadata about a replication target.
type Target struct {
	// Project contains the metadata for the destination to which the replication is established.
	Project *ProjectReference `json:"project,omitempty"`

	// Cluster contains the metadata for the destination to which the replication is established.
	Cluster ClusterReference `json:"cluster"`

	// Bucket contains the metadata for the destination to which the replication is established.
	Bucket BucketReference `json:"bucket"`

	// Scopes contains scope and collection details for the target.
	Scopes []Scope `json:"scopes,omitempty"`

	// Type tells us if the target cluster is capella or external.
	Type string `json:"type,omitempty"`
}

// ProjectReference contains project metadata.
type ProjectReference struct {
	// Id is the project ID.
	Id string `json:"id"`

	// Name is the project name.
	Name string `json:"name"`
}

// ClusterReference contains cluster metadata.
type ClusterReference struct {
	// Id is the cluster ID.
	Id string `json:"id"`

	// Name is the cluster name.
	Name string `json:"name"`
}

// BucketReference contains bucket metadata.
type BucketReference struct {
	// Id is the bucket ID.
	Id string `json:"id"`

	// Name is the bucket name.
	Name string `json:"name"`

	// ConflictResolutionType is the conflict resolution type (lww, seqno, custom).
	ConflictResolutionType string `json:"conflictResolutionType,omitempty"`
}

// Scope contains scope and collection details.
type Scope struct {
	// Name is the scope name.
	Name string `json:"name"`

	// Collections is the list of collection names under this scope.
	Collections []string `json:"collections,omitempty"`
}

// Mapping defines a mapping from source to target scope/collection.
type Mapping struct {
	// SourceScope is the name of the source scope.
	SourceScope string `json:"sourceScope"`

	// TargetScope is the name of the target scope.
	TargetScope string `json:"targetScope"`

	// Collections defines mappings between source and target collections.
	Collections []CollectionMapping `json:"collections,omitempty"`
}

// CollectionMapping defines a mapping between source and target collections.
type CollectionMapping struct {
	// SourceCollection is the name of the source collection.
	SourceCollection string `json:"sourceCollection"`

	// TargetCollection is the name of the target collection.
	TargetCollection string `json:"targetCollection"`
}

// Filter contains the replication settings which are passed to the Couchbase server API.
type Filter struct {
	// DocumentExcludeOptions specifies which document types to filter out.
	DocumentExcludeOptions *DocumentExcludeOptions `json:"documentExcludeOptions,omitempty"`

	// Expressions contains filter expressions.
	Expressions *FilterExpressions `json:"expressions,omitempty"`
}

// DocumentExcludeOptions specifies document types to filter out.
type DocumentExcludeOptions struct {
	// Deletion when true, deletions are filtered out.
	Deletion bool `json:"deletion,omitempty"`

	// Expiration when true, expirations are filtered out.
	Expiration bool `json:"expiration,omitempty"`

	// Ttl when true, TTL value is removed from the replicated documents.
	Ttl bool `json:"ttl,omitempty"`

	// Binary when true, binary documents are filtered.
	Binary bool `json:"binary,omitempty"`
}

// FilterExpressions contains filter expression settings.
type FilterExpressions struct {
	// RegEx is the filter expression to match documents.
	RegEx string `json:"regEx,omitempty"`
}

// ReplicationAuditData contains audit information for a replication.
type ReplicationAuditData struct {
	// CreatedBy is the user who created the replication.
	CreatedBy string `json:"createdBy"`

	// CreatedAt is the RFC3339 timestamp when the replication was created.
	CreatedAt string `json:"createdAt"`

	// ModifiedBy is the user who last modified the replication.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// ModifiedAt is the RFC3339 timestamp when the replication was last modified.
	ModifiedAt string `json:"modifiedAt,omitempty"`

	// Version is the version number of the replication.
	Version int64 `json:"version,omitempty"`
}

// CouchbaseAuditData alias for internal api audit type
// This allows us to use the api.CouchbaseAuditData in the list response
type CouchbaseAuditData api.CouchbaseAuditData
