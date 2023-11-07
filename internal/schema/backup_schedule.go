package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BackupSchedule defines the response as received from V4 Capella Public API when asked to create a new backup schedule.
type BackupSchedule struct {
	// Type represents whether the backup is a Weekly or Daily backup.
	// e.g. 'weekly'
	Type types.String `tfsdk:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule types.Object `tfsdk:"weekly_schedule"`
}

// WeeklySchedule represents the weekly schedule of the backup.
type WeeklySchedule struct {
	// DayOfWeek represents the day of the week for the backup.
	// Enum: "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"
	DayOfWeek types.String `tfsdk:"day_of_week"`

	// StartAt represents the start hour of the backup.
	// Enum: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23
	StartAt types.Int64 `tfsdk:"start_at"`

	// IncrementalEvery represents the interval in hours for incremental backup.
	// Enum: 1, 2, 4, 6, 8, 12, 24
	IncrementalEvery types.Int64 `tfsdk:"incremental_every"`

	// RetentionTime represents the retention time in days.
	// Enum: "30days", "60days", "90days", "180days", "1year", "2years", "3years", "4years", "5years"
	RetentionTime types.String `tfsdk:"retention_time"`

	// CostOptimizedRetention optimizes backup retention to reduce total cost of ownership (TCO).
	CostOptimizedRetention types.Bool `tfsdk:"cost_optimized_retention"`
}
