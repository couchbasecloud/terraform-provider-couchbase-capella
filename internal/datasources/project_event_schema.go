package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectEventBuilder = capellaschema.NewSchemaBuilder("projectEvent")

// ProjectEventSchema returns the schema for the ProjectEvent data source.
func ProjectEventSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", projectEventBuilder, requiredString())
	capellaschema.AddAttr(attrs, "organization_id", projectEventBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", projectEventBuilder, requiredString())
	capellaschema.AddAttr(attrs, "alert_key", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "app_service_id", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "app_service_name", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "cluster_id", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "cluster_name", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "image_url", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "incident_ids", projectEventBuilder, computedStringSet())
	capellaschema.AddAttr(attrs, "key", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "kv", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "occurrence_count", projectEventBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "project_name", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "request_id", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "session_id", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "severity", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "source", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "summary", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "timestamp", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_email", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_id", projectEventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_name", projectEventBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an event for a project. Events represent a trail of actions that users performs within Capella at the project level.",
		Attributes:          attrs,
	}
}
