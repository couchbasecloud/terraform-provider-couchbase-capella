package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AuditLogExportSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to manage audit log exports for a Capella cluster. This allows you to export audit logs for a specific time period and download them for analysis. Audit Logs for the last 30 days can be requested, otherwise they are purged. A pre-signed URL to a s3 bucket location is returned, which is used to download these audit logs.",
		Attributes: map[string]schema.Attribute{
			"id": WithDescription(stringAttribute([]string{computed, useStateForUnknown}),
				"The unique identifier of the audit log export job."),
			"organization_id": WithDescription(stringAttribute([]string{required}),
				"The ID of the Capella organization."),
			"project_id": WithDescription(stringAttribute([]string{required}),
				"The ID of the Capella project that the cluster belongs to."),
			"cluster_id": WithDescription(stringAttribute([]string{required}),
				"The ID of the Capella cluster to export audit logs from."),
			"audit_log_download_url": WithDescription(stringAttribute([]string{computed}),
				"Pre-signed URL to download cluster audit logs."),
			"expiration": WithDescription(stringAttribute([]string{computed}),
				"The timestamp when the download link expires. The timestamp when the audit log export will expire and no longer be available for download."),
			"start": WithDescription(stringAttribute([]string{required}),
				"The start timestamp for the audit log export in RFC3339 format (e.g., '2024-01-01T00:00:00Z'). This defines the beginning of the time period to export logs from."),
			"end": WithDescription(stringAttribute([]string{required}),
				"The end timestamp for the audit log export in RFC3339 format (e.g., '2024-01-02T00:00:00Z'). This defines the end of the time period to export logs from."),
			"created_at": WithDescription(stringAttribute([]string{computed}),
				"The timestamp when this audit log export job was created."),
			"status": WithDescription(stringAttribute([]string{computed}),
				"The current status of the audit log export job. Possible values are 'queued', 'in progress', 'completed', or 'failed'."),
		},
	}
}
