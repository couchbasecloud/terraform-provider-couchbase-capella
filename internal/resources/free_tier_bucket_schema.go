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
		MarkdownDescription: "Manages free-tier bucket resource for a free-tier cluster",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the free-tier bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "Name of the free-tier bucket"),

			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "ID of the Capella organization"),

			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "ID of the Capella project"),

			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "ID of the Capella cluster"),

			"type":            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The bucket type (couchbase or ephemeral)"),
			"storage_backend": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The bucket storage engine type (magma or couchstore)"),
			"memory_allocation_in_mb": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Bucket size allocation in mb",
			},
			"bucket_conflict_resolution": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "Conflict-resolution mechanism of bucket"),
			"durability_level":           WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "Durability level of the bucket"),
			"replicas":                   WithDescription(int64Attribute(computed, useStateForUnknown), "Number of replicas for the data"),
			"flush":                      WithDescription(boolAttribute(computed, useStateForUnknown), "Determines whether the flushing is enabled in the bucket"),
			"time_to_live_in_seconds":    WithDescription(int64Attribute(computed, useStateForUnknown), "Time-to-live (TTL) for items in the bucket, in seconds"),
			"eviction_policy":            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "Eviction policy for the bucket"),
			"stats": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket stats",
				Attributes: map[string]schema.Attribute{
					"item_count":         WithDescription(int64Attribute(computed), "Bucket item count"),
					"ops_per_second":     WithDescription(int64Attribute(computed), "Bucket ops per second value"),
					"disk_used_in_mib":   WithDescription(int64Attribute(computed), "Disk used in mib"),
					"memory_used_in_mib": WithDescription(int64Attribute(computed), "Memory used in mib"),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
