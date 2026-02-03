package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var bucketsBuilder = capellaschema.NewSchemaBuilder("buckets")

// BucketsSchema returns the schema for the Buckets data source.
func BucketsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", bucketsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", bucketsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", bucketsBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "type", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "storage_backend", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "memory_allocation_in_mb", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "bucket_conflict_resolution", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "durability_level", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "replicas", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "flush", bucketsBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "time_to_live_in_seconds", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "eviction_policy", bucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "vbuckets", bucketsBuilder, computedInt64())

	statsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statsAttrs, "item_count", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "ops_per_second", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "disk_used_in_mib", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "memory_used_in_mib", bucketsBuilder, computedInt64())

	capellaschema.AddAttr(dataAttrs, "stats", bucketsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statsAttrs,
	})

	capellaschema.AddAttr(attrs, "data", bucketsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	clusterStatsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(clusterStatsAttrs, "free_memory_in_mb", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(clusterStatsAttrs, "max_replicas", bucketsBuilder, computedInt64())
	capellaschema.AddAttr(clusterStatsAttrs, "total_memory_in_mb", bucketsBuilder, computedInt64())

	capellaschema.AddAttr(attrs, "cluster_stats", bucketsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: clusterStatsAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve buckets for a cluster.",
		Attributes:          attrs,
	}
}
