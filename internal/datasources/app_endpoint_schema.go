package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointsSchema defines the schema for the AppEndpoints datasource.
func AppEndpointsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source retrieves App Endpoint configurations for an App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"app_service_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the App Service.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of App Endpoint configurations.",
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
						"userXattrKey": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user extended attribute key for the App Endpoint.",
						},
						"deltaSyncEnabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Enable or disable delta sync on this App Endpoint.",
						},
						"scopes": schema.MapNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The list of scopes in this App Endpoint. Currently, only one scope can be linked per App Endpoint.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"collections": schema.MapNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The collection configuration defines access control, validation functions, and import filters for a specific collection. The key of the collection configuration object is the name of the collection.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"accessControlFunction": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The Javascript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function.",
												},
												"importFilter": schema.StringAttribute{
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
								"origin": schema.ListAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed origins for CORS.",
								},
								"loginOrigin": schema.ListAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed login origins for CORS.",
								},
								"headers": schema.ListAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed headers for CORS.",
								},
								"maxAge": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Maximum age for CORS preflight requests in seconds.",
								},
								"disabled": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether CORS is disabled for this App Endpoint.",
								},
							},
						},
						"oidc": schema.ListNestedAttribute{
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
									"clientId": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC client ID.",
									},
									"userPrefix": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The user prefix for OIDC users.",
									},
									"discoveryUrl": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC discovery URL.",
									},
									"usernameClaim": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The username claim for OIDC.",
									},
									"rolesClaim": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The roles claim for OIDC.",
									},
									"providerId": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OIDC provider ID.",
									},
									"isDefault": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether this is the default OIDC provider.",
									},
								},
							},
						},
						"requireResync": schema.MapNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of collections that require resync, keyed by scope.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"items": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "List of collections that require resync under this scope.",
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
				},
			},
		},
	}
}
