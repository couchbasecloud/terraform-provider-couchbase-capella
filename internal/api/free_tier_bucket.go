package api

type CreateFreeTierBucketRequest struct {
	Name                 string `json:"name"`
	MemoryAllocationInMb *int64 `json:"memoryAllocationInMb"`
}

type UpdateFreeTierBucketRequest struct {
	MemoryAllocationInMb int64 `json:"memoryAllocationInMb"`
}
