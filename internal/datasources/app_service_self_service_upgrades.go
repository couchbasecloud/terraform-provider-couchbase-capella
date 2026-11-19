package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceSelfServiceUpgrades)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceSelfServiceUpgrades)(nil)
)

// AppServiceSelfServiceUpgrades is the data source implementation for listing the self-service upgrades of an App Service.
type AppServiceSelfServiceUpgrades struct {
	*providerschema.Data
}

// NewAppServiceSelfServiceUpgrades is used in (p *capellaProvider) DataSources for building the provider.
func NewAppServiceSelfServiceUpgrades() datasource.DataSource {
	return &AppServiceSelfServiceUpgrades{}
}

// Metadata returns the App Service self-service upgrades data source type name.
func (d *AppServiceSelfServiceUpgrades) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_self_service_upgrades"
}

// Schema defines the schema for the App Service self-service upgrades data source.
func (d *AppServiceSelfServiceUpgrades) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceSelfServiceUpgradesSchema()
}

// Read refreshes the Terraform state with every self-service upgrade of the App Service.
func (d *AppServiceSelfServiceUpgrades) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceSelfServiceUpgrades
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appServiceId := state.AppServiceId.ValueString()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/selfServiceUpgrades",
		d.HostURL, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString(), appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	// The list order is fixed by the API, which accepts no sortBy.
	response, err := api.GetPaginated[[]appservice_maintenance.GetMaintenanceResponse](ctx, d.ClientV1, d.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Self-Service Upgrades",
			fmt.Sprintf("Could not read self-service upgrades of app service %s, unexpected error: %s", appServiceId, api.ParseError(err)),
		)
		return
	}

	state.Data = make([]providerschema.AppServiceMaintenanceData, 0, len(response))
	for i := range response {
		data, diags := providerschema.NewAppServiceMaintenanceData(ctx, &response[i])
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella App Service Self-Service Upgrades",
				fmt.Sprintf("Could not read self-service upgrades of app service %s, unexpected error: %s", appServiceId, errors.ErrUnableToConvertAuditData),
			)
			return
		}
		state.Data = append(state.Data, data)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the App Service self-service upgrades data source.
func (d *AppServiceSelfServiceUpgrades) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Data = data
}
