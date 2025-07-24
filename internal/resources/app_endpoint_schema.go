package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointSchema defines the schema for the AppEndpoint resource.
func AppEndpointSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage an App Endpoint configuration for a Couchbase Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id":    WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":     WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the App Service."),
			"bucket":             WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the bucket associated with this App Endpoint."),
			"name":               WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the App Endpoint."),
			"user_xattr_key":     WithDescription(stringAttribute([]string{optional}), "The user extended attribute key for the App Endpoint."),
			"delta_sync_enabled": WithDescription(boolAttribute(optional), "States whether delta sync is enabled for this App Endpoint."),

			"scope": stringAttribute([]string{optional}),

			"collections": schema.MapNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for collections within the App Endpoint. The map key is the collection name.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_control_function": WithDescription(stringAttribute([]string{optional}), "The Javascript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function."),
						"import_filter":           WithDescription(stringAttribute([]string{optional}), "The JavaScript function used to filter which documents in the collection that are to be imported by the App Endpoint."),
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
					"max_age": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Maximum age, in seconds, for CORS preflight requests in seconds.",
					},
					"disabled": schema.BoolAttribute{Optional: true, Description: "Disables/Enables CORS for this App Endpoint."},
				},
			},

			"oidc": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of OIDC configurations for the App Endpoint.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer":         WithDescription(stringAttribute([]string{required}), "The OIDC issuer URL."),
						"register":       WithDescription(boolAttribute(optional), "States whether to register the OIDC client."),
						"client_id":      WithDescription(stringAttribute([]string{required}), "The OIDC client ID."),
						"user_prefix":    WithDescription(stringAttribute([]string{optional}), "The user prefix for OIDC users."),
						"discovery_url":  WithDescription(stringAttribute([]string{optional}), "The OIDC discovery URL."),
						"username_claim": WithDescription(stringAttribute([]string{optional}), "The username claim for OIDC."),
						"roles_claim":    WithDescription(stringAttribute([]string{optional}), "The roles claim for OIDC."),
						"provider_id":    WithDescription(stringAttribute([]string{computed}), "The OIDC provider ID."),
						"is_default":     WithDescription(boolAttribute(computed), "States whether this is the default OIDC provider."),
					},
				},
			},

			"require_resync": schema.MapNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of collections that require resync, keyed by scope.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"items": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},

			"admin_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The admin URL for the App Endpoint.",
			},
			"metrics_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The metrics URL for the App Endpoint.",
			},
			"public_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public URL for the App Endpoint.",
			},
		},
	}
}
