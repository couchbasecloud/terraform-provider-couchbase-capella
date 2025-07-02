package resources

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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
	var scopes app_endpoints.AppEndpointScopes
	for scopeName, scopeConfig := range plan.Scopes {
		for collName, collConfig := range scopeConfig.Collections {
			collections[collName] = oapi.CollectionConfig{
				ImportFilter:          collConfig.ImportFilter,
				AccessControlFunction: collConfig.SyncFn,
			}
		}
		scopes[scopeName] = oapi.ScopeConfig{Collections: collections}
	}
	// Create the app endpoint using the API
	createAppEndpointRequest := app_endpoints.CreateAppEndpointRequest{
		Bucket:           plan.Bucket.ValueString(),
		Name:             plan.Name.ValueString(),
		DeltaSyncEnabled: plan.DeltaSyncEnabled.ValueBool(),
		Scopes:           plan.Scopes,
		Cors:             nil,
		Oidc:             nil,
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

	// Validate scopes configuration if provided
	if err := a.validateScopesConfiguration(plan.Scopes); err != nil {
		return err
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

// validateScopesConfiguration validates the scopes configuration.
func (a *AppEndpoint) validateScopesConfiguration(scopes providerschema.AppEndpointScopes) error {
	// Validate access control function if provided
	if !scopes.Default.Collections.Default.AccessControlFunction.IsNull() && !scopes.Default.Collections.Default.AccessControlFunction.IsUnknown() {
		if !providerschema.IsTrimmed(scopes.Default.Collections.Default.AccessControlFunction.ValueString()) {
			return fmt.Errorf("accessControlFunction %s", errors.ErrNotTrimmed)
		}
	}

	// Validate import filter if provided
	if !scopes.Default.Collections.Default.ImportFilter.IsNull() && !scopes.Default.Collections.Default.ImportFilter.IsUnknown() {
		if !providerschema.IsTrimmed(scopes.Default.Collections.Default.ImportFilter.ValueString()) {
			return fmt.Errorf("importFilter %s", errors.ErrNotTrimmed)
		}
	}

	return nil
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
