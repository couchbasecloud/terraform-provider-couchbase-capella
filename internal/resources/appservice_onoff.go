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
func (a *AppServiceOnOffOnDemand) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// Metadata returns the AppServiceOnOffOnDemand cluster resource type name.
func (a *AppServiceOnOffOnDemand) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_onoff_ondemand"
}

// Schema defines the schema for AppServiceOnOffOnDemand.
func (a *AppServiceOnOffOnDemand) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceOnOffOnDemandSchema()
}

// Configure It adds the provider configured api to ClusterOnOff.
func (a *AppServiceOnOffOnDemand) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	a.Data = data
}

// Create allows to switch the cluster to ON or OFF state.
func (a *AppServiceOnOffOnDemand) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := a.validateAppServiceOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandAppServiceOnOff request",
			"Could not switch on the app service, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	if err := a.manageAppServiceActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId, appServiceId); err != nil {
		resp.Diagnostics.AddError(
			"App Service activation failed",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, plan.State.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AppServiceOnOffOnDemand",
			fmt.Sprintf("Could not read Capella AppServiceOnOffOnDemand for the app service: %s associated to cluster: %s: %s", appServiceId, clusterId, errorMessageAfterAppServiceOnOffCreation+app_service_onoff_api.ParseError(err)),
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

func (a *AppServiceOnOffOnDemand) manageAppServiceActivation(ctx context.Context, state, organizationId, projectId, clusterId, appServiceId string) error {
	var (
		url    = fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/activationState", a.HostURL, organizationId, projectId, clusterId, appServiceId)
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

	cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return fmt.Errorf(errorMessageWhileAppServiceOnOffCreation + app_service_onoff_api.ParseError(err))
	}
	return nil
}

func (a *AppServiceOnOffOnDemand) validateAppServiceOnOffRequest(plan providerschema.AppServiceOnOffOnDemand) error {
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
func (a *AppServiceOnOffOnDemand) retrieveAppServiceOnOff(ctx context.Context, organizationId, projectId, clusterId, appServiceId, state string) (*providerschema.AppServiceOnOffOnDemand, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := app_service_onoff_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
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
func (a *AppServiceOnOffOnDemand) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	refreshedState, err := a.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, state.State.String())
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
func (a *AppServiceOnOffOnDemand) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve values from plan
	var plan providerschema.AppServiceOnOffOnDemand
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := a.validateAppServiceOnOffRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create onDemandAppServiceOnOff request",
			"Could not switch on the app service, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	if err := a.manageAppServiceActivation(ctx, plan.State.ValueString(), organizationId, projectId, clusterId, appServiceId); err != nil {
		resp.Diagnostics.AddError(
			"app service activation failed",
			err.Error(),
		)
		return
	}

	refreshedState, err := a.retrieveAppServiceOnOff(ctx, organizationId, projectId, clusterId, appServiceId, plan.State.ValueString())
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

func (a *AppServiceOnOffOnDemand) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Couchbase Capella's v4 does not support a DELETION/destroying resource for cluster on/off.
	// Cluster on/off can only access the POST and DELETE endpoint which are used for switching the cluster to on and off state respectively.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters/operation/clusterOn
}
