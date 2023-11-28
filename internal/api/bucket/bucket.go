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
	Type                     *string `json:"type"`
	StorageBackend           *string `json:"storageBackend"`
	MemoryAllocationInMb     *int64  `json:"memoryAllocationInMb"`
	BucketConflictResolution *string `json:"bucketConflictResolution"`
	DurabilityLevel          *string `json:"durabilityLevel"`
	Replicas                 *int64  `json:"replicas"`
	Flush                    *bool   `json:"flush"`
	TimeToLiveInSeconds      *int64  `json:"timeToLiveInSeconds"`
	EvictionPolicy           *string `json:"evictionPolicy"`
	Name                     string  `json:"name"`
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
	Stats                    *Stats `json:"stats"`
	Id                       string `json:"id"`
	Name                     string `json:"name"`
	Type                     string `json:"type"`
	StorageBackend           string `json:"storageBackend"`
	BucketConflictResolution string `json:"bucketConflictResolution"`
	DurabilityLevel          string `json:"durabilityLevel"`
	EvictionPolicy           string `json:"evictionPolicy"`
	MemoryAllocationInMb     int64  `json:"memoryAllocationInMb"`
	Replicas                 int64  `json:"replicas"`
	TimeToLiveInSeconds      int64  `json:"timeToLiveInSeconds"`
	Flush                    bool   `json:"flush"`
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
	DurabilityLevel      string `json:"durabilityLevel"`
	MemoryAllocationInMb int64  `json:"memoryAllocationInMb"`
	Replicas             int64  `json:"replicas"`
	TimeToLiveInSeconds  int64  `json:"timeToLiveInSeconds"`
	Flush                bool   `json:"flush"`
}
