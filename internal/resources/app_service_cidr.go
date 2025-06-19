package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*AppServiceCidr)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceCidr)(nil)
	_ resource.ResourceWithImportState = (*AppServiceCidr)(nil)
)

type AppServiceCidr struct {
	*providerschema.Data
}

// NewAppServiceCidr is used in (p *capellaProvider) Resources for building the provider.
func NewAppServiceCidr() resource.Resource {
	return &AppServiceCidr{}
}

// Metadata returns the App Service CIDR resource type name.
func (a *AppServiceCidr) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

// Schema defines the schema for the App Service CIDR resource.
func (a *AppServiceCidr) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AllowedCIDRsSchema()
}

// Configure is used to configure the App Service CIDR resource with the provider data.
func (a *AppServiceCidr) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create is used to create a new App Service Allowed CIDR.
func (a *AppServiceCidr) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceCIDR
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cidrReq := api.CreateAllowedCIDRRequest{
		Cidr:      plan.Cidr.ValueString(),
		Comment:   plan.Comment.ValueString(),
		ExpiresAt: plan.ExpiresAt.ValueString(),
	}
	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		cidrReq,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating App Service CIDR",
			"Could not create App Service CIDR, unexpected error: "+api.ParseError(err),
		)
		return
	}

	cidrResp := &api.AppServiceAllowedCIDRResponse{}
	err = json.Unmarshal(response.Body, &cidrResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling App Service CIDR Response",
			"Could not unmarshal App Service CIDR response, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeAllowedCIDRWithPlanAndId(plan, cidrResp.Id))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.refreshAllowedCIDR(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), plan.AppServiceId.ValueString(), cidrResp.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading App Service Allowed CIDR after creation",
			errorMessageAfterAllowListCreation+api.ParseError(err),
		)
		return
	}

	// This is added to workaround any timezone conversions that the API does automatically and
	// may cause an issue in the state file.
	if plan.ExpiresAt != refreshedState.ExpiresAt {
		refreshedState.ExpiresAt = plan.ExpiresAt
	}

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// initializeAllowedCIDRWithPlanAndId initializes an instance of providerschema.AppServiceCIDR
// with the specified plan and ID. It marks all computed fields as null.
func initializeAllowedCIDRWithPlanAndId(plan providerschema.AppServiceCIDR, id string) providerschema.AppServiceCIDR {
	plan.Id = types.StringValue(id)
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	if plan.Comment.IsNull() || plan.Comment.IsUnknown() {
		plan.Comment = types.StringNull()
	}
	return plan
}

// getAllowedCIDR is used to retrieve an existing allow list.
func (a *AppServiceCidr) getAllowedCIDR(ctx context.Context, organizationId, projectId, clusterId, appServiceId, allowListId string) (*api.AppServiceAllowedCIDRResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
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
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	allowListResp := api.ListAppServiceAllowedCIDRResponse{}
	err = json.Unmarshal(response.Body, &allowListResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	for _, allowlist := range allowListResp.Data {
		if allowlist.Id == allowListId {
			// Found the allow list with the matching ID
			return &allowlist, nil
		}
	}
	return nil, errors.ErrNotFound
}

// refreshAllowedCIDR is used to pass an existing Allowed CIDR to the refreshed state.
func (r *AppServiceCidr) refreshAllowedCIDR(ctx context.Context, organizationId, projectId, clusterId, appServiceId, allowedCIDRId string) (*providerschema.AppServiceCIDR, error) {
	allowListResp, err := r.getAllowedCIDR(ctx, organizationId, projectId, clusterId, appServiceId, allowedCIDRId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	// Create audit data object
	audit := providerschema.NewCouchbaseAuditData(allowListResp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting audit data to object")
	}

	refreshedState := providerschema.AppServiceCIDR{
		Id:             types.StringValue(allowListResp.Id),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
		Cidr:           types.StringValue(allowListResp.Cidr),
		Audit:          auditObj,
	}

	// Set optional fields
	if allowListResp.Comment != "" {
		refreshedState.Comment = types.StringValue(allowListResp.Comment)
	}

	if allowListResp.ExpiresAt != "" {
		refreshedState.ExpiresAt = types.StringValue(allowListResp.ExpiresAt)
	}

	return &refreshedState, nil
}

// Read is used to read an existing App Service CIDR and set the state.
func (a *AppServiceCidr) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceCIDR
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	organizationId, projectId, clusterId, appServiceId, allowedCIDRId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating App Service Allowed CIDR state",
			"Could not validate App Service Allowed CIDR: "+err.Error(),
		)
		return
	}

	// refresh the existing allow list
	refreshedState, err := a.refreshAllowedCIDR(ctx, organizationId, projectId, clusterId, appServiceId, allowedCIDRId)
	if err != errors.ErrNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Allowed CIDR",
			"Could not read App Service Allowed CIDR "+allowedCIDRId+": "+err.Error(),
		)
		return
	}

	if state.ExpiresAt.ValueString() != refreshedState.ExpiresAt.ValueString() {
		refreshedState.ExpiresAt = state.ExpiresAt
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update is not supported for App services allowed cidrs.
func (a *AppServiceCidr) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for App Services Allowed CIDRs.
	// App Services Allowed CIDRs can only be created, read and deleted.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/Allowed-CIDRs-(App-Services)
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

// Delete is used to delete an existing App Service Allowed CIDR.
func (a *AppServiceCidr) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceCIDR
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported.
	organizationId, projectId, clusterId, appServiceId, allowedCIDRId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Allowed CIDR",
			"Could not read App Service Allowed CIDR: "+err.Error(),
		)
		return
	}

	// Execute request to delete existing allowed CIDR.
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		allowedCIDRId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
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
			"Error Listing App Service Allowed CIDRs",
			"Could not read App Service Allowed CIDR "+allowedCIDRId+": "+errString,
		)
		return
	}
}

func (a *AppServiceCidr) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
