package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"cloud_provider": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
								"region": schema.StringAttribute{
									Computed: true,
								},
								"cidr": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"couchbase_server": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"version": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"service_groups": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"node": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"cpu": schema.Int64Attribute{
														Computed: true,
													},
													"ram": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
											"disk": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"type": schema.StringAttribute{
														Computed: true,
													},
													"storage": schema.Int64Attribute{
														Computed: true,
													},
													"iops": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
										},
									},
									"num_of_nodes": schema.Int64Attribute{
										Computed: true,
									},
									"services": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"availability": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"support": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"plan": schema.StringAttribute{
									Computed: true,
								},
								"timezone": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"current_state": schema.StringAttribute{
							Computed: true,
						},
						"app_service_id": schema.StringAttribute{
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
					},
				},
			},
		},
	}
}
