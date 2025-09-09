package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
				Description: `The GUID4 ID of the appService. Requires replacement if changed.`,
			},
			"client_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The OpenID Connect provider client ID. Requires replacement if changed.`,
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the cluster. Requires replacement if changed.`,
			},
			"discovery_url": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The URL for the non-standard discovery endpoint. Requires replacement if changed.`,
			},
			"issuer": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The URL for the OpenID Connect issuer. Requires replacement if changed.`,
			},
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the organization. Requires replacement if changed.`,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The GUID4 ID of the project. Requires replacement if changed.`,
			},
			"provider_id": schema.StringAttribute{
				Computed: true,
			},
			"register": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Indicates whether to register a new App Service user account when a user logs in using OpenID Connect. Requires replacement if changed.`,
			},
			"roles_claim": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: `If set, the value(s) of the given OpenID Connect authentication token claim will be added to the user's roles.` + "\n" +
					`The value of this claim in the OIDC token must be either a string or an array of strings, any other type will result in an error.` + "\n" +
					`Requires replacement if changed.`,
			},
			"user_prefix": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Username prefix for all users created for this provider. Requires replacement if changed.`,
			},
			"username_claim": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Allows a different OpenID Connect field to be specified instead of the Subject (sub). Requires replacement if changed.`,
			},
		},
	}
}
