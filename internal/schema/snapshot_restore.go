package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
)

type Filter struct {
	// Name is the attribute to filter by.
	Name types.String `tfsdk:"name"`

	// Values is a set of values for the filter.
	Values types.Set `tfsdk:"values"`
}

type SnapshotRestore struct {
	ClusterID      types.String `tfsdk:"cluster_id"`
	CreatedAt      types.String `tfsdk:"created_at"`
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	RestoreTo      types.String `tfsdk:"restore_to"`
	Snapshot       types.String `tfsdk:"snapshot"`
	Status         types.String `tfsdk:"status"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

type SnapshotRestoreData struct {
	CreatedAt types.String `tfsdk:"created_at"`
	ID        types.String `tfsdk:"id"`
	RestoreTo types.String `tfsdk:"restore_to"`
	Snapshot  types.String `tfsdk:"snapshot"`
	Status    types.String `tfsdk:"status"`
}

type SnapshotRestores struct {
	ClusterID      types.String `tfsdk:"cluster_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	OrganizationID types.String `tfsdk:"organization_id"`

	Data []SnapshotRestoreData `tfsdk:"data"`

	Filters *Filter `tfsdk:"filter"`
}

func (s SnapshotRestore) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cluster_id":      types.StringType,
		"created_at":      types.StringType,
		"id":              types.StringType,
		"project_id":      types.StringType,
		"restore_to":      types.StringType,
		"snapshot":        types.StringType,
		"status":          types.StringType,
		"organization_id": types.StringType,
	}
}

func NewSnapshotRestore(snapshotRestore snapshot_backup.SnapshotRestore, clusterID, projectID, organizationID string) SnapshotRestore {
	return SnapshotRestore{
		ClusterID:      types.StringValue(clusterID),
		CreatedAt:      types.StringValue(snapshotRestore.CreatedAt),
		ID:             types.StringValue(snapshotRestore.ID),
		ProjectID:      types.StringValue(projectID),
		RestoreTo:      types.StringValue(snapshotRestore.RestoreTo),
		Snapshot:       types.StringValue(snapshotRestore.Snapshot),
		Status:         types.StringValue(string(snapshotRestore.Status)),
		OrganizationID: types.StringValue(organizationID),
	}
}

func NewSnapshotRestoreData(snapshotRestore snapshot_backup.SnapshotRestore, clusterID, projectID, organizationID string) SnapshotRestoreData {
	return SnapshotRestoreData{
		CreatedAt: types.StringValue(snapshotRestore.CreatedAt),
		ID:        types.StringValue(snapshotRestore.ID),
		RestoreTo: types.StringValue(snapshotRestore.RestoreTo),
		Snapshot:  types.StringValue(snapshotRestore.Snapshot),
		Status:    types.StringValue(string(snapshotRestore.Status)),
	}
}
