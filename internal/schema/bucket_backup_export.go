package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// Fallback descriptions for the bucket backup export fields, used only until the published OpenAPI
// spec carries them. The leading "\n - " matches the bullet format buildEnhancedDescription emits,
// so the generated docs do not change once the spec takes over.
const (
	ExportIdDescription          = "\n - The ID of the backup export job, returned when the export was created."
	ExportBackupIdDescription    = "\n - The GUID4 ID of the backup to export."
	ExportCycleIdDescription     = "\n - The ID of the backup cycle that was exported."
	ExportBucketNameDescription  = "\n - The name of the bucket the exported backup belongs to."
	ExportStatusDescription      = "\n - The status of the export job: pending, processing, complete, failed or expired."
	ExportCreatedAtDescription   = "\n - The RFC3339 timestamp at which the export was requested."
	ExportSizeInBytesDescription = "\n - The size of the exported archive in bytes. Present once the export is complete."
	ExportChecksumDescription    = "\n - The SHA-256 hash of the exported archive, for verifying the integrity of the downloaded file. Present once the export is complete."
	ExportExpirationDescription  = "\n - The RFC3339 timestamp at which the exported archive is deleted from cloud storage. After this time the backup must be exported again to be downloaded."
	ExportDownloadURLDescription = "\n - Pre-signed URL to download the exported backup archive. A fresh URL is generated on every read and each URL is valid for 1 hour. Present only while the export is complete and the archive has not expired."
)

// BucketBackupExport defines the Terraform schema for a bucket backup export resource.
type BucketBackupExport struct {
	// Id is the ID of the export job, returned when the export was created.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the Capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the Capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the Capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketId is the ID of the bucket the exported backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// BackupId is the ID of the backup to export.
	BackupId types.String `tfsdk:"backup_id"`

	// CycleId is the ID of the backup cycle that was exported.
	CycleId types.String `tfsdk:"cycle_id"`

	// BucketName is the name of the bucket the exported backup belongs to.
	BucketName types.String `tfsdk:"bucket_name"`

	// Status is the status of the export job.
	Status types.String `tfsdk:"status"`

	// CreatedAt is the time at which the export was requested.
	CreatedAt types.String `tfsdk:"created_at"`

	// SizeInBytes is the size of the exported archive.
	SizeInBytes types.Int64 `tfsdk:"size_in_bytes"`

	// Sha256Checksum is the SHA-256 hash of the exported archive.
	Sha256Checksum types.String `tfsdk:"sha256_checksum"`

	// Expiration is the time at which the exported archive is deleted from cloud storage.
	Expiration types.String `tfsdk:"expiration"`

	// BackupDownloadURL is a pre-signed URL to download the exported archive.
	BackupDownloadURL types.String `tfsdk:"backup_download_url"`
}

// BucketBackupExportData defines the Terraform schema for the bucket backup export data source.
// There is no list endpoint, so the export is addressed by ID rather than searched for.
type BucketBackupExportData struct {
	// Id is the ID of the export job to fetch.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the Capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the Capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the Capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketId is the ID of the bucket the exported backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// BackupId is the ID of the exported backup.
	BackupId types.String `tfsdk:"backup_id"`

	// CycleId is the ID of the backup cycle that was exported.
	CycleId types.String `tfsdk:"cycle_id"`

	// BucketName is the name of the bucket the exported backup belongs to.
	BucketName types.String `tfsdk:"bucket_name"`

	// Status is the status of the export job.
	Status types.String `tfsdk:"status"`

	// CreatedAt is the time at which the export was requested.
	CreatedAt types.String `tfsdk:"created_at"`

	// SizeInBytes is the size of the exported archive.
	SizeInBytes types.Int64 `tfsdk:"size_in_bytes"`

	// Sha256Checksum is the SHA-256 hash of the exported archive.
	Sha256Checksum types.String `tfsdk:"sha256_checksum"`

	// Expiration is the time at which the exported archive is deleted from cloud storage.
	Expiration types.String `tfsdk:"expiration"`

	// BackupDownloadURL is a pre-signed URL to download the exported archive.
	BackupDownloadURL types.String `tfsdk:"backup_download_url"`
}

// Validate is used to verify that IDs have been properly imported.
func (b *BucketBackupExport) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: b.OrganizationId,
		ProjectId:      b.ProjectId,
		ClusterId:      b.ClusterId,
		BucketId:       b.BucketId,
		BackupId:       b.BackupId,
		Id:             b.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
