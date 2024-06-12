package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NetworkPeerSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       computedStringAttribute,
						"name":     computedStringAttribute,
						"commands": computedListAttribute,
						"provider_config": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"provider_id": computedStringAttribute,
								"AWS_config": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"accountId": computedStringAttribute,
										"vpcId":     computedStringAttribute,
										"region":    computedStringAttribute,
										"cidr":      computedStringAttribute,
										//"provider_id":    computedStringAttribute,
									},
								},
								"GCP_config": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"cidr":           computedStringAttribute,
										"networkName":    computedStringAttribute,
										"projectId":      computedStringAttribute,
										"serviceAccount": computedStringAttribute,
										//"provider_id":    computedStringAttribute,
									},
								},
							},
						},
						"status": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"reasoning": computedStringAttribute,
								"state":     computedStringAttribute,
							},
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}
