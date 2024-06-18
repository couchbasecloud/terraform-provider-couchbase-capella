package schema

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	//Commands []types.String `tfsdk:"commands"`
	Commands types.Set `tfsdk:"commands"`

	// Name is the Name of the peering relationship.
	Name types.String `tfsdk:"name"`

	//ProviderType is the type of the cloud provider for which the peering connection is created. Which are-
	//     1. aws
	//     2. gcp
	ProviderType types.String `tfsdk:"provider_type"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig ProviderConfig `tfsdk:"provider_config"`

	//AWSConfig AWSConfig `tfsdk:"aws_config"`
	//
	//GCPConfig GCPConfig `tfsdk:"gcp_config"`
	// Status communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
	Status types.Object `tfsdk:"status"`

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
	//AWSConfig AWSConfig `tfsdk:"aws_config"`

	// GCPConfig GCP config data required to establish a VPC peering relationship. Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
	//GCPConfig GCPConfig `tfsdk:"gcp_config"`

	// ProviderId The ID of the VPC peer on AWS or GCP.
	//ProviderId types.String `tfsdk:"provider_id"`

	// AccountId The numeric AWS Account ID or Owner ID.
	AccountId types.String `tfsdk:"account_id"`

	// Cidr The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// Region The AWS region where your VPC is deployed.
	Region types.String `tfsdk:"region"`

	// VpcId The alphanumeric VPC ID which starts with \"vpc-\". This is also known as the networkId.
	VpcId types.String `tfsdk:"vpc_id"`

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

// AWSConfig AWS config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to AWS VPC Peering - [ref](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations).
type AWSConfig struct {
	// AccountId The numeric AWS Account ID or Owner ID.
	AccountId types.String `tfsdk:"account_id"`

	// Cidr The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.
	Cidr types.String `tfsdk:"cidr"`

	// Region The AWS region where your VPC is deployed.
	Region types.String `tfsdk:"region"`

	// VpcId The alphanumeric VPC ID which starts with \"vpc-\". This is also known as the networkId.
	VpcId types.String `tfsdk:"vpc_id"`

	// ProviderId The ID of the VPC peer on GCP.
	ProviderId types.String `tfsdk:"provider_id"`
}

// GCPConfig GCP config data required to establish a VPC peering relationship.
//
//	Refer to the docs for other limitations to GCP VPC Peering - [ref](https://cloud.google.com/vpc/docs/vpc-peering).
type GCPConfig struct {
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
	ProviderId types.String `tfsdk:"provider_id"`
}

// NetworkPeers defines structure based on the response received from V4 Capella Public API when asked to list network peers.
type NetworkPeers struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the cluster
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the network peer needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data It contains the list of resources.
	Data []NetworkPeerData `tfsdk:"data"`
}

// NetworkPeerData defines attributes for a single network peer when fetched from the V4 Capella Public API.
type NetworkPeerData struct {
	//// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	//OrganizationId types.String `tfsdk:"organization_id"`
	//
	//// ProjectId is the ID of the project to which the Capella cluster belongs.
	//ProjectId types.String `tfsdk:"project_id"`
	//
	//// ClusterId is the ID of the cluster for which the network peer needs to be created.
	//ClusterId types.String `tfsdk:"cluster_id"`

	// Id is the ID is the unique UUID generated when a VPC record is created.
	Id types.String `tfsdk:"id"`

	// Commands Commands contains the list of commands that the user must execute in order to complete the association of the network.
	Commands types.List `tfsdk:"commands"`

	// Name is the Name of the peering relationship.
	Name types.String `tfsdk:"name"`

	////ProviderType is the type of the cloud provider for which the peering connection is created. Which are-
	////     1. aws
	////     2. gcp
	//ProviderType types.String `tfsdk:"provider_type"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig ProviderConfig `tfsdk:"provider_config"`

	//AWSConfig AWSConfig `tfsdk:"aws_config"`
	//
	//GCPConfig GCPConfig `tfsdk:"gcp_config"`

	// Status communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
	Status PeeringStatus `tfsdk:"status"`

	Audit types.Object `tfsdk:"audit"`
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
// func NewNetworkPeer(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse, organizationId, projectId, clusterId string, commands []types.String, auditObject basetypes.ObjectValue) (*NetworkPeer, error) {
func NewNetworkPeer(ctx context.Context, networkPeer *network_peer_api.GetNetworkPeeringRecordResponse, organizationId, projectId, clusterId string, auditObject basetypes.ObjectValue) (*NetworkPeer, error) {
	newNetworkPeer := NetworkPeer{
		Id:             types.StringValue(networkPeer.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Name:           types.StringValue(networkPeer.Name),
		Audit:          auditObject,
		//Commands:       commands,
		//AWSConfig: &AWSConfig{
		//	//ProviderId: types.StringValue(networkPeer.AWSConfig.ProviderId),
		//	AccountId: types.StringValue(networkPeer.AWSConfig.AccountId),
		//	VpcId:     types.StringValue(networkPeer.AWSConfig.VpcId),
		//	Region:    types.StringValue(networkPeer.AWSConfig.Region),
		//	Cidr:      types.StringValue(networkPeer.AWSConfig.Cidr),
		//},
		//GCPConfig: &GCPConfig{
		//	//ProviderId:     types.StringValue(networkPeer.GCPConfig.ProviderId),
		//	NetworkName:    types.StringValue(networkPeer.GCPConfig.NetworkName),
		//	ProjectId:      types.StringValue(networkPeer.GCPConfig.ProjectId),
		//	ServiceAccount: types.StringValue(networkPeer.GCPConfig.ServiceAccount),
		//	Cidr:           types.StringValue(networkPeer.AWSConfig.Cidr),
		//},

		//Status: PeeringStatus{
		//	State:     types.StringValue(*networkPeer.Status.State),
		//	Reasoning: types.StringValue(*networkPeer.Status.Reasoning),
		//},
	}

	//newNetworkPeer.Commands = MorphCommands(networkPeer.Commands)

	newProviderConfig, err := morphToProviderConfig(networkPeer)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingProviderConfig, err)
	}
	newNetworkPeer.ProviderConfig = newProviderConfig

	if networkPeer.Status.State != nil {
		state := *networkPeer.Status.State
		reasoning := *networkPeer.Status.Reasoning
		status := PeeringStatus{
			State:     types.StringValue(state),
			Reasoning: types.StringValue(reasoning),
		}
		statusObject, diags := types.ObjectValueFrom(ctx, status.AttributeTypes(), status)
		if diags.HasError() {
			return nil, fmt.Errorf("error while converting peering status")
		}
		newNetworkPeer.Status = statusObject
	}

	newCommands, err := MorphCommands(networkPeer.Commands)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingCidr, err)
	}

	newNetworkPeer.Commands = newCommands

	//newConfig, err := morphToTerraformConfig(networkPeer)
	//if err != nil {
	//	return nil, fmt.Errorf("%s: %w", errors.ErrConvertingProviderConfig, err)
	//}
	//newNetworkPeer.ProviderConfig = *newConfig

	return &newNetworkPeer, nil
}

