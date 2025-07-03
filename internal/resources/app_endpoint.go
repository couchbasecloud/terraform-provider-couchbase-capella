package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
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
func (a *AppEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	if err := a.validateCreateAppEndpointRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create app endpoint request",
			"Could not create app endpoint: "+err.Error(),
		)
		return
	}
	var scopes app_endpoints.ScopesConfig
	var collections map[string]app_endpoints.AppEndpointCollection
	for scopeName, scopeConfig := range plan.Scopes {
		collections = make(map[string]app_endpoints.AppEndpointCollection)
		for collName, collConfig := range scopeConfig.Collections {
			collections[collName.String()] = app_endpoints.AppEndpointCollection{
				ImportFilter:          collConfig.ImportFilter.ValueStringPointer(),
				AccessControlFunction: collConfig.AccessControlFunction.ValueStringPointer(),
			}

		}
		scopes[scopeName.String()] = app_endpoints.ScopeConfig{Collections: collections}
	}

	// Create the app endpoint using the API
	createAppEndpointRequest := app_endpoints.CreateAppEndpointRequest{
		Bucket:           plan.Bucket.ValueString(),
		Name:             plan.Name.ValueString(),
		DeltaSyncEnabled: plan.DeltaSyncEnabled.ValueBool(),
		Scopes:           scopes,
		Cors:             &app_endpoints.AppEndpointCors{Origin: providerschema.BaseStringsToStrings(plan.Cors.Origin), LoginOrigin: providerschema.BaseStringsToStrings(plan.Cors.LoginOrigin), Headers: providerschema.BaseStringsToStrings(plan.Cors.Headers), MaxAge: plan.Cors.MaxAge.ValueInt64Pointer(), Disabled: plan.Cors.Disabled.ValueBoolPointer()},
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

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	response, err := a.Client.ExecuteWithRetry(
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

	createAppServiceResponse := appservice.CreateAppServiceResponse{}
	err = json.Unmarshal(response.Body, &createAppServiceResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app service",
			errorMessageWhileAppServiceCreation+"error during unmarshalling:"+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeAppEndpointWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// validateCreateAppEndpointRequest validates the required fields for creating an app endpoint.
func (a *AppEndpoint) validateCreateAppEndpointRequest(plan providerschema.AppEndpoint) error {
	// Validate required IDs
	if plan.OrganizationId.IsNull() || plan.OrganizationId.IsUnknown() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() || plan.ProjectId.IsUnknown() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() || plan.ClusterId.IsUnknown() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	if plan.AppServiceId.IsNull() || plan.AppServiceId.IsUnknown() {
		return errors.ErrAppServiceIdCannotBeEmpty
	}

	// Validate required bucket name
	if plan.Bucket.IsNull() || plan.Bucket.IsUnknown() {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if !providerschema.IsTrimmed(plan.Bucket.ValueString()) {
		return fmt.Errorf("bucket name %s", errors.ErrNotTrimmed)
	}

	// Validate required endpoint name
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		return fmt.Errorf("app endpoint name cannot be empty")
	}
	if !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("app endpoint name %s", errors.ErrNotTrimmed)
	}
	// Validate endpoint name format
	if !isValidEndpointName(plan.Name.ValueString()) {
		return fmt.Errorf("app endpoint name must be between 1-100 characters and contain only lowercase letters, numbers, hyphens, and underscores")
	}

	// Validate userXattrKey if provided
	if !plan.UserXattrKey.IsNull() && !plan.UserXattrKey.IsUnknown() {
		if !providerschema.IsTrimmed(plan.UserXattrKey.ValueString()) {
			return fmt.Errorf("userXattrKey %s", errors.ErrNotTrimmed)
		}
	}

	// Validate OIDC configurations if provided
	if len(plan.Oidc) > 0 {
		for i, oidc := range plan.Oidc {
			if err := a.validateOidcConfiguration(oidc, i); err != nil {
				return err
			}
		}
	}

	// Validate CORS configuration if provided
	if len(plan.Cors.Origin) > 0 {
		for i, origin := range plan.Cors.Origin {
			if !providerschema.IsTrimmed(origin.ValueString()) {
				return fmt.Errorf("cors origin at index %d %s", i, errors.ErrNotTrimmed)
			}
		}
	}

	if len(plan.Cors.LoginOrigin) > 0 {
		for i, loginOrigin := range plan.Cors.LoginOrigin {
			if !providerschema.IsTrimmed(loginOrigin.ValueString()) {
				return fmt.Errorf("cors loginOrigin at index %d %s", i, errors.ErrNotTrimmed)
			}
		}
	}

	if len(plan.Cors.Headers) > 0 {
		for i, header := range plan.Cors.Headers {
			if !providerschema.IsTrimmed(header.ValueString()) {
				return fmt.Errorf("cors header at index %d %s", i, errors.ErrNotTrimmed)
			}
		}
	}

	// Validate CORS maxAge if provided
	if !plan.Cors.MaxAge.IsNull() && !plan.Cors.MaxAge.IsUnknown() {
		if plan.Cors.MaxAge.ValueInt64() < 0 {
			return fmt.Errorf("cors maxAge cannot be negative")
		}
	}

	return nil
}

// validateOidcConfiguration validates an individual OIDC configuration.
func (a *AppEndpoint) validateOidcConfiguration(oidc providerschema.AppEndpointOidc, index int) error {
	// Validate required issuer
	if oidc.Issuer.IsNull() || oidc.Issuer.IsUnknown() {
		return fmt.Errorf("oidc configuration at index %d: issuer cannot be empty", index)
	}
	if !providerschema.IsTrimmed(oidc.Issuer.ValueString()) {
		return fmt.Errorf("oidc configuration at index %d: issuer %s", index, errors.ErrNotTrimmed)
	}
	// Validate issuer URL format
	if !isValidURL(oidc.Issuer.ValueString()) {
		return fmt.Errorf("oidc configuration at index %d: issuer must be a valid URL", index)
	}

	// Validate required clientId
	if oidc.ClientId.IsNull() || oidc.ClientId.IsUnknown() {
		return fmt.Errorf("oidc configuration at index %d: clientId cannot be empty", index)
	}
	if !providerschema.IsTrimmed(oidc.ClientId.ValueString()) {
		return fmt.Errorf("oidc configuration at index %d: clientId %s", index, errors.ErrNotTrimmed)
	}

	// Validate optional fields if provided
	if !oidc.UserPrefix.IsNull() && !oidc.UserPrefix.IsUnknown() {
		if !providerschema.IsTrimmed(oidc.UserPrefix.ValueString()) {
			return fmt.Errorf("oidc configuration at index %d: userPrefix %s", index, errors.ErrNotTrimmed)
		}
	}

	if !oidc.DiscoveryUrl.IsNull() && !oidc.DiscoveryUrl.IsUnknown() {
		if !providerschema.IsTrimmed(oidc.DiscoveryUrl.ValueString()) {
			return fmt.Errorf("oidc configuration at index %d: discoveryUrl %s", index, errors.ErrNotTrimmed)
		}
		// Validate discovery URL format
		if !isValidURL(oidc.DiscoveryUrl.ValueString()) {
			return fmt.Errorf("oidc configuration at index %d: discoveryUrl must be a valid URL", index)
		}
	}

	if !oidc.UsernameClaim.IsNull() && !oidc.UsernameClaim.IsUnknown() {
		if !providerschema.IsTrimmed(oidc.UsernameClaim.ValueString()) {
			return fmt.Errorf("oidc configuration at index %d: usernameClaim %s", index, errors.ErrNotTrimmed)
		}
	}

	if !oidc.RolesClaim.IsNull() && !oidc.RolesClaim.IsUnknown() {
		if !providerschema.IsTrimmed(oidc.RolesClaim.ValueString()) {
			return fmt.Errorf("oidc configuration at index %d: rolesClaim %s", index, errors.ErrNotTrimmed)
		}
	}

	return nil
}

// isValidURL checks if a string is a valid URL.
func isValidURL(urlString string) bool {
	if urlString == "" {
		return false
	}

	// Basic URL validation - check if it starts with http:// or https://
	return len(urlString) > 7 && (urlString[:7] == "http://" || urlString[:8] == "https://")
}

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

func initializeAppEndpointWithPlan(plan providerschema.AppEndpoint) providerschema.AppEndpoint {
	if plan.UserXattrKey.IsNull() || plan.UserXattrKey.IsUnknown() {
		plan.UserXattrKey = types.StringNull()
	}
	if plan.DeltaSyncEnabled.IsNull() || plan.DeltaSyncEnabled.IsUnknown() {
		plan.DeltaSyncEnabled = types.BoolNull()
	}
	if plan.Cors.Origin == nil {
		plan.Cors.Origin = []types.String{}
	}
	if plan.Cors.LoginOrigin == nil {
		plan.Cors.LoginOrigin = []types.String{}
	}
	if plan.Cors.Headers == nil {
		plan.Cors.Headers = []types.String{}
	}
	if plan.Cors.MaxAge.IsNull() || plan.Cors.MaxAge.IsUnknown() {
		plan.Cors.MaxAge = types.Int64Null()
	}
	if plan.Cors.Disabled.IsNull() || plan.Cors.Disabled.IsUnknown() {
		plan.Cors.Disabled = types.BoolNull()
	}
	if len(plan.Oidc) == 0 {
		plan.Oidc = []providerschema.AppEndpointOidc{}
	}

	return plan
}

func (a *AppEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO: AV-104555: Implement read for App Endpoint
}

// Update allows to update the cluster to ON or OFF state.
func (a *AppEndpoint) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: AV-104552: Implement delete and update for App Endpoint
}

func (a *AppEndpoint) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO: AV-104552: Implement delete and update for App Endpoint
}
