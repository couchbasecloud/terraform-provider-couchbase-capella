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
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = &DataAPI{}
	_ resource.ResourceWithConfigure   = &DataAPI{}
	_ resource.ResourceWithImportState = &DataAPI{}
)

const (
	errorMessageWhileDataAPICreation = "There is an error during Data API enablement. Please check in Capella to see if any hanging resources have been created, unexpected error: "
	errorMessageAfterDataAPICreation = "Data API enablement is initiated, but encountered an error while checking the current state. Please run `terraform plan` after 1-2 minutes to know the current state. Additionally, run `terraform apply --refresh-only` to update the state from remote, unexpected error: "
	errorMessageWhileDataAPIUpdate   = "There is an error during Data API update. Please check in Capella to see if any hanging resources have been created, unexpected error: "
	errorMessageAfterDataAPIUpdate   = "Data API update is initiated, but encountered an error while checking the current state. Please run `terraform plan` after 1-2 minutes to know the current state. Additionally, run `terraform apply --refresh-only` to update the state from remote, unexpected error: "
	errorMessageWhileDataAPIDeletion = "There is an error during Data API disablement, unexpected error: "
	errorMessageAfterDataAPIDeletion = "Data API disablement is initiated, but encountered an error while checking the current state, unexpected error: "
)

// DataAPI is the Data API resource implementation.
type DataAPI struct {
	*providerschema.Data
}

// NewDataAPI is a helper function to simplify the provider implementation.
func NewDataAPI() resource.Resource {
	return &DataAPI{}
}

// Metadata returns the Data API resource type name.
func (d *DataAPI) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_api"
}

// Schema defines the schema for the Data API resource.
func (d *DataAPI) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DataAPISchema()
}

// Configure adds the provider configured client to the Data API resource.
func (d *DataAPI) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.Data = data
}

// ImportState imports a Data API resource.
func (d *DataAPI) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// Create enables Data API on the cluster.
func (d *DataAPI) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.DataAPI
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateDataAPIRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error validating Data API request",
			"Could not validate Data API request, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	requestBody := api.UpdateDataAPIRequest{
		EnableDataApi:        plan.EnableDataApi.ValueBool(),
		EnableNetworkPeering: plan.EnableNetworkPeering.ValueBool(),
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, requestBody, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error enabling Data API",
			errorMessageWhileDataAPICreation+api.ParseError(err),
		)
		return
	}

	refreshedState, err := d.checkDataAPIStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error enabling Data API",
			errorMessageAfterDataAPICreation+api.ParseError(err),
		)
		return
	}

	refreshedState.EnableDataApi = plan.EnableDataApi
	refreshedState.EnableNetworkPeering = plan.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads the Data API status.
func (d *DataAPI) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.DataAPI
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Data API status",
			"Could not read Capella Data API status on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	refreshedState, err := d.getDataAPIState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Data API status",
			"Error reading Data API status, unexpected error: "+errString,
		)
		return
	}

	refreshedState.EnableDataApi = state.EnableDataApi
	refreshedState.EnableNetworkPeering = state.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the Data API configuration.
func (d *DataAPI) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.DataAPI
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	requestBody := api.UpdateDataAPIRequest{
		EnableDataApi:        plan.EnableDataApi.ValueBool(),
		EnableNetworkPeering: plan.EnableNetworkPeering.ValueBool(),
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, requestBody, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Data API",
			errorMessageWhileDataAPIUpdate+api.ParseError(err),
		)
		return
	}

	refreshedState, err := d.checkDataAPIStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Data API",
			errorMessageAfterDataAPIUpdate+api.ParseError(err),
		)
		return
	}

	refreshedState.EnableDataApi = plan.EnableDataApi
	refreshedState.EnableNetworkPeering = plan.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete disables Data API on the cluster.
func (d *DataAPI) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.DataAPI
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Data API state",
			"Could not validate Capella Data API state on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	requestBody := api.UpdateDataAPIRequest{
		EnableDataApi:        false,
		EnableNetworkPeering: false,
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err = d.ClientV1.ExecuteWithRetry(ctx, cfg, requestBody, d.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
		resp.Diagnostics.AddError(
			"Error disabling Data API",
			errorMessageWhileDataAPIDeletion+errString,
		)
		return
	}

	_, err = d.checkDataAPIStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error disabling Data API",
			errorMessageAfterDataAPIDeletion+api.ParseError(err),
		)
	}
}

// checkDataAPIStatus polls the Data API status until it reaches a final state.
func (d *DataAPI) checkDataAPIStatus(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.DataAPI, error) {
	var (
		dataAPIResp *api.GetDataAPIStatusResponse
		err         error
	)

	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = 1 * time.Second

	timer := time.NewTimer(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%s: %w", errors.ErrDataAPIStatusTimeout, err)
		case <-timer.C:
			dataAPIResp, err = d.getDataAPIStatus(ctx, organizationId, projectId, clusterId)
			switch err {
			case nil:
				if api.IsDataAPIFinalState(dataAPIResp.State) && api.IsDataAPIFinalState(dataAPIResp.StateForNetworkPeering) {
					refreshedState := &providerschema.DataAPI{
						OrganizationId:         types.StringValue(organizationId),
						ProjectId:              types.StringValue(projectId),
						ClusterId:              types.StringValue(clusterId),
						State:                  types.StringValue(dataAPIResp.State),
						StateForNetworkPeering: types.StringValue(dataAPIResp.StateForNetworkPeering),
						ConnectionString:       types.StringValue(dataAPIResp.ConnectionString),
					}
					return refreshedState, nil
				}
				tflog.Info(ctx, "waiting for Data API to complete the operation")
			default:
				continue
			}

			timer.Reset(sleep)
		}
	}
}

// getDataAPIStatus retrieves the current Data API status.
func (d *DataAPI) getDataAPIStatus(ctx context.Context, organizationId, projectId, clusterId string) (*api.GetDataAPIStatusResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		return nil, err
	}

	status := api.GetDataAPIStatusResponse{}
	err = json.Unmarshal(response.Body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// getDataAPIState morphs the Data API status into a terraform schema.
func (d *DataAPI) getDataAPIState(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.DataAPI, error) {
	response, err := d.getDataAPIStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	state := &providerschema.DataAPI{
		OrganizationId:         types.StringValue(organizationId),
		ProjectId:              types.StringValue(projectId),
		ClusterId:              types.StringValue(clusterId),
		State:                  types.StringValue(response.State),
		StateForNetworkPeering: types.StringValue(response.StateForNetworkPeering),
		ConnectionString:       types.StringValue(response.ConnectionString),
	}

	return state, nil
}

// validateDataAPIRequest ensures organization id, project id and cluster id are valued.
func validateDataAPIRequest(plan providerschema.DataAPI) error {
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