func morphToProviderConfig(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse) (ProviderConfig, error) {
	var rawConfig map[string]interface{}

	// Unmarshal the provider config into a map
	err := json.Unmarshal(networkPeer.ProviderConfig, &rawConfig)
	if err != nil {
		return ProviderConfig{}, err
	}

	var newProviderConfig ProviderConfig

	// Check for the existence of keys and unmarshal accordingly for aws and gcp config.
	if _, ok := rawConfig["aws_config"]; ok {
		awsConfig, err := networkPeer.AsAWS()
		if err != nil {
			return ProviderConfig{}, fmt.Errorf("%s: %w", errors.ErrReadingAWSConfig, err)
		}
		newProviderConfig.AccountId = types.StringValue(awsConfig.AccountId)
		newProviderConfig.VpcId = types.StringValue(awsConfig.VpcId)
		newProviderConfig.Cidr = types.StringValue(awsConfig.Cidr)
		newProviderConfig.Region = types.StringValue(awsConfig.Region)
		newProviderConfig.ProviderId = types.StringValue(awsConfig.ProviderId)

	} else if _, ok := rawConfig["gcp_config"]; ok {
		gcpConfig, err := networkPeer.AsGCP()
		if err != nil {
			return ProviderConfig{}, fmt.Errorf("%s: %w", errors.ErrReadingGCPConfig, err)
		}
		newProviderConfig.ProjectId = types.StringValue(gcpConfig.ProjectId)
		newProviderConfig.NetworkName = types.StringValue(gcpConfig.NetworkName)
		newProviderConfig.Cidr = types.StringValue(gcpConfig.Cidr)
		newProviderConfig.ServiceAccount = types.StringValue(gcpConfig.ServiceAccount)
		newProviderConfig.ProviderId = types.StringValue(gcpConfig.ProviderId)
	}

	return newProviderConfig, nil

}

// AttributeTypes returns a mapping of field names to their respective attribute types for the CouchbaseServer struct.
// It is used during the conversion of a types.Object field to a CouchbaseServer type.
func (p PeeringStatus) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"state":     types.StringType,
		"reasoning": types.StringType,
	}
}

