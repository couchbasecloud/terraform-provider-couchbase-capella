package resources

import (
	"context"
	"fmt"
	"net/http"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ImportFilter{}
	_ resource.ResourceWithConfigure   = &ImportFilter{}
	_ resource.ResourceWithImportState = &ImportFilter{}
)

// ImportFilter is the resource implementation for managing App Endpoint Import Filters.
type ImportFilter struct {
	*providerschema.Data
}

// NewImportFilter creates a new ImportFilter resource for provider initialization.
func NewImportFilter() resource.Resource {
	return &ImportFilter{}
}

// Metadata returns the resource type name.
func (f *ImportFilter) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_import_filter"
}

// Schema defines the Terraform schema for this resource.
func (f *ImportFilter) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ImportFilterSchema()
}

// Configure sets provider-defined data, clients, etc.
func (f *ImportFilter) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	f.Data = data
}

// buildImportFilterURL builds the URL for the import filter API endpoint.
func buildImportFilterURL(hostURL, organizationId, projectId, clusterId, appServiceId, keyspace string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s.%s.%s/importFilter",
		hostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)
}

// fetchImportFilter fetches the import filter from the API endpoint.
// returns the import filter function as a string.
func (f *ImportFilter) fetchImportFilter(ctx context.Context, organizationId, projectId, clusterId, appServiceId, keyspace string) (string, error) {
	url := buildImportFilterURL(f.HostURL, organizationId, projectId, clusterId, appServiceId, keyspace)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := f.Client.ExecuteWithRetry(ctx, cfg, nil, f.Token, nil)
	if err != nil {
		return "", err
	}
	return string(response.Body), nil
}

// Create upserts the import filter.
func (f *ImportFilter) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ImportFilter
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
			plan.AppEndpointName.ValueString(),
			plan.Scope.ValueString(),
			plan.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/importFilter",
		f.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	headers := map[string]string{"Content-Type": "application/javascript"}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.ImportFilter.ValueString(),
		f.Token,
		headers,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Import Filter",
			"Could not upsert Import Filter: "+api.ParseError(err),
		)
		return
	}

	// Refresh from app service to ensure state matches remote
	body, err := f.fetchImportFilter(ctx, organizationId, projectId, clusterId, appServiceId, keyspace)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Import Filter after creation",
			api.ParseError(err),
		)
	} else {
		plan.ImportFilter = types.StringValue(body)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read fetches the current import filter and refreshes state.
func (f *ImportFilter) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ImportFilter
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.ValidateState()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Import Filter", err.Error())
		return
	}

	keyspace := fmt.Sprintf("%s.%s.%s", IDs["appEndpointName"], IDs["scopeName"], IDs["collectionName"])

	response, err := f.fetchImportFilter(ctx, IDs["organizationId"], IDs["projectId"], IDs["clusterId"], IDs["appServiceId"], keyspace)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Import Filter",
			"Could not read Import Filter: "+errString,
		)
		return
	}

	state.ImportFilter = types.StringValue(response)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the import filter.
func (f *ImportFilter) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ImportFilter
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
			plan.AppEndpointName.ValueString(),
			plan.Scope.ValueString(),
			plan.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/importFilter",
		f.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	headers := map[string]string{"Content-Type": "application/javascript"}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.ImportFilter.ValueString(),
		f.Token,
		headers,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Import Filter",
			"Could not update Import Filter: "+api.ParseError(err),
		)
		return
	}

	// Refresh from server
	body, err := f.fetchImportFilter(ctx, organizationId, projectId, clusterId, appServiceId, keyspace)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Import Filter after update",
			api.ParseError(err),
		)
	} else {
		plan.ImportFilter = types.StringValue(body)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the import filter (resets/removes).
func (f *ImportFilter) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.ImportFilter
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
			state.AppEndpointName.ValueString(),
			state.Scope.ValueString(),
			state.Collection.ValueString(),
		)
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/importFilter",
		f.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := f.Client.ExecuteWithRetry(
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
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Import Filter",
			"Could not delete Import Filter: "+errString,
		)
		return
	}
}

// ImportState imports a resource into Terraform state.
func (f *ImportFilter) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}
