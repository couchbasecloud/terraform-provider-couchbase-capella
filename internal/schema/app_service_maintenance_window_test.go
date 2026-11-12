package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance_window"
)

func TestAppServiceMaintenanceWindowValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       AppServiceMaintenanceWindow
		expectError bool
	}{
		{
			name: "[POSITIVE] IDs passed via terraform apply",
			input: AppServiceMaintenanceWindow{
				OrganizationId: types.StringValue("a"),
				ProjectId:      types.StringValue("b"),
				ClusterId:      types.StringValue("c"),
				AppServiceId:   types.StringValue("d"),
			},
		},
		{
			name: "[POSITIVE] IDs passed via terraform import",
			input: AppServiceMaintenanceWindow{
				AppServiceId: types.StringValue("organization_id=a,project_id=b,cluster_id=c,app_service_id=d"),
			},
		},
		{
			name: "[NEGATIVE] import string without app_service_id",
			input: AppServiceMaintenanceWindow{
				AppServiceId: types.StringValue("organization_id=a,project_id=b,cluster_id=c"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IDs, err := tt.input.Validate()
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, map[Attr]string{
				OrganizationId: "a",
				ProjectId:      "b",
				ClusterId:      "c",
				AppServiceId:   "d",
			}, IDs)
		})
	}
}

func TestNewAppServiceMaintenanceWindow(t *testing.T) {
	got := NewAppServiceMaintenanceWindow(&appservice_maintenance_window.MaintenanceWindow{
		PreferredDay:    "sunday",
		PreferredHour:   2,
		PreferredMinute: 30,
	}, "a", "b", "c", "d")

	assert.Equal(t, AppServiceMaintenanceWindow{
		OrganizationId:  types.StringValue("a"),
		ProjectId:       types.StringValue("b"),
		ClusterId:       types.StringValue("c"),
		AppServiceId:    types.StringValue("d"),
		PreferredDay:    types.StringValue("sunday"),
		PreferredHour:   types.Int64Value(2),
		PreferredMinute: types.Int64Value(30),
	}, *got)
}
