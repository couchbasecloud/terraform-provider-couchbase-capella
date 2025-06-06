package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BucketSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets for an operational cluster.",
		Attributes: map[string]schema.Attribute{
			"id": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The ID of the bucket."),

			"name":            WithDescription(stringAttribute([]string{required, requiresReplace}), "Name of the Capella bucket"),
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization"),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"type": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The bucket type (couchbase or ephemeral).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString("couchbase"),
			},
			"storage_backend": WithDescription(stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}), "The bucket storage engine type (Magma or Couchstore)."),
			"memory_allocation_in_mb": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Bucket size allocation in MB.",
			},
			"bucket_conflict_resolution": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("seqno"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Conflict-resolution mechanism of the bucket.",
			},

			"durability_level": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				MarkdownDescription: "Durability level of the bucket.",
			},

			"replicas": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Number of replicas for the data.",
			},
			"flush": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Determines whether the flushing is enabled in the bucket.",
			},
			"time_to_live_in_seconds": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Time-to-live (TTL) for items in the bucket, in seconds.",
			},
			"eviction_policy": WithDescription(stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}), "Eviction policy for the bucket."),
			"stats": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket stats",
				Attributes: map[string]schema.Attribute{
					"item_count":         WithDescription(int64Attribute(computed), "Bucket item count."),
					"ops_per_second":     WithDescription(int64Attribute(computed), "The value for bucket operations per second."),
					"disk_used_in_mib":   WithDescription(int64Attribute(computed), "Disk used in MiB. "),
					"memory_used_in_mib": WithDescription(int64Attribute(computed), "Memory used in MiB."),
				},
			},
		},
	}
}
