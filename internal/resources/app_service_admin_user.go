package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                   = (*AppServiceAdminUser)(nil)
	_ resource.ResourceWithConfigure      = (*AppServiceAdminUser)(nil)
	_ resource.ResourceWithImportState    = (*AppServiceAdminUser)(nil)
	_ resource.ResourceWithValidateConfig = (*AppServiceAdminUser)(nil)
)

// appServiceAdminUserAccess is the access oneOf. Update sends it as the whole body,
// create nests it under "access". The generated union type has no simple constructor.
type appServiceAdminUserAccess struct {
	AccessAllEndpoints bool     `json:"accessAllEndpoints,omitempty"`
	Endpoints          []string `json:"endpoints,omitempty"`
}

type createAppServiceAdminUserRequest struct {
	Name     string                    `json:"name"`
	Password string                    `json:"password"`
	Access   appServiceAdminUserAccess `json:"access"`
}

type AppServiceAdminUser struct {
	*providerschema.Data
}

// NewAppServiceAdminUser is used in (p *capellaProvider) Resources for building the provider.
func NewAppServiceAdminUser() resource.Resource {
	return &AppServiceAdminUser{}
}

// Metadata returns the App Service admin user resource type name.
func (a *AppServiceAdminUser) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_admin_user"
}

// Schema defines the schema for the App Service admin user resource.
func (a *AppServiceAdminUser) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceAdminUserSchema()
}

// Configure is used to configure the App Service admin user resource with the provider data.
func (a *AppServiceAdminUser) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ValidateConfig ensures exactly one of access_all_endpoints = true or a non-empty endpoints set is configured.
func (a *AppServiceAdminUser) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config providerschema.AppServiceAdminUser
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.AccessAllEndpoints.IsUnknown() || config.Endpoints.IsUnknown() {
		return
	}

	// Read stores null for an empty set, so an empty set would give an inconsistent result.
	if !config.Endpoints.IsNull() && len(config.Endpoints.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoints"),
			"Invalid App Service Admin User Configuration",
			"endpoints must not be empty. Remove it or add at least one App Endpoint.",
		)
		return
	}

	allEndpoints := config.AccessAllEndpoints.ValueBool()
	hasEndpoints := !config.Endpoints.IsNull()
	if allEndpoints == hasEndpoints {
		resp.Diagnostics.AddError(
			"Invalid App Service Admin User Configuration",
			"Exactly one of access_all_endpoints = true or a non-empty endpoints set must be configured.",
		)
	}
}

// Create is used to create a new App Service admin user.
func (a *AppServiceAdminUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceAdminUser
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	access, diags := accessRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	userReq := createAppServiceAdminUserRequest{
		Name:     plan.Name.ValueString(),
		Password: plan.Password.ValueString(),
		Access:   access,
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, userReq, a.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating App Service Admin User",
			"Could not create App Service admin user, unexpected error: "+api.ParseError(err),
		)
		return
	}

	createResp := apigen.CreateAppServiceAdminUserResponse{}
	if err = json.Unmarshal(response.Body, &createResp); err != nil || createResp.Id == nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling App Service Admin User Response",
			fmt.Sprintf("Could not unmarshal App Service admin user response: %v", err),
		)
		return
	}
	userId := *createResp.Id

	plan.Id = types.StringValue(userId)
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.refreshAdminUser(ctx, organizationId, projectId, clusterId, appServiceId, userId, plan)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Reading App Service Admin User After Creation",
			"App Service admin user creation is successful, but encountered an error while reading it: "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read is used to read an existing App Service admin user and set the state.
func (a *AppServiceAdminUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceAdminUser
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, projectId, clusterId, appServiceId, userId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating App Service Admin User State",
			"Could not validate App Service admin user state: "+err.Error(),
		)
		return
	}

	refreshedState, err := a.refreshAdminUser(ctx, organizationId, projectId, clusterId, appServiceId, userId, state)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading App Service Admin User",
			"Could not read App Service admin user "+userId+": "+errString,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update changes the endpoint access of an App Service admin user.
// The API cannot change the password, so a password change updates only the state.
// All other attributes require replacement.
func (a *AppServiceAdminUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.AppServiceAdminUser
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userReq, diags := accessRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
		userId         = state.Id.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	if _, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, userReq, a.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating App Service Admin User",
			"Could not update App Service admin user "+userId+": "+api.ParseError(err),
		)
		return
	}

	refreshedState, err := a.refreshAdminUser(ctx, organizationId, projectId, clusterId, appServiceId, userId, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Admin User",
			"Could not read App Service admin user "+userId+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Delete is used to delete an existing App Service admin user.
func (a *AppServiceAdminUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceAdminUser
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, projectId, clusterId, appServiceId, userId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting App Service Admin User",
			"Could not validate App Service admin user state: "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting App Service Admin User",
			"Could not delete App Service admin user "+userId+": "+errString,
		)
	}
}

// ImportState imports an App Service admin user.
// example: id=user123,organization_id=org123,project_id=proj123,cluster_id=cluster123,app_service_id=app123
func (a *AppServiceAdminUser) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// accessRequest builds the access part of the request body from the plan.
func accessRequest(ctx context.Context, plan providerschema.AppServiceAdminUser) (appServiceAdminUserAccess, diag.Diagnostics) {
	userReq := appServiceAdminUserAccess{AccessAllEndpoints: plan.AccessAllEndpoints.ValueBool()}
	if userReq.AccessAllEndpoints {
		return userReq, nil
	}
	diags := plan.Endpoints.ElementsAs(ctx, &userReq.Endpoints, false)
	return userReq, diags
}

