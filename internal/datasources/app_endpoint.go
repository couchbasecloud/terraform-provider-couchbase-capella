package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppEndpoint)(nil)
	_ datasource.DataSourceWithConfigure = (*AppEndpoint)(nil)
)

// AppEndpoint is the data source implementation for retrieving a single App Endpoint for an App Service.
type AppEndpoint struct {
	*providerschema.Data
}

// NewAppEndpoint is used in (p *capellaProvider) DataSources for building the provider.
func NewAppEndpoint() datasource.DataSource {
	return &AppEndpoint{}
}

// Metadata returns the App Endpoint data source type name.
func (a *AppEndpoint) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint"
}

// Schema defines the schema for the App Endpoint data source.
func (a *AppEndpoint) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointSchema()
}

// Read refreshes the Terraform state with the latest App Endpoint config.
func (a *AppEndpoint) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpoint
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
		endpointName   = config.Name.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointName,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint",
			fmt.Sprintf("Could not read App Endpoint %s: %s", endpointName, api.ParseError(err)),
		)
		return
	}

	var appEndpoint app_endpoints.GetAppEndpointResponse
	if err := json.Unmarshal(response.Body, &appEndpoint); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing App Endpoint response",
			fmt.Sprintf("Could not unmarshal App Endpoint response: %s", err.Error()),
		)
		return
	}

	state := &providerschema.AppEndpoint{}
	if appEndpoint.Cors != nil {
		state = &providerschema.AppEndpoint{
			Cors: &providerschema.AppEndpointCors{},
		}
	}

	// Set computed attributes
	state.AdminURL = types.StringValue(appEndpoint.AdminURL)
	state.PublicURL = types.StringValue(appEndpoint.PublicURL)
	state.MetricsURL = types.StringValue(appEndpoint.MetricsURL)
	state.State = types.StringValue(appEndpoint.State)
	if len(appEndpoint.RequireResync) > 0 {
		requireResyncMap, diags := types.MapValueFrom(
			ctx,
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"items": types.SetType{ElemType: types.StringType},
				},
			},
			appEndpoint.RequireResync)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.RequireResync = requireResyncMap
	} else {
		state.RequireResync = types.MapNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"items": types.SetType{ElemType: types.StringType},
			},
		})
	}

	state.OrganizationId = types.StringValue(organizationId)
	state.ProjectId = types.StringValue(projectId)
	state.ClusterId = types.StringValue(clusterId)
	state.AppServiceId = types.StringValue(appServiceId)

	state.Bucket = types.StringValue(appEndpoint.Bucket)
	state.Name = types.StringValue(appEndpoint.Name)
	state.DeltaSyncEnabled = types.BoolValue(appEndpoint.DeltaSyncEnabled)
	state.UserXattrKey = types.StringValue(appEndpoint.UserXattrKey)

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
				if resp.Diagnostics.HasError() {
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
			if resp.Diagnostics.HasError() {
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
			if resp.Diagnostics.HasError() {
				return
			}
			scopesMapElements[scopeName] = scopeObj
		}

		scopesMap, diags := types.MapValueFrom(
			ctx,
			types.ObjectType{
				AttrTypes: providerschema.
					AppEndpointScope{}.
					AttributeTypes(),
			},
			scopesMapElements,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Scopes = scopesMap
	}

	if appEndpoint.Cors != nil {
		state.Cors.Disabled = types.BoolValue(appEndpoint.Cors.Disabled)
		state.Cors.MaxAge = types.Int64Value(appEndpoint.Cors.MaxAge)

		originSet, diags := types.SetValueFrom(
			ctx,
			types.StringType,
			appEndpoint.Cors.Origin,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Cors.Origin = originSet

		loginOriginSet, diags := types.SetValueFrom(
			ctx,
			types.StringType,
			appEndpoint.Cors.LoginOrigin,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Cors.LoginOrigin = loginOriginSet

		headersSet, diags := types.SetValueFrom(
			ctx,
			types.StringType,
			appEndpoint.Cors.Headers,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Cors.Headers = headersSet
	}
	var oidcSet []providerschema.AppEndpointOidc
	if len(appEndpoint.Oidc) > 0 {
		for _, oidc := range appEndpoint.Oidc {
			oidcSet = append(oidcSet, providerschema.AppEndpointOidc{
				Issuer:        types.StringValue(oidc.Issuer),
				ClientId:      types.StringValue(oidc.ClientId),
				DiscoveryUrl:  types.StringValue(oidc.DiscoveryUrl),
				UsernameClaim: types.StringValue(oidc.UsernameClaim),
				UserPrefix:    types.StringValue(oidc.UserPrefix),
				RolesClaim:    types.StringValue(oidc.RolesClaim),
				ProviderId:    types.StringValue(oidc.ProviderId),
				IsDefault:     types.BoolValue(oidc.IsDefault),
				Register:      types.BoolValue(oidc.Register),
			})
		}
	}
	state.Oidc = oidcSet

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Configure adds the provider configured client to the App Endpoint data source.
func (a *AppEndpoint) Configure(
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
