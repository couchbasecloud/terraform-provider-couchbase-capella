package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var sampleBucketBuilder = capellaschema.NewSchemaBuilder("sampleBucket")

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

	attrs["type"] = WithDescription(stringAttribute([]string{computed}), "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.")
	attrs["storage_backend"] = WithDescription(stringAttribute([]string{computed}), "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.")
	attrs["memory_allocation_in_mb"] = WithDescription(int64Attribute(computed), "The amount of memory to allocate for the bucket memory in MiB. The maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation.")
	attrs["bucket_conflict_resolution"] = WithDescription(stringAttribute([]string{computed}), "The means in which conflicts are resolved during replication. This field may be referred to as conflictResolution in the Couchbase documentation, and seqno and lww may be referred to as sequence Number and Timestamp respectively.")
	attrs["durability_level"] = WithDescription(stringAttribute([]string{computed}), "The durability level of the bucket. This setting defines the minimum durability level required for all writes to a Couchbase bucket. The available durability levels depend on the bucket type. For a Couchbase bucket, the options are: None, Replicate to Majority, Majority and Persist to Active, and Persist to Majority. For an Ephemeral bucket, the options are limited to: None and Replicate to Majority.")
	attrs["replicas"] = WithDescription(int64Attribute(computed), "The number of replicas for the bucket")
	attrs["flush"] = WithDescription(boolAttribute(computed), "Replaced by flushEnabled. Determines whether bucket flush is enabled. Set property to 'true' to be able to delete all items in this bucket using the /flush endpoint. Disable property to avoid inadvertent data loss by calling the the /flush endpoint.")
	attrs["time_to_live_in_seconds"] = WithDescription(int64Attribute(computed), "Specifies the time to live (TTL) value in seconds. This is the maximum time to live for items in the bucket. If specified as 0, TTL is disabled. This is a non-negative value.")
	attrs["eviction_policy"] = WithDescription(stringAttribute([]string{computed}), "The policy which Capella adopts to prevent data loss due to memory exhaustion. This may be also known as Ejection Policy in the Couchbase documentation. For Couchbase bucket, Eviction Policy is fullEviction by default. For Ephemeral buckets, Eviction Policy is a required field, and should be one of the following: noEviction, nruEviction")
	attrs["stats"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"item_count":         WithDescription(int64Attribute(computed), "Number of documents in the bucket."),
			"ops_per_second":     WithDescription(int64Attribute(computed), "Number of operations per second."),
			"disk_used_in_mib":   WithDescription(int64Attribute(computed), "The amount of disk used (in MiB)."),
			"memory_used_in_mib": WithDescription(int64Attribute(computed), "The amount of memory used (in MiB)."),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a sample bucket in Couchbase Capella. Sample buckets are pre-loaded with sample data. Different sample data options include,\"travel-sample\", \"gamesim-sample\", \"beer-sample\".",
		Attributes:          attrs,
	}
}
