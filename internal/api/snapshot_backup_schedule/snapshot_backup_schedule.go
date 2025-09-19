package snapshot_backup_schedule

type SnapshotBackupSchedule struct {
	Interval  int    `json:"interval"`
	Retention int    `json:"retention"`
	StartTime string `json:"startTime"`
}
