package backup_schedule

// CreateBackupScheduleRequest is the request payload sent to the Capella V4 Public API in order to create a new backup schedule.
//
// Couchbase supports a robust scheduled backup and retention time policy as part of an overall disaster recovery plan for production data.
// Couchbase Capella supports scheduled and on-demand backups of bucket data.
// A backup can be restored to the same database where it was created or another database in the same organization.
// On setting up a backup schedule, the bucket automatically backs up the bucket based on the chosen schedule.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateBackupScheduleRequest struct {
	// Type represents whether the backup is a Weekly or Daily backup.
	Type string `json:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule WeeklySchedule `json:"weeklySchedule"`
}

// GetBackupScheduleResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing backup schedule for a bucket.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetBackupScheduleResponse struct {
	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule *WeeklySchedule `json:"weeklySchedule"`

	// Type represents whether the backup is a Weekly or Daily backup.
	Type string `json:"type"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId string `json:"clusterID"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId string `json:"bucketId"`
}

// UpdateBackupScheduleRequest is the request payload sent to the Capella V4 Public API in order to update the existing backup schedule.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type UpdateBackupScheduleRequest struct {
	// Type represents whether the backup is a Weekly or Daily backup.
	Type string `json:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule WeeklySchedule `json:"weeklySchedule"`
}
