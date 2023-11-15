package schema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-capella/internal/api/backup_schedule"
	"terraform-provider-capella/internal/errors"
)

// BackupSchedule defines the response as received from V4 Capella Public API when asked to create a new backup schedule.
type BackupSchedule struct {
	// Id is not present in the backup schedule, it is only present to support import
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// Type represents whether the backup is a Weekly or Daily backup.
	// e.g. 'weekly'
	Type types.String `tfsdk:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule types.Object `tfsdk:"weekly_schedule"`
}

// Validate checks the validity of an API key and extracts associated IDs.
// TODO : add unit testing
func (a *BackupSchedule) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		BucketId:       a.BucketId,
		Id:             a.Id,
	}

	IDs, err := validateBackupScheduleSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
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

func (b WeeklySchedule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"day_of_week":              types.StringType,
		"start_at":                 types.Int64Type,
		"incremental_every":        types.Int64Type,
		"retention_time":           types.StringType,
		"cost_optimized_retention": types.BoolType,
	}
}

// NewWeeklySchedule creates a new WeeklySchedule data object
func NewWeeklySchedule(weeklySchedule backup_schedule.WeeklySchedule) WeeklySchedule {
	return WeeklySchedule{
		DayOfWeek:              types.StringValue(weeklySchedule.DayOfWeek),
		StartAt:                types.Int64Value(weeklySchedule.StartAt),
		IncrementalEvery:       types.Int64Value(weeklySchedule.IncrementalEvery),
		RetentionTime:          types.StringValue(weeklySchedule.RetentionTime),
		CostOptimizedRetention: types.BoolValue(weeklySchedule.CostOptimizedRetention),
	}
}

// NewBackupSchedule creates new backup schedule object
func NewBackupSchedule(backupSchedule *backup_schedule.GetBackupScheduleResponse,
	organizationId, projectId string,
	scheduleObj basetypes.ObjectValue,
) *BackupSchedule {
	newBackup := BackupSchedule{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(backupSchedule.ClusterId),
		BucketId:       types.StringValue(backupSchedule.BucketId),
		Type:           types.StringValue(backupSchedule.Type),
		WeeklySchedule: scheduleObj,
	}
	return &newBackup
}

// validateBackupScheduleSchemaState validates that the IDs passed in as variadic
// parameters were successfully imported for backup schedule.
func validateBackupScheduleSchemaState(state map[Attr]basetypes.StringValue) (map[Attr]string, error) {
	IDs, keyParams := morphState(state)

	// If the state was passed in via terraform import we need to
	// retrieve the individual IDs from the ID string.
	if checkForImportString(state[OrganizationId]) {
		var err error
		IDs, err = splitImportString(state[Id].ValueString(), keyParams)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errors.ErrInvalidImport, err)
		}
	}

	err := checkKeysAndValuesForBackupSchedule(IDs, keyParams)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	delete(IDs, Id)
	return IDs, nil
}

// checkKeysAndValuesForBackupSchedule is used to validate that an ID map
// has been populated with the expected ID keys and that the
// associated values are not empty for backup schedule and skip id key.
func checkKeysAndValuesForBackupSchedule(IDs map[Attr]string, keyParams []Attr) error {
	for _, key := range keyParams {
		if key == "id" {
			continue
		}
		value, ok := IDs[key]
		if !ok {
			return fmt.Errorf("terraform resource was missing: %w: %s", errors.ErrIdMissing, key)

		}
		if value == "" {
			return fmt.Errorf("terraform resource was empty: %w: %s", errors.ErrIdMissing, key)

		}
	}
	return nil
}
