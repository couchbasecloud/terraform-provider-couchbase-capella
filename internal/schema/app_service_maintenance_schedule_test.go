package schema

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestAppServiceMaintenanceScheduleValidate(t *testing.T) {
	expectedIDs := map[Attr]string{
		OrganizationId: "o",
		ProjectId:      "p",
		ClusterId:      "c",
		AppServiceId:   "a",
		MaintenanceId:  "req1",
	}

	tests := []struct {
		name        string
		input       AppServiceMaintenanceSchedule
		expectedErr error
	}{
		{
			name: "[POSITIVE] IDs passed via terraform apply",
			input: AppServiceMaintenanceSchedule{
				OrganizationId: types.StringValue("o"),
				ProjectId:      types.StringValue("p"),
				ClusterId:      types.StringValue("c"),
				AppServiceId:   types.StringValue("a"),
				MaintenanceId:  types.StringValue("req1"),
			},
		},
		{
			name: "[POSITIVE] IDs passed via terraform import",
			input: AppServiceMaintenanceSchedule{
				MaintenanceId: types.StringValue("organization_id=o,project_id=p,cluster_id=c,app_service_id=a,maintenance_id=req1"),
			},
		},
		{
			name: "[NEGATIVE] import string without maintenance_id",
			input: AppServiceMaintenanceSchedule{
				MaintenanceId: types.StringValue("organization_id=o,project_id=p,cluster_id=c,app_service_id=a"),
			},
			expectedErr: errors.ErrInvalidImport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IDs, err := tt.input.Validate()
			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, expectedIDs, IDs)
		})
	}
}

// The maintenance_id import key must not change how existing import strings parse.
func TestAppServiceMaintenanceScheduleImportKeyRegression(t *testing.T) {
	IDs, err := (&AllowList{Id: types.StringValue("id=1,cluster_id=2,project_id=3,organization_id=4")}).Validate()

	require.NoError(t, err)
	assert.Equal(t, map[Attr]string{Id: "1", ClusterId: "2", ProjectId: "3", OrganizationId: "4"}, IDs)
}

func TestNewAppServiceMaintenanceData(t *testing.T) {
	queueAt := time.Date(2026, 9, 5, 2, 30, 0, 0, time.UTC)
	version := "4.0.1"
	audit := types.ObjectNull(CouchbaseAuditData{}.AttributeTypes())

	tests := []struct {
		name     string
		input    appservice_maintenance.GetMaintenanceResponse
		expected func(*testing.T, AppServiceMaintenanceData)
	}{
		{
			name: "[POSITIVE] all fields set",
			input: appservice_maintenance.GetMaintenanceResponse{
				Id: "req1", Type: "upgrade", InitiatedBy: "couchbase", VersionId: &version,
				CurrentState: "pending", QueueAt: &queueAt, CompletedAt: &queueAt,
			},
			expected: func(t *testing.T, got AppServiceMaintenanceData) {
				assert.Equal(t, "req1", got.Id.ValueString())
				assert.Equal(t, "upgrade", got.Kind.ValueString())
				assert.Equal(t, "couchbase", got.InitiatedBy.ValueString())
				assert.Equal(t, "4.0.1", got.VersionId.ValueString())
				assert.Equal(t, "pending", got.CurrentState.ValueString())
				gotQueueAt, diags := got.QueueAt.ValueRFC3339Time()
				require.False(t, diags.HasError())
				assert.True(t, queueAt.Equal(gotQueueAt))
				assert.False(t, got.CompletedAt.IsNull())
			},
		},
		{
			name:  "[POSITIVE] omitted optional fields map to null",
			input: appservice_maintenance.GetMaintenanceResponse{Id: "req1", QueueAt: &queueAt},
			expected: func(t *testing.T, got AppServiceMaintenanceData) {
				assert.True(t, got.VersionId.IsNull())
				assert.True(t, got.CompletedAt.IsNull())
			},
		},
		{
			name:  "[POSITIVE] unknown state value is kept as-is",
			input: appservice_maintenance.GetMaintenanceResponse{Id: "req1", CurrentState: "paused"},
			expected: func(t *testing.T, got AppServiceMaintenanceData) {
				assert.Equal(t, "paused", got.CurrentState.ValueString())
				assert.True(t, got.QueueAt.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expected(t, NewAppServiceMaintenanceData(&tt.input, audit))
		})
	}
}
