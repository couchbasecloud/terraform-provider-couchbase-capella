package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpoint{}
	_ resource.ResourceWithConfigure   = &AppEndpoint{}
	_ resource.ResourceWithImportState = &AppEndpoint{}
)

const errorAppEndpointCreation = "There is an error during App Endpoint creation. unexpected error: "

const errorAppEndpointRefresh = "App Endpoint was created, but terraform encountered an error while checking the current" +
	" state of the App Endpoint. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// AppEndpoint is the AppEndpoint implementation.
type AppEndpoint struct {
	*providerschema.Data
}

// NewAppEndpoint is a helper function to simplify the provider implementation.
func NewAppEndpoint() resource.Resource {
	return &AppEndpoint{}
}

// Metadata returns the App Endpoint resource type name.
func (a *AppEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint"
}

// Schema defines the schema for AppEndpoint.
func (a *AppEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointSchema()
}

// Create creates a new App Endpoint.
func (a *AppEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpoint
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		endpointName   = plan.Name.ValueString()
	)

	appEndpointRequest, diags := morphToAppEndpointRequest(ctx, &plan)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

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
		Method:        http.MethodPost,
		SuccessStatus: http.StatusCreated,
	}

	if _, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		appEndpointRequest,
		a.Token,
		nil,
	); err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorAppEndpointCreation+err.Error(),
		)
		return
	}

	diags = initComputedAttributesToNullBeforeRefresh(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.refreshAppEndpoint(
		ctx,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointName,
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error refreshing App Endpoint",
			errorAppEndpointRefresh+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)

}

// initComputedAttributesToNullBeforeRefresh inits computed attributes to null before refreshing App Endpoint.
func initComputedAttributesToNullBeforeRefresh(ctx context.Context, plan *providerschema.AppEndpoint) diag.Diagnostics {
	var diags diag.Diagnostics

	plan.AdminURL = types.StringNull()
	plan.PublicURL = types.StringNull()
	plan.MetricsURL = types.StringNull()
	plan.State = types.StringNull()

	if plan.UserXattrKey.IsNull() || plan.UserXattrKey.IsUnknown() {
		plan.UserXattrKey = types.StringNull()
	}

	if plan.DeltaSyncEnabled.IsNull() || plan.DeltaSyncEnabled.IsUnknown() {
		plan.DeltaSyncEnabled = types.BoolNull()
	}

	plan.RequireResync = types.MapNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"items": types.SetType{ElemType: types.StringType},
		},
	})

	if !plan.Scopes.IsNull() {
		scopesMap := make(map[string]providerschema.AppEndpointScope)
		diags.Append(plan.Scopes.ElementsAs(ctx, &scopesMap, false)...)
		if diags.HasError() {
			return diags
		}

		for scopeName, scope := range scopesMap {
			if !scope.Collections.IsNull() {
				collectionsMap := make(map[string]providerschema.AppEndpointCollection)
				diags.Append(scope.Collections.ElementsAs(ctx, &collectionsMap, false)...)
				if diags.HasError() {
					return diags
				}

				for collName, collValue := range collectionsMap {
					if collValue.AccessControlFunction.IsNull() || collValue.AccessControlFunction.IsUnknown() {
						collValue.AccessControlFunction = types.StringNull()
					}

					if collValue.ImportFilter.IsNull() || collValue.ImportFilter.IsUnknown() {
						collValue.ImportFilter = types.StringNull()
					}

					collectionsMap[collName] = collValue
				}

				collectionsMapValue, d := types.MapValueFrom(ctx, types.ObjectType{
					AttrTypes: schema.
					AppEndpointCollection{}.
						AttributeTypes(),
				}, collectionsMap)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				scope.Collections = collectionsMapValue
				scopesMap[scopeName] = scope
			}
		}

		scopesMapValue, d := types.MapValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"collections": types.MapType{
					ElemType: types.ObjectType{
						AttrTypes: schema.
						AppEndpointCollection{}.
							AttributeTypes(),
					},
				},
			},
		}, scopesMap)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		plan.Scopes = scopesMapValue
	}

	if plan.Cors != nil {
		if plan.Cors.MaxAge.IsNull() || plan.Cors.MaxAge.IsUnknown() {
			plan.Cors.MaxAge = types.Int64Null()
		}
		if plan.Cors.Disabled.IsNull() || plan.Cors.Disabled.IsUnknown() {
			plan.Cors.Disabled = types.BoolNull()
		}
	}

	if !plan.Oidc.IsNull() {
		var oidcList []providerschema.AppEndpointOidc
		diags.Append(plan.Oidc.ElementsAs(ctx, &oidcList, false)...)
		if diags.HasError() {
			return diags
		}

		for i := range oidcList {
			oidcList[i].ProviderId = types.StringNull()
			oidcList[i].IsDefault = types.BoolNull()

			if oidcList[i].Register.IsNull() || oidcList[i].Register.IsUnknown() {
				oidcList[i].Register = types.BoolNull()
			}
			if oidcList[i].DiscoveryUrl.IsNull() || oidcList[i].DiscoveryUrl.IsUnknown() {
				oidcList[i].DiscoveryUrl = types.StringNull()
			}
			if oidcList[i].UsernameClaim.IsNull() || oidcList[i].UsernameClaim.IsUnknown() {
				oidcList[i].UsernameClaim = types.StringNull()
			}
			if oidcList[i].RolesClaim.IsNull() || oidcList[i].RolesClaim.IsUnknown() {
				oidcList[i].RolesClaim = types.StringNull()
			}
		}

		oidcSet, d := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: schema.
			AppEndpointOidc{}.
				AttributeTypes(),
		}, oidcList)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		plan.Oidc = oidcSet
	}

	return diags
}

