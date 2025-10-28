package snapshot_backup_schedule

type SnapshotBackupSchedule struct {
	Interval      int64    `json:"interval"`
	Retention     int64    `json:"retention"`
	StartTime     string   `json:"startTime"`
	CopyToRegions []string `json:"copyToRegions"`
}
