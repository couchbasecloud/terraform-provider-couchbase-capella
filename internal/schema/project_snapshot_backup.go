package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
)

type ProjectSnapshot struct {
	ClusterID         types.String `tfsdk:"cluster_id"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Expiration        types.String `tfsdk:"expiration"`
	ID                types.String `tfsdk:"id"`
	Progress          types.Object `tfsdk:"progress"`
	ProjectID         types.String `tfsdk:"project_id"`
	AppService        types.String `tfsdk:"app_service"`
	CMEK              types.Set    `tfsdk:"cmek"`
	CrossRegionCopies types.Set    `tfsdk:"cross_region_copies"`
	Retention         types.Int64  `tfsdk:"retention"`
	Server            types.Object `tfsdk:"server"`
	DatabaseSize      types.Int64  `tfsdk:"database_size"`
	OrganizationId    types.String `tfsdk:"organization_id"`
	Type              types.String `tfsdk:"type"`
}

type ProjectSnapshotBackupData struct {
	ClusterID          types.String `tfsdk:"cluster_id"`
	ClusterName        types.String `tfsdk:"cluster_name"`
	CreationDateTime   types.String `tfsdk:"creation_date_time"`
	CreatedBy          types.String `tfsdk:"created_by"`
	CurrentStatus      types.String `tfsdk:"current_status"`
	CloudProvider      types.String `tfsdk:"cloud_provider"`
	Region             types.String `tfsdk:"region"`
	MostRecentSnapshot types.Object `tfsdk:"most_recent_snapshot"`
	OldestSnapshot     types.Object `tfsdk:"oldest_snapshot"`
}

type ProjectSnapshotBackups struct {
	OrganizationId types.String                `tfsdk:"organization_id"`
	ProjectId      types.String                `tfsdk:"project_id"`
	Page           types.Int64                 `tfsdk:"page"`
	PerPage        types.Int64                 `tfsdk:"per_page"`
	SortBy         types.String                `tfsdk:"sort_by"`
	SortDirection  types.String                `tfsdk:"sort_direction"`
	Data           []ProjectSnapshotBackupData `tfsdk:"data"`
	Cursor         *Cursor                     `tfsdk:"cursor"`
}

func (p ProjectSnapshot) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cluster_id":          types.StringType,
		"created_at":          types.StringType,
		"expiration":          types.StringType,
		"id":                  types.StringType,
		"progress":            types.ObjectType{AttrTypes: Progress{}.AttributeTypes()},
		"project_id":          types.StringType,
		"app_service":         types.StringType,
		"cmek":                types.SetType{ElemType: types.ObjectType{AttrTypes: CMEK{}.AttributeTypes()}},
		"cross_region_copies": types.SetType{ElemType: types.ObjectType{AttrTypes: CrossRegionCopy{}.AttributeTypes()}},
		"retention":           types.Int64Type,
		"server":              types.ObjectType{AttrTypes: Server{}.AttributeTypes()},
		"database_size":       types.Int64Type,
		"organization_id":     types.StringType,
		"type":                types.StringType,
	}
}

func NewProjectSnapshot(snapshotBackup snapshot_backup.ProjectSnapshot, organizationID, projectID string, progressObj, serverObj basetypes.ObjectValue, cmekSet, crossRegionCopySet basetypes.SetValue) ProjectSnapshot {
	return ProjectSnapshot{
		ClusterID:         types.StringValue(snapshotBackup.ClusterID),
		CreatedAt:         types.StringValue(snapshotBackup.CreatedAt),
		Expiration:        types.StringValue(snapshotBackup.Expiration),
		ID:                types.StringValue(snapshotBackup.ID),
		Progress:          progressObj,
		ProjectID:         types.StringValue(projectID),
		AppService:        types.StringValue(snapshotBackup.AppService),
		CMEK:              cmekSet,
		CrossRegionCopies: crossRegionCopySet,
		Retention:         types.Int64Value(snapshotBackup.Retention),
		Server:            serverObj,
		DatabaseSize:      types.Int64Value(int64(snapshotBackup.DatabaseSize)),
		OrganizationId:    types.StringValue(organizationID),
		Type:              types.StringValue(snapshotBackup.Type),
	}
}

func NewProjectSnapshotBackupData(projectSnapshotData snapshot_backup.ProjectSnapshotBackupData, mostRecentSnapshot, oldestSnapshot basetypes.ObjectValue) ProjectSnapshotBackupData {
	return ProjectSnapshotBackupData{
		ClusterID:          types.StringValue(projectSnapshotData.ClusterID),
		ClusterName:        types.StringValue(projectSnapshotData.ClusterName),
		CreationDateTime:   types.StringValue(projectSnapshotData.CreationDateTime),
		CreatedBy:          types.StringValue(projectSnapshotData.CreatedBy),
		CurrentStatus:      types.StringValue(projectSnapshotData.CurrentStatus),
		CloudProvider:      types.StringValue(projectSnapshotData.CloudProvider),
		Region:             types.StringValue(projectSnapshotData.Region),
		MostRecentSnapshot: mostRecentSnapshot,
		OldestSnapshot:     oldestSnapshot,
	}

}
