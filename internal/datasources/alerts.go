package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*Alerts)(nil)
	_ datasource.DataSourceWithConfigure = (*Alerts)(nil)
)

// Alerts is the alerts data source implementation.
type Alerts struct {
	*providerschema.Data
}

// NewAlerts is a helper function to simplify the provider implementation.
func NewAlerts() datasource.DataSource {
	return &Alerts{}
}

// Metadata returns the alerts data source type name.
func (d *Alerts) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alerts"
}

// Schema defines the schema for the alerts data source.
func (d *Alerts) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AlertsSchema()
}

// Read refreshes the Terraform state with the alerts in the organization that match the filters.
func (d *Alerts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Alerts
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()

	queryParams, err := alertsQueryParams(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Alerts",
			fmt.Sprintf("Could not read the filters for alerts in organization %s: %s", organizationId, err),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/alerts", d.HostURL, organizationId)
	if len(queryParams) > 0 {
		url += BuildQueryParams(queryParams)
	}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	// Alerts only accept status, name, severity and eventTime as sort keys, so the
	// server default (eventTime, most recent first) is used instead of api.SortById.
	alerts, err := api.GetPaginated[[]alertsapi.Alert](ctx, d.ClientV1, d.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Alerts",
			fmt.Sprintf("Could not read alerts in organization %s, unexpected error: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	state.Data = make([]providerschema.OneAlert, 0, len(alerts))
	for _, alert := range alerts {
		model, diags := providerschema.NewOneAlert(alert)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Data = append(state.Data, model)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// alertsQueryParams maps the optional filters to the list alerts endpoint's query parameters.
// Multi-valued filters are sent as repeated parameters, which the endpoint accepts.
func alertsQueryParams(ctx context.Context, state providerschema.Alerts) (map[string][]string, error) {
	params := make(map[string][]string)

	sets := []struct {
		key   string
		value types.Set
	}{
		{"status", state.Status},
		{"clusterIds", state.ClusterIds},
		{"projectIds", state.ProjectIds},
		{"severityLevels", state.SeverityLevels},
	}
	for _, s := range sets {
		if s.value.IsNull() || s.value.IsUnknown() {
			continue
		}
		values, err := ConvertToList(ctx, s.value)
		if err != nil {
			return nil, err
		}
		params[s.key] = values
	}

	if !state.From.IsNull() && !state.From.IsUnknown() {
		params["from"] = []string{state.From.ValueString()}
	}
	if !state.To.IsNull() && !state.To.IsUnknown() {
		params["to"] = []string{state.To.ValueString()}
	}

	return params, nil
}

// Configure adds the provider configured client to the alerts data source.
func (d *Alerts) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
