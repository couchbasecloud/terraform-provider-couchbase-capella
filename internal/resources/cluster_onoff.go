package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	cluster_onoff_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	cluster_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"

	internal_errors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ClusterOnOffOnDemand{}
	_ resource.ResourceWithConfigure   = &ClusterOnOffOnDemand{}
	_ resource.ResourceWithImportState = &ClusterOnOffOnDemand{}
)

const errorMessageWhileClusterOnOffCreation = "There is an error during switching on/off the cluster. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterClusterOnOffCreation = "Cluster switch on/off is successful, but encountered an error while checking the current" +
	" state of the switched on/off cluster. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// ClusterOnOffOnDemand is the onDemandClusterOnOff resource implementation.
type ClusterOnOffOnDemand struct {
	*providerschema.Data
}

// NewClusterOnOffOnDemand is a helper function to simplify the provider implementation.
func NewClusterOnOffOnDemand() resource.Resource {
	return &ClusterOnOffOnDemand{}
}

// ImportState imports a remote onDemandClusterOnOff cluster that is not created by Terraform.
func (c *ClusterOnOffOnDemand) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// Metadata returns the ClusterOnOffOnDemand cluster resource type name.
func (c *ClusterOnOffOnDemand) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_onoff_ondemand"
}

// Schema defines the schema for ClusterOnOffOnDemand.
func (c *ClusterOnOffOnDemand) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterOnOffOnDemandSchema()
}

// Configure It adds the provider configured api to ClusterOnOff.
func (c *ClusterOnOffOnDemand) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	c.Data = data
}

// Create allows to switch the cluster to ON or OFF state.
func (c *ClusterOnOffOnDemand) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterOnRequest := cluster_onoff_api.CreateClusterOnRequest{}

	// Check for optional fields
	if !plan.TurnOnLinkedAppService.IsNull() && !plan.TurnOnLinkedAppService.IsUnknown() {
		clusterOnRequest.TurnOnLinkedAppService = plan.TurnOnLinkedAppService.ValueBool()
	}

	if err := c.validateClusterOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandClusterOnOff request",
			"Could not switch on the cluster, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	if err := c.manageClusterActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId, clusterOnRequest); err != nil {
		resp.Diagnostics.AddError(
			"Cluster activation failed",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeClusterOnOffWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, plan.State.ValueString(), plan.TurnOnLinkedAppService.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella ClusterOnOffOnDemand",
			"Could not read Capella ClusterOnOffOnDemand for the cluster: %s "+clusterId+"."+errorMessageAfterClusterOnOffCreation+cluster_onoff_api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ClusterOnOffOnDemand) manageClusterActivation(ctx context.Context, state, organizationId, projectId, clusterId string, onPayload any) error {
	var (
		url     = fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/activationState", c.HostURL, organizationId, projectId, clusterId)
		method  string
		payload any
	)

	switch state {
	case "on":
		method = http.MethodPost
		payload = onPayload
	case "off":
		method = http.MethodDelete
	default:
		return errors.New("invalid state value: state must be either 'on' or 'off'")
	}

	cfg := cluster_onoff_api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}
	_, err := c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		payload,
		c.Token,
		nil,
	)
	if err != nil {
		return errors.New(errorMessageWhileClusterOnOffCreation + cluster_onoff_api.ParseError(err))
	}
	return nil
}

// initializeClusterOnOffWithPlan initializes an instance of providerschema.ClusterOnOffOnDemand
// with the specified plan. It marks all computed fields as null.
func initializeClusterOnOffWithPlan(plan providerschema.ClusterOnOffOnDemand) providerschema.ClusterOnOffOnDemand {
	if plan.TurnOnLinkedAppService.IsNull() || plan.TurnOnLinkedAppService.IsUnknown() {
		plan.TurnOnLinkedAppService = types.BoolNull()
	}
	return plan
}

func (c *ClusterOnOffOnDemand) validateClusterOnOffRequest(plan providerschema.ClusterOnOffOnDemand) error {
	if plan.OrganizationId.IsNull() {
		return internal_errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return internal_errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return internal_errors.ErrClusterIdCannotBeEmpty
	}
	if plan.State.IsNull() {
		return internal_errors.ErrOnoffStateCannotBeEmpty
	}
	return nil
}

// retrieveClusterOnOff retrieves onDemandClusterOnOff information from the specified organization and project using the provided cluster ID by Get cluster open-api call.
func (c *ClusterOnOffOnDemand) retrieveClusterOnOff(ctx context.Context, organizationId, projectId, clusterId, state string, linkedApp bool) (*providerschema.ClusterOnOffOnDemand, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", c.HostURL, organizationId, projectId, clusterId)
	cfg := cluster_onoff_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	//There is no GET endpoint so get the cluster response and check current state
	clusterResp := cluster_api.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}

	if validateClusterStateIsSameInPlanAndState(state, string(clusterResp.CurrentState)) {
		clusterResp.CurrentState = cluster_api.State(state)
	}

	refreshedState := providerschema.ClusterOnOffOnDemand{
		ClusterId:              types.StringValue(clusterId),
		ProjectId:              types.StringValue(projectId),
		OrganizationId:         types.StringValue(organizationId),
		State:                  types.StringValue(state),
		TurnOnLinkedAppService: types.BoolValue(linkedApp),
	}

	return &refreshedState, nil
}

func validateClusterStateIsSameInPlanAndState(planClusterState, stateClusterState string) bool {
	return strings.EqualFold(planClusterState, stateClusterState)
}

// Couchbase Capella's v4 does not support a GET endpoint for cluster on/off.
// Cluster on/off can only access the POST and DELETE endpoint for switching the cluster to on and off state respectively.
// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
// This read is calling the retrieveClusterOnOff func to verify the state with the cluster response.
func (c *ClusterOnOffOnDemand) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ClusterOnOffOnDemand
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cluster on/off details in Capella",
			"Could not validate the cluster on/off for cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, state.State.String(), state.TurnOnLinkedAppService.ValueBool())
	if err != nil {
		resourceNotFound, _ := cluster_onoff_api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error parsing read onDemandClusterOnOff request",
			"Could not read the cluster details, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update allows to update the cluster to ON or OFF state.
func (c *ClusterOnOffOnDemand) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve values from plan
	var plan providerschema.ClusterOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterOnRequest := cluster_onoff_api.CreateClusterOnRequest{}
	// Check for optional fields
	if !plan.TurnOnLinkedAppService.IsNull() && !plan.TurnOnLinkedAppService.IsUnknown() {
		clusterOnRequest.TurnOnLinkedAppService = plan.TurnOnLinkedAppService.ValueBool()
	}

	if err := c.validateClusterOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandClusterOnOff request",
			"Could not switch on/off the cluster, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	if err := c.manageClusterActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId, clusterOnRequest); err != nil {
		resp.Diagnostics.AddError(
			"Cluster activation failed",
			err.Error(),
		)
		return
	}

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, plan.State.ValueString(), plan.TurnOnLinkedAppService.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella ClusterOnOffOnDemand",
			"Could not read Capella ClusterOnOffOnDemand for the cluster: %s "+clusterId+"."+errorMessageAfterClusterOnOffCreation+cluster_onoff_api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ClusterOnOffOnDemand) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support a DELETION/destroying resource for cluster on/off.
	// Cluster on/off can only access the POST and DELETE endpoint which are used for switching the cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
}
