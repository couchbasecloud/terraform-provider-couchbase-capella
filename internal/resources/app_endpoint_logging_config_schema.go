package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var loggingConfigBuilder = capellaschema.NewSchemaBuilder("loggingConfig", "consoleLoggingConfig")

func LoggingConfigSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "app_endpoint_name", loggingConfigBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", loggingConfigBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", loggingConfigBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", loggingConfigBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "organization_id", loggingConfigBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "log_level", loggingConfigBuilder, stringAttribute([]string{required}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "log_keys", loggingConfigBuilder, stringSetAttribute(required))

	return schema.Schema{
		MarkdownDescription: "Manages the Log Streaming config for an App Endpoint on the App Service. This config is used to filter what logs to stream to the log collector that is configured on the App Service.",
		Attributes:          attrs,
	}
}
