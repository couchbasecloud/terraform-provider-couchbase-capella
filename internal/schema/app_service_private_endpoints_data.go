package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServicePrivateEndpointsData is the datasource model for the private endpoints of an App Service.
type AppServicePrivateEndpointsData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Status optionally filters the endpoints by status. When null, all endpoints are returned.
	Status types.Set `tfsdk:"status"`

	// Data is the list of private endpoints.
	Data []PrivateEndpointData `tfsdk:"data"`
}
