package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	freeTierClusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/freeTierCluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &FreeTierCluster{}
	_ resource.ResourceWithConfigure   = &FreeTierCluster{}
	_ resource.ResourceWithImportState = &FreeTierCluster{}
)

const errorMessageAfterFreeTierClusterCreationInitiation = "Cluster creation is initiated, but encountered an error while checking the current" +
	" state of the cluster. Please run `terraform plan` after 4-5 minutes to know the" +
	" current status of the cluster. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileFreeTierClusterCreation = "There is an error during cluster creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

type FreeTierCluster struct {
	*providerschema.Data
}

func NewFreeTierCluster() resource.Resource {
	return &FreeTierCluster{}
}

func (f *FreeTierCluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_free_tier_cluster"

}

func (f *FreeTierCluster) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FreeTierClusterSchema()
}

func (f FreeTierCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if err := f.validateFreeTierCreateCluster(plan); err != nil {
		response.Diagnostics.AddError(
			"error while validating create free tier cluster",
			"could not create free tier cluster "+err.Error(),
		)
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
	if plan.OrganizationId.IsNull() {
		response.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()
	if plan.ProjectId.IsNull() {
		response.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
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
			errorMessageWhileFreeTierClusterCreation+"error during unmarshalling"+err.Error(),
		)
		return
	}
	freeTierClusterResponse := freeTierClusterapi.GetFreeTierClusterResponse{}
	err = json.Unmarshal(res.Body, &freeTierClusterResponse)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Cluster",
			errorMessageWhileFreeTierClusterCreation+"error during unmarshalling:"+err.Error(),
		)
		return
	}
	diags = response.State.Set(ctx, initializePendingFreeTierClusterWithPlanAndId(plan, freeTierClusterResponse.ID.String()))
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	err = f.checkFreeTierClusterStatus(ctx, organizationId, projectId, freeTierClusterResponse.ID.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			errorMessageAfterFreeTierClusterCreationInitiation+api.ParseError(err),
		)
		return
	}
	refreshedState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, freeTierClusterResponse.ID.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			errorMessageAfterFreeTierClusterCreationInitiation+api.ParseError(err),
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

func (f FreeTierCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
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

	//get refreshed cluster values from capella
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
	//state.ServiceGroups = refreshedState.ServiceGroups
	diags = response.State.Set(ctx, &refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f FreeTierCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
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

	if err := f.validateFreeTierClusterUpdate(plan, state); err != nil {
		response.Diagnostics.AddError(
			"Error updating cluster",
			"Could not update cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	FreeTierClusterUpdateRequest := freeTierClusterapi.UpdateFreeTierClusterRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	var headers = make(map[string]string)

	// Update existing Cluster
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		FreeTierClusterUpdateRequest,
		f.Token,
		headers,
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

	err = f.checkFreeTierClusterStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating cluster",
			"Could not update cluster id "+state.Id.String()+": "+api.ParseError(err),
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

	// Set state to fully populated data
	diags = response.State.Set(ctx, currentState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f FreeTierCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Retrieve values from state
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

	// Delete existing Cluster
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
		// resourceNotFound as expected
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

func (f FreeTierCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	/// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (f *FreeTierCluster) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (f *FreeTierCluster) validateFreeTierCreateCluster(plan providerschema.FreeTierCluster) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	return nil
}

// initializePendingClusterWithPlanAndId initializes an instance of providerschema.Cluster
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingFreeTierClusterWithPlanAndId(plan providerschema.FreeTierCluster, id string) providerschema.FreeTierCluster {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}

	if plan.EnablePrivateDNSResolution.IsNull() || plan.EnablePrivateDNSResolution.IsUnknown() {
		plan.EnablePrivateDNSResolution = types.BoolNull()
	}

	if plan.CouchbaseServer.IsNull() || plan.CouchbaseServer.IsUnknown() {
		plan.CouchbaseServer = types.ObjectNull(providerschema.CouchbaseServer{}.AttributeTypes())
	}
	plan.AppServiceId = types.StringNull()
	plan.ConnectionString = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Availability = types.ObjectNull(providerschema.Availability{}.AttributeTypes())
	plan.CmekId = types.StringNull()

	if plan.ServiceGroups.IsNull() || plan.ServiceGroups.IsUnknown() {
		plan.ServiceGroups = types.SetNull(types.ObjectType{}.WithAttributeTypes(providerschema.ServiceGroupAttributeTypes()))
	}
	return plan
}

// checkFreeTierClusterStatus monitors the status of a cluster creation, update and deletion operation for a specified
// organization, project, and cluster ID. It periodically fetches the cluster status using the `getCluster`
// function and waits until the cluster reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (c *FreeTierCluster) checkFreeTierClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *freeTierClusterapi.GetFreeTierClusterResponse
		err         error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion
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
			clusterResp, err = c.getFreeTierCluster(ctx, organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterapi.State(clusterResp.CurrentState)) {
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

// getCluster retrieves cluster information from the specified organization and project
// using the provided cluster ID by open-api call.
func (f *FreeTierCluster) getFreeTierCluster(ctx context.Context, organizationId, projectId, clusterId string) (*freeTierClusterapi.GetFreeTierClusterResponse, error) {
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

	clusterResp := freeTierClusterapi.GetFreeTierClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	//clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}

// retrieveCluster retrieves cluster information for a specified organization, project, and cluster ID.
func (f *FreeTierCluster) retrieveFreeTierCluster(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.FreeTierCluster, error) {
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
		return nil, fmt.Errorf("unable to convert availablity data %w", err)
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

// validateClusterUpdate checks if specific fields in a cluster can be updated and returns an error if not.
func (f *FreeTierCluster) validateFreeTierClusterUpdate(plan, state providerschema.FreeTierCluster) error {
	var planOrganizationId, stateOrganizationId string
	if !plan.OrganizationId.IsNull() {
		planOrganizationId = plan.OrganizationId.ValueString()
	}

	if !state.OrganizationId.IsNull() {
		stateOrganizationId = state.OrganizationId.ValueString()
	}

	if planOrganizationId != stateOrganizationId {
		return errors.ErrUnableToUpdateOrganizationId
	}

	var planProjectId, stateProjectId string
	if !plan.ProjectId.IsNull() {
		planProjectId = plan.ProjectId.ValueString()
	}

	if !state.ProjectId.IsNull() {
		stateProjectId = state.ProjectId.ValueString()
	}

	if planProjectId != stateProjectId {
		return errors.ErrUnableToUpdateProjectId
	}

	var planCloudProvider, stateCloudProvider providerschema.CloudProvider
	if plan.CloudProvider != nil {
		planCloudProvider = *plan.CloudProvider
	}
	if state.CloudProvider != nil {
		stateCloudProvider = *state.CloudProvider
	}

	if !reflect.DeepEqual(planCloudProvider, stateCloudProvider) {
		return errors.ErrUnableToUpdateCloudProvider
	}

	return nil
}
