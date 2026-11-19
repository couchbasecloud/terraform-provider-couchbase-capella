package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppEndpointAuditLoggingConfigData is the datasource model for the audit logging config of an App Endpoint.
type AppEndpointAuditLoggingConfigData struct {
	OrganizationId  types.String                   `tfsdk:"organization_id"`
	ProjectId       types.String                   `tfsdk:"project_id"`
	ClusterId       types.String                   `tfsdk:"cluster_id"`
	AppServiceId    types.String                   `tfsdk:"app_service_id"`
	AppEndpointName types.String                   `tfsdk:"app_endpoint_name"`
	AuditEnabled    types.Bool                     `tfsdk:"audit_enabled"`
	EnabledEventIds []types.Int64                  `tfsdk:"enabled_event_ids"`
	DisabledUsers   []AppEndpointAuditDisabledRole `tfsdk:"disabled_users"`
	DisabledRoles   []AppEndpointAuditDisabledRole `tfsdk:"disabled_roles"`
}

// AppEndpointAuditDisabledRole is a user or role excluded from App Endpoint audit logging.
type AppEndpointAuditDisabledRole struct {
	Domain types.String `tfsdk:"domain"`
	Name   types.String `tfsdk:"name"`
}
