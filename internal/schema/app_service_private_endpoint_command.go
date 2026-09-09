package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServiceAWSCommandRequest represents the AWS cli to create a private endpoint
// connected to an App Service's private endpoint service.
type AppServiceAWSCommandRequest struct {
	// ClusterId is the ID of the cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// VpcID The ID of your virtual network.
	VpcID types.String `tfsdk:"vpc_id"`

	// SubnetIDs is a list of subnet ids.
	SubnetIDs []types.String `tfsdk:"subnet_ids"`

	// Command is the AWS command.
	Command types.String `tfsdk:"command"`

	// ServiceName is the AWS VPC endpoint service name, best-effort parsed from Command's
	// `--service-name` flag. It is null when the command's phrasing does not match the
	// expected pattern.
	ServiceName types.String `tfsdk:"service_name"`
}

// AppServiceAzureCommandRequest represents the Azure script to create a private
// endpoint connected to an App Service's private endpoint service.
type AppServiceAzureCommandRequest struct {
	// ClusterId is the ID of the cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// The name of your resource group.
	ResourceGroupName types.String `tfsdk:"resource_group_name"`

	// The virtual network and subnet name.
	VirtualNetwork types.String `tfsdk:"virtual_network"`

	// Command is the Azure script.
	Command types.String `tfsdk:"command"`
}

// AppServiceGCPCommandRequest represents the GCP script to create a private
// endpoint connected to an App Service's private endpoint service.
type AppServiceGCPCommandRequest struct {
	// ClusterId is the ID of the cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// VpcNetworkID The ID of your virtual network.
	VpcNetworkID types.String `tfsdk:"vpc_network_id"`

	// SubnetIDs is a list of subnet ids.
	SubnetIDs []types.String `tfsdk:"subnet_ids"`

	// Command is the GCP command.
	Command types.String `tfsdk:"command"`
}
