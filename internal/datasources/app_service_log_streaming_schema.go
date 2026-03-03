package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceLogStreamingBuilder = capellaschema.NewSchemaBuilder("appServiceLogStreaming", "GetLogStreamingResponse")

// AppServiceLogStreamingSchema returns the schema for the app_service_log_streaming data source.
func AppServiceLogStreamingSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", appServiceLogStreamingBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", appServiceLogStreamingBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceLogStreamingBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceLogStreamingBuilder, requiredStringWithValidator())

	// Computed attributes returned by the API
	capellaschema.AddAttr(attrs, "output_type", appServiceLogStreamingBuilder, computedString())
	capellaschema.AddAttr(attrs, "config_state", appServiceLogStreamingBuilder, computedString())
	capellaschema.AddAttr(attrs, "streaming_state", appServiceLogStreamingBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the log streaming configuration and state for an App Service.",
		Attributes:          attrs,
	}
}
