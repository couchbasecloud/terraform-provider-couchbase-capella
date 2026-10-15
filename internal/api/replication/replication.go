package replication

import "time"

// ModeAsync makes the create call return a job ID at once instead of waiting
// for the replication to be created. The provider always uses it.
const ModeAsync = "async"

// CreateReplicationRequest is the payload sent to the V4 Capella Public API to
// create an XDCR replication from a bucket on the source cluster.
type CreateReplicationRequest struct {
	// SourceBucket is the bucket ID (URL-safe base64 of the bucket name), not the name.
	SourceBucket string                  `json:"sourceBucket"`
	Target       CreateReplicationTarget `json:"target"`

	// Mode is "async" or "sync". The provider always sends "async".
	Mode string `json:"mode"`

	// Direction defaults to "oneWay". Enum: "oneWay" "twoWay".
	Direction *string `json:"direction,omitempty"`

	// Priority defaults to "high". Enum: "low" "medium" "high".
	Priority *string `json:"priority,omitempty"`

	// NetworkUsageLimit is in MiB per second. 0 (the default) means no limit.
	NetworkUsageLimit *int64 `json:"networkUsageLimit,omitempty"`

	// Mappings limits the replication to specific scopes and collections.
	// When omitted, all scopes are replicated.
	Mappings []ScopeMapping `json:"mappings,omitempty"`

	Filter *Filter `json:"filter,omitempty"`

	// MobileSyncMode "active" turns on cross cluster versioning for both buckets. Cross
	// cluster versioning cannot be turned off. Enum: "off" "active".
	MobileSyncMode *string `json:"mobileSyncMode,omitempty"`

	// ConflictLogging needs Couchbase Server 8.0 or later and cross cluster versioning on both buckets.
	ConflictLogging *ConflictLogging `json:"conflictLogging,omitempty"`
}

// CreateReplicationTarget identifies the target of a replication in a create request.
type CreateReplicationTarget struct {
	// Cluster is the target cluster ID for a Capella target,
	// or the cluster reference name for an external target.
	Cluster string `json:"cluster"`

	// Bucket is the target bucket ID for a Capella target,
	// or the bucket name for an external target.
	Bucket string `json:"bucket"`

	// Type defaults to "capella". Enum: "capella" "external".
	Type *string `json:"type,omitempty"`
}

// ScopeMapping maps a source scope to a target scope.
type ScopeMapping struct {
	SourceScope string `json:"sourceScope"`
	TargetScope string `json:"targetScope"`

	// Collections limits the mapping to specific collections.
	// When omitted, all collections in the scope are replicated.
	Collections []CollectionMapping `json:"collections,omitempty"`
}

// CollectionMapping maps a source collection to a target collection.
type CollectionMapping struct {
	SourceCollection string `json:"sourceCollection"`
	TargetCollection string `json:"targetCollection"`
}

// Filter limits which documents are replicated.
//
// On update the server sets only the keys that are present, so a request must
// send the full shape to clear a part that was removed. On GET the server drops
// expressions when regEx is empty, drops documentExcludeOptions when all four
// flags are false, and never returns skipRestream.
type Filter struct {
	DocumentExcludeOptions *DocumentExcludeOptions `json:"documentExcludeOptions,omitempty"`
	Expressions            *FilterExpressions      `json:"expressions,omitempty"`
}

// DocumentExcludeOptions selects document types that are not replicated.
type DocumentExcludeOptions struct {
	Deletion   bool `json:"deletion"`
	Expiration bool `json:"expiration"`
	TTL        bool `json:"ttl"`
	Binary     bool `json:"binary"`
}

// FilterExpressions holds the filter expression for document keys and content.
type FilterExpressions struct {
	// RegEx is the filter expression. An empty value clears the expression.
	RegEx string `json:"regEx"`

	// SkipRestream is only sent; GET does not return it.
	SkipRestream bool `json:"skipRestream"`
}

// UpdateReplicationRequest is the payload sent to the V4 Capella Public API to
// update a replication. A nil field is left unchanged by the server.
type UpdateReplicationRequest struct {
	Priority          *string `json:"priority,omitempty"`
	NetworkUsageLimit *int64  `json:"networkUsageLimit,omitempty"`
	Filter            *Filter `json:"filter,omitempty"`

	// AllScopes resets the replication to all scopes. Omitting Mappings alone leaves
	// the existing mappings unchanged. It must not be sent together with Mappings.
	AllScopes *bool          `json:"allScopes,omitempty"`
	Mappings  []ScopeMapping `json:"mappings,omitempty"`

	// MobileSyncMode "active" fails with 412 when cross cluster versioning is off on a bucket.
	MobileSyncMode *string `json:"mobileSyncMode,omitempty"`

	// ConflictLogging replaces the whole configuration. An empty object clears it.
	ConflictLogging *ConflictLogging `json:"conflictLogging,omitempty"`
}

