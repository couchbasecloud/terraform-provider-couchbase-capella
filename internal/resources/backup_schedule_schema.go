package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BackupScheduleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required),
			"project_id":      stringAttribute(required),
			"cluster_id":      stringAttribute(required),
			"bucket_id":       stringAttribute(required),
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
