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
			},
			"storage_backend": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("couchstore"),
			},
			"memory_allocation_in_mb": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(100),
			},
			"bucket_conflict_resolution": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("seqno"),
			},
			"durability_level": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("none"),
			},
			"replicas": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(1),
			},
			"flush": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"time_to_live_in_seconds": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(0),
			},
			"eviction_policy": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("fullEviction"),
			},
			"stats": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"item_count":         schema.Int64Attribute{},
					"ops_per_second":     schema.Int64Attribute{},
					"disk_used_in_mib":   schema.Int64Attribute{},
					"memory_used_in_mib": schema.Int64Attribute{},
				},
			},
		},
	}
}
