package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NetworkPeerSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source to retrieve network peering records for a Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the network peering records.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the VPC peering record.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the network peering relationship.",
						},
						"provider_config": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The Cloud Service Provider-specific configuration for the network peering. This provides details about the configuration and the ID of the VPC peer on AWS, GCP, or Azure.",
							Attributes: map[string]schema.Attribute{
								"aws_config": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "AWS configuration data required to establish a VPC peering relationship.",
									Attributes: map[string]schema.Attribute{
										"account_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The numeric AWS Account ID or Owner ID.",
										},
										"vpc_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The alphanumeric VPC ID which starts with 'vpc-'. This is also known as the networkId.",
										},
										"region": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The AWS region where your VPC is deployed",
										},
										"cidr": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The AWS VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.",
										},
										"provider_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The ID of the VPC peer on AWS.",
										},
									},
								},
								"gcp_config": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "GCP configuration data required to establish a VPC peering relationship.",
									Attributes: map[string]schema.Attribute{
										"cidr": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The GCP VPC CIDR block of network in which your application runs. This cannot overlap with your Capella CIDR Block.",
										},
										"network_name": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The name of the network that you want to peer with.",
										},
										"project_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The unique identifier for your GCP project.",
										},
										"service_account": schema.StringAttribute{
											Computed: true,
											MarkdownDescription: "The Service Account created or assigned on the external VPC project. The GCP Service Account should have the following permission:" +
												"DNS Admin and Network Admin. It must be in the form of an email address, as shown by the output of the gcloud iam service-accounts list command. [Reference](https://cloud.google.com/iam/docs/service-accounts-create#creating)",
										},
										"provider_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The ID of the VPC peer on GCP.",
										},
									},
								},
								"azure_config": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Azure configuration data required to establish a VNet peering relationship.",
									Attributes: map[string]schema.Attribute{
										"tenant_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The Azure tenant ID where the VNet exists. For more information, see [How to Find Tenant](https://learn.microsoft.com/en-us/entra/fundamentals/how-to-find-tenant)",
										},
										"cidr": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The CIDR block from the virtual network that you created in Azure.",
										},
										"resource_group": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The resource group name holding the resource you’re connecting with Capella.",
										},
										"subscription_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The Azure subscription ID where the VNet exists. For more information, see [Find Your Azure Subscription](https://learn.microsoft.com/en-us/azure/azure-portal/get-subscription-tenant-id#find-your-azure-subscription)",
										},
										"vnet_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The VNet ID is the name of the virtual network peering in Azure.",
										},
										"provider_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The ID of the VNet peer on Azure.",
										},
									},
								},
							},
						},
						"status": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Current status of the network peering connection.",
							Attributes: map[string]schema.Attribute{
								"reasoning": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Detailed reason for the current status of the peering connection.",
								},
								"state": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Current state of the peering connection.",
								},
							},
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}
