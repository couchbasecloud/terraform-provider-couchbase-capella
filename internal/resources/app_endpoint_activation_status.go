package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpointActivationStatus{}
	_ resource.ResourceWithConfigure   = &AppEndpointActivationStatus{}
	_ resource.ResourceWithImportState = &AppEndpointActivationStatus{}
)

const errorMessageWhileAppEndpointActivation = "There is an error during switching online/offline the app endpoint. Please check in Capella to see if any hanging resources have been created, unexpected error: "

// AppEndpointActivationStatus manages activation status (online/offline) of an App Endpoint.
type AppEndpointActivationStatus struct {
	*providerschema.Data
}

// NewAppEndpointActivationStatus is a helper function to simplify the provider implementation.
func NewAppEndpointActivationStatus() resource.Resource {
	return &AppEndpointActivationStatus{}
}

// ImportState imports a remote App Endpoint activation status that is not created by Terraform.
func (r *AppEndpointActivationStatus) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save import name to app_endpoint_name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Metadata returns the resource type name.
func (r *AppEndpointActivationStatus) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_activation_status"
}

// Schema defines the schema for App Endpoint Activation Status.
func (r *AppEndpointActivationStatus) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointActivationStatusSchema()
}

// Configure sets provider-defined data, clients, etc.
func (r *AppEndpointActivationStatus) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.Data = data
}

// Create switches the app endpoint to online/offline state.
func (r *AppEndpointActivationStatus) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpointActivationStatus
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.validate(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create app endpoint activation status request",
			"Could not switch app endpoint online/offline, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
		online          = plan.Online.ValueBool()
	)

	if err := r.manageAppEndpointActivation(ctx, online, organizationId, projectId, clusterId, appServiceId, appEndpointName); err != nil {
		resp.Diagnostics.AddError(
			"App Endpoint activation failed",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.retrieveAppEndpointActivation(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, online)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Endpoint Activation Status",
			fmt.Sprintf("Could not read activation status for the app endpoint: %s associated to app service: %s: %s", appEndpointName, appServiceId, api.ParseError(err)),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AppEndpointActivationStatus) validate(plan providerschema.AppEndpointActivationStatus) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	if plan.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdCannotBeEmpty
	}
	if plan.AppEndpointName.IsNull() {
		return errors.ErrEndpointIdMissing
	}
	if plan.Online.IsNull() {
		return errors.ErrOnoffStateCannotBeEmpty
	}
	return nil
}

func (r *AppEndpointActivationStatus) manageAppEndpointActivation(ctx context.Context, online bool, organizationId, projectId, clusterId, appServiceId, appEndpointName string) error {
	var (
		url    = fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/activationStatus", r.HostURL, organizationId, projectId, clusterId, appServiceId, appEndpointName)
		method string
	)

	if online {
		method = http.MethodPost
	} else {
		method = http.MethodDelete
	}

	cfg := api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return fmt.Errorf(errorMessageWhileAppEndpointActivation + api.ParseError(err))
	}
	return nil
}

// retrieveAppEndpointActivation reads the app endpoint and maps its state to the Online flag.
func (r *AppEndpointActivationStatus) retrieveAppEndpointActivation(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string, desiredOnline bool) (*providerschema.AppEndpointActivationStatus, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s", r.HostURL, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Minimal struct to pick state or online
	var getResp struct {
		State  string `json:"state"`
		Online *bool  `json:"online"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(response.Body, &getResp); err != nil {
		return nil, err
	}

	online := desiredOnline
	switch {
	case getResp.Online != nil:
		online = *getResp.Online
	case getResp.State != "":
		online = strings.EqualFold(getResp.State, "online")
	case getResp.Status != "":
		online = strings.EqualFold(getResp.Status, "online")
	}

	refreshed := providerschema.AppEndpointActivationStatus{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		Online:          types.BoolValue(online),
	}
	return &refreshed, nil
}

// Read verifies the status via GET App Endpoint API and updates state.
// Couchbase Capella's v4 does not support a GET endpoint for activation state directly.
// This read is calling the retrieveAppEndpointActivation func to verify the state with the app endpoint response.
// API reference: Post/Del activation status and Get App Endpoint
// - POST: App Endpoint online
// - DELETE: App Endpoint offline
// - GET: App Endpoint details (contains state)
// See: https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/postAppEndpointActivationStatus
// See: https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/deleteAppEndpointActivationStatus
// See: https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/getAppEndpoint
func (r *AppEndpointActivationStatus) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpointActivationStatus
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading app endpoint activation details in Capella",
			"Could not validate the app endpoint activation for app endpoint "+state.AppEndpointName.String()+": "+err.Error(),
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

	refreshedState, err := r.retrieveAppEndpointActivation(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, state.Online.ValueBool())
	if err != nil {
		resourceNotFound, _ := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error parsing read app endpoint activation request",
			"Could not read the app endpoint details, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update switches the app endpoint to online/offline state.
func (r *AppEndpointActivationStatus) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppEndpointActivationStatus
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.validate(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing update app endpoint activation status request",
			"Could not switch app endpoint online/offline, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
		online          = plan.Online.ValueBool()
	)

	if err := r.manageAppEndpointActivation(ctx, online, organizationId, projectId, clusterId, appServiceId, appEndpointName); err != nil {
		resp.Diagnostics.AddError(
			"App Endpoint activation failed",
			err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveAppEndpointActivation(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, online)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Endpoint Activation Status",
			api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete is a no-op as the activation resource models an action.
func (r *AppEndpointActivationStatus) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support destroying an activation resource.
	// The POST and DELETE endpoints are used to switch the app endpoint online and offline respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/postAppEndpointActivationStatus
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/deleteAppEndpointActivationStatus
}
