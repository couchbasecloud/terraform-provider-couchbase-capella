package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var onOffScheduleBuilder = capellaschema.NewSchemaBuilder("onOffSchedule")

func OnOffScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["timezone"] = WithDescription(stringAttribute([]string{required, requiresReplace}), "Timezone for the schedule. Should be the TZ identifier. For example, 'US/Hawaii', 'Indian/Mauritius'")
	attrs["days"] = schema.ListNestedAttribute{
		Required:            true,
		MarkdownDescription: "List of days the On/Off schedule is active.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"state": WithDescription(stringAttribute([]string{required}), "The cluster state. It can be 'on', 'off', or 'custom'."),
				"day": WithDescription(stringAttribute([]string{required},
					validator.String(stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))),
					"Day of the week for scheduling on/off."),
				"from": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "OnTimeBoundary corresponds to \"from\" and \"to\" time boundaries for when the cluster needs to be in the turned on (healthy) state on a day with \"custom\" scheduling timings.",
					Attributes: map[string]schema.Attribute{
						"hour": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(0),
							MarkdownDescription: "Hour of the time boundary. The valid hour values are from 0 to 23 inclusive.",
						},
						"minute": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(0),
							MarkdownDescription: "Minute of the time boundary. The valid minute values are 0 and 30.",
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
							MarkdownDescription: "Hour of the time boundary. The valid hour values are from 0 to 23 inclusive.",
						},
						"minute": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(0),
							MarkdownDescription: "Minute of the time boundary. The valid minute values are 0 and 30.",
						},
					},
				},
			},
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the On/Off schedule for an operational cluster.",
		Attributes:          attrs,
	}
}
