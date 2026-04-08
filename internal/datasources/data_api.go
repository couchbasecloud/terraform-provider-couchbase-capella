package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataAPI{}
	_ datasource.DataSourceWithConfigure = &DataAPI{}
)

// DataAPI is the Data API data source implementation.
type DataAPI struct {
	*providerschema.Data
}

// NewDataAPI is a helper function to simplify the provider implementation.
func NewDataAPI() datasource.DataSource {
	return &DataAPI{}
}

// Metadata returns the data api data source type name.
func (d *DataAPI) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_api"
}

// Schema defines the schema for the data api data source.
func (d *DataAPI) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataAPISchema()
}

// Read refreshes the Terraform state with the latest data of the Data API status.
func (d *DataAPI) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.DataAPI
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Data API status",
			"Could not read Data API status for cluster with id "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Data API Status",
			"Could not read Data API status in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	dataAPIResp := api.GetDataAPIStatusResponse{}
	err = json.Unmarshal(response.Body, &dataAPIResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Data API status",
			"Could not read Data API status, unexpected error: "+err.Error(),
		)
		return
	}

	newState := providerschema.NewDataAPI(&dataAPIResp, organizationId, projectId, clusterId)

	// Set state
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data api data source.
func (d *DataAPI) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.Data = data
}
