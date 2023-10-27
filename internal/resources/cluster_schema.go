package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			"organization_id": stringAttribute(required),
			"project_id":      stringAttribute(required),
			"name":            stringAttribute(required),
			"description":     stringAttribute(optional),
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type":   stringAttribute(required),
					"region": stringAttribute(required),
					"cidr":   stringAttribute(required),
				},
			},
			"configuration_type": stringAttribute(required),
			"couchbase_server": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": stringAttribute(optional, computed),
				},
			},
			"service_groups": schema.ListNestedAttribute{
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
									Required: true,
									Attributes: map[string]schema.Attribute{
										"type":    stringAttribute(required),
										"storage": int64Attribute(optional, computed),
										"iops":    int64Attribute(optional, computed),
									},
								},
							},
						},
						"num_of_nodes": int64Attribute(required),
						"services":     stringListAttribute(required),
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": stringAttribute(required),
				},
			},
			"support": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"plan":     stringAttribute(required),
					"timezone": stringAttribute(required),
				},
			},
			"current_state":  stringAttribute(computed),
			"app_service_id": stringAttribute(optional, computed),
			"audit":          computedAuditAttribute(),
			"if_match":       stringAttribute(optional),
			"etag":           stringAttribute(computed),
		},
	}
}
