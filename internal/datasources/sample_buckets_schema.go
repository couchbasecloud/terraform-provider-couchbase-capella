package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var sampleBucketsBuilder = capellaschema.NewSchemaBuilder("sampleBuckets")

// SampleBucketsSchema returns the schema for the SampleBuckets data source.
func SampleBucketsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", sampleBucketsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", sampleBucketsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", sampleBucketsBuilder, requiredString())

	statsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statsAttrs, "item_count", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "ops_per_second", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "disk_used_in_mib", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(statsAttrs, "memory_used_in_mib", sampleBucketsBuilder, computedInt64())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "type", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "storage_backend", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "memory_allocation_in_mb", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "bucket_conflict_resolution", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "durability_level", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "replicas", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "flush", sampleBucketsBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "time_to_live_in_seconds", sampleBucketsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "eviction_policy", sampleBucketsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "stats", sampleBucketsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statsAttrs,
	})

	capellaschema.AddAttr(attrs, "data", sampleBucketsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve sample buckets in an operational cluster.",
		Attributes:          attrs,
	}
}
