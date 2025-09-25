package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup_schedule"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type SnapshotBackupSchedule struct {
	OrganizationID types.String   `tfsdk:"organization_id"`
	ProjectID      types.String   `tfsdk:"project_id"`
	ID             types.String   `tfsdk:"id"`
	Interval       types.Int64    `tfsdk:"interval"`
	Retention      types.Int64    `tfsdk:"retention"`
	StartTime      types.String   `tfsdk:"start_time"`
	CopyToRegions  []types.String `tfsdk:"copy_to_regions"`
}

func (s SnapshotBackupSchedule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": types.StringType,
		"project_id":      types.StringType,
		"id":              types.StringType,
		"interval":        types.Int64Type,
		"retention":       types.Int64Type,
		"start_time":      types.StringType,
		"copy_to_regions": types.ListType{ElemType: types.StringType},
	}
}

func NewSnapshotBackupSchedule(snapshotBackupSchedule snapshot_backup_schedule.SnapshotBackupSchedule, organizationID, projectID, clusterID string) SnapshotBackupSchedule {
	return SnapshotBackupSchedule{
		OrganizationID: types.StringValue(organizationID),
		ProjectID:      types.StringValue(projectID),
		ID:             types.StringValue(clusterID),
		Interval:       types.Int64Value(int64(snapshotBackupSchedule.Interval)),
		Retention:      types.Int64Value(int64(snapshotBackupSchedule.Retention)),
		StartTime:      types.StringValue(snapshotBackupSchedule.StartTime),
		CopyToRegions:  ConvertStringList(snapshotBackupSchedule.CopyToRegions),
	}
}

// Validate is used to verify that IDs have been properly imported.
func (s SnapshotBackupSchedule) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: s.OrganizationID,
		ProjectId:      s.ProjectID,
		Id:             s.ID,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

func ConvertStringList(stringList []string) []basetypes.StringValue {
	var stringValueList []basetypes.StringValue
	for _, stringElement := range stringList {
		stringValueList = append(stringValueList, types.StringValue(stringElement))
	}
	return stringValueList
}
