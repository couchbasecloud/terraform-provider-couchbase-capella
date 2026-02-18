package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ClusterStats{}
	_ datasource.DataSourceWithConfigure = &ClusterStats{}
)

// ClusterStats is the ClusterStats data source implementation.
type ClusterStats struct {
	*providerschema.Data
}

// NewClusterStats is a helper function to simplify the provider implementation.
func NewClusterStats() datasource.DataSource {
	return &ClusterStats{}
}

// Metadata returns the cluster stats data source type name.
func (d *ClusterStats) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_stats"
}

// Schema defines the schema for the ClusterStats data source.
func (d *ClusterStats) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterStatsSchema()
}

// Read refreshes the Terraform state with the latest data of cluster stats.
func (d *ClusterStats) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ClusterStats
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/stats", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster Stats",
			fmt.Sprintf(
				"Could not read cluster stats in organization %s, project %s and cluster %s, unexpected error: %s",
				organizationId, projectId, clusterId, api.ParseError(err),
			),
		)
		return
	}

	var statsResponse clusterapi.GetClusterStatsResponse
	err = json.Unmarshal(response.Body, &statsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling Cluster Stats",
			fmt.Sprintf("Could not parse cluster stats response: %s", err.Error()),
		)
		return
	}

	state.FreeMemoryInMb = types.Int64Value(statsResponse.FreeMemoryInMb)
	state.MaxReplicas = types.Int64Value(statsResponse.MaxReplicas)
	state.TotalMemoryInMb = types.Int64Value(statsResponse.TotalMemoryInMb)

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the cluster stats data source.
func (d *ClusterStats) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
