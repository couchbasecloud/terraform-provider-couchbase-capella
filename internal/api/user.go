package api

import "github.com/google/uuid"

// CreateUserRequest defines the model for CreateUserRequest
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

type Resource struct {
	// Id is a GUID4 identifier of the resource.
	Id string `json:"id"`

	// Type is the type of the resource.
	Type string `json:"type"`

	// Roles is an array of strings representing a users project roles
	Roles []string `json:"roles"`
}

type GetUserResponse struct {
	// ID is the ID of the user
	Id uuid.UUID `json:"id"`

	// Name represents the name of the user.
	Name string `json:"name"`

	// Email represents the email of the user.
	Email string `json:"email"`

	// Status depicts whether the user is verified or not
	Status string `json:"status"`

	// Inactive depicts whether the user has accepted the invite for the organization.
	Inactive bool `json:"inactive"`

	// OrganizationId is a GUID4 identifier of the tenant.
	OrganizationIdId uuid.UUID `json:"organizationId"`

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

type CreateUserResponse struct {
	// ID is the ID of the user
	Id uuid.UUID `json:"id"`
}
