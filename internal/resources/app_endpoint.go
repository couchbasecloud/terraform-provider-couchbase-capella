package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppEndpoint{}
	_ resource.ResourceWithConfigure   = &AppEndpoint{}
	_ resource.ResourceWithImportState = &AppEndpoint{}
)

// AppEndpoint is the AppEndpoint implementation.
type AppEndpoint struct {
	*providerschema.Data
}

// NewAppEndpoint is a helper function to simplify the provider implementation.
func NewAppEndpoint() resource.Resource {
	return &AppEndpoint{}
}

// ImportState imports a remote AppEndpoint app service that is not created by Terraform.
func (a *AppEndpoint) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Metadata returns the AppEndpoint cluster resource type name.
func (a *AppEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint"
}

// Configure It adds the provider configured api to ClusterOnOff.
func (a *AppEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.Data = data
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

	scope := plan.Scope.ValueString()
	sc := make(map[string]types.Object)
	plan.Collections.ElementsAs(ctx, &sc, false)

	nestedMap := make(map[string]map[string]map[string]app_endpoints.AppEndpointCollection)
	nestedMap[scope] = make(map[string]map[string]app_endpoints.AppEndpointCollection)

	fieldSetters := map[string]func(*app_endpoints.AppEndpointCollection, string){
		"access_control_function": func(
			c *app_endpoints.AppEndpointCollection, val string,
		) {
			c.AccessControlFunction = &val
		},
		"import_filter": func(c *app_endpoints.AppEndpointCollection, val string) { c.ImportFilter = &val },
	}

	for col, obj := range sc {
		if _, ok := nestedMap[scope]["collections"]; !ok {
			nestedMap[scope]["collections"] = make(map[string]app_endpoints.AppEndpointCollection)
		}

		attr := obj.Attributes()
		endpointCollection := nestedMap[scope]["collections"][col]

		for name, value := range attr {
			if !(value.IsNull() || value.IsUnknown()) {
				fieldSetters[name](&endpointCollection, value.String())
			}
		}

		nestedMap[scope]["collections"][col] = endpointCollection
	}

	createAppEndpointRequest := app_endpoints.CreateAppEndpointRequest{
		Bucket:           plan.Bucket.ValueString(),
		Name:             plan.Name.ValueString(),
		DeltaSyncEnabled: plan.DeltaSyncEnabled.ValueBool(),
		Scopes:           nestedMap,
	}
	if plan.Cors != nil {
		createAppEndpointRequest.Cors = &app_endpoints.AppEndpointCors{
			Origin:      providerschema.BaseStringsToStrings(plan.Cors.Origin),
			LoginOrigin: providerschema.BaseStringsToStrings(plan.Cors.LoginOrigin),
			Headers:     providerschema.BaseStringsToStrings(plan.Cors.Headers),
			MaxAge:      plan.Cors.MaxAge.ValueInt64Pointer(),
			Disabled:    plan.Cors.Disabled.ValueBoolPointer(),
		}
	}
	if len(plan.Oidc) > 0 {
		createAppEndpointRequest.Oidc = make([]app_endpoints.AppEndpointOidc, len(plan.Oidc))
		for i, oidc := range plan.Oidc {
			createAppEndpointRequest.Oidc[i] = app_endpoints.AppEndpointOidc{
				Issuer:        oidc.Issuer.ValueString(),
				ClientId:      oidc.ClientId.ValueString(),
				UserPrefix:    oidc.UserPrefix.ValueStringPointer(),
				DiscoveryUrl:  oidc.DiscoveryUrl.ValueStringPointer(),
				UsernameClaim: oidc.UsernameClaim.ValueStringPointer(),
				RolesClaim:    oidc.RolesClaim.ValueStringPointer(),
				Register:      oidc.Register.ValueBoolPointer(),
			}
		}
	}

	diags = resp.State.Set(ctx, initComputedAttributesToNull(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if jsonData, err := json.MarshalIndent(createAppEndpointRequest, "", "  "); err == nil {
		fmt.Printf("###DEBUG### createAppEndpointRequest: %s\n", string(jsonData))
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var appServiceId = plan.AppServiceId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		createAppEndpointRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileAppServiceCreation+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// validateCreateAppEndpointRequest validates the required fields for creating an app endpoint.
//func (a *AppEndpoint) validateCreateAppEndpointRequest(plan providerschema.AppEndpoint) error {
//	// Validate required IDs
//	if plan.OrganizationId.IsNull() || plan.OrganizationId.IsUnknown() {
//		return errors.ErrOrganizationIdCannotBeEmpty
//	}
//	if plan.ProjectId.IsNull() || plan.ProjectId.IsUnknown() {
//		return errors.ErrProjectIdCannotBeEmpty
//	}
//	if plan.ClusterId.IsNull() || plan.ClusterId.IsUnknown() {
//		return errors.ErrClusterIdCannotBeEmpty
//	}
//	if plan.AppServiceId.IsNull() || plan.AppServiceId.IsUnknown() {
//		return errors.ErrAppServiceIdCannotBeEmpty
//	}
//
//	// Validate required bucket name
//	if plan.Bucket.IsNull() || plan.Bucket.IsUnknown() {
//		return fmt.Errorf("bucket name cannot be empty")
//	}
//	if !providerschema.IsTrimmed(plan.Bucket.ValueString()) {
//		return fmt.Errorf("bucket name %s", errors.ErrNotTrimmed)
//	}
//
//	// Validate required endpoint name
//	if plan.Name.IsNull() || plan.Name.IsUnknown() {
//		return fmt.Errorf("app endpoint name cannot be empty")
//	}
//	if !providerschema.IsTrimmed(plan.Name.ValueString()) {
//		return fmt.Errorf("app endpoint name %s", errors.ErrNotTrimmed)
//	}
//	// Validate endpoint name format
//	if !isValidEndpointName(plan.Name.ValueString()) {
//		return fmt.Errorf("app endpoint name must be between 1-100 characters and contain only lowercase letters, numbers, hyphens, and underscores")
//	}
//
//	// Validate userXattrKey if provided
//	if !plan.UserXattrKey.IsNull() && !plan.UserXattrKey.IsUnknown() {
//		if !providerschema.IsTrimmed(plan.UserXattrKey.ValueString()) {
//			return fmt.Errorf("userXattrKey %s", errors.ErrNotTrimmed)
//		}
//	}
//
//	// Validate OIDC configurations if provided
//	if len(plan.Oidc) > 0 {
//		for i, oidc := range plan.Oidc {
//			if err := a.validateOidcConfiguration(oidc, i); err != nil {
//				return err
//			}
//		}
//	}
//
//	// Validate CORS configuration if provided
//	if len(plan.Cors.Origin) > 0 {
//		for i, origin := range plan.Cors.Origin {
//			if !providerschema.IsTrimmed(origin.ValueString()) {
//				return fmt.Errorf("cors origin at index %d %s", i, errors.ErrNotTrimmed)
//			}
//		}
//	}
//
//	if len(plan.Cors.LoginOrigin) > 0 {
//		for i, loginOrigin := range plan.Cors.LoginOrigin {
//			if !providerschema.IsTrimmed(loginOrigin.ValueString()) {
//				return fmt.Errorf("cors loginOrigin at index %d %s", i, errors.ErrNotTrimmed)
//			}
//		}
//	}
//
//	if len(plan.Cors.Headers) > 0 {
//		for i, header := range plan.Cors.Headers {
//			if !providerschema.IsTrimmed(header.ValueString()) {
//				return fmt.Errorf("cors header at index %d %s", i, errors.ErrNotTrimmed)
//			}
//		}
//	}
//
//	// Validate CORS maxAge if provided
//	if !plan.Cors.MaxAge.IsNull() && !plan.Cors.MaxAge.IsUnknown() {
//		if plan.Cors.MaxAge.ValueInt64() < 0 {
//			return fmt.Errorf("cors maxAge cannot be negative")
//		}
//	}
//
//	return nil
//}

// isValidEndpointName checks if an endpoint name follows the proper naming convention.
func isValidEndpointName(name string) bool {
	if len(name) < 1 || len(name) > 100 {
		return false
	}

	// Check if name contains only lowercase letters, numbers, hyphens, and underscores
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}

	return true
}

func initComputedAttributesToNull(plan providerschema.AppEndpoint) providerschema.AppEndpoint {
	if plan.AdminURL.IsUnknown() || plan.AdminURL.IsNull() {
		plan.AdminURL = types.StringNull()
	}

	if plan.PublicURL.IsUnknown() || plan.PublicURL.IsNull() {
		plan.PublicURL = types.StringNull()

	}

	if plan.MetricsURL.IsUnknown() || plan.MetricsURL.IsNull() {
		plan.MetricsURL = types.StringNull()
	}

	for i := range plan.Oidc {
		if plan.Oidc[i].ProviderId.IsUnknown() || plan.Oidc[i].ProviderId.IsNull() {
			plan.Oidc[i].ProviderId = types.StringNull()
		}

		if plan.Oidc[i].IsDefault.IsUnknown() || plan.Oidc[i].IsDefault.IsNull() {
			plan.Oidc[i].IsDefault = types.BoolNull()
		}
	}

	plan.RequireResync = types.MapNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"items": types.SetType{ElemType: types.StringType},
		},
	})

	return plan
}

