package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointSchema defines the schema for the app endpoint resource.
func AppEndpointSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage an App Endpoint configuration for a Couchbase Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id":    WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the organization."),
			"project_id":         WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the project."),
			"cluster_id":         WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the cluster."),
			"app_service_id":     WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the App Service."),
			"bucket":             WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The name of the bucket associated with this App Endpoint."),
			"name":               WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The name of the App Endpoint."),
			"user_xattr_key":     WithDescription(stringAttribute([]string{optional, computed, useStateForUnknown}), "The user extended attribute key for the App Endpoint."),
			"delta_sync_enabled": WithDescription(boolAttribute(optional, computed, useStateForUnknown), "If delta sync is enabled for this App Endpoint."),
			"scopes": schema.MapNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for scopes within the App Endpoint.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"collections": schema.MapNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Configuration for collections within the App Endpoint.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"access_control_function": WithDescription(
										stringAttribute([]string{optional, computed, useStateForUnknown}),
										"The Javascript function that is used to specify the access control policies to be applied to documents in this collection."),
									"import_filter": WithDescription(
										stringAttribute([]string{optional, computed, useStateForUnknown}),
										"The JavaScript function used to filter which documents in the collection that are to be imported by the App Endpoint."),
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
					"origin": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed origins for CORS.",
					},
					"login_origin": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed login origins for CORS.",
					},
					"headers": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of allowed headers for CORS.",
					},
					"max_age": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The maximum age for CORS preflight requests in seconds. Default is 0, which means no caching.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"disabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Disables/Enables CORS for this App Endpoint.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"oidc": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of OIDC configurations for the App Endpoint.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer":         WithDescription(stringAttribute([]string{required}), "The URL for the OpenID Connect issuer.."),
						"register":       WithDescription(boolAttribute(optional, computed), "Indicates whether to register a new App Service user account when a user logs in using OpenID Connect."),
						"client_id":      WithDescription(stringAttribute([]string{required}), "The OpenID Connect provider client ID."),
						"user_prefix":    WithDescription(stringAttribute([]string{optional, computed}), "Username prefix for all users created for this provider."),
						"discovery_url":  WithDescription(stringAttribute([]string{optional, computed}), "The URL for the non-standard discovery endpoint."),
						"username_claim": WithDescription(stringAttribute([]string{optional, computed}), "Allows a different OpenID Connect field to be specified instead of the Subject (sub)."),
						"roles_claim":    WithDescription(stringAttribute([]string{optional, computed}), "If set, the value(s) of the given OpenID Connect authentication token claim will be added to the user's roles. The value of this claim in the OIDC token must be either a string or an array of strings, any other type will result in an error."),
						"provider_id":    WithDescription(stringAttribute([]string{computed}), "The GUID4 ID of this OpenID Connect provider. Generated by the backend on creation."),
						"is_default":     WithDescription(boolAttribute(computed), "Indicates whether this is the default OpenID Connect provider for this App Endpoint."),
					},
				},
			},
			"require_resync": schema.MapNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of collections that require resync, keyed by scope.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"items": schema.SetAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of collections that require resync.",
							PlanModifiers: []planmodifier.Set{
								setplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current state of the App Endpoint, such as online, offline, resyncing, etc.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The admin URL for the App Endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metrics_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The metrics URL for the App Endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public URL for the App Endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
