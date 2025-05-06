package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierClusterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Couchbase Capella Free Tier Cluster resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the free-tier cluster",
			},
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The GUID4 ID of the organization.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The GUID4 ID of the project.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the free-tier cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the free-tier cluster.",
				Computed:            true,
			},
			"app_service_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the App Service.",
			},
			"enable_private_dns_resolution": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the private DNS resolution is enabled for the cluster.",
			},
			"connection_string": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The connection string of the free-tier cluster.",
			},
			"current_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current state of the free-tier cluster.",
			},
			"cmek_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The customer-managed encryption key (CMEK) ID.",
			},
			"audit": computedAuditAttribute(),
			"etag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The etag of the free-tier cluster, part of the response header",
			},
			"support": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Support information for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"plan": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Support plan for the free-tier cluster. Free tier plan is automatically assigned to free tier clusters.",
					},
					"timezone": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The standard timezone for the cluster. Should be the TZ identifier. For example, 'ET'",
					},
				},
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The cloud provider details for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The cloud provider type. Should be one of 'aws', 'gcp', or 'azure'.",
					},
					"region": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The region for the cloud provider. Should be a valid region for the specified cloud provider. For example 'us-west-2'",
					},
					"cidr": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "CIDR block for Cloud Provider.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"couchbase_server": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Couchbase Server details for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version of Couchbase Server for the free-tier cluster.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"service_groups": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Service groups for the free-tier cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Node details for the service group.",
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Compute details for the node",
									Attributes: map[string]schema.Attribute{
										"cpu": schema.Int64Attribute{
											Computed:            true,
											MarkdownDescription: "The number of CPU cores for the node.",
										},
										"ram": schema.Int64Attribute{
											Computed:            true,
											MarkdownDescription: "The amount of RAM for the node.",
										},
									},
								},
								"disk": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Disk details for the node",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The type of disk for the node.",
										},
										"storage": schema.Int64Attribute{
											Computed:            true,
											MarkdownDescription: "Storage size of the disk.",
										},
										"iops": schema.Int64Attribute{
											Computed:            true,
											MarkdownDescription: "Input/Output Operations Per Second (IOPS) for the disk.",
										},
										"autoexpansion": schema.BoolAttribute{
											Computed:            true,
											MarkdownDescription: "Indicates if auto-expansion is enabled for the disk.",
										},
									},
								},
							},
						},
						"num_of_nodes": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of nodes in the service group. This is 1 for the free-tier cluster.",
						},
						"services": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The services enabled for free-tier cluster.",
						},
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Availability zone details for the free-tier cluster. This is single az for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The availability zone type. Should be 'single' for the free-tier cluster.",
					},
				},
			},
		},
	}

}
