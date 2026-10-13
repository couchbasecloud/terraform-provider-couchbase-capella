package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServicePrivateEndpointServiceData is the datasource model for the private endpoint service of an App Service.
type AppServicePrivateEndpointServiceData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Aws holds the input for the AWS CLI command. When nil, no command is requested.
	Aws *AppServicePrivateEndpointAWSCommandRequest `tfsdk:"aws"`

	// State is the current state of the private endpoint service.
	State types.String `tfsdk:"state"`

	// TargetState is the desired end state of the private endpoint service.
	TargetState types.String `tfsdk:"target_state"`

	// Command is the CLI command to create the private endpoint in the CSP.
	Command types.String `tfsdk:"command"`
}

// AppServicePrivateEndpointAWSCommandRequest is the AWS input for the private endpoint command.
type AppServicePrivateEndpointAWSCommandRequest struct {
	// VpcID is the ID of the AWS VPC.
	VpcID types.String `tfsdk:"vpc_id"`

	// SubnetIDs is the list of AWS subnet IDs.
	SubnetIDs []types.String `tfsdk:"subnet_ids"`
}
