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
	capellaschema.AddAttr(attrs, "name", freeTierClusterBuilder, stringAttribute([]string{required}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "description", freeTierClusterBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "app_service_id", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "current_state", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "cmek_id", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "etag", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "enable_private_dns_resolution", freeTierClusterBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "audit", freeTierClusterBuilder, computedAuditAttribute())

	supportAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(supportAttrs, "plan", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(supportAttrs, "timezone", freeTierClusterBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "support", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: supportAttrs,
	})

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", freeTierClusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(cloudProviderAttrs, "region", freeTierClusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(cloudProviderAttrs, "cidr", freeTierClusterBuilder, stringAttribute([]string{required}))

	capellaschema.AddAttr(attrs, "cloud_provider", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: cloudProviderAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	})

	couchbaseServerAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(couchbaseServerAttrs, "version", freeTierClusterBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "couchbase_server", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: couchbaseServerAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", freeTierClusterBuilder, int64Attribute(computed))
	capellaschema.AddAttr(computeAttrs, "ram", freeTierClusterBuilder, int64Attribute(computed))

	diskAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(diskAttrs, "type", freeTierClusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(diskAttrs, "storage", freeTierClusterBuilder, int64Attribute(computed))
	capellaschema.AddAttr(diskAttrs, "iops", freeTierClusterBuilder, int64Attribute(computed))
	capellaschema.AddAttr(diskAttrs, "autoexpansion", freeTierClusterBuilder, boolAttribute(computed))

	nodeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(nodeAttrs, "compute", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(nodeAttrs, "disk", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: diskAttrs,
	})

	serviceGroupsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serviceGroupsAttrs, "node", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: nodeAttrs,
	})
	capellaschema.AddAttr(serviceGroupsAttrs, "num_of_nodes", freeTierClusterBuilder, int64Attribute(computed))
	capellaschema.AddAttr(serviceGroupsAttrs, "services", freeTierClusterBuilder, stringSetAttribute(computed))

	capellaschema.AddAttr(attrs, "service_groups", freeTierClusterBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: serviceGroupsAttrs,
		},
	})

	availabilityAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(availabilityAttrs, "type", freeTierClusterBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "availability", freeTierClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: availabilityAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a free tier operational cluster.",
		Attributes:          attrs,
	}

}
