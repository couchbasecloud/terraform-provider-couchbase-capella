package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// rest api schema doesn't specify about required field clearly
// is there any other type is supported aparted from weekly?
func BackupScheduleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"bucket_id":       stringAttribute(required, requiresReplace),
			"type":            stringAttribute(required),
			"weekly_schedule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"day_of_week":              stringAttribute(required),
					"start_at":                 int64Attribute(required),
					"incremental_every":        int64Attribute(required),
					"retention_time":           stringAttribute(required),
					"cost_optimized_retention": boolAttribute(required),
				},
			},
		},
	}
}
