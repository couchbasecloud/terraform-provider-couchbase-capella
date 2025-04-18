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
	cluster_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = &FreeTierClusterOnOff{}
	_ resource.ResourceWithConfigure   = &FreeTierClusterOnOff{}
	_ resource.ResourceWithImportState = &FreeTierClusterOnOff{}
)

// FreeTierClusterOnOff is a struct that represents the free tier on/off status of a cluster.
type FreeTierClusterOnOff struct {
	*providerschema.Data
}

func (f *FreeTierClusterOnOff) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	data, ok := request.ProviderData.(*providerschema.Data)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	f.Data = data
}

func (f *FreeTierClusterOnOff) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), request, response)
}

// NewFreeTierClusterOnOff creates a new instance of FreeTierClusterOnOff.
func NewFreeTierClusterOnOff() resource.Resource {
	return &FreeTierClusterOnOff{}
}

// Metadata returns the type name for the resource.
func (f *FreeTierClusterOnOff) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_free_tier_cluster_on_off"
}

func (f *FreeTierClusterOnOff) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FreeTierClusterOnOffSchema()
}

// Create allows to swtich the free-tier cluster on and off.
func (f *FreeTierClusterOnOff) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierClusterOnOff
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	if err := f.manageFreeTierClusterActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId); err != nil {
		response.Diagnostics.AddError(
			"Free-Tier cluster activation failed",
			err.Error(),
		)
		return
	}

	refreshedState, err := f.checkClusterForDesiredStatus(ctx, organizationId, projectId, clusterId, plan.State.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading free-tier cluster on/off",
			"Could not read cluster on/off for the cluster: %s "+clusterId+"."+errorMessageAfterClusterOnOffCreation+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierClusterOnOff) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state providerschema.FreeTierClusterOnOff
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading free-tier cluster on/off details in Capella",
			"Could not validate the free-tier cluster on/off for cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	refreshedState, err := f.retrieveFreeTierClusterOnOff(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, _ := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error parsing read free-tier cluster on/off request",
			"Could not read the cluster details, unexpected error: "+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, &refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierClusterOnOff) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan providerschema.FreeTierClusterOnOff
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	if err := f.manageFreeTierClusterActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId); err != nil {
		response.Diagnostics.AddError(
			"Free-Tier Cluster activation failed",
			err.Error(),
		)
		return
	}

	refreshedState, err := f.checkClusterForDesiredStatus(ctx, organizationId, projectId, clusterId, plan.State.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Capella Free Tier On/Off",
			"Could not read Capella Free Tier on/off for the cluster: %s "+clusterId+"."+errorMessageAfterClusterOnOffCreation+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data.
	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierClusterOnOff) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support a DELETION/destroying resource for free-tier cluster on/off.
	// Free Tier Cluster on/off can only access the POST and DELETE endpoint which are used for switching the free-tier cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/Free-Tier/operation/freeTierClusterOn.
}

// This function is used to manage the activation state of a free-tier cluster.
// It switches off/on the free-tier cluster based on the provided input and uses the appropriate HTTP method (POST for "on" and DELETE for "off").
func (f *FreeTierClusterOnOff) manageFreeTierClusterActivation(ctx context.Context, state, organizationId, projectId, clusterId string) error {
	var (
		url    = fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s/activationState", f.HostURL, organizationId, projectId, clusterId)
		method string
	)

	switch state {
	case "on":
		method = http.MethodPost
	case "off":
		method = http.MethodDelete
	default:
		return fmt.Errorf("invalid state value: state must be either 'on' or 'off'")
	}

	cfg := api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}
	tflog.Info(ctx, url)
	_, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s: %s", errorMessageWhileClusterOnOffCreation, api.ParseError(err))
	}
	return nil
}

// retrieveFreeTierClusterOnOff retrieves the current state of the free-tier cluster.
// if the cluster is in "TurnedOff" state, it returns "off", otherwise it returns "on".
func (f *FreeTierClusterOnOff) retrieveFreeTierClusterOnOff(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.FreeTierClusterOnOff, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	var state string
	clusterResp := cluster_api.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}

	if clusterResp.CurrentState != cluster_api.TurnedOff {
		state = "on"
	} else {
		state = "off"
	}

	refreshedState := providerschema.FreeTierClusterOnOff{
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
		State:          types.StringValue(state),
	}

	return &refreshedState, nil
}

// This function checks the cluster's current state and waits until it matches the desired state.
// When Turning on the cluster it checks for "Healthy" state and when turning off it checks for "TurnedOff" state.
func (f *FreeTierClusterOnOff) checkClusterForDesiredStatus(ctx context.Context, organizationId, projectId, clusterId, state string) (*providerschema.FreeTierClusterOnOff, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	var clusterResp cluster_api.GetClusterResponse
	const timeout = 20 * time.Minute
	const retryInterval = 3 * time.Second // Added retry delay.

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var desiredState cluster_api.State
	switch state {
	case "on":
		desiredState = cluster_api.Healthy
	case "off":
		desiredState = cluster_api.TurnedOff
	default:
		return nil, fmt.Errorf("invalid state: %s", state)
	}

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cluster state check timed out: %w", ctx.Err())
		case <-ticker.C:
			response, err := f.Client.ExecuteWithRetry(ctx, cfg, nil, f.Token, nil)
			if err != nil {
				return nil, fmt.Errorf("API request failed: %w", err)
			}

			if err := json.Unmarshal(response.Body, &clusterResp); err != nil {
				return nil, fmt.Errorf("failed to parse response: %w", err)
			}

			if clusterResp.CurrentState == desiredState {
				refreshedState := providerschema.FreeTierClusterOnOff{
					ClusterId:      types.StringValue(clusterId),
					ProjectId:      types.StringValue(projectId),
					OrganizationId: types.StringValue(organizationId),
					State:          types.StringValue(state),
				}
				return &refreshedState, nil
			}

			tflog.Debug(ctx, fmt.Sprintf("Current cluster state: %s (waiting for %s)", clusterResp.CurrentState, desiredState))
		}
	}
}
