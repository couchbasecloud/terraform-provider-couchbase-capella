package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"name":            stringAttribute([]string{required, requiresReplace}),
			"description":     stringDefaultAttribute("", optional, computed, requiresReplace, useStateForUnknown),
			"expiry":          float64DefaultAttribute(180, optional, computed, requiresReplace, useStateForUnknown),
			"allowed_cidrs": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")})),
			},
			"organization_roles": stringSetAttribute(required, requiresReplace),
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    stringAttribute([]string{required}),
						"roles": stringSetAttribute(required),
						"type":  stringDefaultAttribute("project", optional, computed),
					},
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"rotate": schema.NumberAttribute{
				Optional: true,
				Computed: true,
			},
			"secret": stringAttribute([]string{optional, computed, sensitive}),
			"token":  stringAttribute([]string{computed, sensitive}),
			"audit":  computedAuditAttribute(),
		},
	}
}
