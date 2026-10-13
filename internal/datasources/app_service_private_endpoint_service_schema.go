package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpointService", "GetAppServicePrivateEndpointStateResponse")

// AppServicePrivateEndpointServiceSchema returns the schema for the app_service_private_endpoint_service data source.
func AppServicePrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointServiceBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointServiceBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointServiceBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointServiceBuilder, requiredUUIDString())

	awsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsAttrs, "vpc_id", appServicePrivateEndpointServiceBuilder, &schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The ID of your AWS VPC.",
	}, "CreateVPCEndpointCommandRequest")
	capellaschema.AddAttr(awsAttrs, "subnet_ids", appServicePrivateEndpointServiceBuilder, &schema.SetAttribute{
		Required:            true,
		ElementType:         types.StringType,
		MarkdownDescription: "The IDs of the AWS subnets in which to create the private endpoint.",
	}, "CreateVPCEndpointCommandRequest")
	capellaschema.AddAttr(attrs, "aws", appServicePrivateEndpointServiceBuilder, &schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The AWS network details. Set this to also return the AWS CLI command that creates the private endpoint.",
		Attributes:          awsAttrs,
	})

	capellaschema.AddAttr(attrs, "state", appServicePrivateEndpointServiceBuilder, computedString())
	capellaschema.AddAttr(attrs, "target_state", appServicePrivateEndpointServiceBuilder, computedString())
	capellaschema.AddAttr(attrs, "command", appServicePrivateEndpointServiceBuilder, computedString(), "CreatePrivateEndpointCommandResponse")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the private endpoint service state for an App Service. If the `aws` block is set, it also returns the CLI command to create the private endpoint in your AWS account.",
		Attributes:          attrs,
	}
}
