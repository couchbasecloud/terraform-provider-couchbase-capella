package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLoggingStateBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLoggingState", "CreateAppServiceAuditLogRequest")

// AppServiceAuditLoggingStateSchema returns the schema for the app_service_audit_logging_state data source.
func AppServiceAuditLoggingStateSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLoggingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLoggingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLoggingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLoggingStateBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "audit_enabled", appServiceAuditLoggingStateBuilder, computedBool())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit logging state of an App Service.",
		Attributes:          attrs,
	}
}
