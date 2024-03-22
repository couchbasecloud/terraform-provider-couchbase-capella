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
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"name":            stringAttribute(required),
			"description":     stringAttribute(optional, computed),
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type":   stringAttribute(required),
					"region": stringAttribute(required),
					"cidr":   stringAttribute(required),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"configuration_type": stringAttribute(optional, computed, requiresReplace, useStateForUnknown),
			"couchbase_server": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": stringAttribute(optional, computed),
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
										"type":          stringAttribute(required),
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
					"type": stringAttribute(required),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
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
			"app_service_id": stringAttribute(computed),
			"audit":          computedAuditAttribute(),
			// if_match is only required during update call
			"if_match": stringAttribute(optional),
			"etag":     stringAttribute(computed),
		},
	}
}
