package backup

// WeeklySchedule represents the weekly schedule of the backup.
type WeeklySchedule struct {
	// DayOfWeek represents the day of the week for the backup.
	DayOfWeek string `json:"dayOfWeek"`

	// StartAt represents the start hour of the backup.
	StartAt int64 `json:"startAt"`

	// IncrementalEvery represents the interval in hours for incremental backup.
	IncrementalEvery int64 `json:"incrementalEvery"`

	// RetentionTime represents the retention time in days.
	RetentionTime string `json:"retentionTime"`

	// CostOptimizedRetention optimizes backup retention to reduce total cost of ownership (TCO).
	CostOptimizedRetention bool `json:"costOptimizedRetention"`
}
