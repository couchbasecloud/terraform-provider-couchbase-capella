package bucket

// Stats are the bucket related statistics that are sent by the Capella V4 Public API for any existing bucket.
type Stats struct {
	ItemCount       int64 `json:"itemCount"`
	OpsPerSecond    int64 `json:"opsPerSecond"`
	DiskUsedInMib   int64 `json:"diskUsedInMib"`
	MemoryUsedInMib int64 `json:"memoryUsedInMib"`
}
