package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceAdminUser)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceAdminUser)(nil)
)

// AppServiceAdminUser is the data source implementation for retrieving an App Service admin user.
type AppServiceAdminUser struct {
	*providerschema.Data
}

// NewAppServiceAdminUser is used in (p *capellaProvider) DataSources for building the provider.
func NewAppServiceAdminUser() datasource.DataSource {
	return &AppServiceAdminUser{}
}

// Metadata returns the App Service admin user data source type name.
func (a *AppServiceAdminUser) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_service_admin_user"
}

// Schema defines the schema for the App Service admin user data source.
func (a *AppServiceAdminUser) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAdminUserSchema()
}

// Read refreshes the Terraform state with the App Service admin user.
func (a *AppServiceAdminUser) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceAdminUserData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		a.HostURL,
		state.OrganizationId.ValueString(),
		state.ProjectId.ValueString(),
		state.ClusterId.ValueString(),
		state.AppServiceId.ValueString(),
		state.Id.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Admin User",
			fmt.Sprintf("Could not read App Service admin user %s: %s", state.Id.ValueString(), api.ParseError(err)),
		)
		return
	}

	var user apigen.AppServiceAdminUser
	if err := json.Unmarshal(response.Body, &user); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Admin User",
			"Could not unmarshal App Service admin user response: "+err.Error(),
		)
		return
	}

	one, diags := providerschema.NewOneAppServiceAdminUser(ctx, user)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Name = one.Name
	state.AccessAllEndpoints = one.AccessAllEndpoints
	state.Endpoints = one.Endpoints
	state.Audit = one.Audit

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure is used to configure the AppServiceAdminUser data source with the provider data.
func (a *AppServiceAdminUser) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	a.Data = data
}
