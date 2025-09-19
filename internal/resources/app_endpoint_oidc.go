package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpointOidcProvider{}
	_ resource.ResourceWithConfigure   = &AppEndpointOidcProvider{}
	_ resource.ResourceWithImportState = &AppEndpointOidcProvider{}
)

// AppEndpointOidcProvider is the Access Control Function resource implementation.
type AppEndpointOidcProvider struct {
	*providerschema.Data
}

func NewAppEndpointOidcProvider() resource.Resource {
	return &AppEndpointOidcProvider{}
}

func (r *AppEndpointOidcProvider) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_oidc_provider"
}

// Schema defines the Terraform schema for this resource.
func (r *AppEndpointOidcProvider) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointOidcProviderSchema()
}

// ImportState imports a resource into Terraform state.
func (r *AppEndpointOidcProvider) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// The import ID should be in the format: organizationId,projectId,clusterId,appServiceId,app_endpoint_name,provider_id
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

func (r *AppEndpointOidcProvider) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// Create creates a new OIDC provider and stores provider_id.
func (r *AppEndpointOidcProvider) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpointOidcProvider
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	appEndpointName := plan.AppEndpointName.ValueString()

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	payload := buildAppEndpointOIDCProviderPayload(plan)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	res, err := r.Client.ExecuteWithRetry(ctx, cfg, payload, r.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating OIDC Provider", api.ParseError(err))
		return
	}
	initOidcProviderNullsBeforeRefresh(&plan)

	// Capture providerId from response
	var created api.AppEndpointOIDCProviderResponse
	var providerId string
	if err := json.Unmarshal(res.Body, &created); err != nil {
		resp.Diagnostics.AddError("Error unmarshalling create OIDC Provider response", api.ParseError(err))
		return
	}
	if created.ProviderID != "" {
		providerId = created.ProviderID
		plan.ProviderId = types.StringValue(created.ProviderID)
	} else {
		resp.Diagnostics.AddError(
			"Error Creating App Endpoint OIDC Provider",
			"Empty provider Id returned.",
		)
		return
	}
	// Initialize optional/computed attributes to null before refresh to preserve user intent
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	// Refresh using GET; preserve nulls for optional attributes during create
	details, err := r.getOidcProvider(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, providerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing App Endpoint OIDC Provider after creation",
			fmt.Sprintf("Could not read OIDC provider %s on App Endpoint %s ", providerId, appEndpointName)+". "+api.ParseError(err),
		)
		return
	}
	r.mapResponseToState(&plan, details)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read fetches the OIDC provider and refreshes state.
func (r *AppEndpointOidcProvider) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpointOidcProvider
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading OIDC Provider", "Could not validate state: "+err.Error())
		return
	}

	providerId := IDs[providerschema.ProviderId]
	if providerId == "" {
		tflog.Info(ctx, "providerId missing; removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	details, err := r.getOidcProvider(ctx, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId], IDs[providerschema.ClusterId], IDs[providerschema.AppServiceId], IDs[providerschema.AppEndpointName], providerId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "OIDC provider not found; removing from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading OIDC Provider", "Could not read OIDC provider: "+errString)
		return
	}

	// On Read, populate all fields from remote
	r.mapResponseToState(&state, details)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates an existing OIDC provider for the App Endpoint using the Capella API.
func (r *AppEndpointOidcProvider) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppEndpointOidcProvider
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state to obtain provider ID
	var state providerschema.AppEndpointOidcProvider
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	appEndpointName := plan.AppEndpointName.ValueString()

	providerId := state.ProviderId.ValueString()

	if providerId == "" {
		resp.Diagnostics.AddError("Error Updating OIDC Provider", "provider_id is missing; cannot update")
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/%s",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
		providerId,
	)

	payload := buildAppEndpointOIDCProviderPayload(plan)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(ctx, cfg, payload, r.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating OIDC Provider", api.ParseError(err))
		return
	}

	// Refresh from server to ensure state matches remote
	details, err := r.getOidcProvider(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, providerId)
	if err != nil {
		resp.Diagnostics.AddWarning("Error reading OIDC Provider after update", api.ParseError(err))
		// Preserve providerId even if read failed
		plan.ProviderId = types.StringValue(providerId)
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		return
	}

	r.mapResponseToState(&plan, details)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the OIDC provider.
func (r *AppEndpointOidcProvider) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppEndpointOidcProvider
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()
	appServiceId := state.AppServiceId.ValueString()
	appEndpointName := state.AppEndpointName.ValueString()
	providerId := state.ProviderId.ValueString()

	if providerId == "" {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/%s",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
		providerId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := r.Client.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		resourceNotFound, _ := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError("Error Deleting OIDC Provider", api.ParseError(err))
		return
	}
}

