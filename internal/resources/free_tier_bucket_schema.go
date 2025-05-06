package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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

			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the free-tier bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the Capella organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the Capella project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the Capella cluster",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The bucket type (couchbase or ephemeral)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage_backend": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The bucket storage engine type (magma or couchstore)",
			},
			"memory_allocation_in_mb": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Bucket size allocation in mb",
			},
			"bucket_conflict_resolution": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Conflict-resolution mechanism of bucket",
			},
			"durability_level": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Durability level of the bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
			"eviction_policy": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Eviction policy for the bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
