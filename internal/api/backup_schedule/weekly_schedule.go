package backup_schedule

// WeeklySchedule represents the weekly schedule of the backup.
type WeeklySchedule struct {
	DayOfWeek              string `json:"dayOfWeek"`
	RetentionTime          string `json:"retentionTime"`
	StartAt                int64  `json:"startAt"`
	IncrementalEvery       int64  `json:"incrementalEvery"`
	CostOptimizedRetention bool   `json:"costOptimizedRetention"`
}
