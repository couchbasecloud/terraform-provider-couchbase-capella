package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var loggingConfigBuilder = capellaschema.NewSchemaBuilder("loggingConfig", "consoleLoggingConfig")

func LoggingConfigSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "app_endpoint_name", loggingConfigBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_service_id", loggingConfigBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", loggingConfigBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", loggingConfigBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "organization_id", loggingConfigBuilder, requiredStringWithValidator())

	capellaschema.AddAttr(attrs, "log_level", loggingConfigBuilder, computedString())
	capellaschema.AddAttr(attrs, "log_keys", loggingConfigBuilder, computedStringSet())

	return schema.Schema{
		MarkdownDescription: "Retrieves the App Endpoint Log Streaming config. This config is used to filter what logs to stream to the log collector that is configured on the App Service.",
		Attributes:          attrs,
	}
}
