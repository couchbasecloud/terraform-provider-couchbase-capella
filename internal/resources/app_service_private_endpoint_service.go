package resources

import (
	"context"
	"encoding/json"
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
	_ resource.Resource                = (*AppServicePrivateEndpointService)(nil)
	_ resource.ResourceWithConfigure   = (*AppServicePrivateEndpointService)(nil)
	_ resource.ResourceWithImportState = (*AppServicePrivateEndpointService)(nil)
)

const errorMessageAfterAppServicePrivateEndpointServiceEnable = "Enabling the App Service private endpoint service is initiated, but encountered an error while checking the current" +
	" state. Please run `terraform plan` after a few minutes to know the current state." +
	" Additionally, run `terraform apply --refresh-only` to update the state from remote, unexpected error: "

// AppServicePrivateEndpointService is the App Service private endpoint service resource implementation.
type AppServicePrivateEndpointService struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpointService is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpointService() resource.Resource {
	return &AppServicePrivateEndpointService{}
}

// Metadata returns the App Service private endpoint service resource type name.
func (a *AppServicePrivateEndpointService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoint_service"
}

// Schema defines the schema for the App Service private endpoint service resource.
func (a *AppServicePrivateEndpointService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointServiceSchema()
}

// Configure adds the provider configured client to the App Service private endpoint service resource.
func (a *AppServicePrivateEndpointService) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create enables the private endpoint service on the App Service and waits until it is enabled.
func (a *AppServicePrivateEndpointService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServicePrivateEndpointService
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
	)

	cfg := api.EndpointCfg{
		Url:           a.appServicePrivateEndpointServiceURL(organizationId, projectId, clusterId, appServiceId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusAccepted,
	}
	if _, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error enabling App Service private endpoint service",
			"Could not enable App Service private endpoint service, unexpected error: "+api.ParseError(err),
		)
		return
	}

	// Save the IDs before the wait, so that a failed enable stays in state and is
	// disabled (destroyed) on the next apply.
	plan.State = types.StringNull()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.checkAppServicePrivateEndpointServiceStatus(
		ctx, organizationId, projectId, clusterId, appServiceId, apigen.GetAppServicePrivateEndpointStateResponseStateEnabled,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error enabling App Service private endpoint service",
			errorMessageAfterAppServicePrivateEndpointServiceEnable+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads the state of the private endpoint service on the App Service.
func (a *AppServicePrivateEndpointService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint service",
			"Could not validate App Service private endpoint service state: "+err.Error(),
		)
		return
	}

	refreshedState, err := a.refreshAppServicePrivateEndpointService(
		ctx,
		IDs[providerschema.OrganizationId],
		IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId],
		IDs[providerschema.AppServiceId],
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint service",
			"Could not read App Service private endpoint service, unexpected error: "+errString,
		)
		return
	}

	if refreshedState.State.ValueString() == string(apigen.GetAppServicePrivateEndpointStateResponseStateDisabled) {
		tflog.Info(ctx, "App Service private endpoint service is disabled, removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update is not reachable. Capella has no update endpoint for the App Service
// private endpoint service, and every configurable attribute requires replacement,
// so the framework deletes and recreates the resource instead.
// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
func (a *AppServicePrivateEndpointService) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete disables the private endpoint service on the App Service and waits until it is disabled.
func (a *AppServicePrivateEndpointService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServicePrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error disabling App Service private endpoint service",
			"Could not validate App Service private endpoint service state: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
	)

	cfg := api.EndpointCfg{
		Url:           a.appServicePrivateEndpointServiceURL(organizationId, projectId, clusterId, appServiceId),
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusAccepted,
	}
	_, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
		resp.Diagnostics.AddError(
			"Error disabling App Service private endpoint service",
			"Could not disable App Service private endpoint service, unexpected error: "+errString,
		)
		return
	}

	_, err = a.checkAppServicePrivateEndpointServiceStatus(
		ctx, organizationId, projectId, clusterId, appServiceId, apigen.GetAppServicePrivateEndpointStateResponseStateDisabled,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error disabling App Service private endpoint service",
			"Could not disable App Service private endpoint service, unexpected error: "+errString,
		)
	}
}

// ImportState imports the App Service private endpoint service using the app_service_id attribute.
func (a *AppServicePrivateEndpointService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

func (a *AppServicePrivateEndpointService) appServicePrivateEndpointServiceURL(organizationId, projectId, clusterId, appServiceId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
}

func (a *AppServicePrivateEndpointService) refreshAppServicePrivateEndpointService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.AppServicePrivateEndpointService, error) {
	cfg := api.EndpointCfg{
		Url:           a.appServicePrivateEndpointServiceURL(organizationId, projectId, clusterId, appServiceId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, err
	}

	var stateResp apigen.GetAppServicePrivateEndpointStateResponse
	if err = json.Unmarshal(response.Body, &stateResp); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	refreshedState := providerschema.AppServicePrivateEndpointService{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
		State:          types.StringNull(),
	}
	if stateResp.State != nil {
		refreshedState.State = types.StringValue(string(*stateResp.State))
	}

	return &refreshedState, nil
}

// checkAppServicePrivateEndpointServiceStatus polls the service until it reaches
// targetState. It fails when the service reports failed or the App Service is not
// found, and retries on other errors until statusChangeTimeout.
func (a *AppServicePrivateEndpointService) checkAppServicePrivateEndpointServiceStatus(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	targetState apigen.GetAppServicePrivateEndpointStateResponseState,
) (*providerschema.AppServicePrivateEndpointService, error) {
	ctx, cancel := context.WithTimeout(ctx, statusChangeTimeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("app service private endpoint service did not reach state %s: %w", targetState, ctx.Err())
		case <-ticker.C:
			refreshedState, err := a.refreshAppServicePrivateEndpointService(ctx, organizationId, projectId, clusterId, appServiceId)
			if err != nil {
				if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
					return nil, err
				}
				tflog.Info(ctx, "retrying after error polling App Service private endpoint service state", map[string]interface{}{"error": err.Error()})
				continue
			}

			switch apigen.GetAppServicePrivateEndpointStateResponseState(refreshedState.State.ValueString()) {
			case targetState:
				return refreshedState, nil
			case apigen.GetAppServicePrivateEndpointStateResponseStateFailed:
				return nil, fmt.Errorf("app service private endpoint service reached state %s while waiting for state %s",
					apigen.GetAppServicePrivateEndpointStateResponseStateFailed, targetState)
			}

			tflog.Info(ctx, "waiting for App Service private endpoint service to reach state "+string(targetState))
		}
	}
}
