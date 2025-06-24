package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointSchema defines the schema for the AppEndpoint resource
func AppEndpointSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage an App Endpoint configuration for a Couchbase Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":       WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":       WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"bucket":           WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the bucket associated with this App Endpoint."),
			"name":             WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the App Endpoint."),
			"userXattrKey":     WithDescription(stringAttribute([]string{optional}), "The user extended attribute key for the App Endpoint."),
			"deltaSyncEnabled": WithDescription(boolAttribute(optional), "Whether delta sync is enabled for this App Endpoint."),
			"scopes": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for scopes within the App Endpoint.",
				Attributes: map[string]schema.Attribute{
					"scope": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Configuration for the default scope.",
						Attributes: map[string]schema.Attribute{
							"collections": schema.SetNestedAttribute{
								Computed:            true,
								MarkdownDescription: "The list of collections within this scope.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"accessControlFunction": WithDescription(stringAttribute([]string{optional}), "The access control function for this collection."),
										"importFilter":          WithDescription(stringAttribute([]string{optional}), "The import filter for this collection."),
									},
								},
							},
						},
					},
				},
			},
			"cors": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "CORS configuration for the App Endpoint.",
				Attributes: map[string]schema.Attribute{
					"origin": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed origins for CORS.",
					},
					"loginOrigin": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed login origins for CORS.",
					},
					"headers": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed headers for CORS.",
					},
					"maxAge": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Maximum age for CORS preflight requests in seconds.",
					},
					"disabled": WithDescription(boolAttribute(optional), "Whether CORS is disabled for this App Endpoint."),
				},
			},
			"oidc": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of OIDC configurations for the App Endpoint.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer":        WithDescription(stringAttribute([]string{required}), "The OIDC issuer URL."),
						"register":      WithDescription(boolAttribute(optional), "Whether to register the OIDC client."),
						"clientId":      WithDescription(stringAttribute([]string{required}), "The OIDC client ID."),
						"userPrefix":    WithDescription(stringAttribute([]string{optional}), "The user prefix for OIDC users."),
						"discoveryUrl":  WithDescription(stringAttribute([]string{optional}), "The OIDC discovery URL."),
						"usernameClaim": WithDescription(stringAttribute([]string{optional}), "The username claim for OIDC."),
						"rolesClaim":    WithDescription(stringAttribute([]string{optional}), "The roles claim for OIDC."),
						"providerId":    WithDescription(stringAttribute([]string{optional}), "The OIDC provider ID."),
						"isDefault":     WithDescription(boolAttribute(optional), "Whether this is the default OIDC provider."),
					},
				},
			},
			"requireResync": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for require resync settings.",
				Attributes: map[string]schema.Attribute{
					"_default": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Default require resync configuration.",
						Attributes: map[string]schema.Attribute{
							"items": schema.ListAttribute{
								Optional:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "List of items that require resync.",
							},
						},
					},
				},
			},
			"adminURL": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The admin URL for the App Endpoint.",
			},
			"metricsURL": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The metrics URL for the App Endpoint.",
			},
			"publicURL": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public URL for the App Endpoint.",
			},
		},
	}
}
