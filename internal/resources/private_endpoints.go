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

// PrivateEndpoint is the private endpoint resource implementation.
type PrivateEndpoint struct {
	*providerschema.Data
}

// NewPrivateEndpoint is a helper function to simplify the provider implementation.
func NewPrivateEndpoint() resource.Resource {
	return &PrivateEndpoint{}
}

// Metadata returns the private endpoint resource type name.
func (p *PrivateEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoints"
}

// Schema defines the schema for the private endpoint resource.
func (p *PrivateEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource allows you to manage private endpoints for an operational cluster. Private endpoints allow you to securely connect your Cloud Service Provider's private network (VPC/VNET) to your operational cluster without exposing traffic to the public internet.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(
				stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the organization where the private endpoint will be created.",
			),
			"project_id": WithDescription(
				stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the project containing the cluster where the private endpoint will be created.",
			),
			"cluster_id": WithDescription(
				stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the operational cluster to create the private endpoint for. This enables secure access to the cluster through your Cloud Service Provider's private network.",
			),
			"endpoint_id": WithDescription(
				stringAttribute([]string{required, requiresReplace}),
				"The ID of the private endpoint in your cloud provider.",
			),
			"status": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The current status of the private endpoint. Possible values are:\n" +
					"* `pending` - The endpoint creation is in progress\n" +
					"* `pendingAcceptance` - The endpoint is waiting for acceptance from Capella\n" +
					"* `linked` - The endpoint is successfully connected and active\n" +
					"* `rejected` - The endpoint connection request was rejected\n" +
					"* `unrecognized` - The endpoint state cannot be determined\n" +
					"* `failed` - The endpoint creation or connection attempt failed",
			},
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

// Read reads the private endpoint status.
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
	// From https://developer.hashicorp.com/terraform/plugin/framework/resources/update#caveats
	// If the resource does not support modification and should always be recreated on configuration value updates,
	// the Update logic can be left empty and ensure all configurable schema attributes
	// implement the resource.RequiresReplace() attribute plan modifier.
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

// Configure adds the provider configured client to the private endpoint resource.
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

// ImportState imports a private endpoint to be managed by terraform.
func (p *PrivateEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_id"), req, resp)
}

// validateAcceptPrivateEndpoint ensures organization id, project id, cluster id, and endpoint id are valued.
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

// getPrivateEndpointState morphs private endpoint status to terraform schema.
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

// There is currently no V4 endpoint to get a single private endpoint.  We have to loop through the entire list to find
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
