package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BackupScheduleSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages backup schedule resource associated with a bucket for a Couchbase Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"bucket_id":       WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the bucket."),
			"type":            WithDescription(stringAttribute([]string{required, requiresReplace}), "Type of the backup schedule."),
			"weekly_schedule": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Schedule a full backup once a week with regular incrementals.",
				Attributes: map[string]schema.Attribute{
					"day_of_week":              WithDescription(stringAttribute([]string{required}), "Day of the week for the backup. Values can be \"sunday\" \"monday\" \"tuesday\" \"wednesday\" \"thursday\" \"friday\" \"saturday\""),
					"start_at":                 WithDescription(stringAttribute([]string{required}), "Start at hour (in 24-Hour format). Integer value between 0 and 23."),
					"incremental_every":        WithDescription(stringAttribute([]string{required}), "Interval in hours for incremental backup. Integer value between 1 and 24."),
					"retention_time":           WithDescription(stringAttribute([]string{required}), "Retention time in days, ex: 30days, 1year, 5years"),
					"cost_optimized_retention": WithDescription(stringAttribute([]string{required}), "Optimize backup retention to reduce total cost of ownership (TCO). This gives the option to keep all but the last backup cycle of the month for thirty days; the last cycle will be kept for the defined retention period."),
				},
			},
		},
	}
}
