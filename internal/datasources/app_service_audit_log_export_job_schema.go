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
	// The OpenAPI spec has no descriptions for these fields, so the call site sets them.
	id := requiredStringWithValidator()
	id.MarkdownDescription = "The ID of the audit log export job."
	capellaschema.AddAttr(attrs, "id", appServiceAuditLogExportJobBuilder, id)
	tenantId := computedString()
	tenantId.MarkdownDescription = "The ID of the tenant (organization) that owns the export job."
	capellaschema.AddAttr(attrs, "tenant_id", appServiceAuditLogExportJobBuilder, tenantId)
	status := computedString()
	status.MarkdownDescription = "The status of the export job."
	capellaschema.AddAttr(attrs, "status", appServiceAuditLogExportJobBuilder, status)
	downloadId := computedString()
	downloadId.MarkdownDescription = "The ID used to download the exported audit logs."
	capellaschema.AddAttr(attrs, "download_id", appServiceAuditLogExportJobBuilder, downloadId)
	downloadExpires := computedString()
	downloadExpires.MarkdownDescription = "The RFC3339 timestamp at which the download expires."
	capellaschema.AddAttr(attrs, "download_expires", appServiceAuditLogExportJobBuilder, downloadExpires)
	capellaschema.AddAttr(attrs, "audit", appServiceAuditLogExportJobBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an audit log export job of an App Service.",
		Attributes:          attrs,
	}
}
