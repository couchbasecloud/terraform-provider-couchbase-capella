package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogStreamingStateBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLogStreamingState", "GetAuditLogStreamingResponse")

// AppServiceAuditLogStreamingStateSchema returns the schema for the app_service_audit_log_streaming_state data source.
func AppServiceAuditLogStreamingStateSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogStreamingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogStreamingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogStreamingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogStreamingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "streaming_enabled", appServiceAuditLogStreamingStateBuilder, computedBool())
	capellaschema.AddAttr(attrs, "log_streaming_state", appServiceAuditLogStreamingStateBuilder, computedString())
	capellaschema.AddAttr(attrs, "output_type", appServiceAuditLogStreamingStateBuilder, computedString())
	capellaschema.AddAttr(attrs, "disabled_app_endpoints", appServiceAuditLogStreamingStateBuilder, computedStringSet())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit log streaming state and output type of an App Service.",
		Attributes:          attrs,
	}
}