// Read reads and updates the current state of an App Endpoint.
func (a *AppEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Endpoint state",
			fmt.Sprintf("Could not validate App Endpoint state %s: %s", state.Name.String(), err.Error()),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
		endpointName   = IDs[providerschema.AppEndpointName]
	)

	newstate, err := a.refreshAppEndpoint(
		ctx,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointName,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error refreshing App Endpoint",
			fmt.Sprintf("Could not refresh App Endpoint %s: %s", endpointName, errString),
		)
		return
	}
	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
}

// Update updates an existing App Endpoint.
func (a *AppEndpoint) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppEndpoint
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		endpointName   = plan.Name.ValueString()
	)

	appEndpointRequest, diags := morphToAppEndpointRequest(ctx, &plan)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointName,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		appEndpointRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating App Endpoint",
			fmt.Sprintf("Could not update App Endpoint %s: %s", endpointName, err.Error()),
		)
		return
	}

	// Refresh the plan after update.
	refreshedState, err := a.refreshAppEndpoint(
		ctx,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		endpointName,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing App Endpoint",
			fmt.Sprintf("Could not refresh App Endpoint %s: %s", endpointName, err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes an existing App Endpoint.
func (a *AppEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
		endpointName   = state.Name.ValueString()
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
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}

	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting App Endpoint",
			fmt.Sprintf("Could not delete App Endpoint %s: %s", endpointName, errString),
		)
		return
	}
}

// ImportState imports a remote App Endpoint that was not created by Terraform.
func (a *AppEndpoint) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Configure adds the provider configured api to App Endpoints.
func (a *AppEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.Data = data
}

