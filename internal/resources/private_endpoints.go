package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PrivateEndpoint{}
	_ resource.ResourceWithConfigure   = &PrivateEndpoint{}
	_ resource.ResourceWithImportState = &PrivateEndpoint{}
)

type PrivateEndpoint struct {
	*providerschema.Data
}

func NewPrivateEndpoint() resource.Resource {
	return &PrivateEndpoint{}
}

func (p *PrivateEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoints"
}

func (p *PrivateEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required}),
			"project_id":      stringAttribute([]string{required}),
			"cluster_id":      stringAttribute([]string{required}),
			"endpoint_id":     stringAttribute([]string{required}),
			"status":          stringAttribute([]string{computed}),
		},
	}
}

// Create accepts a private endpoint on the CSP.
func (p *PrivateEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.PrivateEndpoint
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateAcceptPrivateEndpoint(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating private endpoint accept request",
			"Could not validate private endpoint accept request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		endpointId     = plan.EndpointId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpoints/%s/associate",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		endpointId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusNoContent}

	_, err = p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error accepting private endpoint",
			"Could not accept private endpoint "+endpointId+", unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializePrivateEndpointPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := p.getPrivateEndpointState(ctx, organizationId, projectId, clusterId, endpointId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private endpoint service status",
			"Error reading private endpoint service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *PrivateEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.PrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Private Endpoint",
			"Could not validate private endpoint "+state.EndpointId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		endpointId     = IDs[providerschema.EndpointId]
	)

	refreshedState, err := p.getPrivateEndpointState(ctx, organizationId, projectId, clusterId, endpointId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private endpoint status",
			"Error reading private endpoint status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update there is no update API so returns an error.
func (p *PrivateEndpoint) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"No update API for a private endpoint",
		"No update API for a private endpoint",
	)
	return
}

// Delete rejects a private endpoint on the CSP.
func (p *PrivateEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.PrivateEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rejecting private endpoint",
			"Could not reject endpoint due to validation error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		endpointId     = IDs[providerschema.EndpointId]
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpoints/%s/unassociate",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		endpointId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusNoContent}
	_, err = p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rejecting private endpoint",
			"Could not disable private endpoint service for endpoint "+endpointId+" unexpected error: "+err.Error(),
		)
		return
	}
}

func (p *PrivateEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (p *PrivateEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_id"), req, resp)
}

func validateAcceptPrivateEndpoint(plan providerschema.PrivateEndpoint) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.EndpointId.IsNull() {
		return errors.ErrEndpointIdMissing
	}

	return nil
}

// initializePrivateEndpointPlan initializes an instance of providerschema.PrivateEndpoint
// with the specified plan. It marks all computed fields as null.
func initializePrivateEndpointPlan(plan providerschema.PrivateEndpoint) providerschema.PrivateEndpoint {
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.StringNull()
	}
	return plan
}

func (p *PrivateEndpoint) getPrivateEndpointState(ctx context.Context, organizationId, projectId, clusterId, endpointId string) (*providerschema.PrivateEndpoint, error) {
	status, err := p.getPrivateEndpointStatus(ctx, organizationId, projectId, clusterId, endpointId)
	if err != nil {
		return nil, err
	}

	state := providerschema.PrivateEndpoint{
		EndpointId:     types.StringValue(endpointId),
		Status:         types.StringValue(status),
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
	}

	return &state, nil
}

// There is no V4 endpoint to get a single private endpoint.  We have to loop through the entire list to find
// the desired private endpoint.
func (p *PrivateEndpoint) getPrivateEndpointStatus(ctx context.Context, organizationId, projectId, clusterId, endpointId string) (string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpoints", p.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		return "", err
	}

	privateEndpointsResp := api.GetPrivateEndpointsResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointsResp)
	if err != nil {
		return "", err
	}

	for _, e := range privateEndpointsResp.Endpoints {
		if e.Id == endpointId {
			return e.Status, nil
		}
	}

	return "", errors.ErrNotFound
}
