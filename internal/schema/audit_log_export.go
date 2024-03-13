package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type AuditLogExport struct {
	// AuditLogDownloadURL Pre-signed URL to download cluster audit logs.
	AuditLogDownloadURL types.String `tfsdk:"auditlog_download_url"`

	// Expiration The timestamp when the download link expires.
	Expiration types.String `tfsdk:"expiration"`

	// Id is a GUID4 identifier of the export job.
	Id types.String `tfsdk:"id"`

	// Start Specifies the audit log's start date and time.
	Start types.String `tfsdk:"start"`

	// End Specifies the audit log's end date and time.
	End types.String `tfsdk:"end"`

	// CreatedAt The timestamp when the audit logs were exported.
	CreatedAt types.String `tfsdk:"created_at"`

	// Status Indicates status of audit log creation.
	Status types.String `tfsdk:"status"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`
}

type AuditLogExports struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data contains the list of resources.
	Data []AuditLogExport `tfsdk:"data"`
}

// Validate is used to verify that IDs have been properly imported.
func (a *AuditLogExport) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
