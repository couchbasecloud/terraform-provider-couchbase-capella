package bucket

type CreateBucketRequest struct {

	// Name is the name of the cluster (up to 100 characters).
	Name string `json:"name"`

	// Type represents the Bucket Type
	Type string `json:"type"`

	// StorageBackend represents the storage engine used for the bucket.
	StorageBackend string `json:"storageBackend"`

	// MemoryAllocationInMb The amount of memory to allocate for the bucket memory in MiB
	MemoryAllocationInMb int `json:"memoryAllocationInMb"`

	BucketConflictResolution string `json:"bucketConflictResolution"`

	DurabilityLevel string `json:"durabilityLevel"`

	Replicas int `json:"replicas"`

	Flush bool `json:"flush"`

	TimeToLiveInSeconds int `json:"timeToLiveInSeconds"`

	EvictionPolicy string `json:"evictionPolicy"`
}

// CreateBucketResponse defines model for CreateBucketResponse.
type CreateBucketResponse struct {
	// Id The ID of the bucket created.
	Id string `json:"id"`
}

type GetBucketResponse struct {

	// Id is the ID of the bucket created.
	Id string `json:"id"`

	// Name is the name of the cluster (up to 100 characters).
	Name string `json:"name"`

	// Type represents the Bucket Type
	Type string `json:"type"`

	// StorageBackend represents the storage engine used for the bucket.
	StorageBackend string `json:"storageBackend"`

	// MemoryAllocationInMb The amount of memory to allocate for the bucket memory in MiB
	MemoryAllocationInMb int `json:"memoryAllocationInMb"`

	BucketConflictResolution string `json:"bucketConflictResolution"`

	DurabilityLevel string `json:"durabilityLevel"`

	Replicas int `json:"replicas"`

	Flush bool `json:"flush"`

	TimeToLiveInSeconds int `json:"timeToLiveInSeconds"`

	EvictionPolicy string `json:"evictionPolicy"`
}
