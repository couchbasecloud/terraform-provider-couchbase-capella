package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpointResync{}
	_ resource.ResourceWithConfigure   = &AppEndpointResync{}
	_ resource.ResourceWithImportState = &AppEndpointResync{}
)

const errorInitiatingAppEndpointResync = "There is an error initiating app endpoint resync. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

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
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_resync"
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

	var resyncRequest *api.CreateResyncRequest
	if scopes != nil {
		resyncRequest = &api.CreateResyncRequest{
			Scopes: scopes,
		}
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/resync",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodPost,
		SuccessStatus: http.StatusAccepted,
	}

	_, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		resyncRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error initiating app endpoint resync",
			errorInitiatingAppEndpointResync+api.ParseError(err),
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
	url = fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/resync",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	cfg = api.EndpointCfg{
		Url:           url,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app endpoint resync after initiation",
			errorAfterAppEndpointResyncInitiation+api.ParseError(err),
		)
		return
	}

	var resyncResponse api.CreateResyncResponse
	if err = json.Unmarshal(response.Body, &resyncResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing app endpoint resync response",
			errorAfterAppEndpointResyncInitiation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	refreshedState, diags := a.mapResponseToState(ctx, &resyncResponse, &plan)
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

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/resync",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "Resource not found in remote server removing from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading app endpoint resync",
			"Could not read app endpoint resync, unexpected error: "+errString,
		)
		return
	}

	var resyncResponse api.CreateResyncResponse
	if err = json.Unmarshal(response.Body, &resyncResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing app endpoint resync response",
			"error during unmarshalling: "+err.Error(),
		)
		return
	}

	refreshedState, diags := a.mapResponseToState(ctx, &resyncResponse, &state)
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

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/resync",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusAccepted,
	}

	_, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "Resource not found in remote server removing from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error stopping app endpoint resync",
			"Could not stop app endpoint resync, unexpected error: "+errString,
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

// mapResponseToState maps the API response to the Terraform state.
func (a *AppEndpointResync) mapResponseToState(
	ctx context.Context, response *api.CreateResyncResponse, plan *providerschema.AppEndpointResync,
) (*providerschema.AppEndpointResync, diag.Diagnostics) {
	state := &providerschema.AppEndpointResync{
		OrganizationId: plan.OrganizationId,
		ProjectId:      plan.ProjectId,
		ClusterId:      plan.ClusterId,
		AppServiceId:   plan.AppServiceId,
		AppEndpoint:    plan.AppEndpoint,
		Scopes:         plan.Scopes,
		DocsChanged:    types.Int64Value(response.DocsChanged),
		DocsProcessed:  types.Int64Value(response.DocsProcessed),
		LastError:      types.StringValue(response.LastError),
		StartTime:      types.StringValue(response.StartTime.Format("2006-01-02T15:04:05Z")),
		State:          types.StringValue(string(response.State)),
	}

	if len(response.CollectionsProcessing) > 0 {
		mapValue, diags := types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, response.CollectionsProcessing)
		if diags.HasError() {
			return nil, diags
		}

		state.CollectionsProcessing = mapValue
	}

	return state, nil
}
