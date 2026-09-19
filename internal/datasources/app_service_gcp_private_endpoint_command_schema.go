package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceGcpPrivateEndpointCommandBuilder = capellaschema.NewSchemaBuilder("appServiceGcpPrivateEndpointCommand")

// AppServiceGcpPrivateEndpointCommandSchema returns the schema for the
// app_service_gcp_private_endpoint_command data source.
func AppServiceGcpPrivateEndpointCommandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceGcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appServiceGcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceGcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceGcpPrivateEndpointCommandBuilder, requiredString())
	capellaschema.AddAttr(attrs, "vpc_network_id", appServiceGcpPrivateEndpointCommandBuilder, requiredString(), "CreateGCPPrivateEndpointCommandRequest")
	capellaschema.AddAttr(attrs, "subnet_ids", appServiceGcpPrivateEndpointCommandBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(attrs, "command", appServiceGcpPrivateEndpointCommandBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to generate a GCP CLI command for setting up a private endpoint connection to an App Service (Sync Gateway).\n\n" +
			"~> **AWS-only as of this writing:** Couchbase Capella's private endpoints for App Services are AWS-only; GCP support is not confirmed GA. This data source is included for API/schema parity with the cluster-level `gcp_private_endpoint_command` resource, but has not been validated against a live GCP-backed App Service.",
		Attributes: attrs,
	}
}
