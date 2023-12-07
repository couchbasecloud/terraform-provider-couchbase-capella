package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Backup maps Backup resource schema data to the response received from V4 Capella Public API.
type Backup struct {
	// RestoreTimes represents the number of times we have requested a restore.
	// It represents the incremental count each time we request a restore.
	RestoreTimes types.Number `tfsdk:"restore_times"`

	// Method represents the mechanism of the backup.
	// Enum: "incremental" "full"
	// Incremental backups include the data that has changed since the last scheduled backup.
	// Full backup includes all bucket data from the time the backup was created.
	Method types.String `tfsdk:"method"`

	// Date represents the time at which backup was created.
	Date types.String `tfsdk:"date"`

	// BucketName represents the name of the bucket to which the backup belongs to.
	BucketName types.String `tfsdk:"bucket_name"`

	// CycleId is the cycleId to the which the backup belongs to.
	CycleId types.String `tfsdk:"cycle_id"`

	// BucketId is the ID of the bucket to which the backup belongs to.
	BucketId types.String `tfsdk:"bucket_id"`

	// RestoreBefore represents the time at which backup will expire.
	RestoreBefore types.String `tfsdk:"restore_before"`

	// Status represents the status of the backup.
	// Enum: "pending" "ready" "failed"
	Status types.String `tfsdk:"status"`

	// Source represents the way a backup job was initiated.
	// Enum: "manual" "scheduled"
	// Manual represents a manually triggered backup job or on-demand.
	// Scheduled represents a backup job created from a schedule.
	Source types.String `tfsdk:"source"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// Id is a GUID4 identifier of the backup.
	Id types.String `tfsdk:"id"`

	// Provider is the cloud provider where the cluster is hosted.
	CloudProvider types.String `tfsdk:"cloud_provider"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// BackupStats represents various backup level data that couchbase provides.
	BackupStats types.Object `tfsdk:"backup_stats"`

	// ScheduleInfo represents the schedule information of the backup.
	ScheduleInfo types.Object `tfsdk:"schedule_info"`

	// Restore represents information about how to restore the backup.
	Restore types.Object `tfsdk:"restore"`

	// ElapsedTimeInSeconds represents the amount of seconds that have elapsed between the creation and completion of the backup.
	ElapsedTimeInSeconds types.Int64 `tfsdk:"elapsed_time_in_seconds"`
}

// BackupData defines attributes for a single Backup when fetched from the V4 Capella Public API.
type BackupData struct {
	// Method represents the mechanism of the backup.
	// Enum: "incremental" "full"
	// Incremental backups include the data that has changed since the last scheduled backup.
	// Full backup includes all bucket data from the time the backup was created.
	Method types.String `tfsdk:"method"`

	// RestoreBefore represents the time at which backup will expire.
	RestoreBefore types.String `tfsdk:"restore_before"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Id is a GUID4 identifier of the backup.
	Id types.String `tfsdk:"id"`

	// Date represents the time at which backup was created.
	Date types.String `tfsdk:"date"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Status represents the status of the backup.
	// Enum: "pending" "ready" "failed"
	Status types.String `tfsdk:"status"`

	// CycleId is the cycleId to the which the backup belongs to.
	CycleId types.String `tfsdk:"cycle_id"`

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

	// ScheduleInfo represents the schedule information of the backup.
	ScheduleInfo types.Object `tfsdk:"schedule_info"`

	// ElapsedTimeInSeconds represents the amount of seconds that have elapsed between the creation and completion of the backup.
	ElapsedTimeInSeconds types.Int64 `tfsdk:"elapsed_time_in_seconds"`
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

	// Retention represents retention time in days.
	Retention types.String `tfsdk:"retention"`

	// Increment represents interval in hours for incremental backup.
	Increment types.Int64 `tfsdk:"increment"`
}

// Restore provides information about how to restore the backup.
type Restore struct {
	// MapData is specified when you want to restore source data into a different location.
	MapData types.String `tfsdk:"map_data"`

	// SourceClusterId represents the Id of the source cluster the restore is based on.
	SourceClusterId types.String `tfsdk:"source_cluster_id"`

	// FilterKeys represents a regular expression. It is used to selectively
	// restore data, allowing only the restoration of data where the key
	// matches a specific regular expression.
	FilterKeys types.String `tfsdk:"filter_keys"`

	// FilterValues represents a regular expression. It is used to selectively
	// restore data, allowing restoration only when the value matches a
	// specific regular expression.
	FilterValues types.String `tfsdk:"filter_values"`

	// IncludeData when specified restores only the data specified here
	IncludeData types.String `tfsdk:"include_data"`

	// ExcludeData when specified, skips restoring the data specified here.
	ExcludeData types.String `tfsdk:"exclude_data"`

	// TargetClusterId represents the Id of the target cluster to restore to.
	TargetClusterId types.String `tfsdk:"target_cluster_id"`

	// ReplaceTTL sets a new expiration (time-to-live) value for the specified keys.
	ReplaceTTL types.String `tfsdk:"replace_ttl"`

	// ReplaceTTLWith updates the expiration for the keys.
	ReplaceTTLWith types.String `tfsdk:"replace_ttl_with"`

	// Status represents the status of restore.
	Status types.String `tfsdk:"status"`

	// Services represents the array of strings (Services) like data, query.
	Services []types.String `tfsdk:"services"`

	// ForceUpdates when marked true forces data in the Couchbase cluster to
	// be overwritten even if the data in the cluster is newer.
	ForceUpdates types.Bool `tfsdk:"force_updates"`

	// AutoRemoveCollections when marked true automatically delete scopes/collections
	// which are known to be deleted in the backup.
	AutoRemoveCollections types.Bool `tfsdk:"auto_remove_collections"`
}

func (r Restore) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_cluster_id": types.StringType,

		"source_cluster_id": types.StringType,

		"services": types.ListType{ElemType: types.StringType},

		"force_updates": types.BoolType,

		"auto_remove_collections": types.BoolType,

		"filter_keys": types.StringType,

		"filter_values": types.StringType,

		"include_data": types.StringType,

		"exclude_data": types.StringType,

		"map_data": types.StringType,

		"replace_ttl": types.StringType,

		"replace_ttl_with": types.StringType,

		"status": types.StringType,
	}
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

// NewBackupStats creates a new BackupStats data object.
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

// NewScheduleInfo creates a new ScheduleInfo data object.
func NewScheduleInfo(scheduleInfo backup.ScheduleInfo) ScheduleInfo {
	return ScheduleInfo{
		BackupType: types.StringValue(scheduleInfo.BackupType),
		BackupTime: types.StringValue(scheduleInfo.BackupTime),
		Increment:  types.Int64Value(scheduleInfo.Increment),
		Retention:  types.StringValue(scheduleInfo.Retention),
	}
}

// NewBackup creates new backup object.
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
		Restore:              types.ObjectNull(Restore{}.AttributeTypes()),
	}
	return &newBackup
}

// NewBackupData creates new backup data object.
func NewBackupData(backup *backup.GetBackupResponse,
	organizationId, projectId string,
	bStatsObj, sInfoObj basetypes.ObjectValue,
) *BackupData {
	newBackupData := BackupData{
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
	return &newBackupData
}

// Validate is used to verify that IDs have been properly imported.
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
	Data []BackupData `tfsdk:"data"`
}

// Validate is used to verify that IDs have been properly imported.
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
