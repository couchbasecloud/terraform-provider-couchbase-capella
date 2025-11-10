package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogSettingsBuilder = capellaschema.NewSchemaBuilder("auditLogSettings")

// AuditLogSettingsSchema returns the schema for the AuditLogSettings data source.
func AuditLogSettingsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", auditLogSettingsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", auditLogSettingsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", auditLogSettingsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "audit_enabled", auditLogSettingsBuilder, computedBool())
	capellaschema.AddAttr(attrs, "enabled_event_ids", auditLogSettingsBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.Int64Type,
	})

	// Build disabled_users attributes
	disabledUsersAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(disabledUsersAttrs, "domain", auditLogSettingsBuilder, computedString())
	capellaschema.AddAttr(disabledUsersAttrs, "name", auditLogSettingsBuilder, computedString())

	capellaschema.AddAttr(attrs, "disabled_users", auditLogSettingsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: disabledUsersAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve audit log configuration settings for an operational cluster. These settings control which events are logged and which users are excluded.",
		Attributes:          attrs,
	}
}
