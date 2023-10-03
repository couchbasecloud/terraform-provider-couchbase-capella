package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// BucketSchema defines the schema for the terraform provider resource - "Bucket".
// This terraform resource directly maps to the bucket created in a Capella cluster.
// Bucket resource supports Create, Destroy, Read, Import and List operations.
func BucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"audit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("couchbase"),
				Computed: true,
			},
			"storage_backend": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("couchstore"),
				Computed: true,
			},
			"memory_allocation_in_mb": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(100),
				Computed: true,
			},
			"bucket_conflict_resolution": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("seqno"),
				Computed: true,
			},
			"durability_level": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("none"),
				Computed: true,
			},
			"replicas": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(1),
				Computed: true,
			},
			"flush": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(false),
				Computed: true,
			},
			"time_to_live_in_seconds": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(0),
				Computed: true,
			},
			"eviction_policy": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("fullEviction"),
				Computed: true,
			},
			"stats": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"item_count": schema.Int64Attribute{
						Computed: true,
					},
					"ops_per_second": schema.Int64Attribute{
						Computed: true,
					},
					"disk_used_in_mib": schema.Int64Attribute{
						Computed: true,
					},
					"memory_used_in_mib": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}
}