package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BackupScheduleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"bucket_id":       stringAttribute([]string{required, requiresReplace}),
			"type":            stringAttribute([]string{required, requiresReplace}),
			"weekly_schedule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"day_of_week":              stringAttribute([]string{required}),
					"start_at":                 int64Attribute(required),
					"incremental_every":        int64Attribute(required),
					"retention_time":           stringAttribute([]string{required}),
					"cost_optimized_retention": boolAttribute(required),
				},
			},
		},
	}
}
