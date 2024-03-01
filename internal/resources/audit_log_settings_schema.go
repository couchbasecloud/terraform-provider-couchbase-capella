package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"auditenabled":    boolAttribute(computed, optional, requiresReplace),
			"enabledeventids": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"disabledusers": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": stringAttribute(computed, optional, requiresReplace),
						"name":   stringAttribute(computed, optional, requiresReplace),
					},
				},
			},
		},
	}
}
