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

// appServiceAdminUserRequest is the create and update body. The API expects the access
// oneOf flattened into the body, which the generated union type does not do.
// Update sends only the access fields.
type appServiceAdminUserRequest struct {
	Name               string   `json:"name,omitempty"`
	Password           string   `json:"password,omitempty"`
	AccessAllEndpoints bool     `json:"accessAllEndpoints,omitempty"`
	Endpoints          []string `json:"endpoints,omitempty"`
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

	if !config.AccessAllEndpoints.IsNull() && !config.AccessAllEndpoints.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_all_endpoints"),
			"Invalid App Service Admin User Configuration",
			"access_all_endpoints can only be set to true. Remove it and set endpoints to restrict access.",
		)
		return
	}

	allEndpoints := config.AccessAllEndpoints.ValueBool()
	hasEndpoints := len(config.Endpoints.Elements()) > 0
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

	userReq, diags := accessRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	userReq.Name = plan.Name.ValueString()
	userReq.Password = plan.Password.ValueString()

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
func accessRequest(ctx context.Context, plan providerschema.AppServiceAdminUser) (appServiceAdminUserRequest, diag.Diagnostics) {
	userReq := appServiceAdminUserRequest{AccessAllEndpoints: plan.AccessAllEndpoints.ValueBool()}
	if userReq.AccessAllEndpoints {
		return userReq, nil
	}
	diags := plan.Endpoints.ElementsAs(ctx, &userReq.Endpoints, false)
	return userReq, diags
}

// refreshAdminUser fetches the admin user and maps it to state. The password is
// kept from prior, since the API does not return it. The API stores
// access per bucket, so it may return more endpoints than configured. Prior endpoints
// are kept when all of them are in the response, to avoid a perpetual diff.
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
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	user := apigen.AppServiceAdminUser{}
	if err = json.Unmarshal(response.Body, &user); err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
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
		var priorEndpoints []string
		if diags := prior.Endpoints.ElementsAs(ctx, &priorEndpoints, false); diags.HasError() {
			return nil, fmt.Errorf("error reading prior endpoints")
		}
		if len(priorEndpoints) > 0 && !slices.ContainsFunc(priorEndpoints, func(e string) bool {
			return !slices.Contains(*user.Endpoints, e)
		}) {
			state.Endpoints = prior.Endpoints
		} else {
			endpoints, diags := types.SetValueFrom(ctx, types.StringType, *user.Endpoints)
			if diags.HasError() {
				return nil, fmt.Errorf("error converting endpoints to set")
			}
			state.Endpoints = endpoints
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
