package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterOnOffScheduleBuilder = capellaschema.NewSchemaBuilder("clusterOnOffSchedule")

// ClusterOnOffScheduleSchema returns the schema for the ClusterOnOffSchedule data source.
func ClusterOnOffScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterOnOffScheduleBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", clusterOnOffScheduleBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", clusterOnOffScheduleBuilder, requiredString())
	capellaschema.AddAttr(attrs, "timezone", clusterOnOffScheduleBuilder, computedString())

	// Build time attributes (for 'from' and 'to')
	timeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(timeAttrs, "hour", clusterOnOffScheduleBuilder, computedInt64())
	capellaschema.AddAttr(timeAttrs, "minute", clusterOnOffScheduleBuilder, computedInt64())

	// Build days attributes
	daysAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(daysAttrs, "state", clusterOnOffScheduleBuilder, computedString())
	capellaschema.AddAttr(daysAttrs, "day", clusterOnOffScheduleBuilder, computedString())
	capellaschema.AddAttr(daysAttrs, "from", clusterOnOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: timeAttrs,
	})

	// Need separate timeAttrs map for "to" because we can't reuse the same map
	toTimeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(toTimeAttrs, "hour", clusterOnOffScheduleBuilder, computedInt64())
	capellaschema.AddAttr(toTimeAttrs, "minute", clusterOnOffScheduleBuilder, computedInt64())

	capellaschema.AddAttr(daysAttrs, "to", clusterOnOffScheduleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: toTimeAttrs,
	})

	capellaschema.AddAttr(attrs, "days", clusterOnOffScheduleBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: daysAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The On/Off schedule data source allows you to retrieve the on/off schedule for an operational cluster.",
		Attributes:          attrs,
	}
}
