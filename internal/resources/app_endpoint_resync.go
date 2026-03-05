package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpointResync{}
	_ resource.ResourceWithConfigure   = &AppEndpointResync{}
	_ resource.ResourceWithImportState = &AppEndpointResync{}
)

const errorAfterAppEndpointResyncInitiation = "Encountered an error while checking the current" +
	" state of the resync. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// AppEndpointResync is the App Endpoint Resync resource implementation.
type AppEndpointResync struct {
	*providerschema.Data
}

// NewAppEndpointResync is a helper function to simplify the provider implementation.
func NewAppEndpointResync() resource.Resource {
	return &AppEndpointResync{}
}

// Metadata returns the App Endpoint Resync resource type name.
func (a *AppEndpointResync) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_resync_job"
}

// Schema returns the schema for the App Endpoint Resync resource.
func (a *AppEndpointResync) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointResyncSchema()
}

// Create initiates an app endpoint resync.
func (a *AppEndpointResync) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpointResync
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	appEndpointName := plan.AppEndpoint.ValueString()

	var scopes map[string][]string
	diags = plan.Scopes.ElementsAs(ctx, &scopes, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	err := a.startResync(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, scopes)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error starting App Endpoint Resync",
			err.Error(),
		)
		return
	}

	// Set computed attributes to null before refreshing
	plan.CollectionsProcessing = types.MapNull(types.SetType{ElemType: types.StringType})
	plan.DocsChanged = types.Int64Null()
	plan.DocsProcessed = types.Int64Null()
	plan.LastError = types.StringNull()
	plan.StartTime = types.StringNull()
	plan.State = types.StringNull()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh the state by getting the latest data
	resyncResponse, err := a.getResyncStatus(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting App Endpoint Resync status in Capella",
			"Could not get Capella App Endpoint Resync status for app endpoint with name "+appEndpointName+": "+err.Error(),
		)
		return
	}

	refreshedState, diags := a.mapResponseToState(ctx, resyncResponse, &plan, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the current state of the app endpoint resync.
func (a *AppEndpointResync) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpointResync
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint in Capella",
			"Could not read App Endpoint  "+state.AppEndpoint.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId  = IDs[providerschema.OrganizationId]
		projectId       = IDs[providerschema.ProjectId]
		clusterId       = IDs[providerschema.ClusterId]
		appServiceId    = IDs[providerschema.AppServiceId]
		appEndpointName = IDs[providerschema.AppEndpointName]
	)

	resyncResponse, err := a.getResyncStatus(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting App Endpoint Resync status in Capella",
			errorAfterAppEndpointResyncInitiation+err.Error(),
		)
		return
	}

	refreshedState, diags := a.mapResponseToState(ctx, resyncResponse, &state, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update is not supported as app endpoint resync operation cannot be updated.
func (a *AppEndpointResync) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// From https://developer.hashicorp.com/terraform/plugin/framework/resources/update#caveats
	// If the resource does not support modification and should always be recreated on configuration value updates,
	// the Update logic can be left empty and ensure all configurable schema attributes
	// implement the resource.RequiresReplace() attribute plan modifier.
}

// Delete stops the app endpoint resync operation.
func (a *AppEndpointResync) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppEndpointResync
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()
	appServiceId := state.AppServiceId.ValueString()
	appEndpointName := state.AppEndpoint.ValueString()

	err := a.stopResync(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error stopping App Endpoint Resync",
			err.Error(),
		)
		return
	}
}

// ImportState imports a remote app endpoint resync resource.
func (a *AppEndpointResync) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

// Configure adds the provider configured api to AppEndpointResync.
func (a *AppEndpointResync) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
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

