package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogExportBuilder = capellaschema.NewSchemaBuilder("auditLogExport")

// AuditLogExportSchema returns the schema for the AuditLogExport data source.
func AuditLogExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", auditLogExportBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", auditLogExportBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", auditLogExportBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", auditLogExportBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(dataAttrs, "project_id", auditLogExportBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(dataAttrs, "cluster_id", auditLogExportBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(dataAttrs, "audit_log_download_url", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expiration", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "start", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "end", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "created_at", auditLogExportBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", auditLogExportBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", auditLogExportBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve audit log exports for an operational cluster. It will show the pre-signed URL if the export was successful, a failure error if it was unsuccessful, or a message saying no audit logs available if there were no audit logs found.",
		Attributes:          attrs,
	}
}
