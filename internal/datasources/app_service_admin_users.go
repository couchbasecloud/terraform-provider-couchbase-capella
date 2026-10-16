package datasources

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource                   = (*AppServiceAdminUsers)(nil)
	_ datasource.DataSourceWithConfigure      = (*AppServiceAdminUsers)(nil)
	_ datasource.DataSourceWithValidateConfig = (*AppServiceAdminUsers)(nil)
)

// AppServiceAdminUsers is the data source implementation for listing App Service admin users.
type AppServiceAdminUsers struct {
	*providerschema.Data
}

// NewAppServiceAdminUsers is used in (p *capellaProvider) DataSources for building the provider.
func NewAppServiceAdminUsers() datasource.DataSource {
	return &AppServiceAdminUsers{}
}

// Metadata returns the App Service admin users data source type name.
func (a *AppServiceAdminUsers) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_service_admin_users"
}

// Schema defines the schema for the App Service admin users data source.
func (a *AppServiceAdminUsers) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAdminUsersSchema()
}

// Read refreshes the Terraform state with the App Service admin users.
func (a *AppServiceAdminUsers) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceAdminUsers
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers",
		a.HostURL,
		state.OrganizationId.ValueString(),
		state.ProjectId.ValueString(),
		state.ClusterId.ValueString(),
		state.AppServiceId.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	// The list endpoint has no sortBy parameter.
	users, err := api.GetPaginated[[]apigen.AppServiceAdminUser](ctx, a.ClientV1, a.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Admin Users",
			fmt.Sprintf("Could not list admin users of App Service %s: %s", state.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	// The list API has no name query parameter, so the provider filters.
	var names []string
	if state.Filters != nil {
		diags = state.Filters.Values.ElementsAs(ctx, &names, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// A nil slice becomes a null list in state; keep it empty instead.
	state.Data = make([]providerschema.OneAppServiceAdminUser, 0, len(users))
	for _, user := range users {
		if len(names) > 0 && (user.Name == nil || !slices.Contains(names, *user.Name)) {
			continue
		}
		one, diags := providerschema.NewOneAppServiceAdminUser(ctx, user)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Data = append(state.Data, one)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure is used to configure the AppServiceAdminUsers data source with the provider data.
func (a *AppServiceAdminUsers) Configure(
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

// ValidateConfig checks that if 'name' or 'values' is set in the filter block, then both are set.
func (a *AppServiceAdminUsers) ValidateConfig(
	ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse,
) {
	var config providerschema.AppServiceAdminUsers
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Filters != nil && config.Filters.Name.IsNull() != config.Filters.Values.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Filters Configuration",
			"Both 'name' and 'values' in filter block must be configured.",
		)
	}
}
