package appservice

// Compute Following are the supported compute combinations for CPU
// and RAM for different cloud providers. To learn more,
// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type Compute struct {
	// Cpu depicts cpu units (cores).
	Cpu int64 `json:"cpu"`

	// Ram depicts ram units (GB).
	Ram int64 `json:"ram"`
}
