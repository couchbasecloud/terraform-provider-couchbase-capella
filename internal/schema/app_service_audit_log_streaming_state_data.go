package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServiceAuditLogStreamingStateData is the datasource model for the audit log streaming state of an App Service.
type AppServiceAuditLogStreamingStateData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// StreamingEnabled is true when audit log streaming is enabled on the App Service.
	StreamingEnabled types.Bool `tfsdk:"streaming_enabled"`

	// LogStreamingState is the current status of audit log streaming.
	LogStreamingState types.String `tfsdk:"log_streaming_state"`

	// OutputType is the type of output for audit log streaming.
	OutputType types.String `tfsdk:"output_type"`

	// DisabledAppEndpoints is the set of App Endpoints excluded from audit log streaming.
	DisabledAppEndpoints types.Set `tfsdk:"disabled_app_endpoints"`
}
