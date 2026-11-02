package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppEndpointAuditLogEventIDs is the datasource model for the audit log events of an App Endpoint.
type AppEndpointAuditLogEventIDs struct {
	OrganizationId  types.String                 `tfsdk:"organization_id"`
	ProjectId       types.String                 `tfsdk:"project_id"`
	ClusterId       types.String                 `tfsdk:"cluster_id"`
	AppServiceId    types.String                 `tfsdk:"app_service_id"`
	AppEndpointName types.String                 `tfsdk:"app_endpoint_name"`
	Data            []AppEndpointAuditLogEventID `tfsdk:"data"`
}

// AppEndpointAuditLogEventID is an audit event that an App Endpoint can log.
type AppEndpointAuditLogEventID struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Filterable  types.Bool   `tfsdk:"filterable"`
}
