package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogExportBuilder = capellaschema.NewSchemaBuilder("auditLogExport", "clusterAuditLogExport")

func AuditLogExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", auditLogExportBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "project_id", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "cluster_id", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "audit_log_download_url", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "expiration", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "start", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "end", auditLogExportBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "created_at", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "status", auditLogExportBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage audit log exports for an operational cluster. This allows you to export audit logs for a specific time period and download them for analysis. Audit Logs for the last 30 days can be requested, otherwise they are purged. A pre-signed URL to a s3 bucket location is returned, which is used to download these audit logs.",
		Attributes:          attrs,
	}
}
