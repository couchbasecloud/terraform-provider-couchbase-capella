package resources

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func BucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
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
				Optional: true,
			},
		},
	}
}