// refreshAdminUser fetches the admin user and maps it to state. The password is
// kept from prior, since the API does not return it. See reconcileEndpoints for how
// endpoints are mapped.
func (a *AppServiceAdminUser) refreshAdminUser(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId, userId string, prior providerschema.AppServiceAdminUser,
) (*providerschema.AppServiceAdminUser, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	user := apigen.AppServiceAdminUser{}
	if err = json.Unmarshal(response.Body, &user); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	state := providerschema.AppServiceAdminUser{
		Id:                 types.StringValue(userId),
		OrganizationId:     types.StringValue(organizationId),
		ProjectId:          types.StringValue(projectId),
		ClusterId:          types.StringValue(clusterId),
		AppServiceId:       types.StringValue(appServiceId),
		Name:               types.StringPointerValue(user.Name),
		Password:           prior.Password,
		AccessAllEndpoints: types.BoolNull(),
		Endpoints:          types.SetNull(types.StringType),
		Audit:              types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes()),
	}

	if user.AccessAllEndpoints != nil && *user.AccessAllEndpoints {
		state.AccessAllEndpoints = types.BoolValue(true)
	} else if user.Endpoints != nil && len(*user.Endpoints) > 0 {
		state.Endpoints, err = a.reconcileEndpoints(ctx, organizationId, projectId, clusterId, appServiceId, prior.Endpoints, *user.Endpoints)
		if err != nil {
			return nil, err
		}
	}

	if user.Audit != nil {
		audit := providerschema.NewCouchbaseAuditData(api.CouchbaseAuditData(*user.Audit))
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting audit data to object")
		}
		state.Audit = auditObj
	}

	return &state, nil
}

// reconcileEndpoints decides which endpoints set goes into state after a read.
//
// The API grants admin access per bucket, not per App Endpoint. If the user
// configures endpoint "a" on bucket B, the API returns "a" and every other
// App Endpoint on bucket B. A plain copy of the API response would show a
// diff on every plan, so this function filters out that expansion.
//
// The function returns the prior set unchanged when both conditions are true:
//   - The API response contains every prior endpoint.
//   - Each extra endpoint in the response is on a bucket of a prior endpoint.
//
// Otherwise it returns the API response, so Terraform shows the drift. Examples:
//   - A prior endpoint is missing. Someone removed access outside Terraform.
//   - An extra endpoint is on a different bucket. Someone added access outside Terraform.
//   - The prior set is empty, for example after import. The API response is the only source.
//
// The bucket check calls the App Endpoints list API. The function makes this
// call only when extra endpoints exist, so a normal read makes no extra request.
func (a *AppServiceAdminUser) reconcileEndpoints(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId string, prior types.Set, apiEndpoints []string,
) (types.Set, error) {
	// prior is null after import or when the user had access_all_endpoints.
	// ElementsAs gives a nil slice for a null set, without an error.
	var priorEndpoints []string
	if diags := prior.ElementsAs(ctx, &priorEndpoints, false); diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("error reading prior endpoints")
	}

	// allPriorPresent is false for an empty prior set. That forces the API response into state.
	allPriorPresent := len(priorEndpoints) > 0 && !slices.ContainsFunc(priorEndpoints, func(e string) bool {
		return !slices.Contains(apiEndpoints, e)
	})
	// extras are the endpoints in the API response that the user did not configure.
	// Clone first, because DeleteFunc changes the slice in place.
	extras := slices.DeleteFunc(slices.Clone(apiEndpoints), func(e string) bool {
		return slices.Contains(priorEndpoints, e)
	})

	// An exact match needs no bucket lookup.
	keepPrior := allPriorPresent && len(extras) == 0
	if allPriorPresent && len(extras) > 0 {
		var err error
		keepPrior, err = a.extrasSharePriorBuckets(ctx, organizationId, projectId, clusterId, appServiceId, priorEndpoints, extras)
		if err != nil {
			return types.SetNull(types.StringType), err
		}
	}
	if keepPrior {
		return prior, nil
	}

	// Drift or no prior set: the API response is the source of truth.
	endpoints, diags := types.SetValueFrom(ctx, types.StringType, apiEndpoints)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("error converting endpoints to set")
	}
	return endpoints, nil
}

// extrasSharePriorBuckets reports whether every extra endpoint is on a bucket of a prior endpoint.
// It lists the App Endpoints of the App Service to map names to buckets.
func (a *AppServiceAdminUser) extrasSharePriorBuckets(
	ctx context.Context, organizationId, projectId, clusterId, appServiceId string, priorEndpoints, extras []string,
) (bool, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	appEndpoints, err := api.GetPaginated[[]app_endpoints.GetAppEndpointResponse](ctx, a.ClientV1, a.Token, cfg, api.SortByName)
	if err != nil {
		return false, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	bucketOf := make(map[string]string, len(appEndpoints))
	for _, ep := range appEndpoints {
		bucketOf[ep.Name] = ep.Bucket
	}

	priorBuckets := make(map[string]struct{}, len(priorEndpoints))
	for _, e := range priorEndpoints {
		if bucket, ok := bucketOf[e]; ok {
			priorBuckets[bucket] = struct{}{}
		}
	}

	return !slices.ContainsFunc(extras, func(e string) bool {
		bucket, ok := bucketOf[e]
		if !ok {
			return true
		}
		_, shared := priorBuckets[bucket]
		return !shared
	}), nil
}
