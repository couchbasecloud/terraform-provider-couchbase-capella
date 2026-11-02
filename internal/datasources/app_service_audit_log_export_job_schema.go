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
	id := requiredStringWithValidator()
	id.MarkdownDescription = "The ID of the audit log export job."
	capellaschema.AddAttr(attrs, "id", appServiceAuditLogExportJobBuilder, id)
	capellaschema.AddAttr(attrs, "tenant_id", appServiceAuditLogExportJobBuilder, &schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The ID of the tenant (organization) that owns the export job.",
	})
	capellaschema.AddAttr(attrs, "status", appServiceAuditLogExportJobBuilder, &schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The status of the export job.",
	})
	capellaschema.AddAttr(attrs, "download_id", appServiceAuditLogExportJobBuilder, &schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The ID used to download the exported audit logs.",
	})
	capellaschema.AddAttr(attrs, "download_expires", appServiceAuditLogExportJobBuilder, &schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The RFC3339 timestamp at which the download expires.",
	})
	capellaschema.AddAttr(attrs, "audit", appServiceAuditLogExportJobBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an audit log export job of an App Service.",
		Attributes:          attrs,
	}
}
