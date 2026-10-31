package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogExportJobsBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLogExportJobs", "GetAuditExportDocResponse")

// AppServiceAuditLogExportJobsSchema returns the schema for the app_service_audit_log_export_jobs data source.
func AppServiceAuditLogExportJobsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogExportJobsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogExportJobsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogExportJobsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogExportJobsBuilder, requiredUUIDString())

	// export_id filters the list in the provider, as the list endpoint has no filter query parameter.
	attrs["export_id"] = &schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The ID of an export job. When set, only the export job with this ID is returned. When omitted, all export jobs are returned.",
	}

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The ID of the audit log export job."))
	capellaschema.AddAttr(dataAttrs, "cluster_id", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The ID of the cluster for which the App Service is deployed."))
	capellaschema.AddAttr(dataAttrs, "app_service_id", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The ID of the App Service."))
	capellaschema.AddAttr(dataAttrs, "tenant_id", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The ID of the tenant (organization) that owns the export job."))
	capellaschema.AddAttr(dataAttrs, "status", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The status of the export job."))
	capellaschema.AddAttr(dataAttrs, "download_id", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The ID used to download the exported audit logs."))
	capellaschema.AddAttr(dataAttrs, "download_expires", appServiceAuditLogExportJobsBuilder, withDescription(computedString(), "The RFC3339 timestamp at which the download expires."))
	capellaschema.AddAttr(dataAttrs, "audit", appServiceAuditLogExportJobsBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", appServiceAuditLogExportJobsBuilder, &schema.ListNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The audit log export jobs of the App Service.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit log export jobs of an App Service, optionally filtered by export ID.",
		Attributes:          attrs,
	}
}

// withDescription sets a fallback description for fields the OpenAPI spec does not describe.
func withDescription(attr *schema.StringAttribute, description string) *schema.StringAttribute {
	attr.MarkdownDescription = description
	return attr
}
