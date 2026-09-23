package backup

import "time"

// CreateBucketBackupExportResponse is the response received from the Capella V4 Public API after
// queueing a bucket backup export job.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
type CreateBucketBackupExportResponse struct {
	// ExportId is the ID of the export job. Poll the get backup export endpoint with this ID to
	// track progress and obtain the download URL.
	ExportId string `json:"exportId"`

	// CreatedAt is the time at which the export job was created.
	CreatedAt time.Time `json:"createdAt"`
}

// GetBucketBackupExportResponse is the response received from the Capella V4 Public API when asked
// to fetch the details of a bucket backup export job.
//
// The archive details are only populated once the export completes, so they are all pointers.
type GetBucketBackupExportResponse struct {
	// Id is the ID of the backup export job.
	Id string `json:"id"`

	// TenantId is the ID of the organization.
	TenantId string `json:"tenantId"`

	// ProjectId is the ID of the project.
	ProjectId string `json:"projectId"`

	// ClusterId is the ID of the cluster the exported backup was taken on.
	ClusterId string `json:"clusterId"`

	// CycleId is the ID of the backup cycle that was exported.
	CycleId string `json:"cycleId"`

	// BucketName is the name of the bucket the exported backup belongs to.
	BucketName string `json:"bucketName"`

	// Status is the status of the export job.
	Status string `json:"status"`

	// CreatedAt is the time at which the export was requested.
	CreatedAt time.Time `json:"createdAt"`

	// SizeInBytes is the size of the exported archive.
	SizeInBytes *int64 `json:"sizeInBytes,omitempty"`

	// Sha256Checksum is the SHA-256 hash of the exported archive.
	Sha256Checksum *string `json:"sha256Checksum,omitempty"`

	// Expiration is the time at which the exported archive is deleted from cloud storage.
	Expiration *time.Time `json:"expiration,omitempty"`

	// BackupDownloadURL is a pre-signed URL to download the exported archive. Capella generates a
	// fresh one on every request and each is valid for an hour, so it is never worth caching.
	BackupDownloadURL *string `json:"backupDownloadURL,omitempty"`
}
