package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAzurePrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("appServiceAzurePrivateEndpointCommand")

// AppServiceAzurePrivateEndpointCommandSchema returns the schema for the
// app_service_azure_private_endpoint_command data source.
func AppServiceAzurePrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAzurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAzurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAzurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAzurePrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "resource_group_name", appServiceAzurePrivateEndpointCommandBuilder, requiredString(), "CreateAzurePrivateEndpointCommandRequest")
	capellaschema.AddAttr(attrs, "virtual_network", appServiceAzurePrivateEndpointCommandBuilder, requiredString(), "CreateAzurePrivateEndpointCommandRequest")
	capellaschema.AddAttr(attrs, "command", appServiceAzurePrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate an Azure CLI script for setting up a private endpoint connection to an App Service (Sync Gateway). Retrieves the command or script to create the private endpoint, which establishes a private connection between the specified virtual network and the App Service's private endpoint service.\n\n" +
			"~> **AWS-only as of this writing:** Couchbase Capella's private endpoints for App Services are AWS-only; Azure support is not confirmed GA. This data source is included for API/schema parity with the cluster-level `azure_private_endpoint_command` resource, but has not been validated against a live Azure-backed App Service.",
		Attributes: attrs,
	}
}
