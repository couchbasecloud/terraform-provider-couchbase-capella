package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServiceAuditLoggingStateData is the datasource model for the audit logging state of an App Service.
type AppServiceAuditLoggingStateData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AuditEnabled is true when audit logging is enabled on the App Service.
	AuditEnabled types.Bool `tfsdk:"audit_enabled"`
}
