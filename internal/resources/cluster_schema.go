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
				Computed:            true,
				MarkdownDescription: "The ID of the Capella cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the organization that the cluster belongs to."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the project that the cluster belongs to."),
			"name": WithDescription(stringAttribute([]string{required}),
				"The name of the cluster. The name must be unique within the project."),
			"description": stringAttribute([]string{optional, computed}),
			"zones": WithDescription(stringSetAttribute(optional, requiresReplace),
				"The list of availability zones for the cluster. The cluster will be created in the specified zones. If not specified, the cluster will be created in all available zones."),
			"enable_private_dns_resolution": WithDescription(boolDefaultAttribute(false, optional, computed, requiresReplace),
				"Enable private DNS resolution for the cluster. If set to true, the cluster will be accessible via private DNS names. If set to false, the cluster will be accessible via public DNS names."),
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
			"configuration_type": stringAttribute([]string{optional, computed, requiresReplace, useStateForUnknown, deprecated}),
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
										"cpu": WithDescription(int64Attribute(required),
											"The number of CPU cores for the node. The value must be between 1 and 128."),
										"ram": int64Attribute(required),
									},
								},
								"disk": schema.SingleNestedAttribute{
									Description: "The 'storage' and 'IOPS' fields are required for AWS. " +
										"For Azure, only the 'disktype' field is required, and for Ultra disk type, you can provide all 3 - storage, iops and autoexpansion fields. For Premium type, you can only provide the autoexpansion field, others can't be set." +
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
					"timezone": stringAttribute([]string{computed, optional}),
				},
			},
			"current_state":     stringAttribute([]string{computed}),
			"connection_string": stringAttribute([]string{computed}),
			"app_service_id":    stringAttribute([]string{computed}),
			"audit":             computedAuditAttribute(),
			// if_match is only required during update call
			"if_match": stringAttribute([]string{optional}),
			"etag":     stringAttribute([]string{computed}),
		},
	}
}
