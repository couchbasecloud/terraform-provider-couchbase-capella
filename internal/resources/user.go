package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &User{}
	_ resource.ResourceWithConfigure   = &User{}
	_ resource.ResourceWithImportState = &User{}
)

// User is the User resource implementation
type User struct {
	*providerschema.Data
}

func NewUsers() resource.Resource {
	return &User{}
}

// Metadata returns the users resource type name
func (r *User) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the allowlist resource.
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

	createUserRequest := api.CreateUserRequest{
		Name:              plan.Name.ValueString(),
		Email:             plan.Email.ValueString(),
		OrganizationRoles: r.convertOrganizationRoles(plan.OrganizationRoles),
		Resources:         r.convertResources(plan.Resources),
	}

	// Execute request
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/users", r.HostURL, plan.OrganizationId.ValueString()),
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

	refreshedState, err := r.refreshUser(ctx, plan.OrganizationId.String(), plan.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error reading Capella User",
			"Could not read Capella User "+createUserResponse.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error reading Capella User",
			"Could not read Capella User "+createUserResponse.Id.String()+": "+err.Error(),
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

// Read reads user information
func (r *User) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// todo (AV-69625):
}

// Update updates the user
func (r *User) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo (AV-69626):
}

// Delete deletes the user
func (r *User) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo (AV-69627):
}

// convertOrganizationRoles is used to convert all roles
// in an array of basetypes.StringValue to strings.
func (r *User) convertOrganizationRoles(organizationRoles []basetypes.StringValue) []string {
	var convertedRoles []string
	for _, role := range organizationRoles {
		convertedRoles = append(convertedRoles, role.ValueString())
	}
	return convertedRoles
}

// convertResource is used to convert a resource object containing nested fields
// of type basetypes.StringValue to a resource object containing nested fields of type string.
func (r *User) convertResources(resources []providerschema.Resource) []api.Resource {
	var convertedResources []api.Resource
	for _, resource := range resources {
		var convertedResource api.Resource
		convertedResource.Id = resource.Id.ValueString()
		convertedResource.Type = resource.Type.ValueString()

		// Iterate through roles belonging to the user and convert to string
		var convertedRoles []string
		for _, role := range resource.Roles {
			convertedRoles = append(convertedRoles, role.ValueString())
		}
		convertedResource.Roles = convertedRoles

		convertedResources = append(convertedResources, convertedResource)
	}
	return convertedResources
}

// getUser is used to retrieve an existing user
func (r *User) getUser(ctx context.Context, organizationId, userId string) (*api.GetUserResponse, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/userss/%s",
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
		return nil, fmt.Errorf("error executing request: %s", err)
	}

	userResp := api.GetUserResponse{}
	err = json.Unmarshal(response.Body, &userResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %s", err)
	}
	userResp.ETag = response.Response.Header.Get("ETag")
	return &userResp, nil
}

func (r *User) refreshUser(ctx context.Context, organizationId, userId string) (*providerschema.OneUser, error) {
	userResp, err := r.getUser(ctx, organizationId, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %s", err)
	}

	var organizationRoles []basetypes.StringValue
	for _, role := range userResp.OrganizationRoles {
		organizationRoles = append(organizationRoles, types.StringValue(role))
	}

	var resources []providerschema.Resource
	for _, resource := range userResp.Resources {
		var convertedResource providerschema.Resource

		convertedResource.Id = types.StringValue(resource.Id)
		convertedResource.Type = types.StringValue(resource.Type)

		var roles []basetypes.StringValue
		for _, role := range resource.Roles {
			roles = append(roles, types.StringValue(role))
		}

		convertedResource.Roles = roles

		resources = append(resources, convertedResource)
	}

	refreshedState := providerschema.OneUser{
		Id:                  types.StringValue(userResp.Id.String()),
		Name:                types.StringValue(userResp.Name),
		Email:               types.StringValue(userResp.Email),
		Status:              types.StringValue(userResp.Status),
		Inactive:            types.BoolValue(userResp.Inactive),
		OrganizationId:      types.StringValue(organizationId),
		OrganizationRoles:   organizationRoles,
		LastLogin:           types.StringValue(userResp.LastLogin),
		Region:              types.StringValue(userResp.Region),
		TimeZone:            types.StringValue(userResp.TimeZone),
		EnableNotifications: types.BoolValue(userResp.EnableNotifications),
		ExpiresAt:           types.StringValue(userResp.ExpiresAt),
		Resources:           resources,
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(userResp.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(userResp.Audit.CreatedBy),
			ModifiedAt: types.StringValue(userResp.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(userResp.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(userResp.Audit.Version)),
		},
	}

	// TODO: Set optional fields

	return &refreshedState, nil
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
