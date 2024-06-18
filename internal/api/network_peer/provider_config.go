package network_peer

// ProviderConfig provides details about the configuration and the ID of the VPC peer on AWS, GCP.
type ProviderConfig struct {
	// AWSConfig AWS config data required to establish a VPC peering relationship. Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
	AWSConfig *AWSConfig `json:"awsConfig"`

	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
	GCPConfig *GCPConfig `json:"GCPConfig"`

	//ProviderId The ID of the VPC peer on AWS or GCP.
	ProviderId string `json:"providerId"`
}
