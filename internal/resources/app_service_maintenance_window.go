package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance_window"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppServiceMaintenanceWindow{}
	_ resource.ResourceWithConfigure   = &AppServiceMaintenanceWindow{}
	_ resource.ResourceWithImportState = &AppServiceMaintenanceWindow{}
)

// AppServiceMaintenanceWindow is the App Service maintenance window resource implementation.
type AppServiceMaintenanceWindow struct {
	*providerschema.Data
}

// NewAppServiceMaintenanceWindow is a helper function to simplify the provider implementation.
func NewAppServiceMaintenanceWindow() resource.Resource {
	return &AppServiceMaintenanceWindow{}
}

// Metadata returns the App Service maintenance window resource type name.
func (a *AppServiceMaintenanceWindow) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_maintenance_window"
}

// Schema defines the schema for the App Service maintenance window resource.
func (a *AppServiceMaintenanceWindow) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceMaintenanceWindowSchema()
}

// Create sets the maintenance window. It never overwrites a window that already exists:
// the API answers 409 and the user is told to import it instead.
func (a *AppServiceMaintenanceWindow) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceMaintenanceWindow
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	createRequest := appservice_maintenance_window.CreateMaintenanceWindowRequest{
		PreferredDay:    plan.PreferredDay.ValueString(),
		PreferredHour:   plan.PreferredHour.ValueInt64(),
		PreferredMinute: plan.PreferredMinute.ValueInt64(),
	}

	cfg := api.EndpointCfg{Url: a.maintenanceWindowURL(organizationId, projectId, clusterId, appServiceId), Method: http.MethodPost, SuccessStatus: http.StatusNoContent}
	if _, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, createRequest, a.Token, nil); err != nil {
		var apiErr *api.Error
		if stderrors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusConflict {
			resp.Diagnostics.AddError(
				"App Service maintenance window already exists",
				fmt.Sprintf("A maintenance window already exists for app service %s. Import it with "+
					"`terraform import couchbase-capella_app_service_maintenance_window.<name> "+
					"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s` and manage it from Terraform: %s",
					appServiceId, organizationId, projectId, clusterId, appServiceId, apiErr.CompleteError()),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error creating App Service maintenance window",
			"Could not create maintenance window for app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.retrieveMaintenanceWindow(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading App Service maintenance window",
			"The maintenance window was created for app service "+appServiceId+" but could not be read back."+
				" Run `terraform apply --refresh-only` to update the state from remote: "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read refreshes the state with the maintenance window from Capella. A 404 means the window
// was removed outside Terraform; any other error leaves the state untouched.
func (a *AppServiceMaintenanceWindow) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceMaintenanceWindow
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	refreshedState, err := a.retrieveMaintenanceWindow(
		ctx,
		resourceIDs[providerschema.OrganizationId],
		resourceIDs[providerschema.ProjectId],
		resourceIDs[providerschema.ClusterId],
		resourceIDs[providerschema.AppServiceId],
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance window",
			"Could not read maintenance window for app service "+state.AppServiceId.String()+": "+errString,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update replaces the maintenance window with the planned values.
func (a *AppServiceMaintenanceWindow) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppServiceMaintenanceWindow
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating App Service maintenance window",
			"Could not update maintenance window for app service "+plan.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.AppServiceId]
	)

	updateRequest := appservice_maintenance_window.UpdateMaintenanceWindowRequest{
		PreferredDay:    plan.PreferredDay.ValueString(),
		PreferredHour:   plan.PreferredHour.ValueInt64(),
		PreferredMinute: plan.PreferredMinute.ValueInt64(),
	}

	cfg := api.EndpointCfg{Url: a.maintenanceWindowURL(organizationId, projectId, clusterId, appServiceId), Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	if _, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, updateRequest, a.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error updating App Service maintenance window",
			"Could not update maintenance window for app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	currentState, err := a.retrieveMaintenanceWindow(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance window",
			"Could not read maintenance window for app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, currentState)...)
}

// Delete removes the maintenance window. A window that is already gone is not an error.
func (a *AppServiceMaintenanceWindow) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceMaintenanceWindow
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting App Service maintenance window",
			"Could not delete maintenance window for app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	url := a.maintenanceWindowURL(
		resourceIDs[providerschema.OrganizationId],
		resourceIDs[providerschema.ProjectId],
		resourceIDs[providerschema.ClusterId],
		resourceIDs[providerschema.AppServiceId],
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	if _, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil); err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting App Service maintenance window",
			"Could not delete maintenance window for app service "+state.AppServiceId.String()+": "+errString,
		)
	}
}

// ImportState imports an existing maintenance window, for example one set through the Capella UI.
// The import ID is a comma separated string of the IDs, for example:
// "organization_id=<orgId>,project_id=<projId>,cluster_id=<clusterId>,app_service_id=<appServiceId>".
func (a *AppServiceMaintenanceWindow) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// Configure adds the provider configured client to the App Service maintenance window resource.
func (a *AppServiceMaintenanceWindow) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.Data = data
}

func (a *AppServiceMaintenanceWindow) maintenanceWindowURL(organizationId, projectId, clusterId, appServiceId string) string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenanceWindow",
		a.HostURL, organizationId, projectId, clusterId, appServiceId)
}

func (a *AppServiceMaintenanceWindow) retrieveMaintenanceWindow(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.AppServiceMaintenanceWindow, error) {
	cfg := api.EndpointCfg{Url: a.maintenanceWindowURL(organizationId, projectId, clusterId, appServiceId), Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, err
	}

	var windowResp appservice_maintenance_window.GetMaintenanceWindowResponse
	if err = json.Unmarshal(response.Body, &windowResp); err != nil {
		return nil, err
	}

	return providerschema.NewAppServiceMaintenanceWindow(&windowResp, organizationId, projectId, clusterId, appServiceId), nil
}
