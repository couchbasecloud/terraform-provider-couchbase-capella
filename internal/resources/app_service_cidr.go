package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	resp.TypeName = req.ProviderTypeName + "_project"
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

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		a.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
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

	// TODO get latest status of the CIDR

	// TODO save state

}

func (a *AppServiceCidr) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO
}

func (a *AppServiceCidr) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for app services cidr.
	// This resource can only be created, read and deleted.
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

func (a *AppServiceCidr) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}
