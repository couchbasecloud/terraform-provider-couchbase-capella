package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var networkPeerCommandAzureBuilder = capellaschema.NewSchemaBuilder("networkPeerCommandAzure")

// NetworkPeerCommandAzureSchema returns the schema for the NetworkPeerCommandAzure data source.
func NetworkPeerCommandAzureSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "tenant_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "subscription_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "resource_group", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vnet_id", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vnet_peering_service_principal", networkPeerCommandAzureBuilder, requiredString())
	capellaschema.AddAttr(attrs, "command", networkPeerCommandAzureBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate an Azure CLI command for setting up a network peering connection to a cluster.",
		Attributes:          attrs,
	}
}
