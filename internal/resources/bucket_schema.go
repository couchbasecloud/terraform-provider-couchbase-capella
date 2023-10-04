package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
			},
			"storage_backend": schema.StringAttribute{
				Optional: true,
			},
			"memory_allocationinmb": schema.Int64Attribute{
				Optional: true,
			},
			"conflict_resolution": schema.StringAttribute{
				Optional: true,
			},
			"durability_level": schema.StringAttribute{
				Optional: true,
			},
			"replicas": schema.Int64Attribute{
				Optional: true,
			},
			"flush": schema.BoolAttribute{
				Optional: true,
			},
			"ttl": schema.Int64Attribute{
				Optional: true,
			},
			"eviction_policy": schema.StringAttribute{
				Computed: true,
				Optional: true,
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
					"disk_used_inmib": schema.Int64Attribute{
						Computed: true,
					},
					"memory_used_inmib": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}
}
