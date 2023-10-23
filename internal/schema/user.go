package schema

import (
	"fmt"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// User maps User resource schema data
type User struct {
	// Id is a GUID4 identifier of the user.
	Id types.String `tfsdk:"id"`

	// Name represents the name of the user.
	Name types.String `tfsdk:"name"`

	// Name represents the email of the user.
	Email types.String `tfsdk:"email"`

	// Status depicts whether the user is verified or not
	Status types.String `tfsdk:"status"`

	// Inactive depicts whether the user has accepted the invite for the organization.
	Inactive types.Bool `tfsdk:"inactive"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// OrganizationRoles is an array of strings representing the roles granted to the user
	OrganizationRoles []types.String `tfsdk:"organization_roles"`

	// LastLogin is the time(UTC) at which user last logged in.
	LastLogin types.String `tfsdk:"last_login"`

	// Region is the region of the user.
	Region types.String `tfsdk:"region"`

	// TimeZone is the time zone of the user.
	TimeZone types.String `tfsdk:"time_zone"`

	// EnableNotifications represents whether email alerts for databases in projects
	// will be recieved.
	EnableNotifications types.Bool `tfsdk:"enable_notifications"`

	// ExpiresAt is the time at which user expires.
	ExpiresAt types.String `tfsdk:"expires_at"`

	// Resources is an array of objects representing the resources the user has access to
	Resources []Resource `tfsdk:"resources"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// NewUser creates a new instance of a User object
func NewUser(
	Id types.String,
	name types.String,
	email types.String,
	status types.String,
	inactive types.Bool,
	organizationId types.String,
	organizationRoles []types.String,
	lastLogin types.String,
	region types.String,
	timeZone types.String,
	enableNotifications types.Bool,
	expiresAt types.String,
	resources []Resource,
	audit basetypes.ObjectValue,
) *User {
	newUser := User{
		Id:                  Id,
		Name:                name,
		Email:               email,
		Status:              status,
		Inactive:            inactive,
		OrganizationId:      organizationId,
		OrganizationRoles:   organizationRoles,
		LastLogin:           lastLogin,
		Region:              region,
		TimeZone:            timeZone,
		EnableNotifications: enableNotifications,
		ExpiresAt:           expiresAt,
		Resources:           resources,
		Audit:               audit,
	}
	return &newUser
}

type Resource struct {
	// Id is a GUID4 identifier of the resource.
	Id types.String `tfsdk:"id"`

	// Type is the type of the resource.
	Type types.String `tfsdk:"type"`

	// Roles is an array of strings representing a users project roles
	Roles []types.String `tfsdk:"roles"`
}

// Users defines the attributes for a list of users in Capella.
type Users struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Data contains the list of resources.
	Data []User `tfsdk:"data"`
}

// Validate is used to verify that IDs have been properly imported
func (u *User) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: u.OrganizationId,
		Id:             u.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// MorphOrganizationRoles is used to convert nested organizationRoles from
// strings to terraform type.String.
// TODO (AV-53457): add unit testing
func MorphOrganizationRoles(organizationRoles []string) []basetypes.StringValue {
	var morphedRoles []basetypes.StringValue
	for _, role := range organizationRoles {
		morphedRoles = append(morphedRoles, types.StringValue(role))
	}
	return morphedRoles
}

// ConvertOrganizationRoles is used to convert all roles
// in an array of basetypes.StringValue to strings.
// TODO (AV-53457): add unit testing
func ConvertOrganizationRoles(organizationRoles []basetypes.StringValue) []string {
	var convertedRoles []string
	for _, role := range organizationRoles {
		convertedRoles = append(convertedRoles, role.ValueString())
	}
	return convertedRoles
}

// ConvertResource is used to convert a resource object containing nested fields
// of type basetypes.StringValue to a resource object containing nested fields of type string.
// TODO (AV-53457): add unit testing
func ConvertResources(resources []Resource) []api.Resource {
	var convertedResources []api.Resource
	for _, resource := range resources {
		var convertedResource api.Resource
		convertedResource.Id = resource.Id.ValueString()

		resourceType := resource.Type.ValueString()
		convertedResource.Type = &resourceType

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

// MorphResources is used to covert nested resources from strings
// to terraform types.String
// TODO (AV-53457): add unit testing
func MorphResources(resources []api.Resource) []Resource {
	var morphedResources []Resource
	for _, resource := range resources {
		var morphedResource Resource

		morphedResource.Id = types.StringValue(resource.Id)

		// Check for optional field
		if resource.Type != nil {
			resourceType := types.StringValue(*resource.Type)
			morphedResource.Type = resourceType
		}

		var roles []basetypes.StringValue
		for _, role := range resource.Roles {
			roles = append(roles, types.StringValue(role))
		}

		morphedResource.Roles = roles
		morphedResources = append(morphedResources, morphedResource)

	}
	return morphedResources
}
