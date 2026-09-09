package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointsBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpoints")

// AppServicePrivateEndpointsSchema returns the schema for the app_service_private_endpoints data source.
func AppServicePrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointsBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServicePrivateEndpointsBuilder, computedString(), "PrivateEndpoint")
	capellaschema.AddAttr(dataAttrs, "status", appServicePrivateEndpointsBuilder, computedString(), "PrivateEndpoint")
	capellaschema.AddAttr(dataAttrs, "service_name", appServicePrivateEndpointsBuilder, computedString(), "PrivateEndpoint")

	capellaschema.AddAttr(attrs, "data", appServicePrivateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve private endpoints for an App Service (Sync Gateway).",
		Attributes:          attrs,
	}
}
