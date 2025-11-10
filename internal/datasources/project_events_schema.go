package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

	// Build event data attributes
	eventAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(eventAttrs, "alert_key", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "app_service_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "app_service_name", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "cluster_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "cluster_name", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "image_url", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "incident_ids", projectEventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(eventAttrs, "key", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "kv", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "occurrence_count", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(eventAttrs, "project_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "project_name", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "request_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "session_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "severity", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "source", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "summary", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "timestamp", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_email", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_id", projectEventsBuilder, computedString())
	capellaschema.AddAttr(eventAttrs, "user_name", projectEventsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", projectEventsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: eventAttrs,
		},
	})

	// Build cursor attributes for pagination
	hrefsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(hrefsAttrs, "first", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "last", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "next", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "previous", projectEventsBuilder, computedString())

	pagesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(pagesAttrs, "last", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "next", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "page", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "per_page", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "previous", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "total_items", projectEventsBuilder, computedInt64())

	cursorAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cursorAttrs, "hrefs", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: hrefsAttrs,
	})
	capellaschema.AddAttr(cursorAttrs, "pages", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: pagesAttrs,
	})

	capellaschema.AddAttr(attrs, "cursor", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cursorAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all event information for a project. Events represent a trail of actions that users performs within Capella at project level.",
		Attributes:          attrs,
	}
}
