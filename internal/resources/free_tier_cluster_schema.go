package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierClusterSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute([]string{required, requiresReplace},
				validator.String(stringvalidator.LengthAtLeast(1))),
			"project_id": stringAttribute([]string{required, requiresReplace},
				validator.String(stringvalidator.LengthAtLeast(1))),
			"name": stringAttribute([]string{required},
				validator.String(stringvalidator.LengthAtLeast(1))),
			"description":                   stringAttribute([]string{optional, computed}),
			"app_service_id":                stringAttribute([]string{computed}),
			"enable_private_dns_resolution": boolAttribute(computed),
			"connection_string":             stringAttribute([]string{computed}),
			"current_state":                 stringAttribute([]string{computed}),
			"cmek_id":                       stringAttribute([]string{computed}),
			"audit":                         computedAuditAttribute(),
			"etag":                          stringAttribute([]string{computed}),
			"support": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"plan":     stringAttribute([]string{computed}),
					"timezone": stringAttribute([]string{computed}),
				},
			},
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
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": stringAttribute([]string{computed}),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"service_groups": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"cpu": int64Attribute(computed),
										"ram": int64Attribute(computed),
									},
								},
								"disk": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"type":          stringAttribute([]string{computed}),
										"storage":       int64Attribute(computed),
										"iops":          int64Attribute(computed),
										"autoexpansion": boolAttribute(computed),
									},
								},
							},
						},
						"num_of_nodes": int64Attribute(computed),
						"services":     stringSetAttribute(computed),
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": stringAttribute([]string{computed}),
				},
			},
		},
	}

}
