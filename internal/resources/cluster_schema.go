package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterSchema() schema.Schema {
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"region": schema.StringAttribute{
						Required: true,
					},
					"cidr": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"couchbase_server": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"service_groups": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Required: true,
									Attributes: map[string]schema.Attribute{
										"cpu": schema.Int64Attribute{
											Required: true,
										},
										"ram": schema.Int64Attribute{
											Required: true,
										},
									},
								},
								"disk": schema.SingleNestedAttribute{
									Required: true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Required: true,
										},
										"storage": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
										"iops": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
									},
								},
							},
						},
						"num_of_nodes": schema.Int64Attribute{
							Required: true,
						},
						"services": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
						},
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"support": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"plan": schema.StringAttribute{
						Required: true,
					},
					"timezone": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"current_state": schema.StringAttribute{
				Computed: true,
			},
			"app_service_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
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
			"if_match": schema.StringAttribute{
				Optional: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
