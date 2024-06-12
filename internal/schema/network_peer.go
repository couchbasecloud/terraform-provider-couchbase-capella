package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	network_peer_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/network_peer"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// NetworkPeer defines the response as received from V4 Capella Public API when asked to create a new network peer.
type NetworkPeer struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the network peer needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Id is the ID is the unique UUID generated when a VPC record is created.
	Id types.String `tfsdk:"id"`

	// Commands Commands contains the list of commands that the user must execute in order to complete the association of the network.
	Commands []types.String `tfsdk:"commands"`

	// Name is the Name of the peering relationship.
	Name types.String `tfsdk:"name"`

	//ProviderType is the type of the cloud provider for which the peering connection is created. Which are-
	//     1. aws
	//     2. gcp
	ProviderType types.String `tfsdk:"provider_type"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig ProviderConfig `tfsdk:"provider_config"`

	//AWS   AWS  `tfsdk:"aws"`
	//
	//GCP   GCP  `tfsdk:"gcp"`
	// Status communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
	Status PeeringStatus `tfsdk:"status"`

	Audit types.Object `tfsdk:"audit"`
}

// PeeringStatus communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
type PeeringStatus struct {
	Reasoning types.String `tfsdk:"reasoning"`
	State     types.String `tfsdk:"state"`
}

// ProviderConfig provides details about the configuration and the ID of the VPC peer on AWS, GCP.
type ProviderConfig struct {
	// AWSConfig AWS config data required to establish a VPC peering relationship. Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
	AWSConfig *AWSConfigData `tfsdk:"AWS_config"`

	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
	GCPConfig *GCPConfigData `tfsdk:"GCP_config"`

	// ProviderId The ID of the VPC peer on AWS or GCP.
	ProviderId types.String `tfsdk:"provider_id"`
}

//type AWS struct {
//	// AWSConfig AWS config data required to establish a VPC peering relationship. Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
//	AWSConfig *AWSConfigData `json:"AWSConfig"`
//
//	// ProviderId The ID of the VPC peer on AWS.
//	ProviderId types.String `json:"providerId"`
//}
//
//type GCP struct {
//	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
//	GCPConfig *GCPConfigData `json:"GCPConfig"`
//
//	// ProviderId The ID of the VPC peer on GCP.
//	ProviderId types.String `json:"providerId"`
//}

// AWSConfigData AWS config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
type AWSConfigData struct {
	// AccountId The numeric AWS Account ID or Owner ID.
	AccountId types.String `tfsdk:"account_id"`

	// Cidr The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// Region The AWS region where your VPC is deployed.
	Region types.String `tfsdk:"region"`

	// VpcId The alphanumeric VPC ID which starts with \"vpc-\". This is also known as the networkId.
	VpcId types.String `tfsdk:"vpc_id"`

	// ProviderId The ID of the VPC peer on GCP.
	//ProviderId types.String `json:"providerId"`
}

// GCPConfigData GCP config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
type GCPConfigData struct {
	// Cidr The GCP VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// NetworkName The name of the network that you want to peer with.
	NetworkName types.String `tfsdk:"network_name"`

	// ProjectId The unique identifier for your GCP project.
	ProjectId types.String `tfsdk:"project_id"`

	// ServiceAccount is the service account created or assigned on the external VPC project. GCP Service Account with below permissions
	// - DNS Admin
	// - Compute.NetworkAdmin
	// It should be in the form of an email that is shown under `gcloud iam service-accounts list` command.
	// [Reference](https://cloud.google.com/iam/docs/creating-managing-service-accounts#creating)
	ServiceAccount types.String `tfsdk:"service_account"`

	// ProviderId The ID of the VPC peer on GCP.
	//ProviderId types.String `json:"providerId"`
}

func (n *NetworkPeer) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: n.OrganizationId,
		ProjectId:      n.ProjectId,
		ClusterId:      n.ClusterId,
		Id:             n.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}

// NewNetworkPeer create new network peer object.
func NewNetworkPeer(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse, organizationId, projectId, clusterId string, commands []types.String, auditObject basetypes.ObjectValue) (*NetworkPeer, error) {
	newNetworkPeer := NetworkPeer{
		Id:             types.StringValue(networkPeer.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Name:           types.StringValue(networkPeer.Name),
		Commands:       commands,
		Audit:          auditObject,
		Status: PeeringStatus{
			State:     types.StringValue(*networkPeer.Status.State),
			Reasoning: types.StringValue(*networkPeer.Status.Reasoning),
		},
	}

	newConfig, err := morphToTerraformAWSConfig(networkPeer)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingProviderConfig, err)
	}
	newNetworkPeer.ProviderConfig = *newConfig

	return &newNetworkPeer, nil
}

func morphToTerraformAWSConfig(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse) (*ProviderConfig, error) {
	var newConfig ProviderConfig

	newConfig = ProviderConfig{
		ProviderId: types.StringValue(networkPeer.ProviderConfig.ProviderId),
	}
	if networkPeer.ProviderConfig.AWSConfig != nil {
		newConfig.AWSConfig.AccountId = types.StringValue(networkPeer.ProviderConfig.AWSConfig.AccountId)
		newConfig.AWSConfig.VpcId = types.StringValue(networkPeer.ProviderConfig.AWSConfig.VpcId)
		newConfig.AWSConfig.Region = types.StringValue(networkPeer.ProviderConfig.AWSConfig.Region)
		newConfig.AWSConfig.Cidr = types.StringValue(networkPeer.ProviderConfig.AWSConfig.Cidr)
	} else if networkPeer.ProviderConfig.GCPConfig != nil {
		newConfig.GCPConfig.ProjectId = types.StringValue(networkPeer.ProviderConfig.GCPConfig.ProjectId)
		newConfig.GCPConfig.NetworkName = types.StringValue(networkPeer.ProviderConfig.GCPConfig.NetworkName)
		newConfig.GCPConfig.Cidr = types.StringValue(networkPeer.ProviderConfig.GCPConfig.Cidr)
		newConfig.GCPConfig.ServiceAccount = types.StringValue(networkPeer.ProviderConfig.GCPConfig.ServiceAccount)
	} else {
		return nil, errors.ErrUnsupportedCloudProvider
	}

	return &newConfig, nil
}

// MorphCommands is used to convert nested Commands from
// strings to terraform type.String.
func MorphCommands(commands []string) []basetypes.StringValue {
	var morphedCommands []basetypes.StringValue
	for _, command := range commands {
		morphedCommands = append(morphedCommands, types.StringValue(command))
	}
	return morphedCommands
}