// startResync triggers an App Endpoint Resync
func (a *AppEndpointResync) startResync(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string, scopes map[string][]string) error {

	organizationUUID, projectUUID, clusterUUID, appServiceUUID, err := a.mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return err
	}

	convertedScopes := make(map[string]api.ResyncScopes)
	for name, scope := range scopes {
		convertedScopes[name] = scope
	}

	var resyncRequest api.PostAppEndpointResyncJSONRequestBody
	if scopes != nil {
		resyncRequest = api.PostAppEndpointResyncJSONRequestBody{
			Scopes: &convertedScopes,
		}
	}

	postResyncResp, err := a.ClientV2.PostAppEndpointResyncWithResponse(ctx, organizationUUID, projectUUID, clusterUUID, appServiceUUID, appEndpointName, resyncRequest)
	if err != nil {
		tflog.Debug(ctx, "error starting App Endpoint Resync", map[string]interface{}{
			"organizationId":             organizationId,
			"projectId":                  projectId,
			"clusterId":                  clusterId,
			"appServiceId":               appServiceId,
			"appEndpointName":            appEndpointName,
			"updateLoggingConfigRequest": postResyncResp,
			"err":                        err.Error(),
		})
		return err
	}

	if postResyncResp.HTTPResponse.StatusCode != http.StatusAccepted {
		return errors.New("Unexpected status while starting App Endpoint Resync: " + string(postResyncResp.Body))
	}

	return nil
}

// getResyncStatus reads the current App Endpoint Resync status
func (a *AppEndpointResync) getResyncStatus(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string) (*api.ResyncStatus, error) {

	organizationUUID, projectUUID, clusterUUID, appServiceUUID, err := a.mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, err
	}

	getResyncStatusResp, err := a.ClientV2.GetAppEndpointResyncWithResponse(ctx, organizationUUID, projectUUID, clusterUUID, appServiceUUID, appEndpointName)
	if err != nil {
		tflog.Debug(ctx, "error getting App Endpoint Resync status", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
			"err":             err.Error(),
		})
		return nil, err
	}

	if getResyncStatusResp.JSON200 == nil {
		tflog.Debug(ctx, "unexpected status getting app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
		})
		return nil, errors.New("Unexpected status while getting App Endpoint Logging Config: " + string(getResyncStatusResp.Body))
	}

	return getResyncStatusResp.JSON200, err
}

// stopResync stops a running App Endpoint Resync
func (a *AppEndpointResync) stopResync(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string) error {

	organizationUUID, projectUUID, clusterUUID, appServiceUUID, err := a.mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return err
	}

	deleteResyncResp, err := a.ClientV2.DeleteAppEndpointResyncWithResponse(ctx, organizationUUID, projectUUID, clusterUUID, appServiceUUID, appEndpointName)
	if err != nil {
		tflog.Debug(ctx, "error stopping App Endpoint Resync", map[string]interface{}{
			"organizationId":             organizationId,
			"projectId":                  projectId,
			"clusterId":                  clusterId,
			"appServiceId":               appServiceId,
			"appEndpointName":            appEndpointName,
			"updateLoggingConfigRequest": deleteResyncResp,
			"err":                        err.Error(),
		})
		return err
	}

	if deleteResyncResp.HTTPResponse.StatusCode != http.StatusAccepted {
		return errors.New("Unexpected status while starting App Endpoint Resync: " + string(deleteResyncResp.Body))
	}

	return nil
}

// mapResponseToState maps the API response to the Terraform state.
func (a *AppEndpointResync) mapResponseToState(
	ctx context.Context, response *api.ResyncStatus, plan *providerschema.AppEndpointResync, organizationId, projectId, clusterId, appServiceId, appEndpoint string,
) (*providerschema.AppEndpointResync, diag.Diagnostics) {
	state := &providerschema.AppEndpointResync{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
		AppEndpoint:    types.StringValue(appEndpoint),
		Scopes:         plan.Scopes,
		DocsChanged:    types.Int64Value(int64(response.DocsChanged)),
		DocsProcessed:  types.Int64Value(int64(response.DocsProcessed)),
		LastError:      types.StringValue(response.LastError),
		StartTime:      types.StringValue(response.StartTime.Format("2006-01-02T15:04:05Z")),
		State:          types.StringValue(string(response.State)),
	}

	if response.CollectionsProcessing != nil && len(*response.CollectionsProcessing) > 0 {
		mapValue, diags := types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, response.CollectionsProcessing)
		if diags.HasError() {
			return nil, diags
		}

		state.CollectionsProcessing = mapValue
	} else {
		state.CollectionsProcessing = types.MapNull(types.SetType{
			ElemType: types.StringType,
		})
	}

	return state, nil
}

func (a *AppEndpointResync) mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId string) (organizationUUID, projectUUID, clusterUUID, appServiceUUID uuid.UUID, err error) {
	return utils.ParseHierarchyUUIDs(organizationId, projectId, clusterId, appServiceId)
}
