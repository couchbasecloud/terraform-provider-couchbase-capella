package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var eventsBuilder = capellaschema.NewSchemaBuilder("events")

// EventsSchema returns the schema for the Events data source.
func EventsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", eventsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_ids", eventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "cluster_ids", eventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "user_ids", eventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "severity_levels", eventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "tags", eventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "from", eventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "to", eventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "page", eventsBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "per_page", eventsBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "sort_by", eventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "sort_direction", eventsBuilder, optionalString())

	attrs["data"] = computedEventAttributes
	attrs["cursor"] = computedCursorAttribute

	return schema.Schema{
		MarkdownDescription: " Data source to retrieve all events in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes:          attrs,
	}
}
