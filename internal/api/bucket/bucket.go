package bucket

// CreateBucketRequest is the payload passed to V4 Capella Public API to create a bucket in a Capella cluster.
// Creates a new bucket configuration under a cluster.
//
// To learn more about bucket configuration, see https://docs.couchbase.com/server/current/learn/buckets-memory-and-storage/buckets.html.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateBucketRequest struct {
	// Type represents the Bucket Type
	// Default: "couchbase"
	// Enum: "couchbase" "ephemeral"
	//
	// If selected Ephemeral, it is not eligible for imports or App Endpoints creation. This field cannot be changed later.
	// The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Type *string `json:"type"`

	// StorageBackend represents the storage engine used for the bucket.
	// Default: "couchstore"
	// Enum: "couchstore" "magma"
	// Ephemeral buckets do not support StorageBackend, hence not applicable for Ephemeral buckets and throws an error if this field is added.
	// This field is only applicable for a Couchbase bucket. The default value mentioned (Couchstore) is for Couchbase bucket.
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/storage-engines.html
	StorageBackend *string `json:"storageBackend"`

	// MemoryAllocationInMb is the amount of memory to allocate for the bucket memory in MiB
	// Default: 100
	// This is the maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation.
	//
	// The default value (100MiB) mentioned is for Couchbase type buckets with Couchstore as the Storage Backend.
	//
	// For Couchbase buckets, the default and minimum memory allocation changes according to the Storage Backend type as follows:
	// For Couchstore, the default and minimum memory allocation is 100 MiB.
	// For Magma, the default and minimum memory allocation is 1024 MiB.
	// For Ephemeral buckets, the default and minimum memory allocation is 100 MiB.
	MemoryAllocationInMb *int64 `json:"memoryAllocationInMb"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// Default: "seqno"
	// Enum: "seqno" "lww"
	// This field may be referred to as "conflict resolution" in the Couchbase documentation, and seqno and lww may be
	// referred to as "sequence number" and "timestamp" respectively.
	//
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution *string `json:"bucketConflictResolution"`

	// DurabilityLevel is the minimum level at which all writes to the bucket must occur.
	// Default: "none"
	// Enum: "none" "majority" "majorityAndPersistActive" "persistToMajority"
	//
	// The options for Durability level are as follows, according to the bucket type.
	//
	// For a Couchbase bucket:
	// None
	// Replicate to Majority
	// Majority and Persist to Active
	// Persist to Majority
	//
	// For an Ephemeral bucket:
	// None
	// Replicate to Majority
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	DurabilityLevel *string `json:"durabilityLevel"`

	// Replicas states the number of replica nodes for the bucket.
	// Default: 1
	// Enum: 1 2 3
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Replicas *int64 `json:"replicas"`

	// Flush determines whether flushing is enabled on the bucket.
	// Default: false
	// Enable Flush to delete all items in this bucket at the earliest opportunity.
	// Disable Flush to avoid inadvertent data loss.
	Flush *bool `json:"flush"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	// Default: 0
	// This is the maximum time to live for items in the bucket.
	// Default is 0, that means TTL is disabled. This is a non-negative value.
	TimeToLiveInSeconds *int64 `json:"timeToLiveInSeconds"`

	// EvictionPolicy is the policy which Capella adopts to prevent data loss due to memory exhaustion.
	// This may be also known as Ejection Policy in the Couchbase documentation.
	//
	// Default: "fullEviction"
	// Enum: "fullEviction" "noEviction" "nruEviction"
	//
	// For Couchbase bucket, Eviction Policy is fullEviction by default.
	// For Ephemeral buckets, Eviction Policy is a required field, and should be one of the following:
	// noEviction
	// nruEviction
	//
	//To learn more, see https://docs.couchbase.com/server/current/rest-api/rest-bucket-create.html#evictionpolicy
	EvictionPolicy *string `json:"evictionPolicy"`

	// Name is the name of the bucket (up to 100 characters).
	// This field cannot be changed later. The name should be according to the following rules:
	// Characters used for the name should be in the ranges of A-Z, a-z, and 0-9; plus the underscore, period, dash, and percent characters.
	// The name can be a maximum of 100 characters in length.
	// The name cannot have 0 characters or empty. Minimum length of name is 1.
	// The name cannot start with a . (period).
	Name string `json:"name"`
}

// CreateBucketResponse is the response received from Capella V4 Public API on requesting to create a new bucket.
// Common response codes: 201, 403, 422, 429, 500.
type CreateBucketResponse struct {
	// Id is unique ID of the bucket created.
	Id string `json:"id"`
}

// GetBucketResponse is the response received from Capella V4 Public API on requesting to information about an existing bucket.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetBucketResponse struct {
	Stats *Stats `json:"stats"`

	// Id is the ID of the bucket created.
	Id string `json:"id"`

	// Name is the name of the cluster (up to 100 characters).
	Name string `json:"name"`

	// Type represents the Bucket Type
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Type string `json:"type"`

	// StorageBackend represents the storage engine used for the bucket.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/storage-engines.html
	StorageBackend string `json:"storageBackend"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution string `json:"bucketConflictResolution"`

	// DurabilityLevel is the minimum level at which all writes to the bucket must occur.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	DurabilityLevel string `json:"durabilityLevel"`

	// EvictionPolicy is the policy which Capella adopts to prevent data loss due to memory exhaustion.
	//To learn more, see https://docs.couchbase.com/server/current/rest-api/rest-bucket-create.html#evictionpolicy
	EvictionPolicy string `json:"evictionPolicy"`

	// MemoryAllocationInMb is the amount of memory to allocate for the bucket memory in MiB
	MemoryAllocationInMb int64 `json:"memoryAllocationInMb"`

	// Replicas states the number of replica nodes for the bucket.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Replicas int64 `json:"replicas"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	TimeToLiveInSeconds int64 `json:"timeToLiveInSeconds"`

	// Flush determines whether flushing is enabled on the bucket.
	Flush bool `json:"flush"`
}

// PutBucketRequest is the request payload sent to the Capella V4 Public API in order to update an existing bucket.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type PutBucketRequest struct {
	// DurabilityLevel is the minimum level at which all writes to the bucket must occur.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	DurabilityLevel string `json:"durabilityLevel"`

	// MemoryAllocationInMb is the amount of memory to allocate for the bucket memory in MiB
	MemoryAllocationInMb int64 `json:"memoryAllocationInMb"`

	// Replicas states the number of replica nodes for the bucket.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Replicas int64 `json:"replicas"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	TimeToLiveInSeconds int64 `json:"timeToLiveInSeconds"`

	// Flush determines whether flushing is enabled on the bucket.
	Flush bool `json:"flush"`
}
