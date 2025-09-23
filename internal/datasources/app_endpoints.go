package datasources

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource                   = (*AppEndpoints)(nil)
	_ datasource.DataSourceWithConfigure      = (*AppEndpoints)(nil)
	_ datasource.DataSourceWithValidateConfig = (*AppEndpoints)(nil)
)

// AppEndpoints is the data source implementation for retrieving App Endpoints for an App Service.
type AppEndpoints struct {
	*providerschema.Data
}

// NewAppEndpoints is a helper function to simplify the provider implementation.
func NewAppEndpoints() datasource.DataSource {
	return &AppEndpoints{}
}

// Metadata returns the App Endpoints data source type name.
func (a *AppEndpoints) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoints"
}

// Schema defines the schema for the App Endpoints data source.
func (a *AppEndpoints) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointsSchema()
}

// Read refreshes the Terraform state with the latest App Endpoints configs.
func (a *AppEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpoints
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = config.OrganizationId.ValueString()
		projectId      = config.ProjectId.ValueString()
		clusterId      = config.ClusterId.ValueString()
		appServiceId   = config.AppServiceId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	appEndpoints, err := api.GetPaginated[[]app_endpoints.GetAppEndpointResponse](ctx, a.Client, a.Token, cfg, api.SortByName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoints",
			fmt.Sprintf(
				"Could not read App Endpoints in cluster %s and App Service %s, unexpected error: %s",
				clusterId,
				appServiceId,
				api.ParseError(err),
			),
		)
		return
	}

	var filteredAppEndpoints []app_endpoints.GetAppEndpointResponse
	var names []string

	// Since the list API doesn't implement query parameters useful for filtering,
	// filtering is done by provider.
	if config.Filters != nil {
		diags := config.Filters.Values.ElementsAs(ctx, &names, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	for _, appEndpoint := range appEndpoints {
		if slices.Contains(names, appEndpoint.Name) || len(names) == 0 { // If no names are provided, include all app endpoints.
			filteredAppEndpoints = append(filteredAppEndpoints, appEndpoint)
		}
	}

	for _, appEndpoint := range filteredAppEndpoints {
		var requireResyncMap types.Map
		if appEndpoint.RequireResync != nil {
			requireResyncMap, diags = types.MapValueFrom(
				ctx,
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"items": types.SetType{ElemType: types.StringType},
					},
				},
				appEndpoint.RequireResync)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
		} else {
			requireResyncMap = types.MapNull(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"items": types.SetType{ElemType: types.StringType},
					},
				},
			)
		}

		var cors *providerschema.AppEndpointCors
		if appEndpoint.Cors != nil {
			originSet, diags := types.SetValueFrom(
				ctx,
				types.StringType,
				appEndpoint.Cors.Origin,
			)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			loginOriginSet, diags := types.SetValueFrom(
				ctx,
				types.StringType,
				appEndpoint.Cors.LoginOrigin,
			)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			headersSet, diags := types.SetValueFrom(
				ctx,
				types.StringType,
				appEndpoint.Cors.Headers,
			)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			cors = &providerschema.AppEndpointCors{
				Origin:      originSet,
				LoginOrigin: loginOriginSet,
				Headers:     headersSet,
				MaxAge:      types.Int64Value(appEndpoint.Cors.MaxAge),
				Disabled:    types.BoolValue(appEndpoint.Cors.Disabled),
			}
		}

		var oidcSet types.Set
		if len(appEndpoint.Oidc) > 0 {
			oidcSet, diags = types.SetValueFrom(
				ctx,
				types.ObjectType{
					AttrTypes: providerschema.
					AppEndpointOidc{}.
						AttributeTypes(),
				},
				appEndpoint.Oidc,
			)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
		} else {
			oidcSet = types.SetNull(
				types.ObjectType{
					AttrTypes: providerschema.
					AppEndpointOidc{}.
						AttributeTypes(),
				},
			)
		}

		var scopesMap types.Map
		if len(appEndpoint.Scopes) > 0 {
			scopesMapElements := make(map[string]attr.Value)

			for scopeName, scope := range appEndpoint.Scopes {
				collectionsMapElements := make(map[string]attr.Value)

				for collectionName, collection := range scope.Collections {
					collectionObj, diags := types.ObjectValueFrom(
						ctx,
						providerschema.AppEndpointCollection{}.AttributeTypes(),
						providerschema.AppEndpointCollection{
							AccessControlFunction: types.StringValue(collection.AccessControlFunction),
							ImportFilter:          types.StringValue(collection.ImportFilter),
						},
					)
					resp.Diagnostics.Append(diags...)
					if diags.HasError() {
						return
					}

					collectionsMapElements[collectionName] = collectionObj
				}

				collectionsMap, diags := types.MapValueFrom(
					ctx,
					types.ObjectType{
						AttrTypes: providerschema.
						AppEndpointCollection{}.
							AttributeTypes(),
					},
					collectionsMapElements,
				)
				resp.Diagnostics.Append(diags...)
				if diags.HasError() {
					return
				}

				scopeObj, diags := types.ObjectValueFrom(
					ctx,
					providerschema.AppEndpointScope{}.AttributeTypes(),
					providerschema.AppEndpointScope{
						Collections: collectionsMap,
					},
				)
				resp.Diagnostics.Append(diags...)
				if diags.HasError() {
					return
				}
				scopesMapElements[scopeName] = scopeObj
			}

			scopesMap, diags = types.MapValueFrom(
				ctx,
				types.ObjectType{
					AttrTypes: providerschema.
					AppEndpointScope{}.
						AttributeTypes(),
				},
				scopesMapElements,
			)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
		}

		ae := providerschema.OneAppEndpoint{
			Bucket:           types.StringValue(appEndpoint.Bucket),
			Name:             types.StringValue(appEndpoint.Name),
			UserXattrKey:     types.StringValue(appEndpoint.UserXattrKey),
			DeltaSyncEnabled: types.BoolValue(appEndpoint.DeltaSyncEnabled),
			AdminURL:         types.StringValue(appEndpoint.AdminURL),
			MetricsURL:       types.StringValue(appEndpoint.MetricsURL),
			PublicURL:        types.StringValue(appEndpoint.PublicURL),
			State:            types.StringValue(appEndpoint.State),
			RequireResync:    requireResyncMap,
			Cors:             cors,
			Oidc:             oidcSet,
			Scopes:           scopesMap,
		}

		config.AppEndpoints = append(config.AppEndpoints, ae)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (a *AppEndpoints) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.Data = data
}

// ValidateConfig checks that if 'name' or 'values' is set in filter block', then both are set.
func (a *AppEndpoints) ValidateConfig(
	ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse,
) {
	var config providerschema.AppEndpoints
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Filters != nil {
		if (config.Filters.Name.IsNull() && !config.Filters.Values.IsNull()) ||
			(!config.Filters.Name.IsNull() && config.Filters.Values.IsNull()) {
			resp.Diagnostics.AddError(
				"Invalid Filters Configuration",
				"Both 'name' and 'values' in filter block must be configured.",
			)
		}
	}
}
