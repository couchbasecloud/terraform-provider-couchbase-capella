package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterBuilder = capellaschema.NewSchemaBuilder("cluster")

func ClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", clusterBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", clusterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", clusterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "description", clusterBuilder, stringAttribute([]string{optional, computed}))

	attrs["zones"] = WithDescription(stringSetAttribute(optional, requiresReplace),
		"The Cloud Services Provider's availability zones for the cluster."+
			"For single availability zone clusters, only 1 zone is allowed in list.")
	attrs["enable_private_dns_resolution"] = WithDescription(boolDefaultAttribute(false, optional, computed, requiresReplace),
		"EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address. "+
			"This exists to support the use case of customers connecting from their own data centers where it is not possible to make use of a Cloud Service Provider DNS zone.")
	attrs["cloud_provider"] = schema.SingleNestedAttribute{
		Required:            true,
		MarkdownDescription: "The Cloud Service Provider where the cluster will be hosted. ",
		Attributes: map[string]schema.Attribute{
			"type": WithDescription(stringAttribute([]string{required}),
				"The Cloud Service Provider type. Currently supporting AWS, GCP and Azure. For Single Node cluster, only the AWS type Cloud Service Provider is allowed.",
			),
			"region": WithDescription(stringAttribute([]string{required}),
				"The region where the cluster will be hosted."),
			"cidr": WithDescription(stringAttribute([]string{required}),
				"The CIDR block for the cluster's network."),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
	attrs["configuration_type"] = WithDescription(stringAttribute([]string{optional, computed, requiresReplace, useStateForUnknown, deprecated}),
		"The configuration type of the cluster. This field is deprecated.")
	attrs["couchbase_server"] = schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Configuration for the Couchbase Server running on the cluster.",
		Attributes: map[string]schema.Attribute{
			"version": WithDescription(stringAttribute([]string{optional, computed}),
				"The version of Couchbase Server to run on the cluster."),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
			objectplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["service_groups"] = schema.SetNestedAttribute{
		Required:            true,
		MarkdownDescription: "Configuration for the Service Groups in the cluster. Each Service Group represents a set of nodes with the same configuration.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"node": schema.SingleNestedAttribute{
					Required:            true,
					MarkdownDescription: "Node configuration for this Service Group.",
					Attributes: map[string]schema.Attribute{
						"compute": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "Compute resources configuration for the nodes.",
							Attributes: map[string]schema.Attribute{
								"cpu": WithDescription(int64Attribute(required),
									"The number of CPU cores for the node. The value must be between 1 and 128."),
								"ram": WithDescription(int64Attribute(required),
									"The amount of RAM in MB for the node."),
							},
						},
						"disk": schema.SingleNestedAttribute{
							Description: "The 'storage' and 'IOPS' fields are required for AWS. " +
								"For Azure, only the 'disktype' field is required. For the Ultra disk type, you can provide storage, IOPS, and auto-expansion fields. For Premium type, you can only provide the auto-expansion field, others cannot be set." +
								" In the case of GCP, only 'pd ssd' disk type is available, and you cannot set the 'IOPS' field.",
							Required: true,
							Attributes: map[string]schema.Attribute{
								"type": WithDescription(stringAttribute([]string{required}),
									"The type of disk to use. For AWS: gp3 or io2, for Azure: Premium or UltraSSD, for GCP: pd-ssd."),
								"storage": WithDescription(int64Attribute(optional, computed),
									"The size of the disk in GB."),
								"iops": WithDescription(int64Attribute(optional, computed),
									"The number of IOPS for the disk. Only applicable for AWS and Azure."),
								"autoexpansion": WithDescription(boolAttribute(optional, computed),
									"Enable or disable automatic disk expansion.  This can only be set for Azure."),
							},
						},
					},
				},
				"num_of_nodes": WithDescription(int64Attribute(required),
					"The number of nodes in this Service Group."),
				"services": WithDescription(stringSetAttribute(required),
					"The list of Couchbase Services to run on the nodes in this Service Group."),
			},
		},
	}
	attrs["availability"] = schema.SingleNestedAttribute{
		Required:            true,
		MarkdownDescription: "Availability configuration for the cluster.",
		Attributes: map[string]schema.Attribute{
			"type": WithDescription(stringAttribute([]string{required}),
				"The availability type of the cluster. Can be 'single' for Single Node or 'multi' for Multi Node."),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
	attrs["support"] = schema.SingleNestedAttribute{
		Required:            true,
		MarkdownDescription: "Support configuration for the cluster.",
		Attributes: map[string]schema.Attribute{
			"plan": WithDescription(stringAttribute([]string{required}),
				"The support plan options include 'Basic', 'Developer Pro', or 'Enterprise'."),
			"timezone": WithDescription(stringAttribute([]string{computed, optional}),
				"The timezone for support coverage."),
		},
	}

	capellaschema.AddAttr(attrs, "current_state", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "app_service_id", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", clusterBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "if_match", clusterBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(attrs, "etag", clusterBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "Manages the operational cluster resource.",
		Attributes:          attrs,
	}
}
