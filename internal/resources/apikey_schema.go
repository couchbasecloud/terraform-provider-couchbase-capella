package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ApiKeySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute(required, requiresReplace),
			"name":            stringAttribute(required, requiresReplace),
			"description":     stringAttribute(optional, computed, requiresReplace, useStateForUnknown),
			"expiry":          float64Attribute(optional, computed, requiresReplace, useStateForUnknown),
			"allowed_cidrs": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")})),
			},
			"organization_roles": stringListAttribute(required, requiresReplace),
			"resources": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    stringAttribute(required),
						"roles": stringListAttribute(required),
						"type":  stringAttribute(optional, computed),
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"rotate": schema.NumberAttribute{
				Optional: true,
				Computed: true,
			},
			"secret": stringAttribute(optional, computed, sensitive),
			"token":  stringAttribute(computed, sensitive),
			"audit":  computedAuditAttribute(),
		},
	}
}
