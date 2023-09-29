package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

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
	OrganizationRoles []types.String `tfsdk:"organizationRoles"`

	// LastLogin is the time(UTC) at which user last logged in.
	LastLogin types.String `tfsdk:"lastLogin"`

	// Region is the region of the user.
	Region types.String `tfsdk:"region"`

	// TimeZone is the time zone of the user.
	TimeZone types.String `tfsdk:"timeZone"`

	// EnableNotifications represents whether email alerts for databases in projects
	// will be recieved.
	EnableNotifications types.Bool `tfsdk:"enableNotifications"`

	// ExpiresAt is the time at which user expires.
	ExpiresAt types.String `tfsdk:"expiresAt"`

	// Resources is an array of objects representing the resources the user has access to
	Resources []Resource `tfsdk:"resources"`

	// ETag is a unique indentifier which the client uses to determine if the resource has changed.
	ETag types.String `tfsdk:"if_match"`

	// IfMatch is used to check if a request should be made. The request will only proceed if
	// the resources current ETag matches this value.
	IfMatch types.String `tfsdk:"if_match"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

type Resource struct {
	// Id is a GUID4 identifier of the resource.
	Id string `tfsdk:"id"`

	// Type is the type of the resource.
	Type string `tfsdk:"type"`

	// Roles is an array of strings representing a users project roles
	Roles string `tfsdk:"roles"`
}

// OneAllowList maps user resource schema data; there is a separate response object to avoid conversion error for nested fields.
type OneUser struct {
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
	OrganizationRoles []types.String `tfsdk:"organizationRoles"`

	// LastLogin is the time(UTC) at which user last logged in.
	LastLogin types.String `tfsdk:"lastLogin"`

	// Region is the region of the user.
	Region types.String `tfsdk:"region"`

	// TimeZone is the time zone of the user.
	TimeZone types.String `tfsdk:"timeZone"`

	// EnableNotifications represents whether email alerts for databases in projects
	// will be recieved.
	EnableNotifications types.Bool `tfsdk:"enableNotifications"`

	// ExpiresAt is the time at which user expires.
	ExpiresAt types.String `tfsdk:"expiresAt"`

	// Resources is an array of objects representing the resources the user has access to
	Resources types.Object `tfsdk:"resources"`

	// ETag is a unique indentifier which the client uses to determine if the resource has changed.
	ETag types.String `tfsdk:"if_match"`

	// IfMatch is used to check if a request should be made. The request will only proceed if
	// the resources current ETag matches this value.
	IfMatch types.String `tfsdk:"if_match"`

	// Audit represents all audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`
}
