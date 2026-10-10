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
	_ resource.Resource                = &AppServicePrivateEndpoint{}
	_ resource.ResourceWithConfigure   = &AppServicePrivateEndpoint{}
	_ resource.ResourceWithImportState = &AppServicePrivateEndpoint{}
)

// AppServicePrivateEndpoint is the App Service scoped private endpoint resource
// implementation.
type AppServicePrivateEndpoint struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpoint is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpoint() resource.Resource {
	return &AppServicePrivateEndpoint{}
}

// Metadata returns the App Service private endpoint resource type name.
func (p *AppServicePrivateEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoints"
}

// Schema defines the schema for the App Service private endpoint resource.
func (p *AppServicePrivateEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointsSchema()
}

// Create accepts a private endpoint connection on an App Service's private endpoint service.
func (p *AppServicePrivateEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServicePrivateEndpoint
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateAcceptAppServicePrivateEndpoint(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service private endpoint accept request",
			"Could not validate App Service private endpoint accept request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		endpointId     = plan.EndpointId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints/%s",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusNoContent}

	_, err = p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error accepting App Service private endpoint",
			"Could not accept App Service private endpoint "+endpointId+", unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeAppServicePrivateEndpointPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := p.getPrivateEndpointState(ctx, organizationId, projectId, clusterId, appServiceId, endpointId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint status",
			"Error reading App Service private endpoint status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the App Service private endpoint status.
func (p *AppServicePrivateEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoint",
			"Could not validate App Service private endpoint "+state.EndpointId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
		endpointId     = IDs[providerschema.EndpointId]
	)

	refreshedState, err := p.getPrivateEndpointState(ctx, organizationId, projectId, clusterId, appServiceId, endpointId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			resp.State.RemoveResource(ctx)
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint status",
			"Error reading App Service private endpoint status, unexpected error: "+errString,
		)

		return
	}

	// Both rejected and failed associations are terminal and cannot recover in
	// place, so remove them from state to force a clean re-association on the
	// next apply rather than leaving a stuck resource.
	if s := refreshedState.Status.ValueString(); s == "rejected" || s == "failed" {
		tflog.Info(ctx, "App Service private endpoint association is "+s+"; removing from state to force re-association")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update there is no update API so returns an error.
func (p *AppServicePrivateEndpoint) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// From https://developer.hashicorp.com/terraform/plugin/framework/resources/update#caveats
	// If the resource does not support modification and should always be recreated on configuration value updates,
	// the Update logic can be left empty and ensure all configurable schema attributes
	// implement the resource.RequiresReplace() attribute plan modifier.
}

// Delete rejects a private endpoint connection on an App Service's private endpoint service.
func (p *AppServicePrivateEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServicePrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rejecting App Service private endpoint",
			"Could not reject endpoint due to validation error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
		endpointId     = IDs[providerschema.EndpointId]
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints/%s",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rejecting App Service private endpoint",
			"Could not disable App Service private endpoint for endpoint "+endpointId+" unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint resource.
func (p *AppServicePrivateEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	p.Data = data
}

// ImportState imports an App Service private endpoint to be managed by terraform.
func (p *AppServicePrivateEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_id"), req, resp)
}

// validateAcceptAppServicePrivateEndpoint ensures organization id, project id, cluster id,
// app service id, and endpoint id are valued.
func validateAcceptAppServicePrivateEndpoint(plan providerschema.AppServicePrivateEndpoint) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}
	if plan.EndpointId.IsNull() {
		return errors.ErrEndpointIdMissing
	}

	return nil
}

// initializeAppServicePrivateEndpointPlan initializes an instance of
// providerschema.AppServicePrivateEndpoint with the specified plan. It marks all computed
// fields as null.
func initializeAppServicePrivateEndpointPlan(plan providerschema.AppServicePrivateEndpoint) providerschema.AppServicePrivateEndpoint {
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.StringNull()
	}
	if plan.ServiceName.IsNull() || plan.ServiceName.IsUnknown() {
		plan.ServiceName = types.StringNull()
	}
	return plan
}

// getPrivateEndpointState morphs the App Service private endpoint status to a terraform schema.
func (p *AppServicePrivateEndpoint) getPrivateEndpointState(ctx context.Context, organizationId, projectId, clusterId, appServiceId, endpointId string) (*providerschema.AppServicePrivateEndpoint, error) {
	status, serviceName, err := p.getPrivateEndpointStatus(ctx, organizationId, projectId, clusterId, appServiceId, endpointId)
	if err != nil {
		return nil, err
	}

	state := providerschema.AppServicePrivateEndpoint{
		EndpointId:     types.StringValue(endpointId),
		Status:         types.StringValue(status),
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
		AppServiceId:   types.StringValue(appServiceId),
		ServiceName:    types.StringValue(serviceName),
	}

	return &state, nil
}

// There is currently no V4 endpoint to get a single App Service private endpoint. We
// have to loop through the entire list to find the desired private endpoint.
func (p *AppServicePrivateEndpoint) getPrivateEndpointStatus(ctx context.Context, organizationId, projectId, clusterId, appServiceId, endpointId string) (string, string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints", p.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		return "", "", err
	}

	privateEndpointsResp := api.GetAppServicePrivateEndpointsResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointsResp)
	if err != nil {
		return "", "", err
	}

	for _, e := range privateEndpointsResp.Endpoints {
		if e.Id == endpointId {
			return e.Status, e.ServiceName, nil
		}
	}

	return "", "", errors.ErrNotFound
}
