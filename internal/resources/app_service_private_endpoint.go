package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*AppServicePrivateEndpoint)(nil)
	_ resource.ResourceWithConfigure   = (*AppServicePrivateEndpoint)(nil)
	_ resource.ResourceWithImportState = (*AppServicePrivateEndpoint)(nil)
)

// AppServicePrivateEndpoint is the App Service private endpoint resource implementation.
type AppServicePrivateEndpoint struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpoint is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpoint() resource.Resource {
	return &AppServicePrivateEndpoint{}
}

// Metadata returns the App Service private endpoint resource type name.
func (a *AppServicePrivateEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoint"
}

// Schema defines the schema for the App Service private endpoint resource.
func (a *AppServicePrivateEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointSchema()
}

// Configure adds the provider configured client to the App Service private endpoint resource.
func (a *AppServicePrivateEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.Data = data
}

// Create accepts a private endpoint connection request on the App Service.
func (a *AppServicePrivateEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServicePrivateEndpoint
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		endpointId     = plan.EndpointId.ValueString()
	)

	cfg := api.EndpointCfg{
		Url:           a.appServicePrivateEndpointURL(organizationId, projectId, clusterId, appServiceId, endpointId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusNoContent,
	}
	if _, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error accepting App Service private endpoint",
			"Could not accept App Service private endpoint "+endpointId+", unexpected error: "+api.ParseError(err),
		)
		return
	}

	// Save the IDs before the refresh, so that an accepted endpoint stays in
	// state and is removed (destroyed) on the next apply if the refresh fails.
	plan.ServiceName = types.StringNull()
	plan.Status = types.StringNull()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.waitUntilLinked(ctx, organizationId, projectId, clusterId, appServiceId, endpointId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for App Service private endpoint to be linked",
			"App Service private endpoint "+endpointId+" was accepted but did not become linked: "+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads the App Service private endpoint status.
func (a *AppServicePrivateEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint",
			"Could not validate App Service private endpoint "+state.EndpointId.String()+": "+err.Error(),
		)
		return
	}

	refreshedState, err := a.refreshAppServicePrivateEndpoint(
		ctx,
		IDs[providerschema.OrganizationId],
		IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId],
		IDs[providerschema.AppServiceId],
		IDs[providerschema.EndpointId],
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound || stderrors.Is(err, errors.ErrNotFound) {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint",
			"Could not read App Service private endpoint, unexpected error: "+errString,
		)
		return
	}

	// Both rejected and failed associations are terminal and cannot recover in
	// place, so remove them from state to force a clean re-association on the
	// next apply rather than leaving a stuck resource.
	switch s := apigen.PrivateEndpointStatus(refreshedState.Status.ValueString()); s {
	case apigen.PrivateEndpointStatusRejected, apigen.PrivateEndpointStatusFailed:
		tflog.Info(ctx, "App Service private endpoint association is "+string(s)+"; removing from state to force re-association")
		resp.Diagnostics.AddWarning(
			"App Service private endpoint is "+string(s),
			"App Service private endpoint "+IDs[providerschema.EndpointId]+" is "+string(s)+", so it was removed from state and the next apply will accept it again. "+
				"If it stays "+string(s)+", every plan will show a create: fix the connection request on the cloud provider side or remove the resource from the configuration.",
		)
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update is not reachable. Capella has no update endpoint for App Service
// private endpoints, and every configurable attribute requires replacement,
// so the framework deletes and recreates the resource instead.
// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
func (a *AppServicePrivateEndpoint) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete removes the private endpoint from the App Service.
func (a *AppServicePrivateEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServicePrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting App Service private endpoint",
			"Could not validate App Service private endpoint "+state.EndpointId.String()+": "+err.Error(),
		)
		return
	}

	endpointId := IDs[providerschema.EndpointId]
	cfg := api.EndpointCfg{
		Url: a.appServicePrivateEndpointURL(
			IDs[providerschema.OrganizationId],
			IDs[providerschema.ProjectId],
			IDs[providerschema.ClusterId],
			IDs[providerschema.AppServiceId],
			endpointId,
		),
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusNoContent,
	}
	_, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting App Service private endpoint",
			"Could not delete App Service private endpoint "+endpointId+", unexpected error: "+errString,
		)
	}
}

// ImportState imports an App Service private endpoint using the endpoint_id attribute.
func (a *AppServicePrivateEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_id"), req, resp)
}

func (a *AppServicePrivateEndpoint) appServicePrivateEndpointURL(organizationId, projectId, clusterId, appServiceId, endpointId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointId,
	)
}

// refreshAppServicePrivateEndpoint finds the endpoint in the list of App Service
// private endpoints, as there is no endpoint to get a single one. It returns
// errors.ErrNotFound when the endpoint is not in the list.
func (a *AppServicePrivateEndpoint) refreshAppServicePrivateEndpoint(ctx context.Context, organizationId, projectId, clusterId, appServiceId, endpointId string) (*providerschema.AppServicePrivateEndpoint, error) {
	cfg := api.EndpointCfg{
		Url: fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints",
			a.HostURL,
			organizationId,
			projectId,
			clusterId,
			appServiceId,
		),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, err
	}

	var listResp apigen.GetPrivateEndpointsResponse
	if err = json.Unmarshal(response.Body, &listResp); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	for _, e := range listResp.Endpoints {
		if e.Id != endpointId {
			continue
		}
		return &providerschema.AppServicePrivateEndpoint{
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			AppServiceId:   types.StringValue(appServiceId),
			EndpointId:     types.StringValue(endpointId),
			ServiceName:    types.StringPointerValue(e.ServiceName),
			Status:         types.StringValue(string(e.Status)),
		}, nil
	}

	return nil, errors.ErrNotFound
}

// waitUntilLinked polls the endpoint list until the accepted endpoint is
// linked. rejected and failed are terminal, so they fail fast instead of
// waiting for statusChangeTimeout. The caller has already saved the IDs to
// state, so a timeout leaves the endpoint tracked and Read picks it up later.
func (a *AppServicePrivateEndpoint) waitUntilLinked(ctx context.Context, organizationId, projectId, clusterId, appServiceId, endpointId string) (*providerschema.AppServicePrivateEndpoint, error) {
	ctx, cancel := context.WithTimeout(ctx, statusChangeTimeout)
	defer cancel()

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out after %s waiting for status linked", statusChangeTimeout)
		case <-timer.C:
			state, err := a.refreshAppServicePrivateEndpoint(ctx, organizationId, projectId, clusterId, appServiceId, endpointId)
			if err != nil {
				return nil, err
			}
			switch s := apigen.PrivateEndpointStatus(state.Status.ValueString()); s {
			case apigen.PrivateEndpointStatusLinked:
				return state, nil
			case apigen.PrivateEndpointStatusRejected, apigen.PrivateEndpointStatusFailed:
				return nil, fmt.Errorf("status is %s", s)
			}
			timer.Reset(pollInterval)
		}
	}
}
