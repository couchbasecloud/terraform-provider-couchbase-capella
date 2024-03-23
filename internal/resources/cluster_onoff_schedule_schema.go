package resources

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func OnOffScheduleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"timezone":        stringAttribute(required, requiresReplace),
			"days": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": stringAttribute(required),
						"day":   stringAttribute(required),
						"from": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"hour":   int64DefaultAttribute(0, optional, computed),
								"minute": int64DefaultAttribute(0, optional, computed),
							},
						},
						"to": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
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
