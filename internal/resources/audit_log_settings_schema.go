package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogSettingsBuilder = capellaschema.NewSchemaBuilder("auditLogSettings", "clusterAuditSettings")

func AuditLogSettingsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", auditLogSettingsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", auditLogSettingsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", auditLogSettingsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "audit_enabled", auditLogSettingsBuilder, boolAttribute(computed, optional))
	capellaschema.AddAttr(attrs, "enabled_event_ids", auditLogSettingsBuilder, &schema.SetAttribute{
		Computed:    true,
		Optional:    true,
		ElementType: types.Int64Type,
		Validators: []validator.Set{
			setvalidator.ValueInt64sAre(int64validator.AtLeast(1)),
		},
	})

	disabledUserAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(disabledUserAttrs, "domain", auditLogSettingsBuilder, stringAttribute([]string{required}, stringvalidator.LengthAtLeast(1)))
	capellaschema.AddAttr(disabledUserAttrs, "name", auditLogSettingsBuilder, stringAttribute([]string{required}, stringvalidator.LengthAtLeast(1)))

	capellaschema.AddAttr(attrs, "disabled_users", auditLogSettingsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: disabledUserAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage audit log configuration settings for an operational cluster. These settings control which audit events are logged and which users are excluded from audit logging.",
		Attributes:          attrs,
	}
}
