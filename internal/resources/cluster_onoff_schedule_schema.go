package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var onOffScheduleBuilder = capellaschema.NewSchemaBuilder("onOffSchedule", "ClusterOnOffSchedule")

// onOffScheduleHourAttribute returns the hour attribute for an on/off schedule
// time boundary. The V4 API accepts hour values from 0 to 23 inclusive.
func onOffScheduleHourAttribute() *schema.Int64Attribute {
	attribute := int64DefaultAttribute(0, optional, computed)
	attribute.Validators = []validator.Int64{int64validator.Between(0, 23)}
	return attribute
}

// onOffScheduleMinuteAttribute returns the minute attribute for an on/off
// schedule time boundary. The V4 API only accepts minute values 0 and 30.
func onOffScheduleMinuteAttribute() *schema.Int64Attribute {
	attribute := int64DefaultAttribute(0, optional, computed)
	attribute.Validators = []validator.Int64{int64validator.OneOf(0, 30)}
	return attribute
}

func OnOffScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", onOffScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", onOffScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", onOffScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "timezone", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	timeBoundaryAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(timeBoundaryAttrs, "hour", onOffScheduleBuilder, onOffScheduleHourAttribute())
	capellaschema.AddAttr(timeBoundaryAttrs, "minute", onOffScheduleBuilder, onOffScheduleMinuteAttribute())

	dayAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dayAttrs, "state", onOffScheduleBuilder, stringAttribute([]string{required},
		validator.String(stringvalidator.OneOf("on", "off", "custom"))))
	capellaschema.AddAttr(dayAttrs, "day", onOffScheduleBuilder, stringAttribute([]string{required},
		validator.String(stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))))

	capellaschema.AddAttr(dayAttrs, "from", onOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: timeBoundaryAttrs,
	})

	toAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(toAttrs, "hour", onOffScheduleBuilder, onOffScheduleHourAttribute())
	capellaschema.AddAttr(toAttrs, "minute", onOffScheduleBuilder, onOffScheduleMinuteAttribute())

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
