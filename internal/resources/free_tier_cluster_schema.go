package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierClusterBuilder = capellaschema.NewSchemaBuilder("freeTierCluster")

func FreeTierClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", freeTierClusterBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", freeTierClusterBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", freeTierClusterBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "etag", freeTierClusterBuilder, stringAttribute([]string{computed}))

	attrs["name"] = WithDescription(stringAttribute([]string{required}, validator.String(stringvalidator.LengthAtLeast(1))), "Name of the free tier cluster.")
	attrs["description"] = WithDescription(stringAttribute([]string{optional, computed}), "Description of the free tier cluster.")
	attrs["app_service_id"] = WithDescription(stringAttribute([]string{computed}), "The GUID4 ID of the App Service.")
	attrs["connection_string"] = WithDescription(stringAttribute([]string{computed}), "The connection string of the free tier cluster.")
	attrs["current_state"] = WithDescription(stringAttribute([]string{computed}), "The current state of the free tier cluster.")
	attrs["cmek_id"] = WithDescription(stringAttribute([]string{computed}), "The customer-managed encryption key (CMEK) ID.")
	attrs["enable_private_dns_resolution"] = WithDescription(boolAttribute(computed), "Indicates if the private DNS resolution is enabled for the cluster.")
	capellaschema.AddAttr(attrs, "audit", freeTierClusterBuilder, computedAuditAttribute())
	attrs["support"] = schema.SingleNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The Support information for the free tier cluster.",
		Attributes: map[string]schema.Attribute{
			"plan":     WithDescription(stringAttribute([]string{computed}), "The Support plan for the free tier cluster. The free tier plan is automatically assigned to free tier clusters."),
			"timezone": WithDescription(stringAttribute([]string{computed}), "The standard timezone for the cluster. Should be the TZ identifier. For example, 'ET'."),
		},
	}
	attrs["cloud_provider"] = schema.SingleNestedAttribute{
		Required:            true,
		MarkdownDescription: "The Cloud Service Provider details for the free tier cluster.",
		Attributes: map[string]schema.Attribute{
			"type":   WithDescription(stringAttribute([]string{required}), "The Cloud Service Provider type. Currently supporting AWS, GCP and Azure."),
			"region": WithDescription(stringAttribute([]string{required}), "The region for the Cloud Service Provider. This should be a valid region for the specified Cloud Service Provider. For example 'us-west-2'."),
			"cidr":   WithDescription(stringAttribute([]string{required}), "CIDR block for the Cloud Service Provider."),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
	attrs["couchbase_server"] = schema.SingleNestedAttribute{
		Computed:            true,
		MarkdownDescription: "Couchbase Server details for the free tier cluster.",
		Attributes: map[string]schema.Attribute{
			"version": WithDescription(stringAttribute([]string{computed}), "The version of Couchbase Server for the free tier cluster."),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
			objectplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["service_groups"] = schema.SetNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The Service Groups for the free tier cluster.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"node": schema.SingleNestedAttribute{
					Computed:            true,
					MarkdownDescription: "Node details for the Service Group.",
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
								"type":          WithDescription(stringAttribute([]string{computed}), "The type of disk for the node. Should be one of 'ssd', 'hdd', or 'nvme'."),
								"storage":       WithDescription(int64Attribute(computed), "Storage size of the disk."),
								"iops":          WithDescription(int64Attribute(computed), "Input/Output Operations Per Second (IOPS) for the disk."),
								"autoexpansion": WithDescription(boolAttribute(computed), "Indicates if auto-expansion is enabled for the disk."),
							},
						},
					},
				},
				"num_of_nodes": WithDescription(int64Attribute(computed), "The number of nodes in the Service Group."),
				"services":     WithDescription(stringSetAttribute(computed), "The services enabled for the Service Group. Should be a comma-separated list of services. For example, 'data,index,query'."),
			},
		},
	}
	attrs["availability"] = schema.SingleNestedAttribute{
		Computed:            true,
		MarkdownDescription: "Availability zone details for the free tier cluster. Free tier clusters have single availability zones (AZ).",
		Attributes: map[string]schema.Attribute{
			"type": WithDescription(stringAttribute([]string{computed}), "The availability zone type. This should be 'single' for the free tier cluster."),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a free tier operational cluster.",
		Attributes:          attrs,
	}
}
