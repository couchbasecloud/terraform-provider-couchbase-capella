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
	RestoreTimes         types.Number `tfsdk:"restore_times"`
	Method               types.String `tfsdk:"method"`
	Date                 types.String `tfsdk:"date"`
	BucketName           types.String `tfsdk:"bucket_name"`
	CycleId              types.String `tfsdk:"cycle_id"`
	BucketId             types.String `tfsdk:"bucket_id"`
	RestoreBefore        types.String `tfsdk:"restore_before"`
	Status               types.String `tfsdk:"status"`
	Source               types.String `tfsdk:"source"`
	ClusterId            types.String `tfsdk:"cluster_id"`
	ProjectId            types.String `tfsdk:"project_id"`
	Id                   types.String `tfsdk:"id"`
	CloudProvider        types.String `tfsdk:"cloud_provider"`
	BackupStats          types.Object `tfsdk:"backup_stats"`
	OrganizationId       types.String `tfsdk:"organization_id"`
	ScheduleInfo         types.Object `tfsdk:"schedule_info"`
	Restore              types.Object `tfsdk:"restore"`
	ElapsedTimeInSeconds types.Int64  `tfsdk:"elapsed_time_in_seconds"`
}

// BackupData defines attributes for a single Backup when fetched from the V4 Capella Public API.
type BackupData struct {
	Method               types.String `tfsdk:"method"`
	RestoreBefore        types.String `tfsdk:"restore_before"`
	ProjectId            types.String `tfsdk:"project_id"`
	ClusterId            types.String `tfsdk:"cluster_id"`
	Id                   types.String `tfsdk:"id"`
	Date                 types.String `tfsdk:"date"`
	OrganizationId       types.String `tfsdk:"organization_id"`
	Status               types.String `tfsdk:"status"`
	CycleId              types.String `tfsdk:"cycle_id"`
	BucketName           types.String `tfsdk:"bucket_name"`
	BucketId             types.String `tfsdk:"bucket_id"`
	Source               types.String `tfsdk:"source"`
	CloudProvider        types.String `tfsdk:"cloud_provider"`
	BackupStats          types.Object `tfsdk:"backup_stats"`
	ScheduleInfo         types.Object `tfsdk:"schedule_info"`
	ElapsedTimeInSeconds types.Int64  `tfsdk:"elapsed_time_in_seconds"`
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
	BackupType types.String `tfsdk:"backup_type"`
	BackupTime types.String `tfsdk:"backup_time"`
	Retention  types.String `tfsdk:"retention"`
	Increment  types.Int64  `tfsdk:"increment"`
}

// Restore provides information about how to restore the backup.
type Restore struct {
	MapData               types.String   `tfsdk:"map_data"`
	SourceClusterId       types.String   `tfsdk:"source_cluster_id"`
	FilterKeys            types.String   `tfsdk:"filter_keys"`
	FilterValues          types.String   `tfsdk:"filter_values"`
	IncludeData           types.String   `tfsdk:"include_data"`
	ExcludeData           types.String   `tfsdk:"exclude_data"`
	TargetClusterId       types.String   `tfsdk:"target_cluster_id"`
	ReplaceTTL            types.String   `tfsdk:"replace_ttl"`
	ReplaceTTLWith        types.String   `tfsdk:"replace_ttl_with"`
	Status                types.String   `tfsdk:"status"`
	Services              []types.String `tfsdk:"services"`
	ForceUpdates          types.Bool     `tfsdk:"force_updates"`
	AutoRemoveCollections types.Bool     `tfsdk:"auto_remove_collections"`
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
