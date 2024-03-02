package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AuditLogSettingsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute(required),
			"project_id":      stringAttribute(required),
			"cluster_id":      stringAttribute(required),
			"auditenabled":    boolAttribute(computed, optional),
			"enabledeventids": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"disabledusers": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": stringAttribute(computed, optional),
						"name":   stringAttribute(computed, optional),
					},
				},
			},
		},
	}
}
