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
			"organization_id":    WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"bucket":             WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the bucket associated with this app endpoint."),
			"name":               WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the app endpoint."),
			"user_xattr_key":     WithDescription(stringAttribute([]string{optional}), "The user extended attribute key for the app endpoint."),
			"delta_sync_enabled": WithDescription(boolAttribute(optional), "Whether delta sync is enabled for this app endpoint."),
			"scopes": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for scopes within the app endpoint.",
				Attributes: map[string]schema.Attribute{
					"_default": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Configuration for the default scope.",
						Attributes: map[string]schema.Attribute{
							"collections": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Configuration for collections within the default scope.",
								Attributes: map[string]schema.Attribute{
									"_default": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Configuration for the default collection.",
										Attributes: map[string]schema.Attribute{
											"access_control_function": WithDescription(stringAttribute([]string{optional}), "The access control function for the default collection."),
											"import_filter":           WithDescription(stringAttribute([]string{optional}), "The import filter for the default collection."),
										},
									},
								},
							},
						},
					},
				},
			},
			"cors": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "CORS configuration for the app endpoint.",
				Attributes: map[string]schema.Attribute{
					"origin": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed origins for CORS.",
					},
					"login_origin": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed login origins for CORS.",
					},
					"headers": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed headers for CORS.",
					},
					"disabled": WithDescription(boolAttribute(optional), "Whether CORS is disabled for this app endpoint."),
				},
			},
			"oidc": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of OIDC configurations for the app endpoint.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer":         WithDescription(stringAttribute([]string{required}), "The OIDC issuer URL."),
						"register":       WithDescription(boolAttribute(optional), "Whether to register the OIDC client."),
						"client_id":      WithDescription(stringAttribute([]string{required}), "The OIDC client ID."),
						"user_prefix":    WithDescription(stringAttribute([]string{optional}), "The user prefix for OIDC users."),
						"discovery_url":  WithDescription(stringAttribute([]string{optional}), "The OIDC discovery URL."),
						"username_claim": WithDescription(stringAttribute([]string{optional}), "The username claim for OIDC."),
						"roles_claim":    WithDescription(stringAttribute([]string{optional}), "The roles claim for OIDC."),
					},
				},
			},
		},
	}
}
