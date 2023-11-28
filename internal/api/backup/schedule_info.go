package backup

// ScheduleInfo represents the schedule information of the backup.
type ScheduleInfo struct {
	// BackupType represents whether the backup is a Weekly or Daily backup.
	BackupType string `json:"backupType"`

	// BackupTime is the timestamp indicating the backup created time.
	BackupTime string `json:"backupTime"`

	// Retention represents retention time in days.
	Retention string `json:"retention"`

	// Increment represents interval in hours for incremental backup.
	Increment int64 `json:"increment"`
}
