package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AppEndpointDefaultOidcProviderSchema defines the schema for the default OIDC provider resource.
func AppEndpointDefaultOidcProviderSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manage the default OpenID Connect provider for an App Endpoint.",
		Attributes: map[string]schema.Attribute{
			"organization_id":   WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the Capella organization."),
			"project_id":        WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the Capella project."),
			"cluster_id":        WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the Capella cluster."),
			"app_service_id":    WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the Capella App Service."),
			"app_endpoint_name": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The name of the App Endpoint."),
			"provider_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The providerId to set as the default for this App Endpoint.",
			},
		},
	}
}
