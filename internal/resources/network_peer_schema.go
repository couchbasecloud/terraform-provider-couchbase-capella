package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
			"name":            stringAttribute([]string{required}),
			//"commands": schema.SetAttribute{
			//	Computed:    true,
			//	ElementType: types.StringType,
			//},
			"commands":      stringSetAttribute(computed),
			"provider_type": stringAttribute([]string{required}),
			"provider_config": schema.SingleNestedAttribute{
				Description: "The 'accountId', 'vpcId', 'region', and 'cidr' fields are required for AWS VPC peering. " +
					"For GCP, the 'networkName', 'projectId', 'serviceAccount', and 'cidr' fields are required for VPC peering. ",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"provider_id": stringAttribute([]string{computed}),
					"AWS_config": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"accountId": stringAttribute([]string{optional}),
							"vpcId":     stringAttribute([]string{optional}),
							"region":    stringAttribute([]string{optional}),
							"cidr":      stringAttribute([]string{required}),
							//"provider_id":    stringAttribute([]string{computed}),
						},
					},
					"GCP_config": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"cidr":           stringAttribute([]string{required}),
							"networkName":    stringAttribute([]string{optional}),
							"projectId":      stringAttribute([]string{optional}),
							"serviceAccount": stringAttribute([]string{optional}),
							//"provider_id":    stringAttribute([]string{computed}),
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
