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

	capellaschema.AddAttr(attrs, "id", eventBuilder, requiredString())
	capellaschema.AddAttr(attrs, "organization_id", eventBuilder, requiredString())
	capellaschema.AddAttr(attrs, "alert_key", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "app_service_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "app_service_name", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "cluster_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "cluster_name", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "image_url", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "incident_ids", eventBuilder, &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(attrs, "key", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "kv", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "occurrence_count", eventBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "project_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "project_name", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "request_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "session_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "severity", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "source", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "summary", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "timestamp", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_email", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_id", eventBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_name", eventBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an event in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes:          attrs,
	}
}
