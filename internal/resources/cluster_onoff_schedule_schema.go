package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func OnOffScheduleSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The On/Off schedule resource allows you to manage the on/off schedule for a Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timezone": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Timezone for the schedule",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"days": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of days the on/off schedule is active.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Cluster state (on, off, or custom).",
						},
						"day": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"),
							},
							MarkdownDescription: "Day of the week for scheduling on/off.",
						},
						"from": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "OnTimeBoundary corresponds to \"from\" and \"to\" time boundaries for when the cluster needs to be in the turned on (healthy) state on a day with \"custom\" scheduling timings.",
							Attributes: map[string]schema.Attribute{
								"hour": schema.Int64Attribute{
									Optional:            true,
									Computed:            true,
									Default:             int64default.StaticInt64(0),
									MarkdownDescription: "Hour of the time boundary.",
								},
								"minute": schema.Int64Attribute{
									Optional:            true,
									Computed:            true,
									Default:             int64default.StaticInt64(0),
									MarkdownDescription: "Minute of the time boundary.",
								},
							},
						},
						"to": schema.SingleNestedAttribute{
							MarkdownDescription: "OnTimeBoundary corresponds to \"from\" and \"to\" time boundaries for when the cluster needs to be in the turned on (healthy) state on a day with \"custom\" scheduling timings.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"hour": schema.Int64Attribute{
									Optional:            true,
									Computed:            true,
									Default:             int64default.StaticInt64(0),
									MarkdownDescription: "Hour of the time boundary.",
								},
								"minute": schema.Int64Attribute{
									Optional:            true,
									Computed:            true,
									Default:             int64default.StaticInt64(0),
									MarkdownDescription: "Minute of the time boundary.",
								},
							},
						},
					},
				},
			},
		},
	}
}
