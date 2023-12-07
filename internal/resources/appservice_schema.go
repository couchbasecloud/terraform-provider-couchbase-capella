package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AppServiceSchema() schema.Schema {
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
			"name":            stringAttribute(required, requiresReplace),
			"description":     stringDefaultAttribute("", optional, computed, requiresReplace),
			"nodes":           int64Attribute(optional, computed),
			"cloud_provider":  stringAttribute(optional, computed),
			"current_state":   stringAttribute(computed),
			"version":         stringAttribute(optional, computed),
			"compute": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"cpu": int64Attribute(required),
					"ram": int64Attribute(required),
				},
			},
			"audit":    computedAuditAttribute(),
			"if_match": stringAttribute(optional),
			"etag":     stringAttribute(computed),
		},
	}
}
