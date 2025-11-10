package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var eventBuilder = capellaschema.NewSchemaBuilder("event")

// EventSchema returns the schema for the Event data source.
func EventSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", eventBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "organization_id", eventBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "alert_key", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "app_service_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "app_service_name", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "cluster_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "cluster_name", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "image_url", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "incident_ids", eventBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(attrs, "key", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "kv", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "occurrence_count", eventBuilder, &schema.Int64Attribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "project_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "project_name", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "request_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "session_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "severity", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "source", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "summary", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "timestamp", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_email", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_id", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_name", eventBuilder, &schema.StringAttribute{
		Computed: true,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an event in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes:          attrs,
	}
}
