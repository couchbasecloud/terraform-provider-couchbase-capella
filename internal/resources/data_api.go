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
	_ resource.Resource                = &DataAPI{}
	_ resource.ResourceWithConfigure   = &DataAPI{}
	_ resource.ResourceWithImportState = &DataAPI{}
)

const errorMessageAfterDataAPICreation = "Data API creation is successful, but encountered an error while checking the current" +
	" state of the Data API. Please run `terraform plan` after 1-2 minutes to know the" +
	" current Data API state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileDataAPICreation = "There is an error during Data API creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// dataAPIResourceModel is the resource model for the Data API resource.
// It includes both the user-configurable input fields and the computed output fields.
type dataAPIResourceModel struct {
	OrganizationId           types.String `tfsdk:"organization_id"`
	ProjectId                types.String `tfsdk:"project_id"`
	ClusterId                types.String `tfsdk:"cluster_id"`
	EnableDataApi            types.Bool   `tfsdk:"enable_data_api"`
	EnableNetworkPeering     types.Bool   `tfsdk:"enable_network_peering"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	State                    types.String `tfsdk:"state"`
	EnabledForNetworkPeering types.Bool   `tfsdk:"enabled_for_network_peering"`
	StateForNetworkPeering   types.String `tfsdk:"state_for_network_peering"`
	ConnectionString         types.String `tfsdk:"connection_string"`
}

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

// Create creates a new Data API configuration.
func (d *DataAPI) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataAPIResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.validateDataAPIRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create Data API request",
			"Could not create Data API "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	dataAPIRequest := api.UpdateDataAPIRequest{
		EnableDataApi:        plan.EnableDataApi.ValueBool(),
		EnableNetworkPeering: plan.EnableNetworkPeering.ValueBool(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		dataAPIRequest,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileDataAPICreation+api.ParseError(err),
		)
		return
	}

	refreshedState, err := d.retrieveDataAPI(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Reading Capella Data API",
			"Could not read Capella Data API for the cluster: "+clusterId+"."+errorMessageAfterDataAPICreation+api.ParseError(err),
		)
		return
	}

	// Preserve plan input values in the state
	refreshedState.EnableDataApi = plan.EnableDataApi
	refreshedState.EnableNetworkPeering = plan.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the Data API configuration.
func (d *DataAPI) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataAPIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.validateDataAPIRequest(state); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Data API",
			"Could not read Data API for cluster with id "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	refreshedState, err := d.retrieveDataAPI(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Data API",
			"Could not read Data API for cluster with id "+state.ClusterId.String()+": "+errString,
		)
		return
	}

	// Preserve existing plan values for input fields from current state
	refreshedState.EnableDataApi = state.EnableDataApi
	refreshedState.EnableNetworkPeering = state.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the Data API configuration.
func (d *DataAPI) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataAPIResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.validateDataAPIRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Data API",
			"Could not update Data API for cluster with id "+plan.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	dataAPIRequest := api.UpdateDataAPIRequest{
		EnableDataApi:        plan.EnableDataApi.ValueBool(),
		EnableNetworkPeering: plan.EnableNetworkPeering.ValueBool(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		dataAPIRequest,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Data API",
			"Could not update Data API for cluster with id "+plan.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	refreshedState, err := d.retrieveDataAPI(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Data API",
			"Could not read Data API for cluster with id "+plan.ClusterId.String()+": "+errString,
		)
		return
	}

	// Preserve plan input values in the state
	refreshedState.EnableDataApi = plan.EnableDataApi
	refreshedState.EnableNetworkPeering = plan.EnableNetworkPeering

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete disables the Data API by sending a PUT with both flags set to false.
func (d *DataAPI) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataAPIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.validateDataAPIRequest(state); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Data API",
			"Could not delete Data API for cluster with id "+state.ClusterId.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	dataAPIRequest := api.UpdateDataAPIRequest{
		EnableDataApi:        false,
		EnableNetworkPeering: false,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		dataAPIRequest,
		d.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Data API",
			"Could not delete Data API for cluster with id "+state.ClusterId.String()+" unexpected error: "+errString,
		)
		return
	}
}

// ImportState imports an already existing Data API configuration that is not created by Terraform.
func (d *DataAPI) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// Configure adds the provider configured client to the Data API resource.
func (d *DataAPI) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.Data = data
}

// validateDataAPIRequest validates that required IDs are not empty.
func (d *DataAPI) validateDataAPIRequest(model dataAPIResourceModel) error {
	if model.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if model.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if model.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	return nil
}

// retrieveDataAPI retrieves the Data API status from the specified organization, project, and cluster.
func (d *DataAPI) retrieveDataAPI(ctx context.Context, organizationId, projectId, clusterId string) (*dataAPIResourceModel, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	dataAPIResp := api.GetDataAPIStatusResponse{}
	err = json.Unmarshal(response.Body, &dataAPIResp)
	if err != nil {
		return nil, err
	}

	refreshedState := &dataAPIResourceModel{
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		Enabled:                  types.BoolValue(dataAPIResp.Enabled),
		State:                    types.StringValue(dataAPIResp.State),
		EnabledForNetworkPeering: types.BoolValue(dataAPIResp.EnabledForNetworkPeering),
		StateForNetworkPeering:   types.StringValue(dataAPIResp.StateForNetworkPeering),
		ConnectionString:         types.StringValue(dataAPIResp.ConnectionString),
	}

	return refreshedState, nil
}
