package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierBucketSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets of your free tier operational cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the Free Tier bucket.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "Name of the free tier bucket."),

			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the organization."),

			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the project."),

			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), " The GUID4 ID of the cluster."),

			"type":            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The bucket type, which will be either Couchbase or Ephemeral."),
			"storage_backend": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The bucket storage engine type. The type will be either Magma or Couchstore."),
			"memory_allocation_in_mb": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Bucket size allocation in MB.",
			},
			"bucket_conflict_resolution": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "Conflict-resolution mechanism of the bucket. The method used to resolve conflicts during replication. In Couchbase documentation, this field may be referred to as conflictResolution, with seqno and lww representing Sequence Number and Timestamp, respectively. For more details, see [Conflict Resolution](https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution)."),
			"durability_level":           WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The durability level of the bucket. This setting defines the minimum durability level required for all writes to a Couchbase bucket. The available durability levels depend on the bucket type. For a Couchbase bucket, the options are: None, Replicate to Majority, Majority and Persist to Active, and Persist to Majority. For an Ephemeral bucket, the options are limited to: None and Replicate to Majority."),
			"replicas":                   WithDescription(int64Attribute(computed, useStateForUnknown), "Number of replicas for the data."),
			"flush":                      WithDescription(boolAttribute(computed, useStateForUnknown), "Determines whether flushing is enabled in the bucket."),
			"time_to_live_in_seconds":    WithDescription(int64Attribute(computed, useStateForUnknown), "Time-to-live (TTL) for items in the bucket, in seconds."),
			"eviction_policy":            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The eviction policy for the bucket. This is the policy which Capella adopts to prevent data loss due to memory exhaustion. This may be also known as Ejection Policy in the Couchbase documentation."),
			"stats": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket stats",
				Attributes: map[string]schema.Attribute{
					"item_count":         WithDescription(int64Attribute(computed), "Bucket item count."),
					"ops_per_second":     WithDescription(int64Attribute(computed), "The value for bucket operations per second."),
					"disk_used_in_mib":   WithDescription(int64Attribute(computed), "Disk used in MiB."),
					"memory_used_in_mib": WithDescription(int64Attribute(computed), "Memory used in MiB."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
