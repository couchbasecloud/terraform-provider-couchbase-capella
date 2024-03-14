package sample_bucket

// CreateSampleBucketRequest is the payload passed to V4 Capella Public API to load a sample bucket in a Capella cluster.
// Loads a new sample bucket configuration under a cluster.
//
// To learn more about bucket configuration, see https://docs.couchbase.com/server/current/manage/manage-settings/install-sample-buckets.html.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateSampleBucketRequest struct {
	// Name is the name of the sample bucket (up to 100 characters).
	// Enum: travel-sample, beer-sample, gamesim-sample
	// This field cannot be changed later.
	Name string `json:"name"`
}

// CreateBucketSampleResponse is the response received from Capella V4 Public API on requesting to load a new sample bucket.
// Common response codes: 201, 403, 422, 429, 500.
type CreateSampleBucketResponse struct {
	// Id is unique ID of the sample bucket created.
	Id string `json:"bucketId"`

	// Name is the name of the cluster (up to 100 characters).
	Name string `json:"name"`
}

// GetSampleBucketResponse is the response received from Capella V4 Public API on requesting to information about an existing sample bucket.
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
type GetSampleBucketResponse struct {
	Stats *Stats `json:"stats"`

	// Id is the ID of the bucket created.
	Id string `json:"id"`

	// Name is the name of the cluster (up to 100 characters).
	Name string `json:"name"`

	// Type represents the sample Bucket Type
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Type string `json:"type"`

	// StorageBackend represents the storage engine used for the sample bucket.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/storage-engines.html
	StorageBackend string `json:"storageBackend"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution string `json:"bucketConflictResolution"`

	// DurabilityLevel is the minimum level at which all writes to the sample bucket must occur.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	DurabilityLevel string `json:"durabilityLevel"`

	// EvictionPolicy is the policy which Capella adopts to prevent data loss due to memory exhaustion.
	//To learn more, see https://docs.couchbase.com/server/current/rest-api/rest-bucket-create.html#evictionpolicy
	EvictionPolicy string `json:"evictionPolicy"`

	// MemoryAllocationInMb is the amount of memory to allocate for the sample bucket memory in MiB
	MemoryAllocationInMb int64 `json:"memoryAllocationInMb"`

	// Replicas states the number of replica nodes for the sample bucket.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Replicas int64 `json:"replicas"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	TimeToLiveInSeconds int64 `json:"timeToLiveInSeconds"`

	// Flush determines whether flushing is enabled on the sample bucket.
	Flush bool `json:"flush"`
}
