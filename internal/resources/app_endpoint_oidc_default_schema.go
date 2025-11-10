package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointOidcDefaultBuilder = capellaschema.NewSchemaBuilder("appEndpointOidcDefault", "OIDCProvider")

// AppEndpointDefaultOidcProviderSchema defines the schema for the default OIDC provider resource.
func AppEndpointDefaultOidcProviderSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointOidcDefaultBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", appEndpointOidcDefaultBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointOidcDefaultBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointOidcDefaultBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointOidcDefaultBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "provider_id", appEndpointOidcDefaultBuilder, stringAttribute([]string{required}, validator.String(stringvalidator.LengthAtLeast(1))))

	return schema.Schema{
		MarkdownDescription: "Manage the default OpenID Connect provider for an App Endpoint.",
		Attributes:          attrs,
	}
}
