package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogExportJobBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLogExportJob", "GetAuditExportDocResponse")

// AppServiceAuditLogExportJobSchema returns the schema for the App Service audit log export job resource.
func AppServiceAuditLogExportJobSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", appServiceAuditLogExportJobBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogExportJobBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogExportJobBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogExportJobBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogExportJobBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "start", appServiceAuditLogExportJobBuilder, stringAttribute([]string{required, requiresReplace}, rfc3339TimestampValidator{attributeName: "start"}), "CreateClusterAuditLogExportRequest")
	capellaschema.AddAttr(attrs, "end", appServiceAuditLogExportJobBuilder, stringAttribute([]string{required, requiresReplace}, rfc3339TimestampValidator{attributeName: "end"}), "CreateClusterAuditLogExportRequest")
	capellaschema.AddAttr(attrs, "status", appServiceAuditLogExportJobBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "download_id", appServiceAuditLogExportJobBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "download_expires", appServiceAuditLogExportJobBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", appServiceAuditLogExportJobBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "This resource allows you to export the audit logs of an App Service for a specific time period. The export job cannot be updated or deleted. Destroying the resource only removes it from the Terraform state.",
		Attributes:          attrs,
	}
}
