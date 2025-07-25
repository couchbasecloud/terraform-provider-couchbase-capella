package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppEndpoint)(nil)
	_ datasource.DataSourceWithConfigure = (*AppEndpoint)(nil)
)

// AppServiceCidrs is the data source implementation for retrieving allowed CIDRs for an App Service.
type AppEndpoint struct {
	*providerschema.Data
}

func (a *AppEndpoint) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	//TODO implement me
}

func (a *AppEndpoint) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	//TODO implement me
}

// NewAppServiceCidrs is used in (p *capellaProvider) DataSources for building the provider.
func NewAppEndpoint() datasource.DataSource {
	return &AppEndpoint{}
}

// Metadata returns the App Service CIDRs data source type name.
func (a *AppEndpoint) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

// AppEndpointsSchema defines the schema for the AppEndpoints datasource.
func (a *AppEndpoint) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
							MarkdownDescription: "The list of scopes in this App Endpoint. Currently, only one scope can be linked per App Endpoint.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"collections": schema.MapNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The collection configuration defines access control, validation functions, and import filters for a specific collection. The key of the collection configuration object is the name of the collection.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"access_control_function": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The Javascript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function.",
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
								"origin": schema.ListAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed origins for CORS.",
								},
								"login_origin": schema.ListAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of allowed login origins for CORS.",
								},
								"headers": schema.ListAttribute{
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
						"require_resync": schema.MapNestedAttribute{
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
	}
}
