package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AccessFunction{}
	_ resource.ResourceWithConfigure   = &AccessFunction{}
	_ resource.ResourceWithImportState = &AccessFunction{}
)

const errorMessageAfterAccessFunctionCreation = "Access function creation is successful, but encountered an error while checking the current" +
	" state of the access function. Please run `terraform plan` after 1-2 minutes to know the" +
	" current access function state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileAccessFunctionCreation = "There is an error during access function creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// AccessFunction is the AccessFunction resource implementation.
type AccessFunction struct {
	*providerschema.Data
}

func NewAccessFunction() resource.Resource {
	return &AccessFunction{}
}

// Metadata returns the access function resource type name.
func (r *AccessFunction) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_control_function"
}

// Schema defines the schema for the access function resource.
func (r *AccessFunction) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AccessFunctionSchema()
}

// Configure sets provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (r *AccessFunction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create creates a new access function using PUT (upsert).
func (r *AccessFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AccessFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.validateCreateAccessFunctionRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create access function request",
			"Could not create access function, "+err.Error(),
		)
		return
	}

	IDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create access function request",
			"Could not create access function, "+err.Error(),
		)
		return
	}
	organizationId := IDs["organization_id"]
	projectId := IDs["project_id"]
	clusterId := IDs["cluster_id"]
	appServiceId := IDs["app_service_id"]
	keyspace := IDs["keyspace"]

	// Create access function using PUT (upsert)
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		r.HostURL, organizationId, projectId, clusterId, appServiceId, keyspace)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.AccessControlFunction.ValueString(),
		r.Token,
		map[string]string{"Content-Type": "application/javascript"},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileAccessFunctionCreation+api.ParseError(err),
		)
		return
	}

	// Set initial state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh state to get current values
	refreshedState, err := r.refreshAccessFunction(ctx, organizationId, projectId, clusterId, appServiceId, keyspace)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading access function after creation",
			errorMessageAfterAccessFunctionCreation+api.ParseError(err),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *AccessFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AccessFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters
	IDs, err := state.ValidateState()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Access Function",
			"Could not read Capella access function: "+err.Error(),
		)
		return
	}

	// Refresh state
	refreshedState, err := r.refreshAccessFunction(ctx, IDs["organizationId"], IDs["projectId"], IDs["clusterId"],
		IDs["appServiceId"], IDs["keyspace"])
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "access function doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Access Function",
			"Could not read Capella access function: "+errString,
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the access function using PUT (upsert).
func (r *AccessFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AccessFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract IDs from plan
	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	keyspace := plan.Keyspace.ValueString()

	// Update access function using PUT (upsert)
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		r.HostURL, organizationId, projectId, clusterId, appServiceId, keyspace)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.AccessControlFunction.ValueString(),
		r.Token,
		map[string]string{"Content-Type": "application/javascript"},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not update access function: "+api.ParseError(err),
		)
		return
	}

	// Refresh state
	refreshedState, err := r.refreshAccessFunction(ctx, organizationId, projectId, clusterId, appServiceId, keyspace)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading access function after update",
			"Could not read access function after update: "+api.ParseError(err),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the access function.
func (r *AccessFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AccessFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters
	IDs, err := r.validateAccessFunctionState(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella Access Function",
			"Could not delete Capella access function: "+err.Error(),
		)
		return
	}

	// Delete access function
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		r.HostURL, IDs["organizationId"], IDs["projectId"], IDs["clusterId"],
		IDs["appServiceId"], IDs["keyspace"])

	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, _ := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Capella Access Function",
			"Could not delete Capella access function: "+api.ParseError(err),
		)
		return
	}
}

// ImportState imports a resource into Terraform state.
func (r *AccessFunction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID should be in the format: organizationId,projectId,clusterId,appServiceId,keyspace
	resource.ImportStatePassthroughID(ctx, path.Root("keyspace"), req, resp)
}

// Helper functions

func (r *AccessFunction) validateCreateAccessFunctionRequest(plan providerschema.AccessFunction) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return fmt.Errorf("project ID cannot be empty")
	}
	if plan.ClusterId.IsNull() {
		return fmt.Errorf("cluster ID cannot be empty")
	}
	if plan.AppServiceId.IsNull() {
		return fmt.Errorf("app service ID cannot be empty")
	}
	if plan.Keyspace.IsNull() {
		return fmt.Errorf("keyspace cannot be empty")
	}
	if plan.AccessControlFunction.IsNull() {
		return fmt.Errorf("access control function cannot be empty")
	}
	return nil
}

func (r *AccessFunction) validateAccessFunctionState(state providerschema.AccessFunction) (map[string]string, error) {
	if state.OrganizationId.IsNull() {
		return nil, errors.ErrOrganizationIdCannotBeEmpty
	}
	if state.ProjectId.IsNull() {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if state.ClusterId.IsNull() {
		return nil, fmt.Errorf("cluster ID cannot be empty")
	}
	if state.AppServiceId.IsNull() {
		return nil, fmt.Errorf("app service ID cannot be empty")
	}

	if state.Keyspace.IsNull() {
		return nil, fmt.Errorf("keyspace cannot be empty")
	}

	return map[string]string{
		"organizationId": state.OrganizationId.ValueString(),
		"projectId":      state.ProjectId.ValueString(),
		"clusterId":      state.ClusterId.ValueString(),
		"appServiceId":   state.AppServiceId.ValueString(),
		"keyspace":       state.Keyspace.ValueString(),
	}, nil
}

func (r *AccessFunction) getAccessFunction(ctx context.Context, organizationId, projectId, clusterId, appServiceId, keyspace string) (string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		r.HostURL, organizationId, projectId, clusterId, appServiceId, keyspace)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	return string(response.Body), nil
}

func (r *AccessFunction) refreshAccessFunction(ctx context.Context, organizationId, projectId, clusterId, appServiceId, keyspace string) (*providerschema.AccessFunction, error) {
	accessFunction, err := r.getAccessFunction(ctx, organizationId, projectId, clusterId, appServiceId, keyspace)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	refreshedState := &providerschema.AccessFunction{
		OrganizationId:        types.StringValue(organizationId),
		ProjectId:             types.StringValue(projectId),
		ClusterId:             types.StringValue(clusterId),
		AppServiceId:          types.StringValue(appServiceId),
		AccessControlFunction: types.StringValue(accessFunction),
		Keyspace:              types.StringValue(keyspace),
	}
	return refreshedState, nil
}

// API structures for access function operations
type AccessFunctionRequest json.RawMessage

type AccessFunctionResponse struct {
	Function string                 `json:"function"`
	Audit    api.CouchbaseAuditData `json:"audit"`
}
