package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var privateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("privateEndpointService")

// PrivateEndpointServiceSchema returns the schema for the PrivateEndpointService data source.
func PrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", privateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", privateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", privateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "enabled", privateEndpointServiceBuilder, computedBool())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve private endpoint service information for a cluster.",
		Attributes:          attrs,
	}
}
