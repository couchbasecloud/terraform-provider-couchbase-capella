package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	tcslices "github.com/couchbase/tools-common/functional/slices"

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

const errorMessageAfterUserCreation = "User creation is successful, but encountered an error while checking the current" +
	" state of the user. Please run `terraform plan` after 1-2 minutes to know the" +
	" current user state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileUserCreation = "There is an error during user creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// User is the User resource implementation.
type User struct {
	*providerschema.Data
}

func NewUser() resource.Resource {
	return &User{}
}

// Metadata returns the users resource type name.
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

// Create creates a new user.
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
			"Could not create user, "+err.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	createUserRequest := api.CreateUserRequest{
		Email:             plan.Email.ValueString(),
		OrganizationRoles: providerschema.ConvertRoles(plan.OrganizationRoles),
	}

	// check for optional fields
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		createUserRequest.Name = plan.Name.ValueStringPointer()
	}

	if len(plan.Resources) != 0 {
		createUserRequest.Resources = providerschema.ConvertResources(plan.Resources)

	}

	// Execute request
	url := fmt.Sprintf("%s/v4/organizations/%s/users", r.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		createUserRequest,
		r.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileUserCreation+api.ParseError(err),
		)
		return
	}

	createUserResponse := api.CreateUserResponse{}
	err = json.Unmarshal(response.Body, &createUserResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			errorMessageWhileUserCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeUserWithPlanAndId(plan, createUserResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.refreshUser(ctx, organizationId, createUserResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error executing request",
			errorMessageAfterUserCreation+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)

	if checkOrganizationOwner(plan.OrganizationRoles, plan.Resources) {
		attributePath := path.Root("resources")
		diags = resp.State.SetAttribute(ctx, attributePath, plan.Resources)
		resp.Diagnostics.Append(diags...)
	}

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

// Read reads user information.
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
			"Error Reading Capella User",
			"Could not read Capella user: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		userId         = IDs[providerschema.Id]
	)

	// Refresh the existing user
	refreshedState, err := r.refreshUser(ctx, organizationId, userId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella User",
			"Could not read Capella userID "+userId+": "+errString,
		)
		return
	}

	// Set refreshed state
	if checkOrganizationOwner(state.OrganizationRoles, state.Resources) {
		existingResources := state.Resources

		diags = resp.State.Set(ctx, &refreshedState)
		resp.Diagnostics.Append(diags...)
		// overwrite resource values for organization owner. This is needed
		// as the API returns null resources for organization owner.
		attributePath := path.Root("resources")
		diags = resp.State.SetAttribute(ctx, attributePath, existingResources)
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the user.
func (r *User) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state, plan providerschema.User

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			"Could not update user id: "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		userId         = IDs[providerschema.Id]
	)

	patch := constructPatch(state, plan)

	err = r.updateUser(ctx, organizationId, userId, patch)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		resp.Diagnostics.AddError(
			"Error updating user",
			"Could not update Capella user with ID "+userId+": "+errString,
		)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
		}
		return
	}

	refreshedState, err := r.refreshUser(ctx, organizationId, userId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			"Could not update Capella user with ID "+userId+": "+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)

	if checkOrganizationOwner(plan.OrganizationRoles, plan.Resources) {
		attributePath := path.Root("resources")
		diags = resp.State.SetAttribute(ctx, attributePath, plan.Resources)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// constructPatch is used to determine to compare the planned user state with the
// existing user state and populate a Patch struct with the required fields.
func constructPatch(existing, proposed providerschema.User) []api.PatchEntry {
	patch := make([]api.PatchEntry, 0)

	patch = append(patch, handleOrganizationRoles(existing.OrganizationRoles, proposed.OrganizationRoles)...)
	patch = append(patch, handleProjectRoles(existing.Resources, proposed.Resources)...)
	patch = append(patch, compareResources(existing.Resources, proposed.Resources)...)

	return patch
}

// handleOrganizationRoles is used to compare the organizationRoles contained within
// two states and construct patch entries to reflect their differences.
func handleOrganizationRoles(existingRoles, proposedRoles []basetypes.StringValue) []api.PatchEntry {
	entries := make([]api.PatchEntry, 0)

	// Handle changes to organizationRoles
	addRoles, removeRoles := compare(existingRoles, proposedRoles)
	if len(addRoles) > 0 {
		providerschema.ConvertRoles(addRoles)
		entries = append(entries, api.PatchEntry{
			Op:    string(api.Add),
			Path:  "/organizationRoles",
			Value: providerschema.ConvertRoles(addRoles),
		})
	}
	if len(removeRoles) > 0 {
		entries = append(entries, api.PatchEntry{
			Op:    string(api.Remove),
			Path:  "/organizationRoles",
			Value: providerschema.ConvertRoles(removeRoles),
		})
	}

	return entries
}

// handleProjectRoles is used to compare the projectRoles contained within
// two states and construct patch entries to reflect their differences.
func handleProjectRoles(existingResources, proposedResources []providerschema.Resource) []api.PatchEntry {
	entries := make([]api.PatchEntry, 0)

	// populate maps with existing and proposed project roles
	existingMap := make(map[basetypes.StringValue][]basetypes.StringValue)
	proposedMap := make(map[basetypes.StringValue][]basetypes.StringValue)

	for _, resource := range existingResources {
		resourceType := resource.Type.ValueString()
		if resourceType != "" && resourceType != "project" {
			continue
		}
		existingMap[resource.Id] = resource.Roles
	}

	for _, resource := range proposedResources {
		resourceType := resource.Type.ValueString()
		if resourceType != "" && resource.Type.ValueString() != "project" {
			continue
		}
		proposedMap[resource.Id] = resource.Roles
	}

	// compare and construct patch entries
	for _, resource := range proposedResources {
		path := fmt.Sprintf("/resources/%s/roles", resource.Id.ValueString())
		if existing, exists := existingMap[resource.Id]; exists {
			addRoles, removeRoles := compare(existing, resource.Roles)
			if len(addRoles) > 0 {
				entries = append(entries, api.PatchEntry{
					Op:    string(api.Add),
					Path:  path,
					Value: providerschema.ConvertRoles(addRoles),
				})
			}
			if len(removeRoles) > 0 {
				entries = append(entries, api.PatchEntry{
					Op:    string(api.Remove),
					Path:  path,
					Value: providerschema.ConvertRoles(removeRoles),
				})
			}
		}
	}
	return entries
}

// compareResources is used to compare the resources contained within
// two states and construct patch entries to reflect their differences.
func compareResources(existingResources, proposedResources []providerschema.Resource) []api.PatchEntry {
	entries := make([]api.PatchEntry, 0)

	// populate maps with existing and proposed resources
	existingMap := make(map[basetypes.StringValue]providerschema.Resource)
	proposedMap := make(map[basetypes.StringValue]providerschema.Resource)

	for _, resource := range existingResources {
		existingMap[resource.Id] = resource
	}

	for _, resource := range proposedResources {
		proposedMap[resource.Id] = resource
	}

	// compare and construct patch entries
	for _, resource := range proposedResources {
		if _, exists := existingMap[resource.Id]; !exists {
			path := fmt.Sprintf("/resources/%s", resource.Id.ValueString())
			entries = append(entries, api.PatchEntry{
				Op:    "add",
				Path:  path,
				Value: providerschema.ConvertResource(resource),
			})
		}
	}

	for _, resource := range existingResources {
		if _, exists := proposedMap[resource.Id]; !exists {
			path := fmt.Sprintf("/resources/%s", resource.Id.ValueString())
			entries = append(entries, api.PatchEntry{
				Op:    "remove",
				Path:  path,
				Value: providerschema.ConvertResource(resource),
			})
		}
	}

	return entries
}

// checkOrganizationOwner is used to determine whether a list of planned roles for
// a user includes the role 'organizationOwner'.
func checkOrganizationOwner(roles []basetypes.StringValue, resources []providerschema.Resource) bool {
	if resources != nil && slices.Contains(
		roles, basetypes.NewStringValue("organizationOwner")) {
		return true
	}
	return false
}

// compare is used to compare two slices of basetypes.stringvalue
// and determine which values should be added and which should be removed.
func compare(existing, proposed []basetypes.StringValue) ([]basetypes.StringValue, []basetypes.StringValue) {
	// Add values present in the proposed state but not in existing.
	add := tcslices.Difference(proposed, existing)

	// Remove values present in the existing state but not in removed.
	remove := tcslices.Difference(existing, proposed)

	return add, remove
}

// updateUser is used to execute the patch request to update a user.
func (r *User) updateUser(ctx context.Context, organizationId, userId string, patch []api.PatchEntry) error {
	// Update existing user
	url := fmt.Sprintf("%s/v4/organizations/%s/users/%s", r.HostURL, organizationId, userId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPatch, SuccessStatus: http.StatusOK}
	_, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		patch,
		r.Token,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the user.
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
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/users/%s",
		r.HostURL,
		organizationId,
		userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Capella User",
			"Could not delete Capella userId "+userId+": "+errString,
		)
		return
	}
}

