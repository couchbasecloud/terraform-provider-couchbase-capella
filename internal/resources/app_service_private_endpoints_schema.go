package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointsBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpoints")

func AppServicePrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointsBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "endpoint_id", appServicePrivateEndpointsBuilder, stringAttribute(
		[]string{required, requiresReplace},
		validator.String(stringvalidator.LengthAtLeast(1)),
	))
	capellaschema.AddAttr(attrs, "status", appServicePrivateEndpointsBuilder, stringAttribute([]string{computed}), "PrivateEndpoint")
	capellaschema.AddAttr(attrs, "service_name", appServicePrivateEndpointsBuilder, stringAttribute([]string{computed}), "PrivateEndpoint")

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage private endpoints for an App Service (Sync Gateway). Private endpoints allow you to securely connect your Cloud Service Provider's private network (VPC/VNET) to your App Service without exposing traffic to the public internet.\n\n" +
			"~> **AWS-only as of this writing:** Couchbase Capella's private endpoints for App Services are AWS-only; Azure/GCP support is not confirmed GA. This resource's schema does not restrict cloud provider, but it has only been validated against an AWS-backed App Service.",
		Attributes: attrs,
	}
}
