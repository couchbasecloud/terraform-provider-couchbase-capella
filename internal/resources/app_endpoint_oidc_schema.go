package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AppEndpointOidcProviderSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "AppEndpointOidcProvider Resource",
		Attributes: map[string]schema.Attribute{
			"app_endpoint_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The name of the App Endpoint. Requires replacement if changed.`,
			},
			"app_service_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the App Service. Requires replacement if changed.`,
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: `The OpenID Connect provider client ID.`,
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the Capella Cluster. Requires replacement if changed.`,
			},
			"discovery_url": schema.StringAttribute{
				Optional:    true,
				Description: `The URL for the non-standard discovery endpoint.`,
			},
			"issuer": schema.StringAttribute{
				Required:    true,
				Description: `The URL for the OpenID Connect issuer.`,
			},
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the Capella Organization. Requires replacement if changed.`,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the Capella Project. Requires replacement if changed.`,
			},
			"provider_id": schema.StringAttribute{
				Computed: true,
			},
			"register": schema.BoolAttribute{
				Optional:    true,
				Description: `Indicates whether to register a new App Service user account when a user logs in using OpenID Connect.`,
			},
			"roles_claim": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: `If set, the value(s) of the given OpenID Connect authentication token claim will be added to the user's roles.` + "\n" +
					`The value of this claim in the OIDC token must be either a string or an array of strings, any other type will result in an error.`,
			},
			"user_prefix": schema.StringAttribute{
				Optional:    true,
				Description: `Username prefix for all users created for this provider.`,
			},
			"username_claim": schema.StringAttribute{
				Optional:    true,
				Description: `Allows a different OpenID Connect field to be specified instead of the Subject (sub).`,
			},
		},
	}
}
