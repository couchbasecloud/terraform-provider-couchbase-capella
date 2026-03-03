package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
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

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing organization_id", fmt.Sprintf("Invalid organization_id: %s", err.Error()))
		return
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing project_id", fmt.Sprintf("Invalid project_id: %s", err.Error()))
		return
	}
	clusterUUID, err := uuid.Parse(clusterId)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing cluster_id", fmt.Sprintf("Invalid cluster_id: %s", err.Error()))
		return
	}

	response, err := d.ClientV2.GetClusterStatsWithResponse(ctx, orgUUID, projUUID, clusterUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster Stats",
			fmt.Sprintf("Could not read cluster stats in organization %s, project %s and cluster %s: %s: %s",
				organizationId, projectId, clusterId, errors.ErrExecutingRequest, err.Error()),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster Stats",
			fmt.Sprintf("Unexpected response while reading cluster stats: %s", string(response.Body)),
		)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster Stats",
			"API returned an empty response body.",
		)
		return
	}

	tflog.Info(ctx, "read cluster stats", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
	})

	state.FreeMemoryInMb = types.Int64Value(int64(response.JSON200.FreeMemoryInMb))
	state.MaxReplicas = types.Int64Value(int64(response.JSON200.MaxReplicas))
	state.TotalMemoryInMb = types.Int64Value(int64(response.JSON200.TotalMemoryInMb))

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
