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
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"name":            stringAttribute([]string{required, requiresReplace}),
			"provider_type":   stringAttribute([]string{required, requiresReplace}),
			"commands": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"provider_config": schema.SingleNestedAttribute{
				Description: "The 'accountId', 'vpcId', 'region', and 'cidr' fields are required for AWS VPC peering. " +
					"For GCP, the 'networkName', 'projectId', 'serviceAccount', and 'cidr' fields are required for VPC peering. ",
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"aws_config": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"account_id":  stringAttribute([]string{optional}),
							"vpc_id":      stringAttribute([]string{optional}),
							"region":      stringAttribute([]string{optional}),
							"cidr":        stringAttribute([]string{required}),
							"provider_id": stringAttribute([]string{computed}),
						},
					},
					"gcp_config": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"cidr":            stringAttribute([]string{required}),
							"network_name":    stringAttribute([]string{optional}),
							"project_id":      stringAttribute([]string{optional}),
							"service_account": stringAttribute([]string{optional}),
							"provider_id":     stringAttribute([]string{computed}),
						},
					},
					"azure_config": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"tenant_id":       stringAttribute([]string{optional}),
							"cidr":            stringAttribute([]string{required}),
							"resource_group":  stringAttribute([]string{optional}),
							"subscription_id": stringAttribute([]string{optional}),
							"vnet_id":         stringAttribute([]string{optional}),
							"provider_id":     stringAttribute([]string{computed}),
						},
					},
				},
			},
			"status": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"reasoning": stringAttribute([]string{computed}),
					"state":     stringAttribute([]string{computed}),
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
