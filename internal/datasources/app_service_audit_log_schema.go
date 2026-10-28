package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLog", "CreateAppServiceAuditLogRequest")

// AppServiceAuditLogSchema returns the schema for the app_service_audit_log data source.
func AppServiceAuditLogSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "audit_enabled", appServiceAuditLogBuilder, computedBool())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit logging state of an App Service.",
		Attributes:          attrs,
	}
}
