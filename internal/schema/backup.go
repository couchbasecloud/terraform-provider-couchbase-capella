package schema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-capella/internal/errors"
)

// Backup maps Backup resource schema data to the response received from V4 Capella Public API.
type Backup struct {
	// Id is a GUID4 identifier of the backup.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// CycleId is the cycleId to the which the backup belongs to.
	CycleId types.String `tfsdk:"cycle_id"`

	// Date represents the time at which backup was created.
	Date types.String `tfsdk:"date"`

	// RestoreBefore represents the time at which backup will expire.
	RestoreBefore types.String `tfsdk:"restore_before"`

	// Status represents the status of the backup.
	// Enum: "pending" "ready" "failed"
	Status types.String `tfsdk:"status"`

	// Method represents the mechanism of the backup.
	// Enum: "incremental" "full"
	// Incremental backups include the data that has changed since the last scheduled backup.
	// Full backup includes all bucket data from the time the backup was created.
	Method types.String `tfsdk:"method"`

	// BucketName represents the name of the bucket to which the backup belongs to.
	BucketName types.String `tfsdk:"bucket_name"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// Source represents the way a backup job was initiated.
	// Enum: "manual" "scheduled"
	// Manual represents a manually triggered backup job or on-demand.
	// Scheduled represents a backup job created from a schedule.
	Source types.String `tfsdk:"source"`

	// Provider is the cloud provider where the cluster is hosted.
	Provider types.String `tfsdk:"provider"`

	// BackupStats represents various backup level data that couchbase provides.
	BackupStats types.Object `tfsdk:"backup_stats"`

	// ElapsedTimeInSeconds represents the amount of seconds that have elapsed between the creation and completion of the backup.
	ElapsedTimeInSeconds types.Int64 `tfsdk:"elapsed_time_in_seconds"`

	// ScheduleInfo represents the schedule information of the backup.
	ScheduleInfo types.Object `tfsdk:"schedule_info"`

	// Type represents whether the backup is a Weekly or Daily backup.
	Type types.String `tfsdk:"type"`

	// WeeklySchedule represents the weekly schedule of the backup.
	WeeklySchedule types.Object `tfsdk:"weekly_schedule"`
}

// BackupStats has the backup level stats provided by Couchbase.
type BackupStats struct {
	// SizeInMB represents backup size in megabytes.
	SizeInMB types.Int64 `tfsdk:"size_in_mb"`

	// Items is the number of items saved during the backup.
	Items types.Int64 `tfsdk:"items"`

	// Mutations is the number of mutations saved during the backup.
	Mutations types.Int64 `tfsdk:"mutations"`

	// Tombstones is the number of tombstones saved during the backup.
	Tombstones types.Int64 `tfsdk:"tombstones"`

	// GSI is the number of global secondary indexes saved during the backup.
	GSI types.Int64 `tfsdk:"gsi"`

	// FTS is the number of full text search entities saved during the backup.
	FTS types.Int64 `tfsdk:"fts"`

	// CBAS is the number of analytics entities saved during the backup.
	CBAS types.Int64 `tfsdk:"cbas"`

	// Event represents the number of event entities saved during the backup.
	Event types.Int64 `tfsdk:"event"`
}

// ScheduleInfo provides schedule information of the backup
type ScheduleInfo struct {
	// BackupType represents whether the backup is a Weekly or Daily backup.
	BackupType types.String `tfsdk:"backup_type"`

	// BackupTime is the timestamp indicating the backup created time.
	BackupTime types.String `tfsdk:"backup_time"`

	// Increment represents interval in hours for incremental backup.
	Increment types.Int64 `tfsdk:"increment"`

	// Retention represents retention time in days.
	Retention types.String `tfsdk:"retention"`
}

// Validate is used to verify that IDs have been properly imported
func (b Backup) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: b.OrganizationId,
		ProjectId:      b.ProjectId,
		ClusterId:      b.ClusterId,
		Id:             b.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}
