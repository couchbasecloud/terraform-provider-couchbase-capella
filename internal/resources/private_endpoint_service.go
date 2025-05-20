package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	custommodifier "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PrivateEndpointService{}
	_ resource.ResourceWithConfigure   = &PrivateEndpointService{}
	_ resource.ResourceWithImportState = &PrivateEndpointService{}
)

const (
	errorMessageWhileEnablingPrivateEndpointService = "There is an error while enabling private endpoint service. Please check in Capella to see if there are any hanging resources\" +\n\t\" have been created, unexpected error: "
)

// PrivateEndpointService is the scope resource implementation.
type PrivateEndpointService struct {
	*providerschema.Data
}

// NewPrivateEndpointService is a helper function to simplify the provider implementation.
func NewPrivateEndpointService() resource.Resource {
	return &PrivateEndpointService{}
}

// Metadata returns the private endpoint service resource type name.
func (p *PrivateEndpointService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoint_service"
}

// Schema defines the schema for a private endpoint service resource.
func (p *PrivateEndpointService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource allows you to manage the private endpoint service for an operational cluster. The private endpoint service must be enabled before you can create private endpoints to connect your Cloud Service Provider's private network (VPC/VNET) to your operational cluster. This enables secure access to your cluster without exposing traffic to the public internet.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the organization where the private endpoint service will be enabled. This field cannot be changed after the private endpoint service is created.",
			),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the project containing the cluster where the private endpoint service will be enabled. This field cannot be changed after the private endpoint service is created.",
			),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the cluster where the private endpoint service will be enabled. This enables secure access to the cluster through your Cloud Service Provider's private network. This field cannot be changed after the private endpoint service is created.",
			),
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable or disable the private endpoint service for the cluster. When enabled, you can create private endpoints to connect your Cloud Service Provider's private network to the cluster. Note: Setting this to false during creation will result in an error as the service must be enabled to be managed.",
				PlanModifiers:       []planmodifier.Bool{custommodifier.BlockCreateWhenEnabledSetToFalse()},
			},
		},
	}
}

// Create enables private endpoint service.
func (p *PrivateEndpointService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.PrivateEndpointService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateCreateEndpointService(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating private endpoint service request",
			"Could not validate private endpoint service request, unexpected error: "+err.Error(),
		)
		return
	}
	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err = p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error enabling private endpoint service",
			errorMessageWhileEnablingPrivateEndpointService+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializePrivateEndpointServicePlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = p.waitUntilStatusChanges(ctx, true, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error could not enable private endpoint service",
			"Error could not enable private endpoint service, unexpected error: "+err.Error(),
		)
	}

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId)
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

// Read reads the private endpoint service status.
func (p *PrivateEndpointService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.PrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Private Endpoint Service in Capella",
			"Could not read Capella private endpoint service on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading private endpoint service status",
			"Error reading private endpoint service status, unexpected error: "+errString,
		)

		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update will enable/disable the private endpoint service.
func (p *PrivateEndpointService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config providerschema.PrivateEndpointService
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	status := "enabling"

	if !config.Enabled.ValueBool() {
		cfg.Method = http.MethodDelete
		status = "disabling"
	}

	_, err := p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error "+status+" private endpoint service",
			"Error "+status+" private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx,
		config.Enabled.ValueBool(),
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error "+status+" private endpoint service",
			"Error "+status+"private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := p.getServiceState(ctx,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString())
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

// Delete disables private endpoint service on the cluster.
func (p *PrivateEndpointService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.PrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Private Endpoint Service in Capella",
			"Could not validate Capella private endpoint service on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	// If private endpoint service is already disabled, just remove the resource from the state file.
	if !state.Enabled.ValueBool() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error disabling private endpoint service",
			"Could not disable private endpoint service for cluster "+clusterId+" unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx, false, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error could not disable private endpoint service",
			"Error could not disable private endpoint service, unexpected error: "+err.Error(),
		)
	}

}

// Configure It adds the provider configured api to the private endpoint service resource.
func (p *PrivateEndpointService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports a private endpoint service status.
func (p *PrivateEndpointService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// validateCreateEndpointService ensures organization id, project id and cluster id are valued.
func validateCreateEndpointService(plan providerschema.PrivateEndpointService) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}

	return nil
}

// initializePlan initializes an instance of providerschema.PrivateEndpointService
// with the specified plan. It marks all computed fields as null.
func initializePrivateEndpointServicePlan(plan providerschema.PrivateEndpointService) providerschema.PrivateEndpointService {
	if plan.Enabled.IsNull() || plan.Enabled.IsUnknown() {
		plan.Enabled = types.BoolNull()
	}
	return plan
}

// waitUntilStatusChanges terraform will wait until the service status changes on the cluster.
func (p *PrivateEndpointService) waitUntilStatusChanges(ctx context.Context, finalState bool, organizationId, projectId, clusterId string) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Minute*60)
	defer cancel()

	timer := time.NewTimer(time.Minute * 1)

	for {
		select {
		case <-ctx.Done():
			return errors.ErrPrivateEndpointServiceTimeout

		case <-timer.C:
			response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId)
			if err != nil {
				return err
			}

			if response.Enabled == finalState {
				return nil
			}
			timer.Reset(time.Minute * 1)
		}
	}
}

// getServiceStatus retrieves current private endpoint service status.
func (p *PrivateEndpointService) getServiceStatus(ctx context.Context, organizationId, projectId, clusterId string) (*api.GetPrivateEndpointServiceStatusResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService", p.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	status := api.GetPrivateEndpointServiceStatusResponse{}
	err = json.Unmarshal(response.Body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// getServiceState morphs service status into a terraform schema.
func (p *PrivateEndpointService) getServiceState(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.PrivateEndpointService, error) {
	response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	state := providerschema.PrivateEndpointService{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Enabled:        types.BoolValue(response.Enabled),
	}

	return &state, nil
}
