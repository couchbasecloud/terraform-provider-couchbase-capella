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
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
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
func (r *ImportFilter) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_import_filter"
}

// Schema defines the Terraform schema for this resource.
func (r *ImportFilter) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ImportFilterSchema()
}

// Configure sets provider-defined data, clients, etc.
func (r *ImportFilter) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// buildImportFilterURL builds the URL for the import filter API endpoint.
func buildImportFilterURL(hostURL, organizationId, projectId, clusterId, appServiceId, keyspace string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/importFilter",
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
func (r *ImportFilter) fetchImportFilter(ctx context.Context, organizationId, projectId, clusterId, appServiceId, keyspace string) (string, error) {
	url := buildImportFilterURL(r.HostURL, organizationId, projectId, clusterId, appServiceId, keyspace)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		return "", err
	}
	return string(response.Body), nil
}

// Create upserts the import filter.
func (r *ImportFilter) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ImportFilter
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required fields
	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError("Error Creating Import Filter", errors.ErrOrganizationIdCannotBeEmpty.Error())
		return
	}
	if plan.ProjectId.IsNull() {
		resp.Diagnostics.AddError("Error Creating Import Filter", errors.ErrProjectIdMissing.Error())
		return
	}
	if plan.ClusterId.IsNull() {
		resp.Diagnostics.AddError("Error Creating Import Filter", errors.ErrClusterIdMissing.Error())
		return
	}
	if plan.AppServiceId.IsNull() {
		resp.Diagnostics.AddError("Error Creating Import Filter", errors.ErrAppServiceIdMissing.Error())
		return
	}
	if plan.Keyspace.IsNull() {
		resp.Diagnostics.AddError("Error Creating Import Filter", "keyspace cannot be empty")
		return
	}

	url := buildImportFilterURL(
		r.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
		plan.Keyspace.ValueString(),
	)

	headers := map[string]string{"Content-Type": "application/javascript"}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.ImportFilter.ValueString(),
		r.Token,
		headers,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Import Filter",
			"Could not upsert Import Filter: "+api.ParseError(err),
		)
		return
	}

	// Refresh from server to ensure state matches remote
	body, err := r.fetchImportFilter(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), plan.AppServiceId.ValueString(), plan.Keyspace.ValueString())
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
func (r *ImportFilter) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ImportFilter
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.ValidateState()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Import Filter", errors.ErrValidatingResource.Error())
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		appServiceId   = resourceIDs[providerschema.AppServiceId]
		keyspace       = resourceIDs[providerschema.Keyspace]
	)

	url := buildImportFilterURL(
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		keyspace,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
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
			"Error Reading Import Filter",
			"Could not read Import Filter: "+errString,
		)
		return
	}

	state.ImportFilter = types.StringValue(string(response.Body))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the import filter.
func (r *ImportFilter) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ImportFilter
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := buildImportFilterURL(
		r.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
		plan.Keyspace.ValueString(),
	)

	headers := map[string]string{"Content-Type": "application/javascript"}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		plan.ImportFilter.ValueString(),
		r.Token,
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
	body, err := r.fetchImportFilter(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), plan.AppServiceId.ValueString(), plan.Keyspace.ValueString())
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
func (r *ImportFilter) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.ImportFilter
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := buildImportFilterURL(
		r.HostURL,
		state.OrganizationId.ValueString(),
		state.ProjectId.ValueString(),
		state.ClusterId.ValueString(),
		state.AppServiceId.ValueString(),
		state.Keyspace.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
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
func (r *ImportFilter) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("keyspace"), req, resp)
}
