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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	app_service_onoff_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apps_service_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppServiceOnOffOnDemand{}
	_ resource.ResourceWithConfigure   = &AppServiceOnOffOnDemand{}
	_ resource.ResourceWithImportState = &AppServiceOnOffOnDemand{}
)

const errorMessageWhileAppServiceOnOffCreation = "There is an error during switching on/off the cluster. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterAppServiceOnOffCreation = "Cluster switch on/off is successful, but encountered an error while checking the current" +
	" state of the switched on/off cluster. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// AppServiceOnOffOnDemand is the AppServiceOnOffOnDemand implementation.
type AppServiceOnOffOnDemand struct {
	*providerschema.Data
}

// NewAppServiceOnOffOnDemand is a helper function to simplify the provider implementation.
func NewAppServiceOnOffOnDemand() resource.Resource {
	return &AppServiceOnOffOnDemand{}
}

// ImportState imports a remote AppserviceOnOffOnDemand app service that is not created by Terraform.
func (c *AppServiceOnOffOnDemand) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("_app_service_id"), req, resp)
}

// Metadata returns the AppServiceOnOffOnDemand cluster resource type name.
func (c *AppServiceOnOffOnDemand) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_onoff_ondemand"
}

// Schema defines the schema for AppServiceOnOffOnDemand.
func (c *AppServiceOnOffOnDemand) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterOnOffOnDemandSchema()
}

// Configure It adds the provider configured api to ClusterOnOff.
func (c *AppServiceOnOffOnDemand) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (c *AppServiceOnOffOnDemand) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	appServiceOnRequest := app_service_onoff_api.CreateAppServiceOnRequest{}
	appServiceOffRequest := app_service_onoff_api.CreateAppServiceOffRequest{}

	if err := c.validateCreateAppServiceOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandClusterOnOff request",
			"Could not switch on the app service, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var appServiceId = plan.AppServiceId.ValueString()

	if plan.State.ValueString() == "on" {

		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/activationState", c.HostURL, organizationId, projectId, clusterId, appServiceId)
		cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
		_, err := c.Client.ExecuteWithRetry(
			ctx,
			cfg,
			appServiceOnRequest,
			c.Token,
			nil,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error executing the operation to switch on the cluster",
				errorMessageWhileAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
			)
			return
		}
	} else if plan.State.ValueString() == "off" {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/activationState", c.HostURL, organizationId, projectId, clusterId, appServiceId)
		cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
		_, err := c.Client.ExecuteWithRetry(
			ctx,
			cfg,
			appServiceOffRequest,
			c.Token,
			nil,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error executing the operation to switch off the cluster",
				errorMessageWhileAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Invalid state value",
			"State must be either 'on' or 'off'",
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, plan.State.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AppServiceOnOffOnDemand",
			"Could not read Capella AppServiceOnOffOnDemand for the cluster: %s "+clusterId+"."+errorMessageAfterAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
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

func (c *AppServiceOnOffOnDemand) validateCreateAppServiceOnOffRequest(plan providerschema.AppServiceOnOffOnDemand) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	if plan.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdCannotBeEmpty
	}
	if plan.State.IsNull() {
		return errors.ErrOnoffStateCannotBeEmpty
	}
	return nil
}

// retrieveAppServiceOnOff retrieves AppServiceOnOff information from the specified organization and project using the provided cluster ID by Get cluster open-api call.
func (c *AppServiceOnOffOnDemand) retrieveAppServiceOnOff(ctx context.Context, organizationId, projectId, clusterId, appServiceId, state string) (*providerschema.AppServiceOnOffOnDemand, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", c.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
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

	//There is no GET endpoint so get the app service response and check current state
	appServiceResp := apps_service_api.GetAppServiceResponse{}
	err = json.Unmarshal(response.Body, &appServiceResp)
	if err != nil {
		return nil, err
	}

	if validateAppserviceStateIsSameInPlanAndState(state, string(appServiceResp.CurrentState)) {
		appServiceResp.CurrentState = apps_service_api.State(state)
	}

	refreshedState := providerschema.AppServiceOnOffOnDemand{
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
		AppServiceId:   types.StringValue(appServiceId),
		State:          types.StringValue(state),
	}

	return &refreshedState, nil
}

func validateAppserviceStateIsSameInPlanAndState(planAppServiceState, stateAppServiceState string) bool {
	return strings.EqualFold(planAppServiceState, stateAppServiceState)
}

// Couchbase Capella's v4 does not support a GET endpoint for app service on/off.
// App service on/off can only access the POST and DELETE endpoint for switching the app service to on and off state respectively.
// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/appServices/operation/appServiceOn
// This read is calling the retrieveAppServiceOnOff func to verify the state with the cluster response.
func (c *AppServiceOnOffOnDemand) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceOnOffOnDemand
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading app service on/off details in Capella",
			"Could not validate the app service on/off for app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
	)

	refreshedState, err := c.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, state.State.String())
	if err != nil {
		resourceNotFound, _ := app_service_onoff_api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error parsing read AppServiceOnOffOnDemand request",
			"Could not read the app service details, unexpected error: "+err.Error(),
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
func (c *AppServiceOnOffOnDemand) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve values from plan
	var plan providerschema.AppServiceOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	appServiceOnRequest := app_service_onoff_api.CreateClusterOnRequest{}
	appServiceOffRequest := app_service_onoff_api.CreateClusterOffRequest{}

	if err := c.validateCreateAppServiceOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create AppServiceOnOffOnDemand request",
			"Could not switch on/off the app service, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var appServiceId = plan.AppServiceId.ValueString()

	if plan.State.ValueString() == "on" {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/activationState", c.HostURL, organizationId, projectId, clusterId, appServiceId)
		cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
		_, err := c.Client.ExecuteWithRetry(
			ctx,
			cfg,
			appServiceOnRequest,
			c.Token,
			nil,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error executing the operation to switch on the cluster",
				errorMessageWhileAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
			)
			return
		}
	} else if plan.State.ValueString() == "off" {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/activationState", c.HostURL, organizationId, projectId, clusterId, appServiceId)
		cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
		_, err := c.Client.ExecuteWithRetry(
			ctx,
			cfg,
			appServiceOffRequest,
			c.Token,
			nil,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error executing the operation to switch off the cluster",
				errorMessageWhileAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Invalid state value",
			"State must be either 'on' or 'off'",
		)
		return
	}

	refreshedState, err := c.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, plan.State.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AppserviceOnOffOnDemand",
			"Could not read Capella AppserviceOnOffOnDemand for the cluster: %s "+clusterId+"."+errorMessageAfterAppServiceOnOffCreation+app_service_onoff_api.ParseError(err),
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

func (c *AppServiceOnOffOnDemand) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support a DELETION/destroying resource for cluster on/off.
	// Cluster on/off can only access the POST and DELETE endpoint which are used for switching the cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
}