func (a *AppEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var organizationId = state.OrganizationId.ValueString()
	var projectId = state.ProjectId.ValueString()
	var clusterId = state.ClusterId.ValueString()
	var appServiceId = state.AppServiceId.ValueString()
	var endpointName = state.Name.ValueString()

	// Get the app endpoint
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		a.HostURL, organizationId, projectId, clusterId, appServiceId, endpointName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app endpoint",
			fmt.Sprintf("Could not read app endpoint %s: %s", endpointName, err.Error()),
		)
		return
	}

	newstate, err := a.refreshAppEndpoint(ctx, response.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing app endpoint",
			fmt.Sprintf("Could not refresh app endpoint %s: %s", endpointName, err.Error()),
		)
		return
	}
	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
}

// refreshAppEndpoint parses the API response and returns a refreshed AppEndpoint state
func (a *AppEndpoint) refreshAppEndpoint(ctx context.Context, responseBody []byte) (*providerschema.AppEndpoint, error) {
	var appEndpoint app_endpoints.GetAppEndpointResponse
	err := json.Unmarshal(responseBody, &appEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse app endpoint response: %w", err)
	}

	refreshedState := &providerschema.AppEndpoint{}

	// Set basic attributes
	refreshedState.Bucket = types.StringValue(appEndpoint.Bucket)
	refreshedState.Name = types.StringValue(appEndpoint.Name)
	refreshedState.DeltaSyncEnabled = types.BoolValue(appEndpoint.DeltaSyncEnabled)

	if appEndpoint.UserXattrKey != nil {
		refreshedState.UserXattrKey = types.StringValue(*appEndpoint.UserXattrKey)
	} else {
		refreshedState.UserXattrKey = types.StringNull()
	}

	// Set computed attributes
	if appEndpoint.AdminURL != "" {
		refreshedState.AdminURL = types.StringValue(appEndpoint.AdminURL)
	}
	if appEndpoint.PublicURL != "" {
		refreshedState.PublicURL = types.StringValue(appEndpoint.PublicURL)
	}
	if appEndpoint.MetricsURL != "" {
		refreshedState.MetricsURL = types.StringValue(appEndpoint.MetricsURL)
	}

	// Handle scopes and collections
	if len(appEndpoint.Scopes) > 0 {
		for scopeName, scopeData := range appEndpoint.Scopes {
			refreshedState.Scope = types.StringValue(scopeName)

			collectionAttrs := make(map[string]attr.Value)
			for collectionName := range scopeData.Collections {
				collectionAttrs[collectionName] = types.ObjectValueMust(
					map[string]attr.Type{
						"access_control_function": types.StringType,
						"import_filter":           types.StringType,
					},
					map[string]attr.Value{
						"access_control_function": types.StringNull(),
						"import_filter":           types.StringNull(),
					},
				)
			}
			refreshedState.Collections = types.MapValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"access_control_function": types.StringType,
						"import_filter":           types.StringType,
					},
				},
				collectionAttrs,
			)
		}
	}

	// Handle CORS if present
	if appEndpoint.Cors != nil {
		refreshedState.Cors = &providerschema.AppEndpointCors{
			Disabled: types.BoolPointerValue(appEndpoint.Cors.Disabled),
		}
		if appEndpoint.Cors.MaxAge != nil {
			refreshedState.Cors.MaxAge = types.Int64PointerValue(appEndpoint.Cors.MaxAge)
		}
		if len(appEndpoint.Cors.Origin) > 0 {
			origins := make([]types.String, len(appEndpoint.Cors.Origin))
			for i, origin := range appEndpoint.Cors.Origin {
				origins[i] = types.StringValue(origin)
			}
			refreshedState.Cors.Origin = origins
		}
		if len(appEndpoint.Cors.LoginOrigin) > 0 {
			loginOrigins := make([]types.String, len(appEndpoint.Cors.LoginOrigin))
			for i, loginOrigin := range appEndpoint.Cors.LoginOrigin {
				loginOrigins[i] = types.StringValue(loginOrigin)
			}
			refreshedState.Cors.LoginOrigin = loginOrigins
		}
		if len(appEndpoint.Cors.Headers) > 0 {
			headers := make([]types.String, len(appEndpoint.Cors.Headers))
			for i, header := range appEndpoint.Cors.Headers {
				headers[i] = types.StringValue(header)
			}
			refreshedState.Cors.Headers = headers
		}
	}

	// Handle OIDC if present
	if len(appEndpoint.Oidc) > 0 {
		refreshedState.Oidc = make([]providerschema.AppEndpointOidc, len(appEndpoint.Oidc))
		for i, oidc := range appEndpoint.Oidc {
			refreshedState.Oidc[i] = providerschema.AppEndpointOidc{
				Issuer:        types.StringValue(oidc.Issuer),
				ClientId:      types.StringValue(oidc.ClientId),
				UserPrefix:    types.StringPointerValue(oidc.UserPrefix),
				DiscoveryUrl:  types.StringPointerValue(oidc.DiscoveryUrl),
				UsernameClaim: types.StringPointerValue(oidc.UsernameClaim),
				RolesClaim:    types.StringPointerValue(oidc.RolesClaim),
				Register:      types.BoolPointerValue(oidc.Register),
				ProviderId:    types.StringPointerValue(oidc.ProviderId),
				IsDefault:     types.BoolPointerValue(oidc.IsDefault),
			}
		}
	}

	return refreshedState, nil
}

func (a *AppEndpoint) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: AV-104552: Implement delete and update for App Endpoint
}

func (a *AppEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppEndpoint
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var organizationId = state.OrganizationId.ValueString()
	var projectId = state.ProjectId.ValueString()
	var clusterId = state.ClusterId.ValueString()
	var appServiceId = state.AppServiceId.ValueString()
	var endpointName = state.Name.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		a.HostURL, organizationId, projectId, clusterId, appServiceId, endpointName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}

	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting app endpoint",
			fmt.Sprintf("Could not delete app endpoint %s: %s", endpointName, err.Error()),
		)
		return
	}
}
