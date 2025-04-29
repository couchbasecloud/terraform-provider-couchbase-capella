package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func BucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "ID of the Capella bucket",
			},
			"name": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "The name of the bucket"},
			),
			"organization_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the capella teanant"},
			),
			"project_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the Capella project"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"cluster_id": stringAttributeWithValueFields([]string{required, requiresReplace},
				map[string]string{markdownDescription: "ID of the Capella cluster"},
				validator.String(stringvalidator.LengthAtLeast(1)),
			),
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default:             stringdefault.StaticString("couchbase"),
				MarkdownDescription: "The bucket type (couchbase or ephemeral)",
			},
			"storage_backend": stringAttributeWithValueFields([]string{optional, computed, requiresReplace, useStateForUnknown},
				map[string]string{markdownDescription: "The bucket storage engine type (Magma or Couchstore)"},
			),
			"memory_allocation_in_mb": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "bucket size allocation in mb",
			},
			"bucket_conflict_resolution": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default:             stringdefault.StaticString("seqno"),
				MarkdownDescription: "Conflict-resolution mechanism of bucket",
			},
			"durability_level": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("none"),
				MarkdownDescription: "Durability of the bucket",
			},
			"replicas": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Number of replicas for the data",
			},
			"flush": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Determines whether the flushing is enabled in the bucket",
			},
			"time_to_live_in_seconds": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Time-to-live (TTL) for items in the bucket, in seconds.",
			},
			"eviction_policy": stringAttributeWithValueFields([]string{optional, computed, requiresReplace, useStateForUnknown},
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
	}
}
