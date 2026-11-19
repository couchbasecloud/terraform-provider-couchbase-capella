package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppServiceAuditLogExportJobData is the datasource model for an audit log export job of an App Service.
type AppServiceAuditLogExportJobData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Id is the ID of the audit log export job.
	Id types.String `tfsdk:"id"`

	// TenantId is the ID of the tenant (organization) that owns the export job.
	TenantId types.String `tfsdk:"tenant_id"`

	// Status is the status of the export job.
	Status types.String `tfsdk:"status"`

	// DownloadId is the ID used to download the exported audit logs.
	DownloadId types.String `tfsdk:"download_id"`

	// DownloadExpires is the RFC3339 timestamp at which the download expires.
	DownloadExpires types.String `tfsdk:"download_expires"`

	// Audit contains the audit fields of the export job.
	Audit types.Object `tfsdk:"audit"`
}
