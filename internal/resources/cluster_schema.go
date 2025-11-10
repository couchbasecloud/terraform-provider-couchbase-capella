package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterBuilder = capellaschema.NewSchemaBuilder("cluster")

func ClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", clusterBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", clusterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", clusterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "description", clusterBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "zones", clusterBuilder, stringSetAttribute(optional, requiresReplace))
	capellaschema.AddAttr(attrs, "enable_private_dns_resolution", clusterBuilder, boolDefaultAttribute(false, optional, computed, requiresReplace))

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(cloudProviderAttrs, "region", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(cloudProviderAttrs, "cidr", clusterBuilder, stringAttribute([]string{required}))

	capellaschema.AddAttr(attrs, "cloud_provider", clusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: cloudProviderAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	})

	capellaschema.AddAttr(attrs, "configuration_type", clusterBuilder, stringAttribute([]string{optional, computed, requiresReplace, useStateForUnknown, deprecated}))

	couchbaseServerAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(couchbaseServerAttrs, "version", clusterBuilder, stringAttribute([]string{optional, computed}))

	capellaschema.AddAttr(attrs, "couchbase_server", clusterBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Computed:   true,
		Attributes: couchbaseServerAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", clusterBuilder, int64Attribute(required))
	capellaschema.AddAttr(computeAttrs, "ram", clusterBuilder, int64Attribute(required))

	diskAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(diskAttrs, "type", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(diskAttrs, "storage", clusterBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(diskAttrs, "iops", clusterBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(diskAttrs, "autoexpansion", clusterBuilder, boolAttribute(optional, computed))

	nodeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(nodeAttrs, "compute", clusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(nodeAttrs, "disk", clusterBuilder, &schema.SingleNestedAttribute{
		Description: "The 'storage' and 'IOPS' fields are required for AWS. " +
			"For Azure, only the 'disktype' field is required. For the Ultra disk type, you can provide storage, IOPS, and auto-expansion fields. For Premium type, you can only provide the auto-expansion field, others cannot be set." +
			" In the case of GCP, only 'pd ssd' disk type is available, and you cannot set the 'IOPS' field.",
		Required:   true,
		Attributes: diskAttrs,
	})

	serviceGroupAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serviceGroupAttrs, "node", clusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: nodeAttrs,
	})
	capellaschema.AddAttr(serviceGroupAttrs, "num_of_nodes", clusterBuilder, int64Attribute(required))
	capellaschema.AddAttr(serviceGroupAttrs, "services", clusterBuilder, stringSetAttribute(required))

	capellaschema.AddAttr(attrs, "service_groups", clusterBuilder, &schema.SetNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: serviceGroupAttrs,
		},
	})

	availabilityAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(availabilityAttrs, "type", clusterBuilder, stringAttribute([]string{required}))

	capellaschema.AddAttr(attrs, "availability", clusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: availabilityAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	})

	supportAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(supportAttrs, "plan", clusterBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(supportAttrs, "timezone", clusterBuilder, stringAttribute([]string{computed, optional}))

	capellaschema.AddAttr(attrs, "support", clusterBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: supportAttrs,
	})

	capellaschema.AddAttr(attrs, "current_state", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "app_service_id", clusterBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", clusterBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "if_match", clusterBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(attrs, "etag", clusterBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "Manages the operational cluster resource.",
		Attributes:          attrs,
	}
}
