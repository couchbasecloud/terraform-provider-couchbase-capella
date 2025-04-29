package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "ID of the free-tier bucket"},
			),

			"name": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "The name of the free-tier bucket"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"organization_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the capella teanant"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"project_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the Capella project"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"cluster_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the Capella cluster"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"type": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "The bucket type (couchbase or ephemeral)"},
			),
			"storage_backend": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "The bucket storage engine type (magma or couchstore)"},
			),
			"memory_allocation_in_mb": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Bucket size allocation in mb",
			},
			"bucket_conflict_resolution": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "Conflict-resolution mechanism of bucket"},
			),
			"durability_level": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "Durability of the bucket"},
			),
			"replicas": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of replicas for the data",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"flush": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Determines whether the flushing is enabled in the bucket",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"time_to_live_in_seconds": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Time-to-live (TTL) for items in the bucket, in seconds.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"eviction_policy": stringAttributeWithValueFields([]string{computed, useStateForUnknown},
				map[string]string{markdownDescription: "Eviction policy for the bucket."},
			),
			"stats": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"item_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Bucket item count",
					},
					"ops_per_second": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Bucket ops per second value",
					},
					"disk_used_in_mib": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Disk used in mib",
					},
					"memory_used_in_mib": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Memory used in mib",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Bucket stats",
			},
		},

		MarkdownDescription: "Manages free-tier bucket resource for a free-tier cluster",
	}
}
