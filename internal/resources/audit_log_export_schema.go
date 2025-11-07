package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogExportBuilder = capellaschema.NewSchemaBuilder("auditLogExport")

func AuditLogExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", auditLogExportBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "project_id", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "cluster_id", auditLogExportBuilder, stringAttribute([]string{required}))

	attrs["audit_log_download_url"] = WithDescription(stringAttribute([]string{computed}), "Pre-signed URL to download cluster audit logs.")
	attrs["expiration"] = WithDescription(stringAttribute([]string{computed}), "The timestamp when the download link expires. The timestamp when the audit log export will expire and no longer be available for download.")
	attrs["start"] = WithDescription(stringAttribute([]string{required}), "The start timestamp for the audit log export in RFC3339 format (e.g., '2024-01-01T00:00:00Z'). This defines the beginning of the time period to export logs from.")
	attrs["end"] = WithDescription(stringAttribute([]string{required}), "The end timestamp for the audit log export in RFC3339 format (e.g., '2024-01-02T00:00:00Z'). This defines the end of the time period to export logs from.")
	attrs["created_at"] = WithDescription(stringAttribute([]string{computed}), "The timestamp when this audit log export job was created.")
	attrs["status"] = WithDescription(stringAttribute([]string{computed}), "The current status of the audit log export job. Audit log export job statuses are 'queued', 'in progress', 'completed', or 'failed'.")

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage audit log exports for an operational cluster. This allows you to export audit logs for a specific time period and download them for analysis. Audit Logs for the last 30 days can be requested, otherwise they are purged. A pre-signed URL to a s3 bucket location is returned, which is used to download these audit logs.",
		Attributes:          attrs,
	}
}
