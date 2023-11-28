package backup

// ScheduleInfo represents the schedule information of the backup.
type ScheduleInfo struct {
	BackupType string `json:"backupType"`
	BackupTime string `json:"backupTime"`
	Retention  string `json:"retention"`
	Increment  int64  `json:"increment"`
}
