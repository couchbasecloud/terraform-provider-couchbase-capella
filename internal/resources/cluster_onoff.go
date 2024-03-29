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

	cluster_onoff_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	cluster_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ClusterOnOffOnDemand{}
	_ resource.ResourceWithConfigure   = &ClusterOnOffOnDemand{}
	_ resource.ResourceWithImportState = &ClusterOnOffOnDemand{}
)

const errorMessageWhileClusterOnOffCreation = "There is an error during switching on the cluster. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterClusterOnOffCreation = "Cluster switch on is successful, but encountered an error while checking the current" +
	" state of the switched on cluster. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// ClusterOnOffOnDemand is the onDemandClusterOnOff resource implementation.
type ClusterOnOffOnDemand struct {
	*providerschema.Data
}

// NewOnDemandClusterOnOff is a helper function to simplify the provider implementation.
func NewOnDemandClusterOnOff() resource.Resource {
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

// Create allows to switch the cluster to ON state.
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

	if err := c.validateCreateClusterOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandClusterOnOff request",
			"Could not switch on the cluster, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/activationState", c.HostURL, organizationId, projectId, clusterId)
	cfg := cluster_onoff_api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		clusterOnRequest,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing the operation to switch on the cluster",
			errorMessageWhileClusterOnOffCreation+cluster_onoff_api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeClusterOnOffWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, plan.State.ValueString())
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

// initializeClusterOnOffWithPlan initializes an instance of providerschema.ClusterOnOffOnDemand
// with the specified plan. It marks all computed fields as null.
func initializeClusterOnOffWithPlan(plan providerschema.ClusterOnOffOnDemand) providerschema.ClusterOnOffOnDemand {
	if plan.TurnOnLinkedAppService.IsNull() || plan.TurnOnLinkedAppService.IsUnknown() {
		plan.TurnOnLinkedAppService = types.BoolNull()
	}
	return plan
}

func (c *ClusterOnOffOnDemand) validateCreateClusterOnOffRequest(plan providerschema.ClusterOnOffOnDemand) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	if plan.State.IsNull() {
		return errors.ErrOnoffStateCannotBeEmpty
	}
	return nil
}

// retrieveClusterOnOff retrieves onDemandClusterOnOff information from the specified organization and project using the provided cluster ID by open-api call.
func (c *ClusterOnOffOnDemand) retrieveClusterOnOff(ctx context.Context, organizationId, projectId, clusterId, state string) (*providerschema.ClusterOnOffOnDemand, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", c.HostURL, organizationId, projectId, clusterId)
	cfg := cluster_onoff_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.Client.ExecuteWithRetry(
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
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
		State:          types.StringValue(state),
	}

	return &refreshedState, nil
}

func validateClusterStateIsSameInPlanAndState(planClusterState, stateClusterState string) bool {
	return strings.EqualFold(planClusterState, stateClusterState)
}

func (c *ClusterOnOffOnDemand) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Couchbase Capella's v4 does not support a GET endpoint for cluster on/off.
	// Cluster on/off can only access the POST and DELETE endpoint for switching the cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
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

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, state.State.String())
	if err != nil {
		//resourceNotFound, errString := cluster_onoff_api.CheckResourceNotFoundError(err)
		//if resourceNotFound {
		//	tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		//	resp.State.RemoveResource(ctx)
		//	return
		//}
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

// Update allows to switch the cluster to OFF state.
func (c *ClusterOnOffOnDemand) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan providerschema.ClusterOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterOffRequest := cluster_onoff_api.CreateClusterOffRequest{}

	//resourceIDs, err := plan.Validate()
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error updating onDemandClusterOnOff",
	//		"Could not switch off the cluster, unexpected error: "+err.Error(),
	//	)
	//	return
	//}

	if err := c.validateCreateClusterOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandClusterOnOff request",
			"Could not switch on the cluster, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	//var (
	//	organizationId = resourceIDs[providerschema.OrganizationId]
	//	projectId      = resourceIDs[providerschema.ProjectId]
	//	clusterId      = resourceIDs[providerschema.ClusterId]
	//	bucketId       = resourceIDs[providerschema.BucketId]
	//	scopeName      = resourceIDs[providerschema.ScopeName]
	//	collectionName = resourceIDs[providerschema.CollectionName]
	//)

	// Update existing onDemandClusterOnOff
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/activationState", c.HostURL, organizationId, projectId, clusterId)
	cfg := cluster_onoff_api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		clusterOffRequest,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing the operation to switch off the cluster",
			errorMessageWhileClusterOnOffCreation+cluster_onoff_api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveClusterOnOff(ctx, organizationId, projectId, clusterId, plan.State.ValueString())
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

// Delete deletes the onDemandClusterOnOff.
func (c *ClusterOnOffOnDemand) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support a DELETION/destroying resource for cluster on/off.
	// Cluster on/off can only access the POST and DELETE endpoint which are used for switching the cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
}
