package api

import (
	"github.com/google/uuid"
)

// CreateDatabaseRoleRequest represents the POST request body for creating a database role.
type CreateDatabaseRoleRequest struct {
	// Name is the name of the database user role (2-128 characters, no spaces, restricted special characters).
	Name string `json:"name"`

	// Description is an optional description for the database user role.
	Description string `json:"description,omitempty"`

	// Access describes the access information of the database role.
	Access []Access `json:"access"`
}

// CreateDatabaseRoleResponse represents the POST response after creating a database role.
type CreateDatabaseRoleResponse struct {
	Id uuid.UUID `json:"id"`
}

// GetDatabaseRoleResponse represents the GET response for fetching a database role.
type GetDatabaseRoleResponse struct {
	Id          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Audit       CouchbaseAuditData `json:"audit"`
	Access      []Access           `json:"access"`
}

// UpdateDatabaseRoleRequest represents the PUT request body for updating a database role.
type UpdateDatabaseRoleRequest struct {
	// Description is the updated description for the database user role.
	Description string `json:"description,omitempty"`

	// Access describes the updated access information of the database user role.
	Access []Access `json:"access"`
}
