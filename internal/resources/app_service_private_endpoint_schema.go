package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpoint", "PrivateEndpoint")

// AppServicePrivateEndpointSchema returns the schema for the App Service private endpoint resource.
func AppServicePrivateEndpointSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "endpoint_id", appServicePrivateEndpointBuilder, stringAttribute(
		[]string{required, requiresReplace},
		validator.String(stringvalidator.LengthAtLeast(1)),
	))
	capellaschema.AddAttr(attrs, "service_name", appServicePrivateEndpointBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "status", appServicePrivateEndpointBuilder, stringAttribute([]string{computed, useStateForUnknown}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to accept a private endpoint connection request on an App Service. " +
			"Private endpoints allow you to securely connect your Cloud Service Provider's private network (VPC/VNET) to your App Service without exposing traffic to the public internet. " +
			"The private endpoint service must be enabled on the App Service first. " +
			"Destroying this resource removes the private endpoint from the App Service.",
		Attributes: attrs,
	}
}
