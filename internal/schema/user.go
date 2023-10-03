package schema

import (
	"strings"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// User maps User resource schema data
type User struct {
	// Id is a GUID4 identifier of the user.
	Id types.String `tfsdk:"id"`

	// Name represents the name of the user.
	Name *types.String `tfsdk:"name"`

	// Name represents the email of the user.
	Email types.String `tfsdk:"email"`

	// Status depicts whether the user is verified or not
	Status types.String `tfsdk:"status"`

	// Inactive depicts whether the user has accepted the invite for the organization.
	Inactive types.Bool `tfsdk:"inactive"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// OrganizationRoles is an array of strings representing the roles granted to the user
	OrganizationRoles *[]types.String `tfsdk:"organization_roles"`

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
	Resources *[]Resource `tfsdk:"resources"`

	IfMatch types.String `tfsdk:"if_match"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

type Resource struct {
	// Id is a GUID4 identifier of the resource.
	Id types.String `tfsdk:"id"`

	// Type is the type of the resource.
	Type *types.String `tfsdk:"type"`

	// Roles is an array of strings representing a users project roles
	Roles []types.String `tfsdk:"roles"`
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
		Name:                &name,
		Email:               email,
		Status:              status,
		Inactive:            inactive,
		OrganizationId:      organizationId,
		OrganizationRoles:   &organizationRoles,
		LastLogin:           lastLogin,
		Region:              region,
		TimeZone:            timeZone,
		EnableNotifications: enableNotifications,
		ExpiresAt:           expiresAt,
		Resources:           &resources,
		Audit:               audit,
	}
	return &newUser
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import capella_user.new_user id=<uuid>,organization_id=<uuid>`
func (u User) Validate() (userId, organizationId string, err error) {
	const (
		idDelimiter       = ","
		organizationIdSep = "organization_id="
		userIdSep         = "id="
	)

	organizationId = u.OrganizationId.ValueString()
	userId = u.Id.ValueString()
	var found bool

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if u.OrganizationId.IsNull() {
		strs := strings.Split(u.Id.ValueString(), idDelimiter)
		if len(strs) != 4 {
			err = errors.ErrIdMissing
			return
		}
		_, userId, found = strings.Cut(strs[0], userIdSep)
		if !found {
			err = errors.ErrUserIdMissing
			return
		}

		_, organizationId, found = strings.Cut(strs[3], organizationIdSep)
		if !found {
			err = errors.ErrOrganizationIdMissing
			return
		}
	}

	if userId == "" {
		err = errors.ErrUserIdCannotBeEmpty
		return
	}

	if organizationId == "" {
		err = errors.ErrOrganizationIdCannotBeEmpty
		return
	}

	return userId, organizationId, nil
}
