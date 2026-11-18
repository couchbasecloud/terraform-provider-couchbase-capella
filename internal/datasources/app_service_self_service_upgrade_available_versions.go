package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_self_service_upgrade"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceSelfServiceUpgradeAvailableVersions)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceSelfServiceUpgradeAvailableVersions)(nil)
)

// AppServiceSelfServiceUpgradeAvailableVersions is the data source implementation for listing the
// versions an App Service can be upgraded to.
type AppServiceSelfServiceUpgradeAvailableVersions struct {
	*providerschema.Data
}

// NewAppServiceSelfServiceUpgradeAvailableVersions is used in (p *capellaProvider) DataSources for building the provider.
func NewAppServiceSelfServiceUpgradeAvailableVersions() datasource.DataSource {
	return &AppServiceSelfServiceUpgradeAvailableVersions{}
}

// Metadata returns the App Service self-service upgrade available versions data source type name.
func (d *AppServiceSelfServiceUpgradeAvailableVersions) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_self_service_upgrade_available_versions"
}

// Schema defines the schema for the App Service self-service upgrade available versions data source.
func (d *AppServiceSelfServiceUpgradeAvailableVersions) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceSelfServiceUpgradeAvailableVersionsSchema()
}

// Read refreshes the Terraform state with the versions the App Service can be upgraded to.
func (d *AppServiceSelfServiceUpgradeAvailableVersions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceSelfServiceUpgradeAvailableVersions
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appServiceId := state.AppServiceId.ValueString()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/selfServiceUpgrades/availableVersions",
		d.HostURL, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString(), appServiceId)

	// The list is not paginated: it has no cursor and takes no page parameters.
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Self-Service Upgrade Available Versions",
			fmt.Sprintf("Could not read the available versions of app service %s, unexpected error: %s", appServiceId, api.ParseError(err)),
		)
		return
	}

	var versions appservice_self_service_upgrade.ListAvailableVersionsResponse
	if err = json.Unmarshal(response.Body, &versions); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Self-Service Upgrade Available Versions",
			fmt.Sprintf("Could not read the available versions of app service %s, unexpected error: %s: %s", appServiceId, errors.ErrUnmarshallingResponse, err),
		)
		return
	}

	state.Data = make([]providerschema.AppServiceSelfServiceUpgradeAvailableVersion, 0, len(versions.Data))
	for _, version := range versions.Data {
		data, diags := providerschema.NewAppServiceSelfServiceUpgradeAvailableVersion(ctx, version)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Data = append(state.Data, data)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the App Service self-service upgrade available versions data source.
func (d *AppServiceSelfServiceUpgradeAvailableVersions) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
