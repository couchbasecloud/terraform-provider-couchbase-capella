package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source retrieves the details of Couchbase Capella clusters.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the clusters in the project.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the Capella cluster.",
						},
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the cluster (up to 256 characters).",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description of the cluster (up to 1024 characters).",
						},
						"enable_private_dns_resolution": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address.",
						},
						"connection_string": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The connection string used to connect to the cluster.",
						},
						"cloud_provider": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The Cloud Service Provider where the cluster is hosted.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The cluster's Cloud Service Provider (CSP). The supported CSPs include AWS, GCP, and Azure. Only AWS is supported as the cloud provider for Single Node clusters.",
								},
								"region": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The region where the cluster is hosted.",
								},
								"cidr": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "CIDR block for the Cloud Service Provider.",
								},
							},
						},
						"couchbase_server": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the Couchbase Server configuration running on the cluster.",
							Attributes: map[string]schema.Attribute{
								"version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The Couchbase Server version installed in the cluster.",
								},
							},
						},
						"service_groups": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The Couchbase service groups running in the cluster.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"node": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The node configuration for this service group.",
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Computed:            true,
												MarkdownDescription: "The compute resource configuration for the nodes.",
												Attributes: map[string]schema.Attribute{
													"cpu": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "The number of CPU cores for the node.",
													},
													"ram": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "The amount of RAM units in GB for the node.",
													},
												},
											},
											"disk": schema.SingleNestedAttribute{
												Computed:            true,
												MarkdownDescription: "The disk configuration for the nodes.",
												Attributes: map[string]schema.Attribute{
													"type": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "The type of disk to use. For AWS: gp3 or io2, for Azure: Premium (P6, P10, P15, P20, P30, P40, P50, P60) or Ultra, for GCP: pd-ssd.",
													},
													"storage": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "The size of the disk in GB. For AWS: >= 50, for Azure (only required for Ultra Disk types.): >= 64, for GCP: >= 50.",
													},
													"iops": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "The number of IOPS for the disk. Only applicable for certain disk types for AWS and Azure.",
													},
													"autoexpansion": schema.BoolAttribute{
														Computed:            true,
														MarkdownDescription: "Whether to enable automatic disk expansion. Only applicable for Azure.",
													},
												},
											},
										},
									},
									"num_of_nodes": schema.Int64Attribute{
										Computed:            true,
										MarkdownDescription: "The number of nodes in this Service Group. A cluster can have a minimum of 3 nodes and a maximum of 27 nodes.",
									},
									"services": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "The Services running on the nodes. The available Services include 'Data', 'Index', 'Query', and 'Search'.",
									},
								},
							},
						},
						"availability": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Availability configuration for the cluster.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Specifies the availability type of the cluster, either 'single' for Single Node or 'multi' for Multi Node.",
								},
							},
						},
						"support": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Support configuration for the cluster.",
							Attributes: map[string]schema.Attribute{
								"plan": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Plan type, either 'Basic', 'Developer Pro', or 'Enterprise'. Plan type allowed for Single Node cluster is either 'Basic', or 'Developer Pro'. In case of 'Basic' plan timezone field value is ignored.",
								},
								"timezone": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The standard timezone for the cluster. Should be the TZ identifier - ET, GMT, IST, PT.",
								},
							},
						},
						"current_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current cluster status. The cluster statuses are 'deploying', 'scaling', 'destroying', 'peering', 'healthy', and more.",
						},
						"app_service_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the App Service associated with this cluster.",
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}
