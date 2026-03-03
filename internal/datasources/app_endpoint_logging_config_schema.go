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
		MarkdownDescription: "Retrieves the Logging Config associated with an App Endpoint",
		Attributes:          attrs,
	}
}
