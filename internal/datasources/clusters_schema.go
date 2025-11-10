package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterBuilder = capellaschema.NewSchemaBuilder("cluster")

func ClusterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", clusterBuilder, requiredString())

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", clusterBuilder, computedInt64())
	capellaschema.AddAttr(computeAttrs, "ram", clusterBuilder, computedInt64())

	diskAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(diskAttrs, "type", clusterBuilder, computedString())
	capellaschema.AddAttr(diskAttrs, "storage", clusterBuilder, computedInt64())
	capellaschema.AddAttr(diskAttrs, "iops", clusterBuilder, computedInt64())
	capellaschema.AddAttr(diskAttrs, "autoexpansion", clusterBuilder, computedBool())

	nodeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(nodeAttrs, "compute", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(nodeAttrs, "disk", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: diskAttrs,
	})

	serviceGroupAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serviceGroupAttrs, "node", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: nodeAttrs,
	})
	capellaschema.AddAttr(serviceGroupAttrs, "num_of_nodes", clusterBuilder, computedInt64())
	capellaschema.AddAttr(serviceGroupAttrs, "services", clusterBuilder, &schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", clusterBuilder, computedString())
	capellaschema.AddAttr(cloudProviderAttrs, "region", clusterBuilder, computedString())
	capellaschema.AddAttr(cloudProviderAttrs, "cidr", clusterBuilder, computedString())

	couchbaseServerAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(couchbaseServerAttrs, "version", clusterBuilder, computedString())

	availabilityAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(availabilityAttrs, "type", clusterBuilder, computedString())

	supportAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(supportAttrs, "plan", clusterBuilder, computedString())
	capellaschema.AddAttr(supportAttrs, "timezone", clusterBuilder, computedString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "enable_private_dns_resolution", clusterBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "connection_string", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cloudProviderAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "couchbase_server", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: couchbaseServerAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "service_groups", clusterBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: serviceGroupAttrs,
		},
	})
	capellaschema.AddAttr(dataAttrs, "availability", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: availabilityAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "support", clusterBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: supportAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "current_state", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "app_service_id", clusterBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", clusterBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", clusterBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source retrieves the details of Couchbase Capella clusters.",
		Attributes:          attrs,
	}
}
