package cluster

// Defines values for CloudProviderType.
const (
	Aws   CloudProviderType = "aws"
	Azure CloudProviderType = "azure"
	Gcp   CloudProviderType = "gcp"
)

// CloudProvider depicts where the cluster will be hosted.
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type CloudProvider struct {
	// Cidr block for Cloud Provider.
	Cidr string `json:"cidr"`

	// Region is cloud provider region, e.g. 'us-west-2'. For information about supported regions, see
	// [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	// [Google Cloud Provider](https://docs.couchbase.com/cloud/reference/gcp.html).
	// [Azure Cloud](https://docs.couchbase.com/cloud/reference/azure.html).
	Region string `json:"region"`

	// Type is cloud provider type, either 'AWS', 'GCP', or 'Azure'.
	Type CloudProviderType `json:"type"`
}

// CloudProviderType is cloud provider type, either 'AWS', 'GCP', or 'Azure'.
type CloudProviderType string
