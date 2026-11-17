package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpointService", "GetAppServicePrivateEndpointStateResponse")

// AppServicePrivateEndpointServiceSchema returns the schema for the App Service private endpoint service resource.
func AppServicePrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "state", appServicePrivateEndpointServiceBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to enable the private endpoint service on an App Service. " +
			"The private endpoint service must be enabled before you can create private endpoints to connect your Cloud Service Provider's private network (VPC/VNET) to your App Service. " +
			"Destroying this resource disables the private endpoint service.",
		Attributes: attrs,
	}
}
