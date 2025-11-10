package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectEventBuilder = capellaschema.NewSchemaBuilder("projectEvent")

// ProjectEventSchema returns the schema for the ProjectEvent data source.
func ProjectEventSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", projectEventBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "organization_id", projectEventBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "project_id", projectEventBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(attrs, "alert_key", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "app_service_id", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "app_service_name", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "cluster_id", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "cluster_name", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "image_url", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "incident_ids", projectEventBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(attrs, "key", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "kv", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "occurrence_count", projectEventBuilder, &schema.Int64Attribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "project_name", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "request_id", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "session_id", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "severity", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "source", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "summary", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "timestamp", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_email", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_id", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(attrs, "user_name", projectEventBuilder, &schema.StringAttribute{
		Computed: true,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an event for a project. Events represent a trail of actions that users performs within Capella at the project level.",
		Attributes:          attrs,
	}
}
