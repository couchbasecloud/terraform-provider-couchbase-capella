package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type AWSCommandRequest struct {
	// ClusterId is the ID of the cluster for which the collection needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// VpcID The ID of your virtual network
	VpcID types.String `tfsdk:"vpcID"`

	SubnetIDs []types.String `tfsdk:"subnetIDs"`
}

type AzureCommandRequest struct {
	// The name of your resource group
	ResourceGroupName types.String `tfsdk:"resourceGroupName"`

	// The virtual network and subnet name
	VirtualNetwork types.String `tfsdk:"virtualNetwork"`
}
