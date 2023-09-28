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
	Roles string `json:"roles"`
}

type GetUserResponse struct {
	// ID is the ID of the user
	Id uuid.UUID `json:"id"`
}

type CreateUserResponse struct {
}
