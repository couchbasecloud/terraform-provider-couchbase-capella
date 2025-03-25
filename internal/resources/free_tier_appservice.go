package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"time"
)

// Ensure the FreeTierAppService implements the required interfaces
var (
	_ resource.Resource                = &FreeTierAppService{}
	_ resource.ResourceWithConfigure   = &FreeTierAppService{}
	_ resource.ResourceWithImportState = &FreeTierAppService{}
)

// FreeTierAppService represents the free-tier app service resource.
type FreeTierAppService struct {
	*providerschema.Data
}

// NewFreeTierAppService is a helpfer function to simplify the provider implementation
func NewFreeTierAppService() resource.Resource {
	return &FreeTierAppService{}
}

// Metadata returns the metadata for the free-tier app service.
func (f *FreeTierAppService) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_free_tier_app_service"
}

// Schema defines the schema for the free-tier app service.
func (f *FreeTierAppService) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = FreeTierAppServiceSchema()
}

// Create the free-tier app service.
func (f *FreeTierAppService) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierAppService
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	createFreeTierAppServiceRequest := appservice.CreateFreeTierAppServiceRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createFreeTierAppServiceRequest.Description = plan.Description.ValueStringPointer()
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/freeTier", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	apiResp, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		createFreeTierAppServiceRequest,
		f.Token,
		nil,
	)
	if err != nil {
		response.Diagnostics.AddError("Failed to create free tier app service",
			errors.ErrFreeTierCreateAppServiceError.Error()+api.ParseError(err))
		return
	}
	freeTierCreateAppServiceResponse := appservice.CreateAppServiceResponse{}
	err = json.Unmarshal(apiResp.Body, &freeTierCreateAppServiceResponse)
	if err != nil {
		response.Diagnostics.AddError("Failed to create free tier app service",
			errors.ErrFreeTierCreateAppServiceError.Error()+" error during unmarshalling "+api.ParseError(err))
		return
	}

	diags = response.State.Set(ctx, initializePendingFreeTierAppServiceWithPlanAndId(plan, freeTierCreateAppServiceResponse.Id.String()))
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err = f.checkFreeTierAppServiceStatus(ctx, organizationId, projectId, clusterId, freeTierCreateAppServiceResponse.Id.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating app service",
			errors.ErrFreeTierAppServiceAfterCreation.Error()+api.ParseError(err),
		)
		return
	}

	refreshedState, err := f.refreshFreeTierAppService(ctx, organizationId, projectId, clusterId, freeTierCreateAppServiceResponse.Id.String())
	if err != nil {
		response.Diagnostics.AddError("Failed to create free tier app service",
			errors.ErrFreeTierAppServiceAfterCreation.Error()+" error during refreshing "+api.ParseError(err))
		return
	}

	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

// Read the free-tier app service.
func (f *FreeTierAppService) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var currState providerschema.FreeTierAppService
	diags := request.State.Get(ctx, &currState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := currState.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading free tier app service",
			"Could not read Capella free-tier app service id: "+currState.Id.String()+err.Error(),
		)
		return
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.Id]
	)

	refreshedState, err := f.refreshFreeTierAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading free tier app service",
			"Could not read the free tier app service with ID "+currState.Id.String()+": "+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

// Update the free-tier app service.
func (f *FreeTierAppService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state providerschema.FreeTierAppService
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()

	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free-tier app service",
			"Could not update free-tier app service id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	if err := f.validateAppServiceAttributesTrimmed(plan); err != nil {
		response.Diagnostics.AddError(
			"Error updating free-tier app service",
			"Could not update free-tier app service id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.Id]
	)

	updateFreeTierAppServiceRequest := appservice.UpdateFreeTierAppServiceRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		updateFreeTierAppServiceRequest,
		f.Token,
		nil,
	)
	if err != nil {
		response.Diagnostics.AddError("Failed to update free tier app service",
			"could not update app service with id "+state.Id.String()+api.ParseError(err))
		return
	}
	err = f.checkFreeTierAppServiceStatus(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier app service",
			"could not update app service with id "+state.Id.String()+api.ParseError(err),
		)
	}

	refreshedState, err := f.refreshFreeTierAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		response.Diagnostics.AddError("Failed to update free tier app service",
			"could not update app service with id "+state.Id.String()+api.ParseError(err))
		return
	}

	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the free-tier app service.
