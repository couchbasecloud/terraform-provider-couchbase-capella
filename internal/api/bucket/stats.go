package bucket

// Stats are the bucket related statistics that are sent by the Capella V4 Public API for any existing bucket.
type Stats struct {
	ItemCount       int `json:"itemCount"`
	OpsPerSecond    int `json:"opsPerSecond"`
	DiskUsedInMib   int `json:"diskUsedInMib"`
	MemoryUsedInMib int `json:"memoryUsedInMib"`
}
