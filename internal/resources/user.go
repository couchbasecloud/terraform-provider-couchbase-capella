package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &User{}
	_ resource.ResourceWithConfigure   = &User{}
	_ resource.ResourceWithImportState = &User{}
)

const (
	organizationIdKey = "organizationId"
	userIdKey         = "userId"
)

// User is the User resource implementation
type User struct {
	*providerschema.Data
}

func NewUser() resource.Resource {
	return &User{}
}

// Metadata returns the users resource type name
func (r *User) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the user resource.
func (r *User) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UserSchema()
}

// Configure sets provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (r *User) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create creates a new user
func (r *User) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.validateCreateUserRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create user request",
			"Could not create user "+err.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	createUserRequest := api.CreateUserRequest{
		Name:              plan.Name.ValueString(),
		Email:             plan.Email.ValueString(),
		OrganizationRoles: providerschema.ConvertOrganizationRoles(plan.OrganizationRoles),
		Resources:         providerschema.ConvertResources(plan.Resources),
	}

	// Execute request
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/users", r.HostURL, organizationId),
		http.MethodPost,
		createUserRequest,
		r.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not execute request, unexpected error: "+err.Error(),
		)
		return
	}

	createUserResponse := api.CreateUserResponse{}
	err = json.Unmarshal(response.Body, &createUserResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.refreshUser(ctx, organizationId, createUserResponse.Id.String(), plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			"Could not read user, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *User) validateCreateUserRequest(plan providerschema.User) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.Email.IsNull() {
		return errors.ErrEmailCannotBeEmpty
	}
	if plan.OrganizationRoles == nil {
		return errors.ErrOrganizationRolesCannotBeEmpty
	}
	return nil
}

// Read reads user information
func (r *User) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allow list: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		userId         = IDs[providerschema.Id]
	)

	// Refresh the existing user
	refreshedState, err := r.refreshUser(ctx, organizationId, userId, state)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella User",
				"Could not read Capella userID "+userId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella User",
			"Could not read Capella userID "+userId+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the user
func (r *User) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for users.
	// Users are instead updated via a PATCH request.
	// http://cbc-cp-api.s3-website-us-east-1.amazonaws.com/#tag/allowedCIDRs(Cluster)
	//
	// The update logic has been therefore been left blank. In this situation, terraform apply
	// will default to deleting and executing a new create.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
	//
	// TODO (AV-63471): Implement logic to parse and execute a PATCH request
}

// Delete deletes the user
func (r *User) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve existing state
	var state providerschema.User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella User",
			"Could not delete Capella user: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		userId         = IDs[providerschema.Id]
	)

	// Execute request to delete existing user
	_, err = r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/users/%s",
			r.HostURL,
			organizationId,
			userId,
		),
		http.MethodDelete,
		nil,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Deleting Capella User",
				"Could not delete Capella userId "+userId+": "+err.CompleteError(),
			)
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Deleting Capella User",
			"Could not delete Capella userId "+userId+": "+err.Error(),
		)
		return
	}
}

// getUser is used to retrieve an existing user
func (r *User) getUser(ctx context.Context, organizationId, userId string) (*api.GetUserResponse, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/users/%s",
			r.HostURL,
			organizationId,
			userId,
		),
		http.MethodGet,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	userResp := api.GetUserResponse{}
	err = json.Unmarshal(response.Body, &userResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errors.ErrUnmarshallingResponse, err)
	}
	return &userResp, nil
}

func (r *User) refreshUser(ctx context.Context, organizationId, userId string, plan providerschema.User) (*providerschema.User, error) {
	userResp, err := r.getUser(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	audit := providerschema.NewCouchbaseAuditData(userResp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, errors.ErrUnableToConvertAuditData
	}

	// Set optional fields - these may be left blank
	var name basetypes.StringValue
	if userResp.Name != nil {
		name = types.StringValue(*userResp.Name)
	}

	refreshedState := providerschema.NewUser(
		types.StringValue(userResp.Id.String()),
		name,
		types.StringValue(userResp.Email),
		types.StringValue(userResp.Status),
		types.BoolValue(userResp.Inactive),
		types.StringValue(userResp.OrganizationId.String()),
		providerschema.MorphOrganizationRoles(userResp.OrganizationRoles),
		types.StringValue(userResp.LastLogin),
		types.StringValue(userResp.Region),
		types.StringValue(userResp.TimeZone),
		types.BoolValue(userResp.EnableNotifications),
		types.StringValue(userResp.ExpiresAt),
		plan.Resources,
		auditObj,
	)
	return refreshedState, nil
}

// ImportState imports a remote user that was not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=cluster123,project_id=proj123,organization_id=org123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *User) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// handleCapellaUserError is used to differentiate between error types which
// may be returned during requests to capella.
func handleCapellaUserError(err error) error {
	switch err := err.(type) {
	case nil:
	case api.Error:
		return fmt.Errorf("%w: %s", errors.ErrUnableToReadCapellaUser, err.CompleteError())
	default:
		return fmt.Errorf("%w: %s", errors.ErrUnableToReadCapellaUser, err.Error())
	}
	return nil
}
