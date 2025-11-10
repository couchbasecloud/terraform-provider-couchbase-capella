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

	// Build event data attributes
	eventAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(eventAttrs, "alert_key", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "app_service_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "app_service_name", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "cluster_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "cluster_name", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "image_url", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "incident_ids", eventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(eventAttrs, "key", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "kv", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "occurrence_count", eventsBuilder, computedInt64())
	capellaschema.AddAttr(eventAttrs, "project_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "project_name", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "request_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "session_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "severity", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "source", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "summary", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "timestamp", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_email", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_id", eventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_name", eventsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", eventsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: eventAttrs,
		},
	})

	// Build cursor attributes for pagination
	hrefsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(hrefsAttrs, "first", eventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "last", eventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "next", eventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "previous", eventsBuilder, computedString())

	pagesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(pagesAttrs, "last", eventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "next", eventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "page", eventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "per_page", eventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "previous", eventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "total_items", eventsBuilder, computedInt64())

	cursorAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cursorAttrs, "hrefs", eventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: hrefsAttrs,
	})
	capellaschema.AddAttr(cursorAttrs, "pages", eventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: pagesAttrs,
	})

	capellaschema.AddAttr(attrs, "cursor", eventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cursorAttrs,
	})

	return schema.Schema{
		MarkdownDescription: " Data source to retrieve all events in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes:          attrs,
	}
}
