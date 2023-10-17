package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/api/appservice"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppService{}
	_ resource.ResourceWithConfigure   = &AppService{}
	_ resource.ResourceWithImportState = &AppService{}
)

type AppService struct {
	*providerschema.Data
}

func NewAppService() resource.Resource {
	return &AppService{}
}

func (a *AppService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service"

}

func (a *AppService) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceSchema()
}

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
			"Error parsing create_app_service.tf app service request",
			"Could not create_app_service.tf app service "+err.Error(),
		)
		return
	}

	appServiceRequest := appservice.CreateAppServiceRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Nodes:       plan.Nodes.ValueInt64(),
		Compute: appservice.Compute{
			Cpu: plan.Compute.Cpu.ValueInt64(),
			Ram: plan.Compute.Ram.ValueInt64(),
		},
	}

	if !plan.Version.IsNull() && !plan.Version.IsUnknown() {
		version := plan.Version.ValueString()
		appServiceRequest.Version = &version
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices", a.HostURL, organizationId, projectId, clusterId),
		http.MethodPost,
		appServiceRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not execute request, unexpected error: "+err.Error(),
		)
		return
	}

	createAppServiceResponse := appservice.CreateAppServiceResponse{}
	err = json.Unmarshal(response.Body, &createAppServiceResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app service",
			"Could not create app service, unexpected error: "+err.Error(),
		)
		return
	}

	err = a.checkAppServiceStatus(ctx, organizationId, projectId, clusterId, createAppServiceResponse.Id.String())
	_, err = handleAppServiceError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app service",
			"Could not create app service, unexpected error: "+err.Error(),
		)
		return
	}
	refreshedState, err := a.refreshAppService(ctx, organizationId, projectId, clusterId, createAppServiceResponse.Id.String())
	_, err = handleAppServiceError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app service",
			"Could not create app service, unexpected error: "+err.Error(),
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

func (a *AppService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	appServiceId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service",
			"Could not read Capella app service list: "+err.Error(),
		)
		return
	}

	// Refresh the existing user
	refreshedState, err := a.refreshAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella App Service",
				"Could not read Capella appServiceID "+appServiceId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service",
			"Could not read Capella appServiceID "+appServiceId+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *AppService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (a *AppService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
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
// example: id=user123,organization_id=org123,project_id=proj123,cluster_id=cluster123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (a *AppService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

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

func (a *AppService) refreshAppService(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.OneAppService, error) {
	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId),
		http.MethodGet,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	appServiceResponse := appservice.GetAppServiceResponse{}
	err = json.Unmarshal(response.Body, &appServiceResponse)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.OneAppService{
		Id:            types.StringValue(appServiceId),
		Name:          types.StringValue(appServiceResponse.Name),
		Description:   types.StringValue(appServiceResponse.Description),
		CloudProvider: types.StringValue(appServiceResponse.CloudProvider),
		Nodes:         types.Int64Value(int64(appServiceResponse.Nodes)),
		Compute: providerschema.Compute{
			Cpu: types.Int64Value(appServiceResponse.Compute.Cpu),
			Ram: types.Int64Value(appServiceResponse.Compute.Ram),
		},
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		CurrentState:   types.StringValue(string(appServiceResponse.CurrentState)),
		Version:        types.StringValue(appServiceResponse.Version),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(appServiceResponse.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(appServiceResponse.Audit.CreatedBy),
			ModifiedAt: types.StringValue(appServiceResponse.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(appServiceResponse.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(appServiceResponse.Audit.Version)),
		},
	}
	return &refreshedState, nil
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
			appServiceResp, err = a.getAppService(organizationId, projectId, clusterId, appServiceId)
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
func (a *AppService) getAppService(organizationId, projectId, clusterId, appServiceId string) (*appservice.GetAppServiceResponse, error) {
	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", a.HostURL, organizationId, projectId, clusterId, appServiceId),
		http.MethodGet,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	appServiceResp := appservice.GetAppServiceResponse{}
	err = json.Unmarshal(response.Body, &appServiceResp)
	if err != nil {
		return nil, err
	}
	return &appServiceResp, nil
}

// handleAppServiceError extracts error message if error is api.Error and also checks whether error is
// resource not found
func handleAppServiceError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}
