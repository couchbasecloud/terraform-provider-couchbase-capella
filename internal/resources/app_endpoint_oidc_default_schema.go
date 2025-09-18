package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// AppEndpointDefaultOidcProviderSchema defines the schema for the default OIDC provider resource.
func AppEndpointDefaultOidcProviderSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manage the default OpenID Connect provider for an App Endpoint.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The GUID4 ID of the Capella Organization. Requires replacement if changed.",
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The GUID4 ID of the Capella Project. Requires replacement if changed.",
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The GUID4 ID of the Capella Cluster. Requires replacement if changed.",
			},
			"app_service_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The GUID4 ID of the App Service. Requires replacement if changed.",
			},
			"app_endpoint_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The name of the App Endpoint. Requires replacement if changed.",
			},
			"provider_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The providerId to set as the default for this App Endpoint.",
			},
		},
	}
}
