package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
			"organization_id":    stringAttribute(required),
			"name":               stringAttribute(required),
			"description":        stringAttribute(optional, computed),
			"expiry":             float64Attribute(optional, computed),
			"allowed_cidrs":      stringListAttribute(optional, computed),
			"organization_roles": stringListAttribute(required),
			"resources": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    stringAttribute(required),
						"roles": stringAttribute(required),
						"type":  stringAttribute(optional, computed),
					},
				},
			},
			"rotate": boolAttribute(optional, computed),
			"secret": stringAttribute(optional, computed, sensitive),
			"token":  stringAttribute(computed, sensitive),
			"audit":  computedAuditAttribute(),
		},
	}
}
