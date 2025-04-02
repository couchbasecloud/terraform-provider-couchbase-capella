package datasources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
		Clusters: &Clusters{
			FreeTierClusterFilter: true,
		},
	}
}

// Metadata returns the cluster data source type name.
func (f *FreeTierClusters) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_free_tier_clusters"
}
