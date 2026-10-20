package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAWSPrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("appServiceAWSPrivateEndpointCommand")

// AppServiceAWSPrivateEndpointCommandSchema returns the schema for the
// app_service_aws_private_endpoint_command data source.
func AppServiceAWSPrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAWSPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appServiceAWSPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAWSPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAWSPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vpc_id", appServiceAWSPrivateEndpointCommandBuilder, requiredString(), "CreateVPCEndpointCommandRequest")
	capellaschema.AddAttr(attrs, "subnet_ids", appServiceAWSPrivateEndpointCommandBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(attrs, "command", appServiceAWSPrivateEndpointCommandBuilder, computedString())
	capellaschema.AddAttr(attrs, "service_name", appServiceAWSPrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate an AWS CLI command for setting up a private endpoint connection to an App Service (Sync Gateway). Retrieves the command or script to create the private endpoint, which establishes a private connection between the specified VPC and the App Service's private endpoint service.\n\n" +
			"~> **`service_name`** is parsed best-effort from the `--service-name` flag in `command`. It is not returned by the API directly (unlike the operational cluster's private endpoint service, the App Service private endpoint status API has no service name field) — if Couchbase changes the command's phrasing, this attribute is left null rather than failing the read; `command` is unaffected either way.",
		Attributes: attrs,
	}
}
