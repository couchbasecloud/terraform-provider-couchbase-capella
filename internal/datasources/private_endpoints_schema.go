package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var privateEndpointsBuilder = capellaschema.NewSchemaBuilder("privateEndpoints")

// PrivateEndpointsSchema returns the schema for the PrivateEndpoints data source.
func PrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", privateEndpointsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", privateEndpointsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", privateEndpointsBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", privateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", privateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", privateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", privateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", privateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", privateEndpointsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", privateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve private endpoints for a cluster.",
		Attributes:          attrs,
	}
}