// MorphCommands is used to convert nested Commands from
// strings to terraform type.String.
func MorphCommands(commands []string) (basetypes.SetValue, error) {
	//var morphedCommands []types.String
	//for _, command := range commands {
	//	morphedCommands = append(morphedCommands, types.StringValue(command))
	//}
	//return morphedCommands
	var newCommand []attr.Value
	for _, command := range commands {
		newCommand = append(newCommand, types.StringValue(command))
	}

	newCommands, diags := types.SetValue(types.StringType, newCommand)
	if diags.HasError() {
		return types.SetUnknown(types.StringType), fmt.Errorf("error while converting commands")
	}

	return newCommands, nil
}

// NewNetworkPeerData create new network peer data object.
// func NewNetworkPeerData(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse, organizationId, projectId, clusterId string, commands []types.String, auditObject basetypes.ObjectValue) (*NetworkPeer, error) {
func NewNetworkPeerData(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse, organizationId, projectId, clusterId string, auditObject basetypes.ObjectValue) (*NetworkPeerData, error) {
	newNetworkPeerData := NetworkPeerData{
		Id: types.StringValue(networkPeer.Id.String()),
		//OrganizationId: types.StringValue(organizationId),
		//ProjectId:      types.StringValue(projectId),
		//ClusterId:      types.StringValue(clusterId),
		Name: types.StringValue(networkPeer.Name),
		//Commands:       commands,
		Audit: auditObject,
		Status: PeeringStatus{
			State:     types.StringValue(*networkPeer.Status.State),
			Reasoning: types.StringValue(*networkPeer.Status.Reasoning),
		},
	}

	//if networkPeer.AWSConfig != nil {
	//	newNetworkPeerData.AWSConfig = &AWSConfig{
	//		AccountId:  types.StringValue(networkPeer.AWSConfig.AccountId),
	//		VpcId:      types.StringValue(networkPeer.AWSConfig.VpcId),
	//		Region:     types.StringValue(networkPeer.AWSConfig.Region),
	//		Cidr:       types.StringValue(networkPeer.AWSConfig.Cidr),
	//		ProviderId: types.StringValue(networkPeer.AWSConfig.ProviderId),
	//	}
	//}
	//
	//if networkPeer.GCPConfig != nil {
	//	newNetworkPeerData.GCPConfig = &GCPConfig{
	//		NetworkName:    types.StringValue(networkPeer.GCPConfig.NetworkName),
	//		ProjectId:      types.StringValue(networkPeer.GCPConfig.ProjectId),
	//		ServiceAccount: types.StringValue(networkPeer.GCPConfig.ServiceAccount),
	//		Cidr:           types.StringValue(networkPeer.GCPConfig.Cidr),
	//		ProviderId:     types.StringValue(networkPeer.GCPConfig.ProviderId),
	//	}
	//}

	var newCommand []attr.Value
	for _, command := range networkPeer.Commands {
		newCommand = append(newCommand, types.StringValue(command))
	}

	commands, diags := types.ListValue(types.StringType, newCommand)
	if diags.HasError() {
		return &NetworkPeerData{}, fmt.Errorf("error while converting commands")
	}

	newNetworkPeerData.Commands = commands

	//var newCommands []types.String
	//for _, command := range networkPeer.Commands {
	//	newCommands = append(newCommands, types.StringValue(command))
	//}
	//newNetworkPeerData.Commands = newCommands

	newProviderConfig, err := morphToProviderConfig(networkPeer)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingProviderConfig, err)
	}
	newNetworkPeerData.ProviderConfig = newProviderConfig

	return &newNetworkPeerData, nil
}

//func morphToTerraformConfig(networkPeer *network_peer_api.GetNetworkPeeringRecordResponse) (*ProviderConfig, error) {
//	var newConfig ProviderConfig
//	newConfig = ProviderConfig{
//		ProviderId: types.StringValue(networkPeer.ProviderConfig.ProviderId),
//	}
//
//	if networkPeer.ProviderConfig.AWSConfig != nil {
//		newConfig.AWSConfig = &AWSConfig{
//			AccountId: types.StringValue(networkPeer.ProviderConfig.AWSConfig.AccountId),
//			VpcId:     types.StringValue(networkPeer.ProviderConfig.AWSConfig.VpcId),
//			Region:    types.StringValue(networkPeer.ProviderConfig.AWSConfig.Region),
//			Cidr:      types.StringValue(networkPeer.ProviderConfig.AWSConfig.Cidr),
//		}
//	}
//
//	if networkPeer.ProviderConfig.GCPConfig != nil {
//		newConfig.GCPConfig = &GCPConfig{
//			NetworkName:    types.StringValue(networkPeer.ProviderConfig.GCPConfig.NetworkName),
//			ProjectId:      types.StringValue(networkPeer.ProviderConfig.GCPConfig.ProjectId),
//			Cidr:           types.StringValue(networkPeer.ProviderConfig.GCPConfig.Cidr),
//			ServiceAccount: types.StringValue(networkPeer.ProviderConfig.GCPConfig.ServiceAccount),
//		}
//	}
//	return &newConfig, nil
//}
