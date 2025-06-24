package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpointsSchema defines the schema for the AppEndpoints datasource
func AppEndpointsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source retrieves App Endpoint configurations for a Couchbase Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
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
							MarkdownDescription: "Whether delta sync is enabled for this App Endpoint.",
						},
						"scopes": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration for scopes within the App Endpoint.",
							Attributes: map[string]schema.Attribute{
								"_default": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Configuration for the default scope.",
									Attributes: map[string]schema.Attribute{
										"collections": schema.SingleNestedAttribute{
											Computed:            true,
											MarkdownDescription: "Configuration for collections within the default scope.",
											Attributes: map[string]schema.Attribute{
												"_default": schema.SingleNestedAttribute{
													Computed:            true,
													MarkdownDescription: "Configuration for the default collection.",
													Attributes: map[string]schema.Attribute{
														"accessControlFunction": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The access control function for the default collection.",
														},
														"importFilter": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The import filter for the default collection.",
														},
													},
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
						"requireResync": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration for require resync settings.",
							Attributes: map[string]schema.Attribute{
								"_default": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Default require resync configuration.",
									Attributes: map[string]schema.Attribute{
										"items": schema.ListAttribute{
											Computed:            true,
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
				},
			},
		},
	}
}