func (f *FreeTierAppService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state providerschema.FreeTierAppService
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	resourceIDs, err := state.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting free-tier app service",
			"Could not delete free-tier app service id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appserviceId   = resourceIDs[providerschema.Id]
	)
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, appserviceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		resourceNotfound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotfound {
			tflog.Info(ctx, "resource does not exist in remote server removing resource from the state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Failed to delete free tier app service",
			"could not delete app service with id "+state.Id.String()+errString,
		)
		return
	}

	err = f.checkFreeTierAppServiceStatus(ctx, organizationId, projectId, clusterId, appserviceId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			response.Diagnostics.AddError(
				"Failed to delete free tier app service",
				"could not delete app service with id "+state.Id.String()+errString,
			)
			return
		}
		return
	}

	// This will only be reached when app service deletion has failed,
	// and the app service record still exists in the cp metadata. Therefore,
	// no error will be returned when performing a GET call.
	freeTierAppService, err := f.refreshFreeTierAppService(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString(), state.Id.ValueString())
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error deleting app service",
			"Could not delete app service id "+state.Id.String()+": "+errString,
		)
		return
	}
	response.Diagnostics.AddError(
		"Error deleting app service",
		fmt.Sprintf("Could not delete free-tier app service id %s, as current app service state: %s", state.Id.String(), freeTierAppService.CurrentState),
	)
}

// Configure adds the provider configured client to the free-tier app service resource.
func (f *FreeTierAppService) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

// ImportState imports a remote free-tier app service that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=appService123,organization_id=org123,project_id=proj123,cluster_id=cluster123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (f *FreeTierAppService) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

// initializePendingFreeTierAppServiceWithPlanAndId initializes an instance of providerschema.AppService
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingFreeTierAppServiceWithPlanAndId(plan providerschema.FreeTierAppService, id string) providerschema.FreeTierAppService {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}
	plan.Compute = types.ObjectNull(providerschema.AppServiceCompute{}.AttributeTypes())
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Etag = types.StringNull()
	return plan
}

// getAppService retrieves app service information from the specified organization, project and cluster
// using the provided app service ID by open-api call.
func (f *FreeTierAppService) getFreeTierAppService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*appservice.GetAppServiceResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, appServiceId)
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

	appServiceResp := appservice.GetAppServiceResponse{}
	err = json.Unmarshal(response.Body, &appServiceResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	appServiceResp.Etag = response.Response.Header.Get("ETag")
	return &appServiceResp, nil
}

func (f *FreeTierAppService) refreshFreeTierAppService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.FreeTierAppService, error) {
	appServiceResponse, err := f.getFreeTierAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(appServiceResponse.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)

	}
	compute := providerschema.NewAppServiceCompute(appServiceResponse.Compute)
	computeObj, diags := types.ObjectValueFrom(ctx, compute.AttributeTypes(), compute)

	refreshedState := providerschema.NewFreeTierAppService(
		appServiceResponse,
		organizationId,
		projectId,
		auditObj,
		computeObj,
	)
	return refreshedState, nil
}

// checkFreeTierAppServiceStatus monitors the status of an app service creation, update and deletion operation for a specified
// organization, project, cluster and appService ID. It periodically fetches the app service status using the `getFreeTierAppService`
// function and waits until the app service reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (f *FreeTierAppService) checkFreeTierAppServiceStatus(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) error {
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

	timer := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return errors.ErrAppServiceCreationStatusTimeout

		case <-timer.C:
			appServiceResp, err = f.getFreeTierAppService(ctx, organizationId, projectId, clusterId, appServiceId)
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

func (f *FreeTierAppService) validateAppServiceAttributesTrimmed(plan providerschema.FreeTierAppService) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	if (!plan.Description.IsNull() && !plan.Description.IsUnknown()) && !providerschema.IsTrimmed(plan.Description.ValueString()) {
		return fmt.Errorf("description %s", errors.ErrNotTrimmed)
	}
	return nil
}
