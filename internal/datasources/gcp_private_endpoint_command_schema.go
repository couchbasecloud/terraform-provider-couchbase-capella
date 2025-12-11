package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var gcpPrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("gcpPrivateEndpointCommand")

// GcpPrivateEndpointCommandSchema returns the schema for the GcpPrivateEndpointCommand data source.
func GcpPrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", gcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", gcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", gcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vpc_network_id", gcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "subnet_ids", gcpPrivateEndpointCommandBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(attrs, "command", gcpPrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate a GCP CLI command for setting up a private endpoint connection to an operational cluster.",
		Attributes:          attrs,
	}
}
