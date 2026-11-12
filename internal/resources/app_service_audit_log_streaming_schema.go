package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogStreamingBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLogStreaming", "PutAuditLogStreamingRequest")

// AppServiceAuditLogStreamingSchema returns the schema for the app_service_audit_log_streaming resource.
func AppServiceAuditLogStreamingSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogStreamingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogStreamingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogStreamingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogStreamingBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "streaming_enabled", appServiceAuditLogStreamingBuilder, boolAttribute(required))
	capellaschema.AddAttr(attrs, "output_type", appServiceAuditLogStreamingBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "credentials", appServiceAuditLogStreamingBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Sensitive:  true,
		Attributes: buildCredentialsAttributes(),
	})
	capellaschema.AddAttr(attrs, "disabled_app_endpoints", appServiceAuditLogStreamingBuilder, stringSetAttribute(optional))

	capellaschema.AddAttr(attrs, "log_streaming_state", appServiceAuditLogStreamingBuilder, stringAttribute([]string{computed}), "GetAuditLogStreamingResponse")

	return schema.Schema{
		Attributes: attrs,
		MarkdownDescription: `Manages audit log streaming for an App Service.
Destroying this resource disables audit log streaming and removes the collector configuration, including credentials.
Note: The API does not return credentials. After an import, credentials are null in the state, so the next apply sends the configured credentials again.`,
	}
}
