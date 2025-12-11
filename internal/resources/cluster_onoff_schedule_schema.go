package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var onOffScheduleBuilder = capellaschema.NewSchemaBuilder("onOffSchedule", "ClusterOnOffSchedule")

func OnOffScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "timezone", onOffScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	timeBoundaryAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(timeBoundaryAttrs, "hour", onOffScheduleBuilder, int64DefaultAttribute(0, optional, computed))
	capellaschema.AddAttr(timeBoundaryAttrs, "minute", onOffScheduleBuilder, int64DefaultAttribute(0, optional, computed))

	dayAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dayAttrs, "state", onOffScheduleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(dayAttrs, "day", onOffScheduleBuilder, stringAttribute([]string{required},
		validator.String(stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))))

	capellaschema.AddAttr(dayAttrs, "from", onOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: timeBoundaryAttrs,
	})

	toAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(toAttrs, "hour", onOffScheduleBuilder, int64DefaultAttribute(0, optional, computed))
	capellaschema.AddAttr(toAttrs, "minute", onOffScheduleBuilder, int64DefaultAttribute(0, optional, computed))

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
