package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ClusterSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id":               stringAttribute([]string{required, requiresReplace}),
			"project_id":                    stringAttribute([]string{required, requiresReplace}),
			"name":                          stringAttribute([]string{required}),
			"description":                   stringAttribute([]string{optional, computed}),
			"enable_private_dns_resolution": boolDefaultAttribute(false, optional, requiresReplace),
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type":   stringAttribute([]string{required}),
					"region": stringAttribute([]string{required}),
					"cidr":   stringAttribute([]string{required}),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"couchbase_server": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": stringAttribute([]string{optional, computed}),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"service_groups": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Required: true,
									Attributes: map[string]schema.Attribute{
										"cpu": int64Attribute(required),
										"ram": int64Attribute(required),
									},
								},
								"disk": schema.SingleNestedAttribute{
									Description: "The 'storage' and 'IOPS' fields are required for AWS. " +
										"For Azure, only the 'disktype' field is required, and for Ultra, you can provide all three fields. " +
										"In the case of GCP, only 'pd ssd' disk type is available, and you cannot set the 'IOPS' field.",
									Required: true,
									Attributes: map[string]schema.Attribute{
										"type":          stringAttribute([]string{required}),
										"storage":       int64Attribute(optional, computed),
										"iops":          int64Attribute(optional, computed),
										"autoexpansion": boolAttribute(optional, computed),
									},
								},
							},
						},
						"num_of_nodes": int64Attribute(required),
						"services":     stringSetAttribute(required),
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": stringAttribute([]string{required}),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"support": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"plan":     stringAttribute([]string{required}),
					"timezone": stringAttribute([]string{required}),
				},
			},
			"current_state":  stringAttribute([]string{computed}),
			"app_service_id": stringAttribute([]string{computed}),
			"audit":          computedAuditAttribute(),
			// if_match is only required during update call
			"if_match": stringAttribute([]string{optional}),
			"etag":     stringAttribute([]string{computed}),
		},
	}
}
