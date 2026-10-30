package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointAuditLoggingConfigBuilder = capellaschema.NewSchemaBuilder("appEndpointAuditLoggingConfig", "GetAppEndpointAuditLogResponse")

// AppEndpointAuditLoggingConfigSchema returns the schema for the app_endpoint_audit_logging_config data source.
func AppEndpointAuditLoggingConfigSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointAuditLoggingConfigBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appEndpointAuditLoggingConfigBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointAuditLoggingConfigBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointAuditLoggingConfigBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointAuditLoggingConfigBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "audit_enabled", appEndpointAuditLoggingConfigBuilder, computedBool())
	capellaschema.AddAttr(attrs, "enabled_event_ids", appEndpointAuditLoggingConfigBuilder, &schema.SetAttribute{
		Computed:            true,
		ElementType:         types.Int64Type,
		MarkdownDescription: "The IDs of the audit events that are enabled.",
	})
	capellaschema.AddAttr(attrs, "disabled_users", appEndpointAuditLoggingConfigBuilder, computedAppEndpointAuditDisabledRoles("The users that are excluded from audit logging."))
	capellaschema.AddAttr(attrs, "disabled_roles", appEndpointAuditLoggingConfigBuilder, computedAppEndpointAuditDisabledRoles("The roles that are excluded from audit logging."))

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit logging config of an App Endpoint.",
		Attributes:          attrs,
	}
}

func computedAppEndpointAuditDisabledRoles(description string) *schema.SetNestedAttribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "domain", appEndpointAuditLoggingConfigBuilder, computedString(), "DisabledUserRole")
	capellaschema.AddAttr(attrs, "name", appEndpointAuditLoggingConfigBuilder, computedString(), "DisabledUserRole")

	return &schema.SetNestedAttribute{
		Computed:            true,
		MarkdownDescription: description,
		NestedObject: schema.NestedAttributeObject{
			Attributes: attrs,
		},
	}
}
