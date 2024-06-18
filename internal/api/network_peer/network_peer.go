package network_peer

import (
	"github.com/google/uuid"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// CreateNetworkPeeringRequest Creates a network peering record for Capella. Capella does not support peering of networks between different cloud providers. For example, you cannot peer GCP VPC that hosts Capella cluster with an AWS VPC hosting an application.
//
// - Create configures a Couchbase Capella private networking with the cloud provider. Setting up a private network enables your application to interact with Couchbase Capella over a private connection by co-locating them through network peering.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateNetworkPeeringRequest struct {
	// Name is the name of the peering relationship. -  The name of the peering relationship must be at least 2 characters long. -  The name can not exceed 128 characters.
	Name string `json:"name"`

	// ProviderConfig The config data for a peering relationship for a cluster on AWS, GCP.
	//ProviderConfig json.RawMessage `json:"providerConfig"`

	// AWSConfig AWS config data required to establish a VPC peering relationship. Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
	AWSConfig AWSConfig `json:"awsConfig"`

	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
	GCPConfig GCPConfig `json:"gcpConfig"`

	//// ProviderType Type of the cloud provider for which the peering connection is created. Which are- 1. aws 2. gcp
	ProviderType string `json:"providerType"`
}

// CreateNetworkPeeringResponse is the response received from the Capella V4 Public API when asked to create a new network peering connection.
type CreateNetworkPeeringResponse struct {
	// Id is the ID is the unique UUID generated when a VPC record is created.
	Id uuid.UUID `json:"id"`
}

// GetNetworkPeeringRecordResponse Fetches the details of the network peering meta data based on the peerID provided.
//
//	In order to access this endpoint, the provided API key must have at least one of the following roles:
//	 - Organization Owner
//	 - Project Owner
//	 - Project Manager
//
//	To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type GetNetworkPeeringRecordResponse struct {
	Audit api.CouchbaseAuditData `json:"audit"`

	// Commands contains the list of commands that the user must execute in order to complete the association of the network.
	Commands []string `json:"commands"`

	// Id The ID is the unique UUID generated when a VPC record is created.
	Id uuid.UUID `json:"id"`

	// Name is the name of the peering relationship.
	Name string `json:"name"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig ProviderConfig `json:"providerConfig"`

	AWSConfig *AWSConfig `json:"awsConfig"`

	GCPConfig *GCPConfig `json:"gcpConfig"`

	Status PeeringStatus `json:"status"`
}

// AWSConfig is the AWS config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
type AWSConfig struct {
	// AccountId The numeric AWS Account ID or Owner ID.
	AccountId string `json:"accountId"`

	// Cidr The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr string `json:"cidr"`

	// Region The AWS region where your VPC is deployed.
	Region string `json:"region"`

	// VpcId The alphanumeric VPC ID which starts with \"vpc-\". This is also known as the networkId.
	VpcId string `json:"vpcId"`

	// ProviderId The ID of the VPC peer on GCP.
	ProviderId string `json:"providerId"`
}

// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
type GCPConfig struct {
	// Cidr The GCP VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr string `json:"cidr"`

	// NetworkName The name of the network that you want to peer with.
	NetworkName string `json:"networkName"`

	// ProjectId The unique identifier for your GCP project.
	ProjectId string `json:"projectId"`

	// ServiceAccount is the ServiceAccount created or assigned on the external VPC project. GCP Service Account with below permissions
	// - DNS Admin
	// - Compute.NetworkAdmin
	// It should be in the form of an email that is shown under `gcloud iam service-accounts list` command.
	// [Reference](https://cloud.google.com/iam/docs/creating-managing-service-accounts#creating)
	ServiceAccount string `json:"serviceAccount"`

	// ProviderId The ID of the VPC peer on GCP.
	//ProviderId string `json:"providerId"`
}

//
//// AsAWS returns the union data inside the GetNetworkPeeringRecordResponse_ProviderConfig as a AWS
//func (t *CreateNetworkPeeringRequest) AsAWS() (AWSConfig, error) {
//	var body AWSConfig
//	err := json.Unmarshal(t.ProviderConfig, &body)
//	return body, err
//}
//
//// FromAWS overwrites any union data inside the GetNetworkPeeringRecordResponse_ProviderConfig as the provided AWS
//func (t *CreateNetworkPeeringRequest) FromAWS(v AWSConfig) error {
//	b, err := json.Marshal(v)
//	t.ProviderConfig = b
//	return err
//}
//
//// AsGCP returns the union data inside the GetNetworkPeeringRecordResponse_ProviderConfig as a GCP
//func (t *CreateNetworkPeeringRequest) AsGCP() (GCPConfig, error) {
//	var body GCPConfig
//	err := json.Unmarshal(t.ProviderConfig, &body)
//	return body, err
//}
//
//// FromGCP overwrites any union data inside the GetNetworkPeeringRecordResponse_ProviderConfig as the provided GCP
//func (t *CreateNetworkPeeringRequest) FromGCP(v GCPConfig) error {
//	b, err := json.Marshal(v)
//	t.ProviderConfig = b
//	return err
//}
