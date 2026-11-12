package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogExportJobBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLogExportJob", "GetAuditExportDocResponse")

// AppServiceAuditLogExportJobSchema returns the schema for the app_service_audit_log_export_job data source.
// The output field names match the nested fields of the app_service_audit_log_export_jobs data source.
func AppServiceAuditLogExportJobSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogExportJobBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogExportJobBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogExportJobBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogExportJobBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "id", appServiceAuditLogExportJobBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "tenant_id", appServiceAuditLogExportJobBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", appServiceAuditLogExportJobBuilder, computedString())
	capellaschema.AddAttr(attrs, "download_id", appServiceAuditLogExportJobBuilder, computedString())
	capellaschema.AddAttr(attrs, "download_expires", appServiceAuditLogExportJobBuilder, computedString())
	capellaschema.AddAttr(attrs, "audit", appServiceAuditLogExportJobBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an audit log export job of an App Service.",
		Attributes:          attrs,
	}
}
