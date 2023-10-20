package api

import "github.com/google/uuid"

// CreateUserRequest is the request payload sent to the Capella V4 Public API when asked to invite a new user to an organization.
// This request simply invites a new user under the organization.
// An invitation email is triggered and sent to the user.
// Upon receiving the invitation email, the user is required to click on a provided URL,
// which will redirect them to a page with a user interface (UI) where they can set their username and password.
//
// The modification of any personal information related to a user can only be performed by the user through the UI.
// Similarly, the user can solely conduct password updates through the UI.
//
// The "caller" possessing Organization Owner access rights retains the exclusive user creation capability.
// They hold the authority to assign roles at the organization and project levels.
// At present, our support is limited to the resourceType of "project" exclusively.
//
// In order to access this endpoint, the provided API key must have the following role:
// Organization Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-user-roles.html
type CreateUserRequest struct {
	// Name represents the name of the user.
	Name string `json:"name"`

	// Email represents the email of the user.
	Email string `json:"email"`

	// OrganizationRoles is an array of strings representing the roles granted to the user.
	OrganizationRoles []string `json:"organizationRoles"`

	// Resources is an array of objects representing the resources the user has access to.
	Resources []Resource `json:"resources"`
}

// CreateUserResponse is the response received from the Capella V4 Public API when asked to invite a new user to an organization.
type CreateUserResponse struct {
	// ID is the ID of the user
	Id uuid.UUID `json:"id"`
}

// Resource defines either a project or cluster to which the newly invited user should have access.
type Resource struct {
	// Id is a GUID4 identifier of the resource.
	Id string `json:"id"`

	// Type is the type of the resource.
	Type *string `json:"type"`

	// Roles is an array of strings representing a users project roles
	Roles []string `json:"roles"`
}

// GetUserResponse is the response received from the Capella V4 Public API when asked to get existing user's details.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Organization Member
// Project Creator
// The results are always limited by the role and scope of the caller's privileges.
//
// When performing a GET request for a user with an organization owner role,
// the response will exclude project-level permissions for that user.
// This is because organization owners have access to all resources at the organization level, rendering project-level permissions unnecessary for them.
//
// To learn more about the roles, see:
// Organization Roles: https://docs.couchbase.com/cloud/organizations/organization-user-roles.html
// Project Roles: https://docs.couchbase.com/cloud/projects/project-roles.html
type GetUserResponse struct {
	// ID is the ID of the user
	Id uuid.UUID `json:"id"`

	// Name represents the name of the user.
	Name *string `json:"name"`

	// Email represents the email of the user.
	Email string `json:"email"`

	// Status depicts whether the user is verified or not
	Status string `json:"status"`

	// Inactive depicts whether the user has accepted the invite for the organization.
	Inactive bool `json:"inactive"`

	// OrganizationId is a GUID4 identifier of the tenant.
	OrganizationId uuid.UUID `json:"organizationId"`

	// OrganizationRoles is an array of strings representing the roles granted to the user.
	OrganizationRoles []string `json:"organizationRoles"`

	// LastLogin is the time(UTC) at which user last logged in.
	LastLogin string `json:"lastLogin"`

	// Region is the region of the user.
	Region string `json:"region"`

	// TimeZone is the time zone of the user.
	TimeZone string `json:"timeZone"`

	// EnableNotifications represents whether email alerts for databases in projects
	// will be recieved.
	EnableNotifications bool `json:"enableNotifications"`

	// ExpiresAt is the time at which user expires.
	ExpiresAt string `json:"expiresAt"`

	// Resources is an array of objects representing the resources the user has access to.
	Resources []Resource `json:"resources"`

	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`
}

// UpdateUserRequest is the payload sent to the Capella V4 Public API when asked to update a user in an organization.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organizations.html
type UpdateUserRequest struct {
	// Op is the type of operation
	//
	// Enum: "add" "remove"
	op string

	// Path is the path of the resource that needs to be updated
	//
	// Organization Roles: /organizationRoles
	// Resources: /resources/{resourceId}
	// Resource Roles: /resources/{resourceId}/roles
	path string

	// Value is an array of OrganizationRoles (strings) or an
	// Array of ProjectRoles (strings) or a Resource (object)
	value interface{}
}

// GetUsersResponse is the response received from the Capella V4 Public API when asked to list all users that have access to an organization.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Organization Member
// Project Creator
// The results are always limited by the role and scope of the caller's privileges.
// When retrieving a list of users through a GET request, if a user holds the organization owner role,
// the response will exclude project-level permissions for those users.
// This is because organization owners have full access to all resources within the organization, making project-level permissions irrelevant for them.
//
// To learn more about the roles, see:
// Organization Roles: https://docs.couchbase.com/cloud/organizations/organization-user-roles.html
// Project Roles: https://docs.couchbase.com/cloud/projects/project-roles.html
type GetUsersResponse struct {
	Data []GetUserResponse `json:"data"`
}
