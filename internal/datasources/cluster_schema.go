package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var oneClusterBuilder = capellaschema.NewSchemaBuilder("oneCluster")

func OneClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", oneClusterBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", oneClusterBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "id", oneClusterBuilder, requiredStringWithValidator())

	capellaschema.AddAttr(attrs, "name", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "description", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "enable_private_dns_resolution", oneClusterBuilder, computedBool())
	capellaschema.AddAttr(attrs, "connection_string", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "configuration_type", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "current_state", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "app_service_id", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "etag", oneClusterBuilder, computedString())
	capellaschema.AddAttr(attrs, "if_match", oneClusterBuilder, computedString())

	capellaschema.AddAttr(attrs, "zones", oneClusterBuilder, &schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	// cloud_provider
	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", oneClusterBuilder, computedString(), "CloudProvider")
	capellaschema.AddAttr(cloudProviderAttrs, "region", oneClusterBuilder, computedString(), "CloudProvider")
	capellaschema.AddAttr(cloudProviderAttrs, "cidr", oneClusterBuilder, computedString(), "CloudProvider")

	capellaschema.AddAttr(attrs, "cloud_provider", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cloudProviderAttrs,
	})

	// couchbase_server
	couchbaseServerAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(couchbaseServerAttrs, "version", oneClusterBuilder, computedString())

	capellaschema.AddAttr(attrs, "couchbase_server", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: couchbaseServerAttrs,
	})

	// service_groups
	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", oneClusterBuilder, computedInt64())
	capellaschema.AddAttr(computeAttrs, "ram", oneClusterBuilder, computedInt64())

	diskAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(diskAttrs, "type", oneClusterBuilder, computedString())
	capellaschema.AddAttr(diskAttrs, "storage", oneClusterBuilder, computedInt64())
	capellaschema.AddAttr(diskAttrs, "iops", oneClusterBuilder, computedInt64())
	capellaschema.AddAttr(diskAttrs, "autoexpansion", oneClusterBuilder, computedBool())

	nodeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(nodeAttrs, "compute", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(nodeAttrs, "disk", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: diskAttrs,
	})

	serviceGroupAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serviceGroupAttrs, "node", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: nodeAttrs,
	})
	capellaschema.AddAttr(serviceGroupAttrs, "num_of_nodes", oneClusterBuilder, computedInt64())
	capellaschema.AddAttr(serviceGroupAttrs, "services", oneClusterBuilder, &schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	capellaschema.AddAttr(attrs, "service_groups", oneClusterBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: serviceGroupAttrs,
		},
	})

	// availability
	availabilityAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(availabilityAttrs, "type", oneClusterBuilder, computedString(), "Availability")

	capellaschema.AddAttr(attrs, "availability", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: availabilityAttrs,
	})

	// support
	supportAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(supportAttrs, "plan", oneClusterBuilder, computedString(), "Support")
	capellaschema.AddAttr(supportAttrs, "timezone", oneClusterBuilder, computedString(), "Support")

	capellaschema.AddAttr(attrs, "support", oneClusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: supportAttrs,
	})

	capellaschema.AddAttr(attrs, "audit", oneClusterBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "The data source retrieves the details of a single Couchbase Capella cluster.",
		Attributes:          attrs,
	}
}