// getUser is used to retrieve an existing user.
func (r *User) getUser(ctx context.Context, organizationId, userId string) (*api.GetUserResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/users/%s",
		r.HostURL,
		organizationId,
		userId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	userResp := api.GetUserResponse{}
	err = json.Unmarshal(response.Body, &userResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &userResp, nil
}

// refreshUser retrieves user information for a specified organization and and user.
// It returns a schema representing the current user state.
func (r *User) refreshUser(ctx context.Context, organizationId, userId string) (*providerschema.User, error) {
	userResp, err := r.getUser(ctx, organizationId, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(userResp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
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
		providerschema.MorphRoles(userResp.OrganizationRoles),
		types.StringValue(userResp.LastLogin),
		types.StringValue(userResp.Region),
		types.StringValue(userResp.TimeZone),
		types.BoolValue(userResp.EnableNotifications),
		types.StringValue(userResp.ExpiresAt),
		providerschema.MorphResources(userResp.Resources),
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

// initializeUserWithPlanAndId initializes an instance of providerschema.User
// with the specified plan and ID. It marks all computed fields as null.
func initializeUserWithPlanAndId(plan providerschema.User, id string) providerschema.User {
	plan.Id = types.StringValue(id)
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		plan.Name = types.StringNull()
	}
	plan.Status = types.StringNull()
	plan.Inactive = types.BoolNull()
	plan.LastLogin = types.StringNull()
	plan.Region = types.StringNull()
	plan.TimeZone = types.StringNull()
	plan.EnableNotifications = types.BoolNull()
	plan.ExpiresAt = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return plan
}
