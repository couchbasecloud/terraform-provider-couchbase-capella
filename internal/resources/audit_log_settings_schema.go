package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogSettingsBuilder = capellaschema.NewSchemaBuilder("auditLogSettings")

func AuditLogSettingsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", auditLogSettingsBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "project_id", auditLogSettingsBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "cluster_id", auditLogSettingsBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "audit_enabled", auditLogSettingsBuilder, boolAttribute(computed, optional))
	capellaschema.AddAttr(attrs, "enabled_event_ids", auditLogSettingsBuilder, &schema.SetAttribute{
		Computed:    true,
		Optional:    true,
		ElementType: types.Int64Type,
	})

	disabledUserAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(disabledUserAttrs, "domain", auditLogSettingsBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(disabledUserAttrs, "name", auditLogSettingsBuilder, stringAttribute([]string{required}))

	attrs["disabled_users"] = &schema.SetNestedAttribute{
		Computed: true,
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: disabledUserAttrs,
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage audit log configuration settings for an operational cluster. These settings control which audit events are logged and which users are excluded from audit logging.",
		Attributes:          attrs,
	}
}
