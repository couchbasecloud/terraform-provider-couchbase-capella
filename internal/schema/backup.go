package schema

import (
	"fmt"
	"terraform-provider-capella/internal/api/backup"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	CloudProvider types.String `tfsdk:"cloud_provider"`

	// BackupStats represents various backup level data that couchbase provides.
	BackupStats types.Object `tfsdk:"backup_stats"`

	// ElapsedTimeInSeconds represents the amount of seconds that have elapsed between the creation and completion of the backup.
	ElapsedTimeInSeconds types.Int64 `tfsdk:"elapsed_time_in_seconds"`

	// ScheduleInfo represents the schedule information of the backup.
	ScheduleInfo types.Object `tfsdk:"schedule_info"`
}

// BackupStats has the backup level stats provided by Couchbase.
type BackupStats struct {
	// SizeInMB represents backup size in megabytes.
	SizeInMB types.Float64 `tfsdk:"size_in_mb"`

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

// ScheduleInfo provides schedule information of the backup.
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

func (b BackupStats) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"size_in_mb": types.Float64Type,
		"items":      types.Int64Type,
		"mutations":  types.Int64Type,
		"tombstones": types.Int64Type,
		"gsi":        types.Int64Type,
		"fts":        types.Int64Type,
		"cbas":       types.Int64Type,
		"event":      types.Int64Type,
	}
}

// NewBackupStats creates a new BackupStats data object
func NewBackupStats(backupStats backup.BackupStats) BackupStats {
	return BackupStats{
		SizeInMB:   types.Float64Value(backupStats.SizeInMB),
		Items:      types.Int64Value(backupStats.Items),
		Mutations:  types.Int64Value(backupStats.Mutations),
		Tombstones: types.Int64Value(backupStats.Tombstones),
		GSI:        types.Int64Value(backupStats.GSI),
		FTS:        types.Int64Value(backupStats.FTS),
		CBAS:       types.Int64Value(backupStats.CBAS),
		Event:      types.Int64Value(backupStats.Event),
	}
}

func (b ScheduleInfo) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_type": types.StringType,
		"backup_time": types.StringType,
		"increment":   types.Int64Type,
		"retention":   types.StringType,
	}
}

// NewScheduleInfo creates a new ScheduleInfo data object
func NewScheduleInfo(scheduleInfo backup.ScheduleInfo) ScheduleInfo {
	return ScheduleInfo{
		BackupType: types.StringValue(scheduleInfo.BackupType),
		BackupTime: types.StringValue(scheduleInfo.BackupTime),
		Increment:  types.Int64Value(scheduleInfo.Increment),
		Retention:  types.StringValue(scheduleInfo.Retention),
	}
}

// NewBackup creates new backup object
func NewBackup(backup *backup.GetBackupResponse,
	organizationId, projectId string,
	bStatsObj, sInfoObj basetypes.ObjectValue,
) *Backup {
	newBackup := Backup{
		Id:                   types.StringValue(backup.Id),
		OrganizationId:       types.StringValue(organizationId),
		ProjectId:            types.StringValue(projectId),
		ClusterId:            types.StringValue(backup.ClusterId),
		CycleId:              types.StringValue(backup.CycleId),
		Date:                 types.StringValue(backup.Date),
		RestoreBefore:        types.StringValue(backup.RestoreBefore),
		Status:               types.StringValue(string(backup.Status)),
		Method:               types.StringValue(backup.Method),
		BucketName:           types.StringValue(backup.BucketName),
		BucketId:             types.StringValue(backup.BucketId),
		Source:               types.StringValue(backup.Source),
		CloudProvider:        types.StringValue(backup.CloudProvider),
		BackupStats:          bStatsObj,
		ScheduleInfo:         sInfoObj,
		ElapsedTimeInSeconds: types.Int64Value(backup.ElapsedTimeInSeconds),
	}
	return &newBackup
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

// Backups defines structure based on the response received from V4 Capella Public API when asked to list backups.
type Backups struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// Data contains the list of resources.
	Data []Backup `tfsdk:"data"`
}

// Validate is used to verify that IDs have been properly imported
func (b Backups) Validate() (bucketId, clusterId, projectId, organizationId string, err error) {
	if b.BucketId.IsNull() {
		return "", "", "", "", errors.ErrBucketIdMissing
	}
	if b.OrganizationId.IsNull() {
		return "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if b.ProjectId.IsNull() {
		return "", "", "", "", errors.ErrProjectIdMissing
	}
	if b.ClusterId.IsNull() {
		return "", "", "", "", errors.ErrClusterIdMissing
	}

	return b.BucketId.ValueString(), b.ClusterId.ValueString(), b.ProjectId.ValueString(), b.OrganizationId.ValueString(), nil
}
