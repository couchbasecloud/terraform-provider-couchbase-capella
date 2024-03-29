package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OnOffScheduleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"timezone":        stringAttribute([]string{required, requiresReplace}),
			"days": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": stringAttribute([]string{required}),
						"day":   stringAttribute([]string{required}, stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday")),
						"from": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"hour":   int64DefaultAttribute(0, optional, computed),
								"minute": int64DefaultAttribute(0, optional, computed),
							},
						},
						"to": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"hour":   int64DefaultAttribute(0, optional, computed),
								"minute": int64DefaultAttribute(0, optional, computed),
							},
						},
					},
				},
			},
		},
	}
}
