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

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*Cors)(nil)
	_ resource.ResourceWithConfigure   = (*Cors)(nil)
	_ resource.ResourceWithImportState = (*Cors)(nil)
)

// Cors is the resource implementation for configuring App Endpoint CORS.
type Cors struct {
	*providerschema.Data
}

// NewCors is used in (p *capellaProvider) Resources for building the provider.
func NewCors() resource.Resource {
	return &Cors{}
}

// Metadata returns the CORS resource type name.
func (c *Cors) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_cors"
}

// Schema defines the schema for the CORS resource.
func (c *Cors) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CorsSchema()
}

// Create is used to create a new CORS configuration.
func (c *Cors) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Cors
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var origin []string
	diags = plan.Origin.ElementsAs(ctx, &origin, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loginOrigin []string
	if !plan.LoginOrigin.IsNull() {
		diags = plan.LoginOrigin.ElementsAs(ctx, &loginOrigin, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var headers []string
	if !plan.Headers.IsNull() {
		diags = plan.Headers.ElementsAs(ctx, &headers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	corsReq := api.CorsRequest{
		Origin:      origin,
		LoginOrigin: loginOrigin,
		Headers:     headers,
		MaxAge:      plan.MaxAge.ValueInt64(),
		Disabled:    plan.Disabled.ValueBool(),
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/cors",
		c.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)
	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodPut,
		SuccessStatus: http.StatusNoContent,
	}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		corsReq,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating CORS Configuration",
			"Could not create CORS configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}

	// Initialize computed attributes to null if not set
	if plan.MaxAge.IsNull() || plan.MaxAge.IsUnknown() {
		plan.MaxAge = types.Int64Null()
	}
	if plan.Disabled.IsNull() || plan.Disabled.IsUnknown() {
		plan.Disabled = types.BoolNull()
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the created CORS configuration
	refreshedState, err := c.refreshCors(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading CORS configuration after creation",
			"CORS was created successfully but could not be read back: "+api.ParseError(err),
		)
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// getCors is used to retrieve an existing CORS configuration.
func (c *Cors) getCors(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string,
) (*api.CorsResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/cors",
		c.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	corsResp := &api.CorsResponse{}
	err = json.Unmarshal(response.Body, corsResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	return corsResp, nil
}

// refreshCors is used to pass an existing CORS configuration to the refreshed state.
func (c *Cors) refreshCors(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string,
) (*providerschema.Cors, error) {
	corsResp, err := c.getCors(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		return nil, err
	}

	originSet, diags := types.SetValueFrom(ctx, types.StringType, corsResp.Origin)
	if diags.HasError() {
		var allErrs []string
		for _, d := range diags.Errors() {
			allErrs = append(allErrs, d.Detail())
		}
		return nil, fmt.Errorf("failed to convert origin to set: %s", strings.Join(allErrs, "; "))
	}

	loginOriginSet, diags := types.SetValueFrom(ctx, types.StringType, corsResp.LoginOrigin)
	if diags.HasError() {
		var allErrs []string
		for _, d := range diags.Errors() {
			allErrs = append(allErrs, d.Detail())
		}
		return nil, fmt.Errorf("failed to convert login origin to set: %s", strings.Join(allErrs, "; "))
	}

	headersSet, diags := types.SetValueFrom(ctx, types.StringType, corsResp.Headers)
	if diags.HasError() {
		var allErrs []string
		for _, d := range diags.Errors() {
			allErrs = append(allErrs, d.Detail())
		}
		return nil, fmt.Errorf("failed to convert headers to set: %s", strings.Join(allErrs, "; "))
	}

	refreshedState := &providerschema.Cors{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		Origin:          originSet,
		LoginOrigin:     loginOriginSet,
		Headers:         headersSet,
		MaxAge:          types.Int64Value(corsResp.MaxAge),
		Disabled:        types.BoolValue(corsResp.Disabled),
	}

	return refreshedState, nil
}

// Read is used to read an existing CORS configuration and set the state.
func (c *Cors) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Cors
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, projectId, clusterId, appServiceId, appEndpointName, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating CORS state",
			"Could not validate CORS: "+err.Error(),
		)
		return
	}

	// Refresh the existing CORS configuration
	refreshedState, err := c.refreshCors(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading CORS Configuration",
			"Could not read CORS configuration: "+errString,
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update is used to update an existing CORS configuration.
func (c *Cors) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.Cors
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert list attributes to string slices
	var origin []string
	if !plan.Origin.IsNull() && !plan.Origin.IsUnknown() {
		diags = plan.Origin.ElementsAs(ctx, &origin, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var loginOrigin []string
	if !plan.LoginOrigin.IsNull() && !plan.LoginOrigin.IsUnknown() {
		diags = plan.LoginOrigin.ElementsAs(ctx, &loginOrigin, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var headers []string
	if !plan.Headers.IsNull() && !plan.Headers.IsUnknown() {
		diags = plan.Headers.ElementsAs(ctx, &headers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	corsReq := api.CorsRequest{
		Origin:      origin,
		LoginOrigin: loginOrigin,
		Headers:     headers,
		MaxAge:      plan.MaxAge.ValueInt64(),
		Disabled:    plan.Disabled.ValueBool(),
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/cors",
		c.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		corsReq,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating CORS Configuration",
			"Could not update CORS configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}

	// Read the updated CORS configuration
	refreshedState, err := c.refreshCors(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddWarning(
			"Error reading CORS configuration after update",
			"CORS was updated successfully but could not be read back: "+errString,
		)
		// Set state with plan data since update was successful
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete is used to delete an existing CORS configuration.
// There's no actual delete endpoint for CORS so just remove it from the state file.
func (c *Cors) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// As per the requirements, there's no actual delete endpoint for CORS
	// We just remove it from the state file
	tflog.Info(ctx, "CORS configuration removed from state.")
}

// ImportState is used to import a remote CORS configuration into Terraform state.
func (c *Cors) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

// Configure is used to configure the CORS resource with the provider data.
func (c *Cors) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
