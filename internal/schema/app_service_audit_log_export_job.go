package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServiceAuditLogExportJobResource is the resource model for an App Service audit log export job.
type AppServiceAuditLogExportJobResource struct {
	// Id is the ID of the export job.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the ID of the organization to which the App Service belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Service belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Service belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Start is the start date and time of the exported audit logs.
	Start types.String `tfsdk:"start"`

	// End is the end date and time of the exported audit logs.
	End types.String `tfsdk:"end"`

	// Status is the status of the export job.
	Status types.String `tfsdk:"status"`

	// DownloadId is the ID used to download the exported audit logs.
	DownloadId types.String `tfsdk:"download_id"`

	// DownloadExpires is the time when the download expires.
	DownloadExpires types.String `tfsdk:"download_expires"`

	// Audit contains the audit information for the export job.
	Audit types.Object `tfsdk:"audit"`
}

// Validate is used to verify that IDs have been properly imported.
func (a *AppServiceAuditLogExportJobResource) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
