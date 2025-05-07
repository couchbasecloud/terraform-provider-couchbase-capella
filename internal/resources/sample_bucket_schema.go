package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func SampleBucketSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to manage a sample bucket in Couchbase Capella. Sample buckets are pre-loaded with sample data - \"travel-sample\", \"gamesim-sample\", \"beer-sample\".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the bucket. This is the base64 encoding of the bucket name.",
			},
			"name":                       WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the sample dataset to be loaded. The name has to be one of the following sample datasets: \"travel-sample\", \"gamesim-sample\" or \"beer-sample\"."),
			"organization_id":            WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Capella organization."),
			"project_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"type":                       WithDescription(stringAttribute([]string{computed}), "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation."),
			"storage_backend":            WithDescription(stringAttribute([]string{computed}), "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation."),
			"memory_allocation_in_mb":    WithDescription(int64Attribute(computed), "The amount of memory to allocate for the bucket memory in MiB. The maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation."),
			"bucket_conflict_resolution": WithDescription(stringAttribute([]string{computed}), "The means in which conflicts are resolved during replication. This field may be referred to as conflictResolution in the Couchbase documentation, and seqno and lww may be referred to as sequence Number and Timestamp respectively."),
			"durability_level":           WithDescription(stringAttribute([]string{computed}), "This is the minimum level at which all writes to the Couchbase bucket must occur. The options for Durability level are as follows, according to the bucket type. For a Couchbase bucket: None, Replicate to Majority, Majority and Persist to Active, Persist to Majority. For an Ephemeral bucket: None, Replicate to Majority"),
			"replicas":                   WithDescription(int64Attribute(computed), "The number of replicas for the bucket"),
			"flush":                      WithDescription(boolAttribute(computed), "Replaced by flushEnabled. Determines whether bucket flush is enabled. Set property to true to be able to delete all items in this bucket using the /flush endpoint. Disable property to avoid inadvertent data loss by calling the the /flush endpoint."),
			"time_to_live_in_seconds":    WithDescription(int64Attribute(computed), "Specifies the time to live (TTL) value in seconds. This is the maximum time to live for items in the bucket. If specified as 0, TTL is disabled. This is a non-negative value."),
			"eviction_policy":            WithDescription(stringAttribute([]string{computed}), "The policy which Capella adopts to prevent data loss due to memory exhaustion. This may be also known as Ejection Policy in the Couchbase documentation. For Couchbase bucket, Eviction Policy is fullEviction by default. For Ephemeral buckets, Eviction Policy is a required field, and should be one of the following: noEviction, nruEviction"),
			"stats": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"item_count":         WithDescription(int64Attribute(computed), "Number of documents in the bucket."),
					"ops_per_second":     WithDescription(int64Attribute(computed), "Number of operations per second."),
					"disk_used_in_mib":   WithDescription(int64Attribute(computed), "The amount of disk used (in MiB)."),
					"memory_used_in_mib": WithDescription(int64Attribute(computed), "The amount of memory used (in MiB)."),
				},
			},
		},
	}
}
