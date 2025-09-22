package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointsSchema defines the schema for the AppEndpoints datasource.
func AppEndpointSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source retrieves App Endpoint configurations for an App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"app_service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the App Service.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"app_endpoints": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of App Endpoints.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bucket": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the bucket associated with this App Endpoint.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the App Endpoint.",
						},
						"user_xattr_key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user extended attribute key for the App Endpoint.",
						},
						"delta_sync_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Enable or disable delta sync on this App Endpoint.",
						},
						"scopes": schema.MapNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration for scopes within the App Endpoint.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"collections": schema.MapNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Configuration for collections within the App Endpoint.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"access_control_function": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The Javascript function that is used to specify the access control policies to be applied to documents in this collection.",
												},
												"import_filter": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The JavaScript function used to filter which documents in the collection that are to be imported by the App Endpoint.",
												},
											},
										},
									},
								},
							},
						},
						"cors": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "CORS configuration for the App Endpoint.",
							Attributes: map[string]schema.Attribute{
								"origin": schema.SetAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed origins for CORS.",
								},
								"login_origin": schema.SetAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed login origins for CORS.",
								},
								"headers": schema.SetAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed headers for CORS.",
								},
								"max_age": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Maximum age for CORS preflight requests in seconds.",
								},
								"disabled": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether CORS is disabled for this App Endpoint.",
								},
							},
						},
						"oidc": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of OIDC configurations for the App Endpoint.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"issuer": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC issuer URL.",
									},
									"register": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether to register the OIDC client.",
									},
									"client_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC client ID.",
									},
									"user_prefix": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The user prefix for OIDC users.",
									},
									"discovery_url": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC discovery URL.",
									},
									"username_claim": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The username claim for OIDC.",
									},
									"roles_claim": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The roles claim for OIDC.",
									},
									"provider_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC provider ID.",
									},
									"is_default": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether this is the default OIDC provider.",
									},
								},
							},
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The state of the App Endpoint. Possible values include `online`, `offline` and `resyncing`.",
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
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				MarkdownDescription: "Filter criteria for App Endpoints.  Only filtering by App Endpoint name is supported.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the attribute to filter.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("name"),
						},
					},
					"values": schema.SetAttribute{
						MarkdownDescription: "List of values to match against.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
		},
	}
}
