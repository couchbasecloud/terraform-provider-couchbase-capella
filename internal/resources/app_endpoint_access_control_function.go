package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AccessControlFunction{}
	_ resource.ResourceWithConfigure   = &AccessControlFunction{}
	_ resource.ResourceWithImportState = &AccessControlFunction{}
)

const errorUpsertingAccessFunction = "There was an error upserting the access control function.  Error: "

// AccessControlFunction is the Access Control Function resource implementation.
type AccessControlFunction struct {
	*providerschema.Data
}

func NewAccessControlFunction() resource.Resource {
	return &AccessControlFunction{}
}

// Metadata returns the access function resource type name.
func (a *AccessControlFunction) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_access_control_function"
}

// Schema defines the schema for the access function resource.
func (a *AccessControlFunction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AccessControlFunctionSchema()
}

// Create upserts an access control function.
func (a *AccessControlFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AccessControlFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		keyspace       = fmt.Sprintf(
			"%s.%s.%s",
			plan.AppEndpoint.ValueString(),
			plan.Scope.ValueString(),
			plan.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.AccessControlFunction.ValueString(),
		a.Token,
		map[string]string{"Content-Type": "application/javascript"},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorUpsertingAccessFunction+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (a *AccessControlFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AccessControlFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Access Function",
			"Could not read app endpoint access function: "+err.Error(),
		)
		return
	}

	keyspace := fmt.Sprintf("%s.%s.%s", IDs["appEndpointName"], IDs["scopeName"], IDs["collectionName"])

	accessControlFunction, err := a.getAccessFunction(
		ctx,
		IDs["organizationId"],
		IDs["projectId"],
		IDs["clusterId"],
		IDs["appServiceId"],
		keyspace,
	)
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

	diags = resp.State.Set(ctx, providerschema.AccessControlFunction{
		OrganizationId:        types.StringValue(IDs["organizationId"]),
		ProjectId:             types.StringValue(IDs["projectId"]),
		ClusterId:             types.StringValue(IDs["clusterId"]),
		AppServiceId:          types.StringValue(IDs["appServiceId"]),
		AppEndpoint:           types.StringValue(IDs["appEndpointName"]),
		Scope:                 types.StringValue(IDs["scopeName"]),
		Collection:            types.StringValue(IDs["collectionName"]),
		AccessControlFunction: types.StringValue(accessControlFunction),
	})
	resp.Diagnostics.Append(diags...)
}

// Update updates the access control function used by app endpoint.
func (a *AccessControlFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AccessControlFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		keyspace       = fmt.Sprintf(
			"%s.%s.%s",
			plan.AppEndpoint.ValueString(),
			plan.Scope.ValueString(),
			plan.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.AccessControlFunction.ValueString(),
		a.Token,
		map[string]string{"Content-Type": "application/javascript"},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorUpsertingAccessFunction+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the access function.
func (a *AccessControlFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AccessControlFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
		keyspace       = fmt.Sprintf(
			"%s.%s.%s",
			state.AppEndpoint.ValueString(),
			state.Scope.ValueString(),
			state.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
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
func (a *AccessControlFunction) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// The import ID should be in the format: organizationId,projectId,clusterId,appServiceId,keyspace
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint"), req, resp)
}

// Configure sets provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (a *AccessControlFunction) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
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
	a.Data = data
}

func (a *AccessControlFunction) getAccessFunction(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId, keyspace string,
) (string, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/accessControlFunction",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	return string(response.Body), nil
}
