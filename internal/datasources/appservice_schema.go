package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func AppServiceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              computedStringAttribute,
						"organization_id": computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"name":            computedStringAttribute,
						"description":     computedStringAttribute,
						"nodes":           computedInt64Attribute,
						"cloud_provider":  computedStringAttribute,
						"current_state":   computedStringAttribute,
						"compute": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"cpu": computedInt64Attribute,
								"ram": computedInt64Attribute,
							},
						},
						"version": computedStringAttribute,
						"audit":   computedAuditAttribute,
					},
				},
			},
		},
	}
}
