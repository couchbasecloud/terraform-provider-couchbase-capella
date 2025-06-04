package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

var (
	_ resource.Resource                = (*AppServiceCidr)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceCidr)(nil)
	_ resource.ResourceWithImportState = (*AppServiceCidr)(nil)
)

type AppServiceCidr struct {
	*providerschema.Data
}

func NewAppServiceCidr() resource.Resource {
	return &AppServiceCidr{}
}

func (a *AppServiceCidr) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

func (a *AppServiceCidr) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	// TODO
}

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
			"", // TODO: Add error message detail
		)
		return
	}

	cidrResp := &api.AppServiceAllowedCIDRResponse{}
	err = json.Unmarshal(response.Body, &cidrResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling App Service CIDR Response",
			"", // TODO : Add error message detail
		)
		return
	}

	// TODO set computed attributes to null

	// TODO get latest status of the CIDR

	// TODO save state

}

// getAllowList is used to retrieve an existing allow list.
func (r *AllowList) getAllowList(ctx context.Context, organizationId, projectId, clusterId, appServiceId, allowListId string) (*api.GetAllowListResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
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
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	allowListResp := api.GetAllowListResponse{}
	err = json.Unmarshal(response.Body, &allowListResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &allowListResp, nil
}

func (a *AppServiceCidr) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AllowList
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allow list: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		allowListId    = IDs[providerschema.Id]
	)

	// refresh the existing allow list
	refreshedState, err := a.refreshAllowList(ctx, organizationId, projectId, clusterId, allowListId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allowListID "+allowListId+": "+errString,
		)
		return
	}

	if state.ExpiresAt != refreshedState.ExpiresAt {
		refreshedState.ExpiresAt = state.ExpiresAt
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *AppServiceCidr) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO return error
}

func (a *AppServiceCidr) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}

func (a *AppServiceCidr) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// TODO
}
