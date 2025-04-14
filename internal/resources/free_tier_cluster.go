package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	freeTierClusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/freeTierCluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &FreeTierCluster{}
	_ resource.ResourceWithConfigure   = &FreeTierCluster{}
	_ resource.ResourceWithImportState = &FreeTierCluster{}
)

type FreeTierCluster struct {
	*providerschema.Data
}

func NewFreeTierCluster() resource.Resource {
	return &FreeTierCluster{}
}

func (f *FreeTierCluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_free_tier_cluster"

}

func (f *FreeTierCluster) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FreeTierClusterSchema()
}

func (f *FreeTierCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	freeTierClusterCreateRequest := freeTierClusterapi.CreateFreeTierClusterRequest{
		Name: plan.Name.ValueString(),
		CloudProvider: clusterapi.CloudProvider{
			Cidr:   plan.CloudProvider.Cidr.ValueString(),
			Region: plan.CloudProvider.Region.ValueString(),
			Type:   clusterapi.CloudProviderType(plan.CloudProvider.Type.ValueString()),
		},
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		freeTierClusterCreateRequest.Description = plan.Description.ValueStringPointer()
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier", f.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	res, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		freeTierClusterCreateRequest,
		f.Token,
		nil,
	)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating cluster",
			errors.ErrorMessageWhileFreeTierClusterCreation.Error()+api.ParseError(err),
		)
		return
	}
	freeTierClusterResponse := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(res.Body, &freeTierClusterResponse)
	if err != nil {
		response.Diagnostics.AddError(
			"Error unmarshalling the response",
			errors.ErrorMessageWhileFreeTierClusterCreation.Error()+"error during unmarshalling:"+err.Error(),
		)
		return
	}
	diags = response.State.Set(ctx, initializePendingFreeTierClusterWithPlanAndId(plan, freeTierClusterResponse.Id.String()))
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	err = f.checkFreeTierClusterStatus(ctx, organizationId, projectId, freeTierClusterResponse.Id.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			errors.ErrorMessageAfterFreeTierClusterCreationInitiation.Error()+api.ParseError(err),
		)
		return
	}
	refreshedState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, freeTierClusterResponse.Id.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error fetching the cluster info",
			errors.ErrorMessageAfterFreeTierClusterCreationInitiation.Error()+api.ParseError(err),
		)
	}

	if clusterapi.State(refreshedState.CurrentState.ValueString()) != clusterapi.Healthy {
		response.Diagnostics.AddError(
			"Error creating cluster",
			fmt.Sprintf("Could not create cluster id %s, as current Cluster state: %s", refreshedState.Id.String(), refreshedState.CurrentState),
		)

	}

	// Set state to fully populated data.
	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f *FreeTierCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state providerschema.FreeTierCluster
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()

	if err != nil {
		response.Diagnostics.AddError(
			"error reading free tier cluster",
			"could not read free tier cluster "+state.Id.String()+": "+err.Error(),
		)
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	//get refreshed cluster values from capella.
	refreshedState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			"Could Not Read Capella Cluster "+state.Id.String()+": "+errString,
		)
		return
	}

	diags = response.State.Set(ctx, &refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f *FreeTierCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	freeTierClusterUpdateRequest := freeTierClusterapi.UpdateFreeTierClusterRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Update existing Cluster.
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		freeTierClusterUpdateRequest,
		f.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+": "+errString,
		)
		return
	}

	currentState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data.
	diags = response.State.Set(ctx, currentState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Retrieve values from state.
	var state providerschema.FreeTierCluster
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting free tier cluster",
			"Could not delete cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	// Delete existing Cluster.
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error Deleting Free Tier Cluster",
			"Could not delete cluster id "+state.Id.String()+": "+errString,
		)
		return
	}

	err = f.checkFreeTierClusterStatus(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.Id.ValueString())
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			response.Diagnostics.AddError(
				"Error Deleting Capella Cluster",
				"Could not delete cluster id "+state.Id.String()+": "+errString,
			)
			return
		}
		// resourceNotFound as expected.
		return
	}

	// This case will only occur when cluster deletion has failed,
	// and the cluster record still exists in the cp metadata. Therefore,
	// no error will be returned when performing a GET call.
	cluster, err := f.retrieveFreeTierCluster(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.Id.ValueString())
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error Deleting Capella Cluster",
			"Could not delete cluster id "+state.Id.String()+": "+errString,
		)
		return
	}
	response.Diagnostics.AddError(
		"Error deleting cluster",
		fmt.Sprintf("Could not delete cluster id %s, as current Cluster state: %s", state.Id.String(), cluster.CurrentState),
	)
}

