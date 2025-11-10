package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointOidcBuilder = capellaschema.NewSchemaBuilder("appEndpointOidc")

func AppEndpointOidcProviderSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointOidcBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", appEndpointOidcBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointOidcBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointOidcBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointOidcBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "issuer", appEndpointOidcBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "register", appEndpointOidcBuilder, boolAttribute(optional, computed))
	capellaschema.AddAttr(attrs, "client_id", appEndpointOidcBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "user_prefix", appEndpointOidcBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "discovery_url", appEndpointOidcBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "username_claim", appEndpointOidcBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "roles_claim", appEndpointOidcBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "provider_id", appEndpointOidcBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "is_default", appEndpointOidcBuilder, boolAttribute(computed))

	return schema.Schema{
		MarkdownDescription: "App Endpoint OpenID Connect Provider Resource",
		Attributes:          attrs,
	}
}
