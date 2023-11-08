package appservice

// AppServiceCompute depicts the couchbase compute, following are the supported compute combinations for CPU
// and RAM for different cloud providers.
// To learn more, see:
// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
type AppServiceCompute struct {
	// Cpu depicts cpu units (cores).
	Cpu int64 `json:"cpu"`

	// Ram depicts ram units (GB).
	Ram int64 `json:"ram"`
}
