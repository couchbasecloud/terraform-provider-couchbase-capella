package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"time"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppService{}
	_ resource.ResourceWithConfigure   = &AppService{}
	_ resource.ResourceWithImportState = &AppService{}
)

const errorMessageAfterAppServiceCreationInitiation = "App Service creation is initiated, but encountered an error while checking the current" +
	" state of the app service. Please run `terraform plan` after 4-5 minutes to know the" +
	" current status of the app service. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileAppServiceCreation = "There is an error during app service creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// AppService is the AppService resource implementation.
type AppService struct {
	*providerschema.Data
}

// NewAppService is a helper function to simplify the provider implementation.
func NewAppService() resource.Resource {
	return &AppService{}
}

// Metadata returns the AppService resource type name.
func (a *AppService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service"

}

// Schema defines the schema for the AppService resource.
func (a *AppService) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceSchema()
}

// Create creates a new AppService.
func (a *AppService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validateCreateAppServiceRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create app service request",
			"Could not create app service "+err.Error(),
		)
		return
	}

	appServiceRequest := appservice.CreateAppServiceRequest{
		Name: plan.Name.ValueString(),
		Compute: appservice.AppServiceCompute{
			Cpu: plan.Compute.Cpu.ValueInt64(),
			Ram: plan.Compute.Ram.ValueInt64(),
		},
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		appServiceRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Nodes.IsNull() && !plan.Nodes.IsUnknown() {
		appServiceRequest.Nodes = plan.Nodes.ValueInt64Pointer()
	}

	if !plan.Version.IsNull() && !plan.Version.IsUnknown() {
		version := plan.Version.ValueString()
		appServiceRequest.Version = &version
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		appServiceRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileAppServiceCreation+err.Error(),
		)
		return
	}

	createAppServiceResponse := appservice.CreateAppServiceResponse{}
	err = json.Unmarshal(response.Body, &createAppServiceResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app service",
			errorMessageWhileAppServiceCreation+"error during unmarshalling:"+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializePendingAppServiceWithPlanAndId(plan, createAppServiceResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = a.checkAppServiceStatus(ctx, organizationId, projectId, clusterId, createAppServiceResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating app service",
			errorMessageAfterAppServiceCreationInitiation+api.ParseError(err),
		)
		return
	}
	refreshedState, err := a.refreshAppService(ctx, organizationId, projectId, clusterId, createAppServiceResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating app service",
			errorMessageAfterAppServiceCreationInitiation+api.ParseError(err),
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

// Read reads the app service project information.
func (a *AppService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service",
			"Could not read Capella app service id: "+err.Error(),
		)
		return
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.Id]
	)

	// Refresh the existing app service
	refreshedState, err := a.refreshAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service",
			"Could not read Capella appServiceID "+appServiceId+": "+errString,
		)
		return
	}

	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		refreshedState.IfMatch = state.IfMatch
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the AppService.
func (a *AppService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state providerschema.AppService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app service",
			"Could not update app service id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.Id]
	)

	// Added temporarily until https://couchbasecloud.atlassian.net/browse/AV-65838 is fixed
	if plan.Name != state.Name {
		resp.Diagnostics.AddError(
			"Error updating app service",
			"Could not update app service id "+state.Id.String()+" unexpected error: "+errors.ErrUnableToUpdateAppServiceName.Error(),
		)
		return
	}

	appServiceRequest := appservice.UpdateAppServiceRequest{
		Nodes: plan.Nodes.ValueInt64(),
		Compute: appservice.AppServiceCompute{
			Cpu: plan.Compute.Cpu.ValueInt64(),
			Ram: plan.Compute.Ram.ValueInt64(),
		},
	}

	var headers = make(map[string]string)
	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		headers["If-Match"] = state.IfMatch.ValueString()
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		appServiceRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app service",
			"Could not update app service id "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	err = a.checkAppServiceStatus(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app service",
			"Could not update app service id "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	currentState, err := a.refreshAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app service",
			"Could not update app service id "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	if !plan.IfMatch.IsUnknown() && !plan.IfMatch.IsNull() {
		currentState.IfMatch = plan.IfMatch
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the app service.
func (a *AppService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting app service",
			"Could not delete app service id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.Id]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	// Delete existing App Service
	_, err = a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
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
			"Error deleting app service",
			"Could not delete app service id "+state.Id.String()+": "+errString,
		)
		return
	}

	err = a.checkAppServiceStatus(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString(), state.Id.ValueString())
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			resp.Diagnostics.AddError(
				"Error deleting app service",
				"Could not delete app service id "+state.Id.String()+": "+errString,
			)
			return
		}
		// resourceNotFound as expected
		return
	}

	// This will only be reached when app service deletion has failed,
	// and the app service record still exists in the cp metadata. Therefore,
	// no error will be returned when performing a GET call.
	appService, err := a.refreshAppService(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString(), state.Id.ValueString())
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting app service",
			"Could not delete app service id "+state.Id.String()+": "+errString,
		)
		return
	}
	resp.Diagnostics.AddError(
		"Error deleting app service",
		fmt.Sprintf("Could not delete app service id %s, as current app service state: %s", state.Id.String(), appService.CurrentState),
	)
}

// Configure adds the provider configured client to the app service resource.
func (a *AppService) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports a remote app service that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=appService123,organization_id=org123,project_id=proj123,cluster_id=cluster123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (a *AppService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// validateCreateAppServiceRequest validates the payload of create app service request.
func (a *AppService) validateCreateAppServiceRequest(plan providerschema.AppService) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	return nil
}

// refreshAppService is used to pass an existing AppService to the refreshed state.
func (a *AppService) refreshAppService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.AppService, error) {
	appServiceResponse, err := a.getAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(appServiceResponse.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)

	}

	refreshedState := providerschema.NewAppService(
		appServiceResponse,
		organizationId,
		projectId,
		auditObj,
	)
	return refreshedState, nil
}

// checkAppServiceStatus monitors the status of an app service creation, update and deletion operation for a specified
// organization, project, cluster and appService ID. It periodically fetches the app service status using the `getAppService`
// function and waits until the app service reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (a *AppService) checkAppServiceStatus(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) error {
	var (
		appServiceResp *appservice.GetAppServiceResponse
		err            error
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
			const msg = "app service creation status transition timed out after initiation"
			return fmt.Errorf(msg)

		case <-timer.C:
			appServiceResp, err = a.getAppService(ctx, organizationId, projectId, clusterId, appServiceId)
			switch err {
			case nil:
				if appservice.IsFinalState(appServiceResp.CurrentState) {
					return nil
				}
				const msg = "waiting for app service to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

// getAppService retrieves app service information from the specified organization, project and cluster
// using the provided app service ID by open-api call
func (a *AppService) getAppService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*appservice.GetAppServiceResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	appServiceResp := appservice.GetAppServiceResponse{}
	err = json.Unmarshal(response.Body, &appServiceResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	appServiceResp.Etag = response.Response.Header.Get("ETag")
	return &appServiceResp, nil
}

// initializePendingAppServiceWithPlanAndId initializes an instance of providerschema.AppService
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingAppServiceWithPlanAndId(plan providerschema.AppService, id string) providerschema.AppService {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}
	if plan.Nodes.IsNull() || plan.Nodes.IsUnknown() {
		plan.Nodes = types.Int64Null()
	}
	if plan.CloudProvider.IsNull() || plan.CloudProvider.IsUnknown() {
		plan.CloudProvider = types.StringNull()
	}
	if plan.Version.IsNull() || plan.Version.IsUnknown() {
		plan.Version = types.StringNull()
	}
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Etag = types.StringNull()
	return plan
}
