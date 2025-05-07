package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func FreeTierClusterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Couchbase Capella Free Tier Cluster resource",
		Attributes: map[string]schema.Attribute{
			"id":                            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the free-tier cluster."),
			"organization_id":               WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the organization."),
			"project_id":                    WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the project."),
			"name":                          WithDescription(stringAttribute([]string{required}, stringvalidator.LengthAtLeast(1)), "Name of the free-tier cluster."),
			"description":                   WithDescription(stringAttribute([]string{optional, computed}), "Description of the free-tier cluster."),
			"app_service_id":                WithDescription(stringAttribute([]string{computed}), "The GUID4 ID of the App Service."),
			"connection_string":             WithDescription(stringAttribute([]string{computed}), "The connection string of the free-tier cluster."),
			"current_state":                 WithDescription(stringAttribute([]string{computed}), "The current state of the free-tier cluster."),
			"cmek_id":                       WithDescription(stringAttribute([]string{computed}), "The customer-managed encryption key (CMEK) ID."),
			"etag":                          WithDescription(stringAttribute([]string{computed}), "The etag of the free-tier cluster, part of the response header"),
			"enable_private_dns_resolution": WithDescription(boolAttribute(computed), "Indicates if the private DNS resolution is enabled for the cluster."),
			"audit":                         computedAuditAttribute(),
			"support": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Support information for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"plan":     WithDescription(stringAttribute([]string{computed}), "Support plan for the free-tier cluster. Free tier plan is automatically assigned to free tier clusters."),
					"timezone": WithDescription(stringAttribute([]string{computed}), "The standard timezone for the cluster. Should be the TZ identifier. For example, 'ET'"),
				},
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The cloud provider details for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"type":   WithDescription(stringAttribute([]string{required}), "The cloud provider type. Should be one of 'aws', 'gcp', or 'azure'."),
					"region": WithDescription(stringAttribute([]string{required}), "The region for the cloud provider. Should be a valid region for the specified cloud provider. For example 'us-west-2'"),
					"cidr":   WithDescription(stringAttribute([]string{required}), "CIDR block for Cloud Provider."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"couchbase_server": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Couchbase Server details for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"version": WithDescription(stringAttribute([]string{computed}), "The version of Couchbase Server for the free-tier cluster."),
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
										"cpu": WithDescription(int64Attribute(computed), "The number of CPU cores for the node."),
										"ram": WithDescription(int64Attribute(computed), "The amount of RAM for the node."),
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
										"storage":       WithDescription(int64Attribute(computed), "Storage size of the disk."),
										"iops":          WithDescription(int64Attribute(computed), "Input/Output Operations Per Second (IOPS) for the disk."),
										"autoexpansion": WithDescription(boolAttribute(computed), "Indicates if auto-expansion is enabled for the disk."),
									},
								},
							},
						},
						"num_of_nodes": WithDescription(int64Attribute(computed), "The number of nodes in the service group."),
						"services":     WithDescription(stringAttribute([]string{computed}), "The services enabled for the service group. Should be a comma-separated list of services. For example, 'data,index,query'"),
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Availability zone details for the free-tier cluster. This is single az for the free-tier cluster.",
				Attributes: map[string]schema.Attribute{
					"type": WithDescription(stringAttribute([]string{computed}), "The availability zone type. Should be 'single' for the free-tier cluster."),
				},
			},
		},
	}

}
