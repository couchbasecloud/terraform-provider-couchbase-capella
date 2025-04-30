package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the Couchbase Capella cluster resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("The unique identifier of the Capella organization that owns this cluster.")),
			"project_id": stringAttribute([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("The unique identifier of the Capella project where this cluster will be created.")),
			"name": stringAttribute([]string{required},
				withMarkdown[*schema.StringAttribute]("The name of the cluster. This must be unique within the project.")),
			"description": stringAttribute([]string{optional, computed},
				withMarkdown[*schema.StringAttribute]("A description of the cluster's purpose or characteristics.")),
			"zones": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The cloud provider availability zones where the cluster will be deployed. Currently only supports single AZ clusters.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"enable_private_dns_resolution": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enables private DNS resolution for the cluster.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The cloud provider configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"type":   stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The type of cloud provider (e.g., 'aws', 'azure', 'gcp').")),
					"region": stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The region where the cluster will be deployed.")),
					"cidr":   stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The CIDR block for the cluster's network.")),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"configuration_type": stringAttribute([]string{optional, computed, requiresReplace, useStateForUnknown, deprecated}, withMarkdown[*schema.StringAttribute]("The type of cluster configuration.")),
			"couchbase_server": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Couchbase Server configuration settings.",
				Attributes: map[string]schema.Attribute{
					"version": stringAttribute([]string{optional, computed}, withMarkdown[*schema.StringAttribute]("The version of Couchbase Server to deploy.")),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"service_groups": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "Configuration for service groups in the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "Node configuration for the service group.",
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Required:            true,
									MarkdownDescription: "Compute resources configuration for the node.",
									Attributes: map[string]schema.Attribute{
										"cpu": int64Attribute([]string{required}, withMarkdown[*schema.Int64Attribute]("Number of CPU cores allocated to the node.")),
										"ram": int64Attribute([]string{required}, withMarkdown[*schema.Int64Attribute]("Amount of RAM in MB allocated to the node.")),
									},
								},
								"disk": schema.SingleNestedAttribute{
									Required:            true,
									MarkdownDescription: "Disk configuration for the node. Storage and IOPS are required for AWS. For Azure, only disktype is required, with additional options for Ultra disk type. GCP only supports pd-ssd disk type.",
									Attributes: map[string]schema.Attribute{
										"type":          stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The type of disk to use (e.g., 'pd-ssd', 'premium', 'ultra').")),
										"storage":       int64Attribute([]string{optional, computed}, withMarkdown[*schema.Int64Attribute]("The size of the disk in GB.")),
										"iops":          int64Attribute([]string{optional, computed}, withMarkdown[*schema.Int64Attribute]("The number of IOPS for the disk.")),
										"autoexpansion": boolAttribute([]string{optional, computed}, withMarkdown[*schema.BoolAttribute]("Whether to enable automatic disk expansion.")),
									},
								},
							},
						},
						"num_of_nodes": int64Attribute([]string{required}, withMarkdown[*schema.Int64Attribute]("The number of nodes in this service group.")),
						"services":     stringSetAttribute([]string{required}, withMarkdown[*schema.SetAttribute]("The services to run on this service group (e.g., 'data', 'index', 'query', 'search').")),
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Availability configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"type": stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The type of availability configuration (e.g., 'single', 'multi').")),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"support": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Support configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"plan":     stringAttribute([]string{required}, withMarkdown[*schema.StringAttribute]("The support plan for the cluster.")),
					"timezone": stringAttribute([]string{computed, optional}, withMarkdown[*schema.StringAttribute]("The timezone for support operations.")),
				},
			},
			"current_state":     stringAttribute([]string{computed}, withMarkdown[*schema.StringAttribute]("The current state of the cluster.")),
			"connection_string": stringAttribute([]string{computed}, withMarkdown[*schema.StringAttribute]("The connection string for accessing the cluster.")),
			"app_service_id":    stringAttribute([]string{computed}, withMarkdown[*schema.StringAttribute]("The ID of the associated application service.")),
			"audit":             computedAuditAttribute(),
			"if_match":          stringAttribute([]string{optional}, withMarkdown[*schema.StringAttribute]("The ETag value for optimistic concurrency control during updates.")),
			"etag":              stringAttribute([]string{computed}, withMarkdown[*schema.StringAttribute]("The current ETag value of the cluster resource.")),
		},
	}
}
