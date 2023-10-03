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
				Required: true,
			},
			"storage_backend": schema.StringAttribute{
				Required: true,
			},
			"memory_allocationinmb": schema.Int64Attribute{
				Required: true,
			},
			"conflict_resolution": schema.StringAttribute{
				Required: true,
			},
			"durability_level": schema.StringAttribute{
				Required: true,
			},
			"replicas": schema.Int64Attribute{
				Required: true,
			},
			"flush": schema.BoolAttribute{
				Required: true,
			},
			"ttl": schema.Int64Attribute{
				Required: true,
			},
			"eviction_policy": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
