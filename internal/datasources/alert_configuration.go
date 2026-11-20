package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*OrganizationAlertConfiguration)(nil)
	_ datasource.DataSourceWithConfigure = (*OrganizationAlertConfiguration)(nil)
	_ datasource.DataSource              = (*ProjectAlertConfiguration)(nil)
	_ datasource.DataSourceWithConfigure = (*ProjectAlertConfiguration)(nil)
	_ datasource.DataSource              = (*ClusterAlertConfiguration)(nil)
	_ datasource.DataSourceWithConfigure = (*ClusterAlertConfiguration)(nil)
)

// OrganizationAlertConfiguration is the organization alert configuration data source implementation.
type OrganizationAlertConfiguration struct {
	*providerschema.Data
}

// NewOrganizationAlertConfiguration is a helper function to simplify the provider implementation.
func NewOrganizationAlertConfiguration() datasource.DataSource {
	return &OrganizationAlertConfiguration{}
}

// Metadata returns the organization alert configuration data source type name.
func (d *OrganizationAlertConfiguration) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_alert_configuration"
}

// Schema defines the schema for the organization alert configuration data source.
func (d *OrganizationAlertConfiguration) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OrganizationAlertConfigurationSchema()
}

// Read refreshes the Terraform state with the organization's alert configuration.
func (d *OrganizationAlertConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.OrganizationAlertConfiguration
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	values, diags := readAlertConfiguration(ctx, d.Data, state.OrganizationId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.AlertConfigurationValues = values

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the organization alert configuration data source.
func (d *OrganizationAlertConfiguration) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ProjectAlertConfiguration is the project alert configuration data source implementation.
type ProjectAlertConfiguration struct {
	*providerschema.Data
}

// NewProjectAlertConfiguration is a helper function to simplify the provider implementation.
func NewProjectAlertConfiguration() datasource.DataSource {
	return &ProjectAlertConfiguration{}
}

// Metadata returns the project alert configuration data source type name.
func (d *ProjectAlertConfiguration) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_alert_configuration"
}

// Schema defines the schema for the project alert configuration data source.
func (d *ProjectAlertConfiguration) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectAlertConfigurationSchema()
}

// Read refreshes the Terraform state with the project's alert configuration.
func (d *ProjectAlertConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ProjectAlertConfiguration
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := fmt.Sprintf("%s/projects/%s", state.OrganizationId.ValueString(), state.ProjectId.ValueString())
	values, diags := readAlertConfiguration(ctx, d.Data, path)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.AlertConfigurationValues = values

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the project alert configuration data source.
func (d *ProjectAlertConfiguration) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ClusterAlertConfiguration is the cluster alert configuration data source implementation.
type ClusterAlertConfiguration struct {
	*providerschema.Data
}

// NewClusterAlertConfiguration is a helper function to simplify the provider implementation.
func NewClusterAlertConfiguration() datasource.DataSource {
	return &ClusterAlertConfiguration{}
}

// Metadata returns the cluster alert configuration data source type name.
func (d *ClusterAlertConfiguration) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_alert_configuration"
}

// Schema defines the schema for the cluster alert configuration data source.
func (d *ClusterAlertConfiguration) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterAlertConfigurationSchema()
}

// Read refreshes the Terraform state with the cluster's alert configuration.
func (d *ClusterAlertConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ClusterAlertConfiguration
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := fmt.Sprintf("%s/projects/%s/clusters/%s",
		state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.ClusterId.ValueString())
	values, diags := readAlertConfiguration(ctx, d.Data, path)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.AlertConfigurationValues = values

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the cluster alert configuration data source.
func (d *ClusterAlertConfiguration) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// readAlertConfiguration fetches the alert configuration at path, which is relative to
// /v4/organizations. It is shared by the organization, project and cluster data sources.
func readAlertConfiguration(ctx context.Context, data *providerschema.Data, path string) (providerschema.AlertConfigurationValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration", data.HostURL, path)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(ctx, cfg, nil, data.Token, nil)
	if err != nil {
		diags.AddError(
			"Error Reading Capella Alert Configuration",
			fmt.Sprintf("Could not read the alert configuration at %s, unexpected error: %s", path, api.ParseError(err)),
		)
		return providerschema.AlertConfigurationValues{}, diags
	}

	var configResp alertsapi.GetAlertConfigurationResponse
	if err := json.Unmarshal(response.Body, &configResp); err != nil {
		diags.AddError(
			"Error Reading Capella Alert Configuration",
			fmt.Sprintf("Could not read the alert configuration at %s: %s: %s", path, errors.ErrUnmarshallingResponse, err),
		)
		return providerschema.AlertConfigurationValues{}, diags
	}

	return providerschema.NewAlertConfigurationValues(ctx, configResp)
}
