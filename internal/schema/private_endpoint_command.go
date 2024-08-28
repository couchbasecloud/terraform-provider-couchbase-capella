package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AWSCommandRequest represents the AWS cli to create a private endpoint.
type AWSCommandRequest struct {
	// ClusterId is the ID of the cluster associated with the private endpoint.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// VpcID The ID of your virtual network.
	VpcID types.String `tfsdk:"vpc_id"`

	// SubnetIDs is a list of subnet ids.
	SubnetIDs []types.String `tfsdk:"subnet_ids"`

	// Command is the AWS command.
	Command types.String `tfsdk:"command"`
}

// AzureCommandRequest represents the Azure script to create a private endpoint.
type AzureCommandRequest struct {
	// ClusterId is the ID of the cluster associated with the private endpoint.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// The name of your resource group.
	ResourceGroupName types.String `tfsdk:"resource_group_name"`

	// The virtual network and subnet name.
	VirtualNetwork types.String `tfsdk:"virtual_network"`

	// Command is the Azure script.
	Command types.String `tfsdk:"command"`
}
