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

	capellaschema.AddAttr(attrs, "organization_id", projectEventsBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "project_id", projectEventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "cluster_ids", projectEventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "user_ids", projectEventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "severity_levels", projectEventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "tags", projectEventsBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	})
	capellaschema.AddAttr(attrs, "from", projectEventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "to", projectEventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "page", projectEventsBuilder, &schema.Int64Attribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "per_page", projectEventsBuilder, &schema.Int64Attribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "sort_by", projectEventsBuilder, &schema.StringAttribute{
		Optional: true,
	})
	capellaschema.AddAttr(attrs, "sort_direction", projectEventsBuilder, &schema.StringAttribute{
		Optional: true,
	})

	attrs["data"] = computedEventAttributes
	attrs["cursor"] = computedCursorAttribute

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all event information for a project. Events represent a trail of actions that users performs within Capella at project level.",
		Attributes:          attrs,
	}
}