// refreshAppEndpoint parses the API response and returns a refreshed AppEndpoint state.
func (a *AppEndpoint) refreshAppEndpoint(
	ctx context.Context, orgId, projId, clusterId, appServiceId, endpointName string,
) (*providerschema.AppEndpoint, error) {

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		a.HostURL,
		orgId,
		projId,
		clusterId,
		appServiceId,
		endpointName,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var appEndpoint app_endpoints.GetAppEndpointResponse
	if err := json.Unmarshal(response.Body, &appEndpoint); err != nil {
		return nil, fmt.Errorf("could not unmarshal App Endpoint response: %w", err)
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
		if diags.HasError() {
			return nil, fmt.Errorf("error converting require_resync: %s", diags.Errors())
		}

		state.RequireResync = requireResyncMap
	} else {
		state.RequireResync = types.MapNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"items": types.SetType{ElemType: types.StringType},
			},
		})
	}

	state.OrganizationId = types.StringValue(orgId)
	state.ProjectId = types.StringValue(projId)
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
				if diags.HasError() {
					return nil, fmt.Errorf("error converting collection %s: %v", collectionName, diags.Errors())
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
			if diags.HasError() {
				return nil, fmt.Errorf(
					"error converting collections for scope %s: %v",
					scopeName,
					diags.Errors(),
				)
			}

			scopeObj, diags := types.ObjectValueFrom(
				ctx,
				providerschema.AppEndpointScope{}.AttributeTypes(),
				providerschema.AppEndpointScope{
					Collections: collectionsMap,
				},
			)
			if diags.HasError() {
				return nil, fmt.Errorf(
					"error converting scope %s: %v",
					scopeName,
					diags.Errors(),
				)
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
		if diags.HasError() {
			return nil, fmt.Errorf("error converting scopes: %v", diags.Errors())
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
		if diags.HasError() {
			return nil, fmt.Errorf("error converting CORS origins: %v", diags.Errors())
		}
		state.Cors.Origin = originSet

		loginOriginSet, diags := types.SetValueFrom(
			ctx,
			types.StringType,
			appEndpoint.Cors.LoginOrigin,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting CORS login origins: %v", diags.Errors())
		}
		state.Cors.LoginOrigin = loginOriginSet

		headersSet, diags := types.SetValueFrom(
			ctx,
			types.StringType,
			appEndpoint.Cors.Headers,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting CORS headers: %v", diags.Errors())
		}
		state.Cors.Headers = headersSet
	}

	if len(appEndpoint.Oidc) > 0 {
		oidcSet, diags := types.SetValueFrom(
			ctx,
			types.ObjectType{
				AttrTypes: providerschema.
				AppEndpointOidc{}.
					AttributeTypes(),
			},
			appEndpoint.Oidc,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting OIDC configurations: %v", diags.Errors())
		}
		state.Oidc = oidcSet

	} else {
		state.Oidc = types.SetNull(types.ObjectType{
			AttrTypes: providerschema.
			AppEndpointOidc{}.
				AttributeTypes(),
		})
	}

	return state, nil
}

// morphToAppEndpointRequest converts the Terraform plan to the request format expected by the Capella API for creating/updating an App Endpoint.
func morphToAppEndpointRequest(
	ctx context.Context, plan *providerschema.AppEndpoint,
) (*app_endpoints.AppEndpointRequest, diag.Diagnostics) {
	var tfScopes map[string]providerschema.AppEndpointScope

	if diags := plan.Scopes.ElementsAs(ctx, &tfScopes, false); diags.HasError() {
		return nil, diags
	}

	// Convert from schema structs to API structs
	apiScopes := make(app_endpoints.Scopes)
	for scopeName, tfScope := range tfScopes {
		var tfCollections map[string]providerschema.AppEndpointCollection
		if diags := tfScope.Collections.ElementsAs(ctx, &tfCollections, false); diags.HasError() {
			return nil, diags
		}

		apiCollections := make(map[string]app_endpoints.Collection)
		for collName, tfColl := range tfCollections {
			apiCollections[collName] = app_endpoints.Collection{
				AccessControlFunction: tfColl.AccessControlFunction.ValueString(),
				ImportFilter:          tfColl.ImportFilter.ValueString(),
			}
		}

		apiScopes[scopeName] = app_endpoints.Scope{
			Collections: apiCollections,
		}
	}

	appEndpointRequest := &app_endpoints.AppEndpointRequest{
		Bucket: plan.Bucket.ValueString(),
		Name:   plan.Name.ValueString(),
		Scopes: apiScopes,
	}

	if !plan.DeltaSyncEnabled.IsNull() || !plan.DeltaSyncEnabled.IsUnknown() {
		appEndpointRequest.DeltaSyncEnabled = plan.DeltaSyncEnabled.ValueBool()
	}

	if !plan.UserXattrKey.IsNull() {
		appEndpointRequest.UserXattrKey = plan.UserXattrKey.ValueString()
	}

	if plan.Cors != nil {
		corsRequest := &app_endpoints.AppEndpointCors{
			MaxAge:   plan.Cors.MaxAge.ValueInt64(),
			Disabled: plan.Cors.Disabled.ValueBool(),
		}

		if !plan.Cors.Origin.IsNull() {
			origins := []string{}
			if diags := plan.Cors.Origin.ElementsAs(ctx, &origins, false); diags.HasError() {
				return nil, diags
			}

			corsRequest.Origin = origins
		}

		if !plan.Cors.LoginOrigin.IsNull() {
			loginOrigins := []string{}
			if diags := plan.Cors.LoginOrigin.ElementsAs(ctx, &loginOrigins, false); diags.HasError() {
				return nil, diags
			}

			corsRequest.LoginOrigin = loginOrigins
		}

		if !plan.Cors.Headers.IsNull() {
			headers := []string{}
			if diags := plan.Cors.Headers.ElementsAs(ctx, &headers, false); diags.HasError() {
				return nil, diags
			}

			corsRequest.Headers = headers
		}

		appEndpointRequest.Cors = corsRequest
	}

	if !plan.Oidc.IsNull() {
		oidc := []app_endpoints.AppEndpointOidc{}
		if diags := plan.Oidc.ElementsAs(ctx, &oidc, false); diags.HasError() {
			return nil, diags
		}
		appEndpointRequest.Oidc = oidc
	}

	return appEndpointRequest, nil
}
