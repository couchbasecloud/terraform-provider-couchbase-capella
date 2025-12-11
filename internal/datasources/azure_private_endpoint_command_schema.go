package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var azurePrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("azurePrivateEndpointCommand")

// AzurePrivateEndpointCommandSchema returns the schema for the AzurePrivateEndpointCommand data source.
func AzurePrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", azurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", azurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", azurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "resource_group_name", azurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "virtual_network", azurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "command", azurePrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate an Azure CLI command for setting up a private endpoint connection to an operational cluster.",
		Attributes:          attrs,
	}
}
