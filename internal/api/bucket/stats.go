package bucket

type Stats struct {
	ItemCount       int `json:"itemCount"`
	OpsPerSecond    int `json:"opsPerSecond"`
	DiskUsedInMib   int `json:"diskUsedInMib"`
	MemoryUsedInMib int `json:"memoryUsedInMib"`
}
