package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterStatsBuilder = capellaschema.NewSchemaBuilder("cluster_stats")

func ClusterStatsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterStatsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", clusterStatsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", clusterStatsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "free_memory_in_mb", clusterStatsBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "max_replicas", clusterStatsBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "total_memory_in_mb", clusterStatsBuilder, computedInt64())

	return schema.Schema{
		MarkdownDescription: "The data source retrieves the statistics of a Couchbase Capella cluster.",
		Attributes:          attrs,
	}
}
