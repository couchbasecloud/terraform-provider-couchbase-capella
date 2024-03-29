package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type AppServiceOnOffOnDemand struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the appServiceId of the capella tenant.
	AppServiceId types.String `tfsdk:"appService_id"`

	// State is the state to which the app service needs to be turned to i.e. on or off.
	State types.String `tfsdk:"state"`
}