func (f *FreeTierCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	/// Retrieve import ID and save to id attribute.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (f *FreeTierCluster) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	data, ok := request.ProviderData.(*providerschema.Data)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("expected *providerschema.FreeTierCluster, got %T", request.ProviderData),
		)
		return
	}
	f.Data = data
}

// initializePendingClusterWithPlanAndId initializes an instance of providerschema.Cluster.
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingFreeTierClusterWithPlanAndId(plan providerschema.FreeTierCluster, id string) providerschema.FreeTierCluster {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}

	plan.EnablePrivateDNSResolution = types.BoolNull()
	plan.CouchbaseServer = types.ObjectNull(providerschema.CouchbaseServer{}.AttributeTypes())
	plan.AppServiceId = types.StringNull()
	plan.ConnectionString = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Availability = types.ObjectNull(providerschema.Availability{}.AttributeTypes())
	plan.CmekId = types.StringNull()

	plan.ServiceGroups = types.SetNull(types.ObjectType{}.WithAttributeTypes(providerschema.ServiceGroupAttributeTypes()))
	plan.Support = types.ObjectNull(providerschema.Support{}.AttributeTypes())
	plan.Etag = types.StringNull()
	return plan
}

// checkFreeTierClusterStatus monitors the status of a cluster creation, update and deletion operation for a specified,
// organization, project, and cluster ID. It periodically fetches the cluster status using the `getCluster`
// function and waits until the cluster reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (f *FreeTierCluster) checkFreeTierClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *clusterapi.GetClusterResponse
		err         error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion.
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 3

	timer := time.NewTimer(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cluster creation status transition timed out after initiation, unexpected error: %w", err)
		case <-timer.C:
			clusterResp, err = f.getFreeTierCluster(ctx, organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterResp.CurrentState) {
					return nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

// getFreeTierCluster retrieves cluster information from the specified organization and project.
// using the provided cluster ID by open-api call.
func (f *FreeTierCluster) getFreeTierCluster(ctx context.Context, organizationId, projectId, clusterId string,
) (*clusterapi.GetClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	clusterResp := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}

// retrieveFreeTierCluster retrieves cluster information for a specified.
// organization, project, and cluster ID.
func (f *FreeTierCluster) retrieveFreeTierCluster(
	ctx context.Context,
	organizationId,
	projectId,
	clusterId string,
) (*providerschema.FreeTierCluster, error) {
	freeTierClusterResp, err := f.getFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(freeTierClusterResp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
	}

	availability := providerschema.NewAvailability(freeTierClusterResp.Availability)
	availabilityObj, diags := types.ObjectValueFrom(ctx, availability.AttributeTypes(), availability)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert availability data %w", err)
	}

	support := providerschema.NewSupport(freeTierClusterResp.Support)
	supportObj, diags := types.ObjectValueFrom(ctx, support.AttributeTypes(), support)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert support data %w", err)
	}

	serviceGroups, err := providerschema.NewTerraformServiceGroups(freeTierClusterResp)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert service groups data %w", err)
	}
	serviceGroupObjList, err, diag := providerschema.NewServiceGroups(ctx, serviceGroups)
	if err != nil {
		if diag.HasError() {
			return nil, err
		}
	}

	serviceGroupsObj, diags := types.SetValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(providerschema.ServiceGroupAttributeTypes()), serviceGroupObjList)
	if diags.HasError() {
		return nil, fmt.Errorf("error while converting servicegroups to service group object ")
	}

	refreshedState, err := providerschema.NewFreeTierCluster(ctx, freeTierClusterResp, organizationId, projectId, auditObj, availabilityObj, supportObj, serviceGroupsObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}
