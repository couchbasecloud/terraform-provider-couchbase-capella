package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type Progress struct {
	Status types.String `tfsdk:"status"`
	Time   types.String `tfsdk:"time"`
}

type CMEK struct {
	ID         types.String `tfsdk:"id"`
	ProviderID types.String `tfsdk:"provider_id"`
}

type Server struct {
	Version types.String `tfsdk:"version"`
}

type CrossRegionCopy struct {
	RegionCode types.String `tfsdk:"region_code"`
	Status     types.String `tfsdk:"status"`
	Time       types.String `tfsdk:"time"`
}

type SnapshotBackup struct {
	ClusterID         types.String   `tfsdk:"cluster_id"`
	CreatedAt         types.String   `tfsdk:"created_at"`
	Expiration        types.String   `tfsdk:"expiration"`
	ID                types.String   `tfsdk:"id"`
	Retention         types.Int64    `tfsdk:"retention"`
	RegionsToCopy     []types.String `tfsdk:"regions_to_copy"`
	CrossRegionCopies types.Set      `tfsdk:"cross_region_copies"`
	Progress          types.Object   `tfsdk:"progress"`
	CMEK              types.Set      `tfsdk:"cmek"`
	ProjectID         types.String   `tfsdk:"project_id"`
	Server            types.Object   `tfsdk:"server"`
	Size              types.Int64    `tfsdk:"size"`
	OrganizationId    types.String   `tfsdk:"organization_id"`
	Type              types.String   `tfsdk:"type"`
}

type SnapshotBackupData struct {
	CreatedAt         types.String `tfsdk:"created_at"`
	Expiration        types.String `tfsdk:"expiration"`
	ID                types.String `tfsdk:"id"`
	Retention         types.Int64  `tfsdk:"retention"`
	CrossRegionCopies types.Set    `tfsdk:"cross_region_copies"`
	Progress          types.Object `tfsdk:"progress"`
	CMEK              types.Set    `tfsdk:"cmek"`
	Server            types.Object `tfsdk:"server"`
	Size              types.Int64  `tfsdk:"size"`
	Type              types.String `tfsdk:"type"`
}

// SnapshotBackups defines structure based on the response received from V4 Capella Public API when asked to list snapshot backups.
type SnapshotBackups struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	// Data contains the list of resources.
	Data []SnapshotBackupData `tfsdk:"data"`
}

func (p Progress) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status": types.StringType,
		"time":   types.StringType,
	}
}

func NewProgress(progress snapshot_backup.Progress) Progress {
	return Progress{
		Status: types.StringValue(string(progress.Status)),
		Time:   types.StringValue(progress.Time),
	}
}

func (c CMEK) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"provider_id": types.StringType,
	}
}

func NewCMEK(cmek snapshot_backup.CMEK) CMEK {
	return CMEK{
		ID:         types.StringValue(cmek.ID),
		ProviderID: types.StringValue(cmek.ProviderID),
	}
}

func (s Server) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"version": types.StringType,
	}
}

func NewServer(server snapshot_backup.Server) Server {
	return Server{
		Version: types.StringValue(server.Version),
	}
}

func (c CrossRegionCopy) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region_code": types.StringType,
		"status":      types.StringType,
		"time":        types.StringType,
	}
}

func NewCrossRegionCopy(crossRegionCopy snapshot_backup.CrossRegionCopy) CrossRegionCopy {
	return CrossRegionCopy{
		RegionCode: types.StringValue(crossRegionCopy.RegionCode),
		Status:     types.StringValue(string(crossRegionCopy.Status)),
		Time:       types.StringValue(crossRegionCopy.Time),
	}
}

func (s SnapshotBackup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cluster_id":          types.StringType,
		"created_at":          types.StringType,
		"expiration":          types.StringType,
		"id":                  types.StringType,
		"progress":            types.ObjectType{AttrTypes: Progress{}.AttributeTypes()},
		"project_id":          types.StringType,
		"retention":           types.Int64Type,
		"regions_to_copy":     types.SetType{ElemType: types.StringType},
		"cross_region_copies": types.SetType{ElemType: types.ObjectType{AttrTypes: CrossRegionCopy{}.AttributeTypes()}},
		"cmek":                types.SetType{ElemType: types.ObjectType{AttrTypes: CMEK{}.AttributeTypes()}},
		"server":              types.ObjectType{AttrTypes: Server{}.AttributeTypes()},
		"size":                types.Int64Type,
		"organization_id":     types.StringType,
		"type":                types.StringType,
	}
}

func NewSnapshotBackup(ctx context.Context, snapshotBackup snapshot_backup.SnapshotBackup, ID, clusterID, projectID, organizationID string, progressObj, serverObj basetypes.ObjectValue, cmekSet, crossRegionCopySet basetypes.SetValue) SnapshotBackup {
	return SnapshotBackup{
		ID:                types.StringValue(ID),
		ClusterID:         types.StringValue(clusterID),
		Expiration:        types.StringValue(snapshotBackup.Expiration),
		ProjectID:         types.StringValue(projectID),
		OrganizationId:    types.StringValue(organizationID),
		CreatedAt:         types.StringValue(snapshotBackup.CreatedAt),
		Retention:         types.Int64Value(snapshotBackup.Retention),
		CrossRegionCopies: crossRegionCopySet,
		Progress:          progressObj,
		CMEK:              cmekSet,
		Server:            serverObj,
		Size:              types.Int64Value(int64(snapshotBackup.Size)),
		Type:              types.StringValue(snapshotBackup.Type),
	}
}

func NewSnapshotBackupData(snapshotBackup snapshot_backup.SnapshotBackup, ID, clusterID, projectID, organizationID string, progressObj, serverObj basetypes.ObjectValue, cmekSet, crossRegionCopySet basetypes.SetValue) SnapshotBackupData {
	return SnapshotBackupData{
		ID:                types.StringValue(ID),
		Expiration:        types.StringValue(snapshotBackup.Expiration),
		CreatedAt:         types.StringValue(snapshotBackup.CreatedAt),
		Retention:         types.Int64Value(snapshotBackup.Retention),
		CrossRegionCopies: crossRegionCopySet,
		Progress:          progressObj,
		CMEK:              cmekSet,
		Server:            serverObj,
		Size:              types.Int64Value(int64(snapshotBackup.Size)),
		Type:              types.StringValue(snapshotBackup.Type),
	}
}

// Validate is used to verify that IDs have been properly imported.
func (s SnapshotBackup) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: s.OrganizationId,
		ProjectId:      s.ProjectID,
		ClusterId:      s.ClusterID,
		Id:             s.ID,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

// ConvertRegionsToCopy is used to convert all regionsToCopy
// in an array of basetypes.StringValue to strings.
func ConvertRegionsToCopy(regionsToCopy []basetypes.StringValue) []string {
	var convertedRegionsToCopy []string
	for _, region := range regionsToCopy {
		convertedRegionsToCopy = append(convertedRegionsToCopy, region.ValueString())
	}
	return convertedRegionsToCopy
}
