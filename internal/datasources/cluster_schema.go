package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Retrieves the details of Couchbase Capella clusters.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the organization that the cluster and project belongs to.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the project that the cluster belongs to.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of clusters in the project.",
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
							MarkdownDescription: "Description of the cluster (up to 1024 characters).",
						},
						"enable_private_dns_resolution": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address.",
						},
						"connection_string": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The connection string to use to connect to the cluster.",
						},
						"cloud_provider": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The cloud provider where the cluster is hosted.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Cloud provider type. Currently supported values are AWS, GCP and Azure. Note: For singleNode cluster, only AWS type cloud provider is allowed.",
								},
								"region": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The region where the cluster is hosted.",
								},
								"cidr": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "CIDR block for Cloud Provider.",
								},
							},
						},
						"couchbase_server": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration for the Couchbase Server running on the cluster.",
							Attributes: map[string]schema.Attribute{
								"version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Version of the Couchbase Server installed in the cluster.",
								},
							},
						},
						"service_groups": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The couchbase service groups running in the cluster.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"node": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Node configuration for this service group.",
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Computed:            true,
												MarkdownDescription: "Compute resources configuration for the nodes.",
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
														MarkdownDescription: "The size of the disk in GB. For AWS: >= 50, for Azure (Only required for Ultra Disk types.): >= 64, for GCP: >= 50.",
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
										MarkdownDescription: "The number of nodes in this service group. The minimum number of nodes for the cluster can be 3 and maximum can be 27 nodes.",
									},
									"services": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "The couchbase services running on the nodes. The allowed services for singleNode cluster are one or all of - data, index, query and search.",
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
									MarkdownDescription: "The availability type of the cluster. Can be 'single' for Single Node or 'multi' for Multi Node.",
								},
							},
						},
						"support": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Support configuration for the cluster.",
							Attributes: map[string]schema.Attribute{
								"plan": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Plan type, either 'Basic', 'Developer Pro', or 'Enterprise'. Plan type allowed for singleNode cluster is either 'Basic', or 'Developer Pro'. In case of 'Basic' plan timezone field value is ignored.",
								},
								"timezone": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The standard timezone for the cluster. Should be the TZ identifier - ET, GMT, IST, PT.",
								},
							},
						},
						"current_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current state of the cluster. Example- deploying, scaling, destroying, peering,etc.",
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
