package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceDSBuilder = capellaschema.NewSchemaBuilder("appService")

func AppServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceDSBuilder, &schema.StringAttribute{
		Required: true,
	})

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", appServiceDSBuilder, &schema.Int64Attribute{
		Computed: true,
	})
	capellaschema.AddAttr(computeAttrs, "ram", appServiceDSBuilder, &schema.Int64Attribute{
		Computed: true,
	})

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "organization_id", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "cluster_id", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "name", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "description", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "nodes", appServiceDSBuilder, &schema.Int64Attribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "cloud_provider", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "current_state", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "compute", appServiceDSBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computeAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "version", appServiceDSBuilder, &schema.StringAttribute{
		Computed: true,
	})
	dataAttrs["audit"] = computedAuditAttribute

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
