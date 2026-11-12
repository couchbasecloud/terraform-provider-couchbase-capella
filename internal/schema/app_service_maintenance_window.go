package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance_window"
)

// AppServiceMaintenanceWindow is the preferred weekly maintenance window of an App Service.
type AppServiceMaintenanceWindow struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the appServiceId of the capella tenant.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// PreferredDay is the lowercase day of the week the maintenance window starts on.
	PreferredDay types.String `tfsdk:"preferred_day"`

	// PreferredHour is the hour of the day (0-23) the maintenance window starts at.
	PreferredHour types.Int64 `tfsdk:"preferred_hour"`

	// PreferredMinute is the minute of the hour (0-59) the maintenance window starts at.
	PreferredMinute types.Int64 `tfsdk:"preferred_minute"`
}

// Validate is used to verify that IDs have been properly imported.
func (a *AppServiceMaintenanceWindow) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %w", err)
	}

	return IDs, nil
}

// NewAppServiceMaintenanceWindow creates a new AppServiceMaintenanceWindow from the V4 API response.
func NewAppServiceMaintenanceWindow(
	resp *appservice_maintenance_window.MaintenanceWindow,
	organizationId, projectId, clusterId, appServiceId string,
) *AppServiceMaintenanceWindow {
	return &AppServiceMaintenanceWindow{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		PreferredDay:    types.StringValue(resp.PreferredDay),
		PreferredHour:   types.Int64Value(resp.PreferredHour),
		PreferredMinute: types.Int64Value(resp.PreferredMinute),
	}
}