// getOidcProvider gets OIDC provider details.
func (r *AppEndpointOidcProvider) getOidcProvider(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName, providerId string) (api.AppEndpointOIDCProviderResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/%s",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
		providerId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	res, err := r.Client.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		return api.AppEndpointOIDCProviderResponse{}, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	var out api.AppEndpointOIDCProviderResponse
	if err := json.Unmarshal(res.Body, &out); err != nil {
		return api.AppEndpointOIDCProviderResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return out, nil
}

// mapResponseToState maps response fields to state.
// If preserveNulls is true, optional attributes that are null in state will not be populated.
func (r *AppEndpointOidcProvider) mapResponseToState(state *providerschema.AppEndpointOidcProvider, resp api.AppEndpointOIDCProviderResponse) {
	// Required fields
	if resp.Issuer != "" {
		state.Issuer = types.StringValue(resp.Issuer)
	}
	if resp.ClientID != "" {
		state.ClientId = types.StringValue(resp.ClientID)
	}

	// Optional fields with null-preservation on create
	if resp.DiscoveryURL != "" && !state.DiscoveryUrl.IsNull() {
		state.DiscoveryUrl = types.StringValue(resp.DiscoveryURL)
	}
	if resp.UserPrefix != "" && !state.UserPrefix.IsNull() {
		state.UserPrefix = types.StringValue(resp.UserPrefix)
	}
	if resp.UsernameClaim != "" && !state.UsernameClaim.IsNull() {
		state.UsernameClaim = types.StringValue(resp.UsernameClaim)
	}
	if resp.RolesClaim != "" && !state.RolesClaim.IsNull() {
		state.RolesClaim = types.StringValue(resp.RolesClaim)
	}
	// Preserve null for optional bool on create if not set in config
	if !state.Register.IsNull() {
		state.Register = types.BoolValue(resp.Register)
	}

	// Computed ID
	if resp.ProviderID != "" {
		state.ProviderId = types.StringValue(resp.ProviderID)
	}
	state.IsDefault = types.BoolValue(resp.IsDefault)

}

// initOidcProviderNullsBeforeRefresh initializes computed attributes to null
// prior to the first refresh after create.
func initOidcProviderNullsBeforeRefresh(plan *providerschema.AppEndpointOidcProvider) {
	plan.IsDefault = types.BoolNull()
	plan.ProviderId = types.StringNull()
}

func buildAppEndpointOIDCProviderPayload(plan providerschema.AppEndpointOidcProvider) api.AppEndpointOIDCProviderRequest {
	payload := api.AppEndpointOIDCProviderRequest{
		Issuer:   plan.Issuer.ValueString(),
		ClientID: plan.ClientId.ValueString(),
	}
	if !plan.DiscoveryUrl.IsNull() && !plan.DiscoveryUrl.IsUnknown() {
		v := plan.DiscoveryUrl.ValueString()
		payload.DiscoveryURL = &v
	}
	if !plan.Register.IsNull() && !plan.Register.IsUnknown() {
		v := plan.Register.ValueBool()
		payload.Register = &v
	}
	if !plan.RolesClaim.IsNull() && !plan.RolesClaim.IsUnknown() {
		v := plan.RolesClaim.ValueString()
		payload.RolesClaim = &v
	}
	if !plan.UserPrefix.IsNull() && !plan.UserPrefix.IsUnknown() {
		v := plan.UserPrefix.ValueString()
		payload.UserPrefix = &v
	}
	if !plan.UsernameClaim.IsNull() && !plan.UsernameClaim.IsUnknown() {
		v := plan.UsernameClaim.ValueString()
		payload.UsernameClaim = &v
	}
	return payload
}
