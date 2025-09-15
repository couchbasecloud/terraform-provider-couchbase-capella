package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpointDefaultOidcProvider{}
	_ resource.ResourceWithConfigure   = &AppEndpointDefaultOidcProvider{}
	_ resource.ResourceWithImportState = &AppEndpointDefaultOidcProvider{}
)

// AppEndpointDefaultOidcProvider manages the default OIDC provider selection for an App Endpoint.
type AppEndpointDefaultOidcProvider struct {
	*providerschema.Data
}

func NewAppEndpointDefaultOidcProvider() resource.Resource {
	return &AppEndpointDefaultOidcProvider{}
}

func (r *AppEndpointDefaultOidcProvider) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_default_oidc_provider"
}

func (r *AppEndpointDefaultOidcProvider) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointDefaultOidcProviderSchema()
}

func (r *AppEndpointDefaultOidcProvider) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports a resource into Terraform state.
func (r *AppEndpointDefaultOidcProvider) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is the app endpoint name to align with similar resources.
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

// Create sets the default OIDC provider for the App Endpoint.
func (r *AppEndpointDefaultOidcProvider) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpointDefaultOidcProvider
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
	providerId := plan.ProviderId.ValueString()

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/defaultProvider",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	payload := api.AppEndpointOIDCDefaultProviderRequest{ProviderID: providerId}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(ctx, cfg, payload, r.Token, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		resp.Diagnostics.AddError("Error Setting Default OIDC Provider", api.ParseError(err))
		return
	}

	// Refresh state by listing providers and picking the default
	selected, err := r.getDefaultProvider(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddWarning("Error reading default OIDC provider after create", api.ParseError(err))
	} else if selected.ProviderID != "" {
		plan.ProviderId = types.StringValue(selected.ProviderID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read finds the current default provider by listing available providers.
func (r *AppEndpointDefaultOidcProvider) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpointDefaultOidcProvider
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Default OIDC Provider", "Could not validate state: "+err.Error())
		return
	}

	selected, err := r.getDefaultProvider(ctx, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId], IDs[providerschema.ClusterId], IDs[providerschema.AppServiceId], IDs[providerschema.AppEndpointName])
	if err != nil {
		// If list returns 404, remove from state
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Default OIDC Provider", "Could not list OIDC providers: "+errString)
		return
	}

	if selected.ProviderID == "" {
		// No default currently set; remove from state so plan can set it
		tflog.Info(ctx, "No default OIDC provider found; removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	state.ProviderId = types.StringValue(selected.ProviderID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update sets the default provider (idempotent with Create).
func (r *AppEndpointDefaultOidcProvider) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppEndpointDefaultOidcProvider
	var state providerschema.AppEndpointDefaultOidcProvider
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	appEndpointName := plan.AppEndpointName.ValueString()
	providerId := plan.ProviderId.ValueString()
	if providerId == "" {
		// fall back to state if plan doesn't include it
		providerId = state.ProviderId.ValueString()
	}
	if providerId == "" {
		resp.Diagnostics.AddError("Error Updating Default OIDC Provider", "provider_id is required to set default provider")
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/defaultProvider",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)

	payload := api.AppEndpointOIDCDefaultProviderRequest{ProviderID: providerId}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(ctx, cfg, payload, r.Token, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Default OIDC Provider", api.ParseError(err))
		return
	}

	selected, err := r.getDefaultProvider(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddWarning("Error reading default OIDC provider after update", api.ParseError(err))
		plan.ProviderId = types.StringValue(providerId)
	} else if selected.ProviderID != "" {
		plan.ProviderId = types.StringValue(selected.ProviderID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete is a no-op as the API does not support unsetting default; remove from state.
func (r *AppEndpointDefaultOidcProvider) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppEndpointDefaultOidcProvider
	_ = req.State.Get(ctx, &state)
	// Best-effort: there is no API to clear default; just remove from state
	tflog.Info(ctx, "Deleting default OIDC provider is not supported by API; removing from state")
}

// getDefaultProvider lists providers and returns the one marked as default.
func (r *AppEndpointDefaultOidcProvider) getDefaultProvider(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string) (api.AppEndpointOIDCProviderResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	res, err := r.Client.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		return api.AppEndpointOIDCProviderResponse{}, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	var list api.AppEndpointOIDCProviderListResponse
	if err := json.Unmarshal(res.Body, &list); err != nil {
		return api.AppEndpointOIDCProviderResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}
	for _, p := range list.Data {
		if p.IsDefault {
			return p, nil
		}
	}
	return api.AppEndpointOIDCProviderResponse{}, nil
}
