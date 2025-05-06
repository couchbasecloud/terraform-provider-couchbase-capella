package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BackupScheduleSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages backup schedule resource associated with a bucket for a Couchbase Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"bucket_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name.",
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Type of the backup schedule.",
			},
			"weekly_schedule": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Schedule a full backup once a week with regular incrementals.",
				Attributes: map[string]schema.Attribute{
					"day_of_week": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Day of the week for the backup. Values can be \"sunday\" \"monday\" \"tuesday\" \"wednesday\" \"thursday\" \"friday\" \"saturday\"",
					},
					"start_at": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Start at hour (in 24-Hour format). Integer value between 0 and 23.",
					},
					"incremental_every": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Interval in hours for incremental backup. Integer value between 1 and 24.",
					},
					"retention_time": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Retention time in days, ex: 30days, 1year, 5years",
					},
					"cost_optimized_retention": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Optimize backup retention to reduce total cost of ownership (TCO). This gives the option to keep all but the last backup cycle of the month for thirty days; the last cycle will be kept for the defined retention period.",
					},
				},
			},
		},
	}
}
