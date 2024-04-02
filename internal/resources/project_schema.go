package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ProjectSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"name":            stringAttribute([]string{required}),
			"description":     stringAttribute([]string{optional, computed}),
			"if_match":        stringAttribute([]string{optional}),
			"etag":            stringAttribute([]string{computed}),
			"audit":           computedAuditAttribute(),
		},
	}
}