// ConflictLogging is the conflict logging configuration in a create or update request.
// Disabled, Bucket and Collection must be sent together. An empty object means no conflict logging.
type ConflictLogging struct {
	// Disabled is true to keep the configuration without logging conflicts.
	Disabled *bool `json:"disabled,omitempty"`

	// Bucket is the bucket ID of the log bucket, not the name.
	Bucket *string `json:"bucket,omitempty"`

	// Collection is the default log target, in "scope.collection" format.
	Collection *string `json:"collection,omitempty"`

	// LoggingRules is keyed by a source scope or a source "scope.collection". A nil rule
	// is sent as null and excludes the namespace from conflict logging.
	LoggingRules map[string]*ConflictLoggingRule `json:"loggingRules,omitempty"`
}

// ConflictLoggingRule is the log target for one source namespace. An empty rule logs to the
// default bucket and collection. Bucket and Collection must be set together.
type ConflictLoggingRule struct {
	Bucket     *string `json:"bucket,omitempty"`
	Collection *string `json:"collection,omitempty"`
}

// GetConflictLogging is the conflict logging configuration in a GET response.
type GetConflictLogging struct {
	Disabled     bool                            `json:"disabled"`
	Bucket       string                          `json:"bucket"`
	Collection   string                          `json:"collection"`
	LoggingRules map[string]*ConflictLoggingRule `json:"loggingRules"`
}

// CreateReplicationAsyncResponse is the 202 response of an async create call.
type CreateReplicationAsyncResponse struct {
	JobId string `json:"jobId"`
}

// GetReplicationJobResponse is the status of an async replication create job.
type GetReplicationJobResponse struct {
	JobId string   `json:"jobId"`
	State JobState `json:"state"`

	// RetryNumber is 0 on the first attempt.
	RetryNumber int    `json:"retryNumber"`
	LastError   string `json:"lastError"`

	// ReplicationId and ReverseReplicationId are set only when the state is complete.
	// ReverseReplicationId is set only for a twoWay replication.
	ReplicationId        string `json:"replicationId"`
	ReverseReplicationId string `json:"reverseReplicationId"`

	LastUpdatedTimestamp time.Time `json:"lastUpdatedTimestamp"`
}

// GetReplicationResponse is the response of a GET on one replication.
type GetReplicationResponse struct {
	Id          string `json:"id"`
	Status      Status `json:"status"`
	ChangesLeft int64  `json:"changesLeft"`
	Error       string `json:"error"`

	Source SourceDetails `json:"source"`
	Target TargetDetails `json:"target"`

	Mappings          []ScopeMapping `json:"mappings"`
	Direction         string         `json:"direction"`
	Priority          string         `json:"priority"`
	NetworkUsageLimit int64          `json:"networkUsageLimit"`

	Filter *Filter `json:"filter"`

	MobileSyncMode string `json:"mobileSyncMode"`

	ConflictLogging *GetConflictLogging `json:"conflictLogging"`

	Audit ReplicationAuditData `json:"audit"`
}

// SourceDetails describes the source of a replication.
type SourceDetails struct {
	Project Reference       `json:"project"`
	Cluster Reference       `json:"cluster"`
	Bucket  BucketReference `json:"bucket"`
}

// TargetDetails describes the target of a replication.
//
// For an external target, Cluster.Name is the cluster reference name and
// Bucket.Name is the bucket name. Project is empty.
type TargetDetails struct {
	Project Reference       `json:"project"`
	Cluster Reference       `json:"cluster"`
	Bucket  BucketReference `json:"bucket"`
	Type    string          `json:"type"`
}

// Reference is the ID and name of a project or a cluster.
type Reference struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// BucketReference is the ID, name and conflict resolution type of a bucket.
type BucketReference struct {
	Id                     string `json:"id"`
	Name                   string `json:"name"`
	ConflictResolutionType string `json:"conflictResolutionType"`
}

// ReplicationAuditData is the audit data of a replication. Unlike
// api.CouchbaseAuditData, the replication endpoints return only these two fields.
type ReplicationAuditData struct {
	// CreatedBy is a user UUID, or a name such as "internal-support" for internal users.
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// ReplicationSummary is one entry of a replication list.
type ReplicationSummary struct {
	Id            string               `json:"id"`
	SourceCluster string               `json:"sourceCluster"`
	TargetCluster string               `json:"targetCluster"`
	Status        Status               `json:"status"`
	Direction     string               `json:"direction"`
	Audit         ReplicationAuditData `json:"audit"`
}
