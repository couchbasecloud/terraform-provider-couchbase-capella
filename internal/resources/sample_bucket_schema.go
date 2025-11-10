package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var sampleBucketBuilder = capellaschema.NewSchemaBuilder("sampleBucket", "PostSampleBucket")

func SampleBucketSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", sampleBucketBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "name", sampleBucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "organization_id", sampleBucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", sampleBucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", sampleBucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "type", sampleBucketBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "storage_backend", sampleBucketBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "memory_allocation_in_mb", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "bucket_conflict_resolution", sampleBucketBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "durability_level", sampleBucketBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "replicas", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "flush", sampleBucketBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "time_to_live_in_seconds", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "eviction_policy", sampleBucketBuilder, stringAttribute([]string{computed}))

	statsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statsAttrs, "item_count", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "ops_per_second", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "disk_used_in_mib", sampleBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "memory_used_in_mib", sampleBucketBuilder, int64Attribute(computed))

	capellaschema.AddAttr(attrs, "stats", sampleBucketBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statsAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a sample bucket in Couchbase Capella. Sample buckets are pre-loaded with sample data. Different sample data options include,\"travel-sample\", \"gamesim-sample\", \"beer-sample\".",
		Attributes:          attrs,
	}
}
