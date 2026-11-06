package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AlertConfigurationOverrides)(nil)
	_ datasource.DataSourceWithConfigure = (*AlertConfigurationOverrides)(nil)
)

// AlertConfigurationOverrides is the alert configuration overrides data source implementation.
type AlertConfigurationOverrides struct {
	*providerschema.Data
}

// NewAlertConfigurationOverrides is a helper function to simplify the provider implementation.
func NewAlertConfigurationOverrides() datasource.DataSource {
	return &AlertConfigurationOverrides{}
}

// Metadata returns the alert configuration overrides data source type name.
func (d *AlertConfigurationOverrides) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_configuration_overrides"
}

// Schema defines the schema for the alert configuration overrides data source.
func (d *AlertConfigurationOverrides) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AlertConfigurationOverridesSchema()
}

// Read refreshes the Terraform state with the alert configuration overrides in the organization.
func (d *AlertConfigurationOverrides) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AlertConfigurationOverrides
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration/overrides", d.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	// Overrides only accept updatedAt and scope as sort keys, so the server default
	// (updatedAt, most recent first) is used instead of api.SortById.
	overrides, err := api.GetPaginated[[]alertsapi.AlertConfigurationOverride](ctx, d.ClientV1, d.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Alert Configuration Overrides",
			fmt.Sprintf("Could not read alert configuration overrides in organization %s, unexpected error: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	state.Data = make([]providerschema.OneAlertConfigurationOverride, 0, len(overrides))
	for _, override := range overrides {
		model, diags := providerschema.NewOneAlertConfigurationOverride(ctx, override)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Data = append(state.Data, model)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the alert configuration overrides data source.
func (d *AlertConfigurationOverrides) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
