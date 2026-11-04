package datasources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance_window"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppServiceMaintenanceWindow{}
	_ datasource.DataSourceWithConfigure = &AppServiceMaintenanceWindow{}
)

// errorCodeForMaintenanceWindowNotFound is the code of the V4 API error AppServiceMaintenanceWindowNotFound,
// returned when an App Service has no maintenance window set. It tells that 404 apart from the
// BackendNotFound 404 for an unknown App Service. The value has not been confirmed against the
// published error catalog; a wrong value makes a no-window read fail instead of returning nulls.
const errorCodeForMaintenanceWindowNotFound = 11074

// AppServiceMaintenanceWindow is the App Service maintenance window data source implementation.
type AppServiceMaintenanceWindow struct {
	*providerschema.Data
}

// NewAppServiceMaintenanceWindow is a helper function to simplify the provider implementation.
func NewAppServiceMaintenanceWindow() datasource.DataSource {
	return &AppServiceMaintenanceWindow{}
}

// Metadata returns the App Service maintenance window data source type name.
func (a *AppServiceMaintenanceWindow) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_maintenance_window"
}

// Schema defines the schema for the App Service maintenance window data source.
func (a *AppServiceMaintenanceWindow) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceMaintenanceWindowSchema()
}

// Read refreshes the Terraform state with the App Service maintenance window. When no window
// is set the preferred_* attributes are left null.
func (a *AppServiceMaintenanceWindow) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceMaintenanceWindow
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance window",
			"Could not read maintenance window for app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.AppServiceId]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenanceWindow",
		a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusNotFound && apiErr.Code == errorCodeForMaintenanceWindowNotFound {
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance window",
			"Could not read maintenance window for app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	var windowResp appservice_maintenance_window.GetMaintenanceWindowResponse
	if err = json.Unmarshal(response.Body, &windowResp); err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance window",
			"Could not read maintenance window, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, providerschema.NewAppServiceMaintenanceWindow(
		&windowResp, organizationId, projectId, clusterId, appServiceId))...)
}

// Configure adds the provider configured client to the App Service maintenance window data source.
func (a *AppServiceMaintenanceWindow) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.Data = data
}
