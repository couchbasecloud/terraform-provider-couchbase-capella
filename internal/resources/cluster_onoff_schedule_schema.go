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
	capellaschema.AddAttr(attrs, "timezone", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	timeBoundaryAttrs := make(map[string]schema.Attribute)
	timeBoundaryAttrs["hour"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(0),
	}
	timeBoundaryAttrs["minute"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(0),
	}

	dayAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dayAttrs, "state", onOffScheduleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(dayAttrs, "day", onOffScheduleBuilder, stringAttribute([]string{required},
		validator.String(stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))))

	capellaschema.AddAttr(dayAttrs, "from", onOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: timeBoundaryAttrs,
	})

	toAttrs := make(map[string]schema.Attribute)
	toAttrs["hour"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(0),
	}
	toAttrs["minute"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(0),
	}

	capellaschema.AddAttr(dayAttrs, "to", onOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: toAttrs,
	})

	capellaschema.AddAttr(attrs, "days", onOffScheduleBuilder, &schema.ListNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dayAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the On/Off schedule for an operational cluster.",
		Attributes:          attrs,
	}
}
