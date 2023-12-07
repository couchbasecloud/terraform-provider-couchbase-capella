package backup_schedule

// WeeklySchedule represents the weekly schedule of the backup.
type WeeklySchedule struct {
	// DayOfWeek represents the day of the week for the backup.
	DayOfWeek string `json:"dayOfWeek"`

	// RetentionTime represents the retention time in days.
	RetentionTime string `json:"retentionTime"`

	// StartAt represents the start hour of the backup.
	StartAt int64 `json:"startAt"`

	// IncrementalEvery represents the interval in hours for incremental backup.
	IncrementalEvery int64 `json:"incrementalEvery"`

	// CostOptimizedRetention optimizes backup retention to reduce total cost of ownership (TCO).
	CostOptimizedRetention bool `json:"costOptimizedRetention"`
}
