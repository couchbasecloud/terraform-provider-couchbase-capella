package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointAuditLoggingBuilder = capellaschema.NewSchemaBuilder("appEndpointAuditLogging", "GetAppEndpointAuditLogResponse")

// AppEndpointAuditLoggingSchema returns the schema for the App Endpoint audit logging resource.
func AppEndpointAuditLoggingSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointAuditLoggingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appEndpointAuditLoggingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointAuditLoggingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointAuditLoggingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointAuditLoggingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))

	capellaschema.AddAttr(attrs, "audit_enabled", appEndpointAuditLoggingBuilder, boolAttribute(computed, optional, useStateForUnknown))
	capellaschema.AddAttr(attrs, "enabled_event_ids", appEndpointAuditLoggingBuilder, &schema.SetAttribute{
		Computed:      true,
		Optional:      true,
		ElementType:   types.Int64Type,
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	})
	capellaschema.AddAttr(attrs, "disabled_users", appEndpointAuditLoggingBuilder, disabledUserRoleSetAttribute())
	capellaschema.AddAttr(attrs, "disabled_roles", appEndpointAuditLoggingBuilder, disabledUserRoleSetAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the audit logging configuration of an App Endpoint. Destroying this resource removes it from the Terraform state only. The audit logging configuration in Capella stays in its last state.",
		Attributes:          attrs,
	}
}

func disabledUserRoleSetAttribute() *schema.SetNestedAttribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "domain", appEndpointAuditLoggingBuilder, stringAttribute([]string{required}, stringvalidator.LengthAtLeast(1)), "DisabledUserRole")
	capellaschema.AddAttr(attrs, "name", appEndpointAuditLoggingBuilder, stringAttribute([]string{required}, stringvalidator.LengthAtLeast(1)), "DisabledUserRole")

	return &schema.SetNestedAttribute{
		Computed:      true,
		Optional:      true,
		NestedObject:  schema.NestedAttributeObject{Attributes: attrs},
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}
