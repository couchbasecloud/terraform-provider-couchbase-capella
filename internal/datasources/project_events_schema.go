package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectEventsBuilder = capellaschema.NewSchemaBuilder("projectEvents")

// ProjectEventsSchema returns the schema for the ProjectEvents data source.
func ProjectEventsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", projectEventsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", projectEventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "cluster_ids", projectEventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "user_ids", projectEventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "severity_levels", projectEventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "tags", projectEventsBuilder, optionalStringSet())
	capellaschema.AddAttr(attrs, "from", projectEventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "to", projectEventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "page", projectEventsBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "per_page", projectEventsBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "sort_by", projectEventsBuilder, optionalString())
	capellaschema.AddAttr(attrs, "sort_direction", projectEventsBuilder, optionalString())

	attrs["data"] = computedEventAttributes
	attrs["cursor"] = computedCursorAttribute

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all event information for a project. Events represent a trail of actions that users performs within Capella at project level.",
		Attributes:          attrs,
	}
}
