package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierClusterBuilder = capellaschema.NewSchemaBuilder("freeTierCluster")

func FreeTierClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", freeTierClusterBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", freeTierClusterBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", freeTierClusterBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "etag", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "name", freeTierClusterBuilder, stringAttribute([]string{required}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "description", freeTierClusterBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "app_service_id", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "current_state", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "cmek_id", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "enable_private_dns_resolution", freeTierClusterBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "audit", freeTierClusterBuilder, computedAuditAttribute())
	attrs["support"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"plan":     stringAttribute([]string{computed}),
			"timezone": stringAttribute([]string{computed}),
		},
	}
	attrs["cloud_provider"] = schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"type":   stringAttribute([]string{required}),
			"region": stringAttribute([]string{required}),
			"cidr":   stringAttribute([]string{required}),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
	attrs["couchbase_server"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"version": stringAttribute([]string{computed}),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
			objectplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["service_groups"] = schema.SetNestedAttribute{
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
	}
	attrs["availability"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"type": stringAttribute([]string{computed}),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a free tier operational cluster.",
		Attributes:          attrs,
	}
}
