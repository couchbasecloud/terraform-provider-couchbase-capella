package backup

type GetBackupResponse struct {
	// Id is a GUID4 identifier of the backup.
	Id string `json:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId string `json:"organizationID"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId string `json:"projectID"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId string `json:"clusterID"`

	// CycleId is the cycleId to the which the backup belongs to.
	CycleId string `json:"cycleID"`

	// Date represents the time at which backup was created.
	Date string `json:"date"`

	// RestoreBefore represents the time at which backup will expire.
	RestoreBefore string `json:"restoreBefore"`

	// Status represents the status of the backup.
	// Enum: "pending" "ready" "failed"
	Status string `json:"status"`

	// Method represents the mechanism of the backup.
	// Enum: "incremental" "full"
	// Incremental backups include the data that has changed since the last scheduled backup.
	// Full backup includes all bucket data from the time the backup was created.
	Method string `json:"method"`

	// BucketName represents the name of the bucket to which the backup belongs to.
	BucketName string `json:"bucketName"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId string `json:"bucketID"`

	// Source represents the way a backup job was initiated.
	// Enum: "manual" "scheduled"
	// Manual represents a manually triggered backup job or on-demand.
	// Scheduled represents a backup job created from a schedule.
	Source string `json:"source"`

	// Provider is the cloud provider where the cluster is hosted.
	Provider string `json:"provider"`

	// BackupStats represents various backup level data that couchbase provides.
	BackupStats BackupStats `json:"stats"`

	// ElapsedTimeInSeconds represents the amount of seconds that have elapsed between the creation and completion of the backup.
	ElapsedTimeInSeconds int `json:"elapsedTimeInSeconds"`

	// ScheduleInfo represents the schedule information of the backup.
	ScheduleInfo ScheduleInfo `json:"scheduleInfo"`

	// Type represents whether the backup is a Weekly or Daily backup.
	Type string `json:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule WeeklySchedule `json:"weeklySchedule"`
}

type CreateBackupRequest struct {
	// Type represents whether the backup is a Weekly or Daily backup.
	Type string `json:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule WeeklySchedule `json:"weeklySchedule"`
}
