package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// NetworkPeer defines the response as received from V4 Capella Public API when asked to create a new network peer.
type NetworkPeer struct {
	Audit types.Object `tfsdk:"audit"`

	// Commands Commands contains the list of commands that the user must execute in order to complete the association of the network.
	Commands []types.String `tfsdk:"commands"`

	// Id is the ID is the unique UUID generated when a VPC record is created.
	Id types.String `tfsdk:"id"`

	// Name is the Name of the peering relationship.
	Name types.String `tfsdk:"name"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig ProviderConfig `tfsdk:"providerConfig"`
	Status         PeeringStatus  `tfsdk:"status"`
}

// PeeringStatus communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
type PeeringStatus struct {
	Reasoning types.String `tfsdk:"reasoning"`
	State     types.String `tfsdk:"state"`
}

// ProviderConfig provides details about the configuration and the ID of the VPC peer on AWS, GCP.
type ProviderConfig struct {
	// AWSConfig AWS config data required to establish a VPC peering relationship. Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
	AWSConfig *AWSConfigData `tfsdk:"AWSConfig"`

	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
	GCPConfig *GCPConfigData `tfsdk:"GCPConfig"`

	// ProviderId The ID of the VPC peer on AWS or GCP.
	ProviderId *types.String `tfsdk:"providerId"`
}

// AWSConfigData AWS config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
type AWSConfigData struct {
	// AccountId The numeric AWS Account ID or Owner ID.
	AccountId types.String `tfsdk:"accountId"`

	// Cidr The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// Region The AWS region where your VPC is deployed.
	Region types.String `tfsdk:"region"`

	// VpcId The alphanumeric VPC ID which starts with \"vpc-\". This is also known as the networkId.
	VpcId types.String `tfsdk:"vpcId"`
}

// GCPConfigData GCP config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
type GCPConfigData struct {
	// Cidr The GCP VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// NetworkName The name of the network that you want to peer with.
	NetworkName types.String `tfsdk:"networkName"`

	// ProjectId The unique identifier for your GCP project.
	ProjectId types.String `tfsdk:"projectId"`

	// ServiceAccount is the service account created or assigned on the external VPC project. GCP Service Account with below permissions
	// - DNS Admin
	// - Compute.NetworkAdmin
	// It should be in the form of an email that is shown under `gcloud iam service-accounts list` command.
	// [Reference](https://cloud.google.com/iam/docs/creating-managing-service-accounts#creating)
	ServiceAccount types.String `tfsdk:"serviceAccount"`
}
