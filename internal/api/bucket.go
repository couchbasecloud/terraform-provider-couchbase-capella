package api

type CreateBucketRequest struct {
	Name                     string `json:"name"`
	Type                     string `json:"type,omitempty"`
	StorageBackend           string `json:"storageBackend,omitempty"`
	MemoryAllocationInMb     int64  `json:"memoryAllocationInMb,omitempty"`
	BucketConflictResolution string `json:"bucketConflictResolution,omitempty"`
	DurabilityLevel          string `json:"durabilityLevel,omitempty"`
	Replicas                 int64  `json:"replicas,omitempty"`
	Flush                    bool   `json:"flush,omitempty"`
	TimeToLiveInSeconds      int64  `json:"timeToLiveInSeconds,omitempty"`
}

type CreateBucketResponse struct {
	Id string `json:"id"`
}

type GetBucketResponse struct {
	// Audit contains all audit-related fields.
	Audit                    CouchbaseAuditData `json:"audit"`
	Id                       string             `json:"id"`
	Name                     string             `json:"name"`
	Type                     string             `json:"type"`
	StorageBackend           string             `json:"storageBackend"`
	MemoryAllocationInMb     int64              `json:"memoryAllocationInMb"`
	BucketConflictResolution string             `json:"bucketConflictResolution"`
	DurabilityLevel          string             `json:"durabilityLevel"`
	Replicas                 int64              `json:"replicas"`
	Flush                    bool               `json:"flush,omitempty"`
	TimeToLiveInSeconds      int64              `json:"timeToLiveInSeconds"`
	EvictionPolicy           string             `json:"evictionPolicy"`
	Stats                    *Stats             `json:"stats"`
}

type Stats struct {
	ItemCount       int64 `json:"itemCount"`
	OpsPerSecond    int64 `json:"opsPerSecond"`
	DiskUsedInMiB   int64 `json:"diskUsedInMiB"`
	MemoryUsedInMiB int64 `json:"memoryUsedInMiB"`
}

type PutBucketRequest struct {
	MemoryAllocationInMb int64  `json:"memoryAllocationInMb"`
	DurabilityLevel      string `json:"durabilityLevel"`
	Replicas             int64  `json:"replicas"`
	Flush                bool   `json:"flush"`
	TimeToLiveInSeconds  int64  `json:"timeToLiveInSeconds"`
}

// GetBucketsResponse defines the model for a GetBucketsResponse.
type GetBucketsResponse struct {
	Data []GetBucketResponse `json:"data"`
}
