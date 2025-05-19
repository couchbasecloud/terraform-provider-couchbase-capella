package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NetworkPeerSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to manage network peering for a Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The unique identifier for the network peering record.",
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the organization."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the project."),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the cluster to set up network peering."),
			"name": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The name of the network peering relationship."),
			"provider_type": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The cloud provider type for the network peering (aws, gcp, or azure)."),
			"commands": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The list of commands required to set up network peering.",
			},
			"provider_config": schema.SingleNestedAttribute{
				Description: "Cloud provider-specific configuration for the network peering." +
					"The 'accountId', 'vpcId', 'region', and 'cidr' fields are required for AWS VPC peering. " +
					"For GCP, the 'networkName', 'projectId', 'serviceAccount', and 'cidr' fields are required for VPC peering. ",
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"aws_config": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "AWS-specific configuration for VPC peering.",
						Attributes: map[string]schema.Attribute{
							"account_id": WithDescription(stringAttribute([]string{optional}),
								"The numeric AWS Account ID or Owner ID."),
							"vpc_id": WithDescription(stringAttribute([]string{optional}),
								"The alphanumeric VPC ID which starts with 'vpc-'. This is also known as the networkId."),
							"region": WithDescription(stringAttribute([]string{optional}),
								"The AWS region where your VPC is deployed."),
							"cidr": WithDescription(stringAttribute([]string{required}),
								"The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block."),
							"provider_id": WithDescription(stringAttribute([]string{computed}),
								"The ID of the VPC peer on AWS."),
						},
					},
					"gcp_config": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "GCP-specific configuration for VPC network peering.",
						Attributes: map[string]schema.Attribute{
							"cidr": WithDescription(stringAttribute([]string{required}),
								"The GCP VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block."),
							"network_name": WithDescription(stringAttribute([]string{optional}),
								"The name of the network that you want to peer with."),
							"project_id": WithDescription(stringAttribute([]string{optional}),
								"The unique identifier for your GCP project."),
							"service_account": WithDescription(stringAttribute([]string{optional}),
								"ServiceAccount created or assigned on the external VPC project. GCP Service Account with DNS Admin and Compute.NetworkAdmin permissions. Must be in email form shown by 'gcloud iam service-accounts list'."),
							"provider_id": WithDescription(stringAttribute([]string{computed}),
								"The ID of the VPC peer on GCP."),
						},
					},
					"azure_config": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Azure-specific configuration for VNet peering.",
						Attributes: map[string]schema.Attribute{
							"tenant_id": WithDescription(stringAttribute([]string{optional}),
								"The Azure tenant ID where the VNet exists."),
							"cidr": WithDescription(stringAttribute([]string{required}),
								"The CIDR block from the virtual network that you created in Azure. This cannot overlap with your Capella CIDR Block."),
							"resource_group": WithDescription(stringAttribute([]string{optional}),
								"The resource group name holding the resource you're connecting with Capella."),
							"subscription_id": WithDescription(stringAttribute([]string{optional}),
								"The Azure subscription ID where the VNet exists."),
							"vnet_id": WithDescription(stringAttribute([]string{optional}),
								"The VNet ID is the name of the virtual network peering in Azure."),
							"provider_id": WithDescription(stringAttribute([]string{computed}),
								"The ID of the VNet peer on Azure."),
						},
					},
				},
			},
			"status": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Current status of the network peering connection.",
				Attributes: map[string]schema.Attribute{
					"reasoning": WithDescription(stringAttribute([]string{computed}),
						"Detailed reason for the current status of the peering connection."),
					"state": WithDescription(stringAttribute([]string{computed}),
						"Current state of the peering connection (e.g., pending, active, failed)."),
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
