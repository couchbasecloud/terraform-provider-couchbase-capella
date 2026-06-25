package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &FreeTierClusters{}
	_ datasource.DataSourceWithConfigure = &FreeTierClusters{}
)

type FreeTierClusters struct {
	*Clusters
}

func NewFreeTierClusters() datasource.DataSource {
	return &FreeTierClusters{
		Clusters: &Clusters{},
	}
}

// Metadata returns the cluster data source type name.
func (f *FreeTierClusters) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_free_tier_clusters"
}

// Read refreshes the Terraform state with the latest data of free-tier clusters.
func (f *FreeTierClusters) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Clusters
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading free tier clusters",
			"Could not read free tier clusters, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading free tier clusters",
			"Could not read free tier clusters, unexpected error: project ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier", f.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]clusterapi.GetClusterResponse](ctx, f.ClientV1, f.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Free Tier Clusters",
			fmt.Sprintf(
				"Could not read free tier clusters in organization %s and project %s, unexpected error: %s",
				organizationId, projectId, api.ParseError(err),
			),
		)
		return
	}

	state.Data = make([]providerschema.ClusterData, 0, len(response))
	for i := range response {
		cluster := response[i]
		audit := providerschema.NewCouchbaseAuditData(cluster.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella Free Tier Clusters",
				fmt.Sprintf("Could not read free tier clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, errors.ErrUnableToConvertAuditData),
			)
			return
		}

		newClusterData, err := providerschema.NewClusterData(&cluster, organizationId, projectId, auditObj)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Capella Free Tier Clusters",
				fmt.Sprintf("Could not read free tier clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.Error()),
			)
			return
		}
		state.Data = append(state.Data, *newClusterData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
