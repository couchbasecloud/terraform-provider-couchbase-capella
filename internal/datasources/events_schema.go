package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var eventsBuilder = capellaschema.NewSchemaBuilder("events")

// EventsSchema returns the schema for the Events data source.
func EventsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", eventsBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "project_ids", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "cluster_ids", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "user_ids", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "severity_levels", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "tags", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "from", eventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "to", eventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "page", eventsBuilder, &schema.Int64Attribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "per_page", eventsBuilder, &schema.Int64Attribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "sort_by", eventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "sort_direction", eventsBuilder, &schema.StringAttribute{
		Optional: true,
	})

	attrs["data"] = computedEventAttributes
	attrs["cursor"] = computedCursorAttribute

	return schema.Schema{
		MarkdownDescription: " Data source to retrieve all events in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes:          attrs,
	}
}
