package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AuditLogSettingsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required}),
			"project_id":      stringAttribute([]string{required}),
			"cluster_id":      stringAttribute([]string{required}),
			"audit_enabled":   boolAttribute(computed, optional),
			"enabled_event_ids": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"disabled_users": schema.SetNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": stringAttribute([]string{required}),
						"name":   stringAttribute([]string{required}),
					},
				},
			},
		},
	}
}
