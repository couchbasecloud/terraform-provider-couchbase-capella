package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppEndpointAuditLogging defines the Terraform state for the App Endpoint audit logging resource.
type AppEndpointAuditLogging struct {
	OrganizationId  types.String `tfsdk:"organization_id"`
	ProjectId       types.String `tfsdk:"project_id"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`
	AuditEnabled    types.Bool   `tfsdk:"audit_enabled"`
	EnabledEventIds types.Set    `tfsdk:"enabled_event_ids"`
	DisabledUsers   types.Set    `tfsdk:"disabled_users"`
	DisabledRoles   types.Set    `tfsdk:"disabled_roles"`
}

// AppEndpointAuditLogDisabledUserRole is an element of disabled_users and disabled_roles.
type AppEndpointAuditLogDisabledUserRole struct {
	Domain types.String `tfsdk:"domain"`
	Name   types.String `tfsdk:"name"`
}

// AppEndpointAuditLogDisabledUserRoleType is the object type of disabled_users and disabled_roles elements.
var AppEndpointAuditLogDisabledUserRoleType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"domain": types.StringType,
	"name":   types.StringType,
}}

// Validate validates the resource state and returns the IDs.
// It also parses the import string when the resource is imported.
func (a *AppEndpointAuditLogging) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  a.OrganizationId,
		ProjectId:       a.ProjectId,
		ClusterId:       a.ClusterId,
		AppServiceId:    a.AppServiceId,
		AppEndpointName: a.AppEndpointName,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
