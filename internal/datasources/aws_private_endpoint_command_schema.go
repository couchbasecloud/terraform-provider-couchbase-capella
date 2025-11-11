package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var awsPrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("awsPrivateEndpointCommand")

// AwsPrivateEndpointCommandSchema returns the schema for the AwsPrivateEndpointCommand data source.
func AwsPrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", awsPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", awsPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", awsPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vpc_id", awsPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "subnet_ids", awsPrivateEndpointCommandBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(attrs, "command", awsPrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate an AWS CLI command for setting up a private endpoint connection to an operational cluster. Retrieves the command or script to create the private endpoint, which establishes a private connection between the specified VPC and the Capella private endpoint service.",
		Attributes:          attrs,
	}
}
