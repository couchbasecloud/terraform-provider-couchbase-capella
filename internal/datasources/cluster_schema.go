package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                            computedStringAttribute,
						"organization_id":               computedStringAttribute,
						"project_id":                    computedStringAttribute,
						"name":                          computedStringAttribute,
						"description":                   computedStringAttribute,
						"enable_private_dns_resolution": computedBoolAttribute,
						"connection_string":             computedStringAttribute,
						"cloud_provider": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type":   computedStringAttribute,
								"region": computedStringAttribute,
								"cidr":   computedStringAttribute,
							},
						},
						"couchbase_server": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"version": computedStringAttribute,
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
													"cpu": computedInt64Attribute,
													"ram": computedInt64Attribute,
												},
											},
											"disk": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"type":          computedStringAttribute,
													"storage":       computedInt64Attribute,
													"iops":          computedInt64Attribute,
													"autoexpansion": computedBoolAttribute,
												},
											},
										},
									},
									"num_of_nodes": computedInt64Attribute,
									"services":     computedListAttribute,
								},
							},
						},
						"availability": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type": computedStringAttribute,
							},
						},
						"support": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"plan":     computedStringAttribute,
								"timezone": computedStringAttribute,
							},
						},
						"current_state":  computedStringAttribute,
						"app_service_id": computedStringAttribute,
						"audit":          computedAuditAttribute,
					},
				},
			},
		},
	}
}
