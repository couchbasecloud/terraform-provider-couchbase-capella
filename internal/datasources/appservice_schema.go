package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceDSBuilder = capellaschema.NewSchemaBuilder("appService")

func AppServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceDSBuilder, requiredString())

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", appServiceDSBuilder, computedInt64(), "AppServiceCompute")
	capellaschema.AddAttr(computeAttrs, "ram", appServiceDSBuilder, computedInt64(), "AppServiceCompute")

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "nodes", appServiceDSBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "current_state", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "compute", appServiceDSBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "version", appServiceDSBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", appServiceDSBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", appServiceDSBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source retrieves information for an App Service in Capella. App Service is a fully managed application backend designed to provide data synchronization between mobile or IoT applications running Couchbase Lite and your Couchbase Capella database.",
		Attributes:          attrs,
	}
}
