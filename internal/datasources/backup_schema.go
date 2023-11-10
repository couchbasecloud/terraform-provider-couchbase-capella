package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func BackupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"bucket_id":       requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              computedStringAttribute,
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"bucket_id":       computedStringAttribute,
						"cycle_id":        computedStringAttribute,
						"date":            computedStringAttribute,
						"restore_before":  computedStringAttribute,
						"status":          computedStringAttribute,
						"method":          computedStringAttribute,
						"bucket_name":     computedStringAttribute,
						"source":          computedStringAttribute,
						"cloud_provider":  computedStringAttribute,
						"backup_stats": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"size_in_mb": computedFloat64Attribute,
								"items":      computedInt64Attribute,
								"mutations":  computedInt64Attribute,
								"tombstones": computedInt64Attribute,
								"gsi":        computedInt64Attribute,
								"fts":        computedInt64Attribute,
								"cbas":       computedInt64Attribute,
								"event":      computedInt64Attribute,
							},
						},
						"elapsed_time_in_seconds": computedInt64Attribute,
						"schedule_info": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"backup_type": computedStringAttribute,
								"backup_time": computedStringAttribute,
								"increment":   computedInt64Attribute,
								"retention":   computedStringAttribute,
							},
						},
					},
				},
			},
		},
	}
}
