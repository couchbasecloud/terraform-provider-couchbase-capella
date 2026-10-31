package schema

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// AppServiceAuditLogExportJobsData is the datasource model for the audit log export jobs of an App Service.
type AppServiceAuditLogExportJobsData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the App Service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// ExportId optionally filters the export jobs by ID. When null, all export jobs are returned.
	ExportId types.String `tfsdk:"export_id"`

	// Data is the list of audit log export jobs.
	Data []AppServiceAuditLogExportJob `tfsdk:"data"`
}

// AppServiceAuditLogExportJob is a single audit log export job of an App Service.
type AppServiceAuditLogExportJob struct {
	Id              types.String `tfsdk:"id"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	TenantId        types.String `tfsdk:"tenant_id"`
	Status          types.String `tfsdk:"status"`
	DownloadId      types.String `tfsdk:"download_id"`
	DownloadExpires types.String `tfsdk:"download_expires"`
	Audit           types.Object `tfsdk:"audit"`
}

// NewAppServiceAuditLogExportJob maps an API audit log export job to its Terraform model.
func NewAppServiceAuditLogExportJob(ctx context.Context, job apigen.GetAuditExportDocResponse) (AppServiceAuditLogExportJob, diag.Diagnostics) {
	var diags diag.Diagnostics

	auditObj := types.ObjectNull(CouchbaseAuditData{}.AttributeTypes())
	if job.Audit != nil {
		audit := NewCouchbaseAuditData(api.CouchbaseAuditData(*job.Audit))
		auditObj, diags = types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	}

	downloadExpires := types.StringNull()
	if job.DownloadExpires != nil {
		downloadExpires = types.StringValue(job.DownloadExpires.Format(time.RFC3339))
	}

	return AppServiceAuditLogExportJob{
		Id:              types.StringPointerValue(job.Id),
		ClusterId:       types.StringPointerValue(job.ClusterId),
		AppServiceId:    types.StringPointerValue(job.AppServiceId),
		TenantId:        types.StringPointerValue(job.TenantId),
		Status:          types.StringPointerValue(job.Status),
		DownloadId:      types.StringPointerValue(job.DownloadId),
		DownloadExpires: downloadExpires,
		Audit:           auditObj,
	}, diags
}
