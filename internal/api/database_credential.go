package api

import (
	"github.com/google/uuid"
)

// Access defines the level of access that the database credential will have across buckets and scopes.
// This access is currently defined for all buckets and all scopes in the cluster.
// todo: Support for granular access per bucket and per scope will be added in AV-62864
type Access struct {
	Privileges []string `json:"privileges"`
}

// CreateDatabaseCredentialRequest represents the schema for the POST Capella V4 API request that creates the database credential.
// Password is an optional field, if not passed, the password for the database credential is auto-generated.
type CreateDatabaseCredentialRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Access   []Access `json:"access"`
}

// CreateDatabaseCredentialResponse represents the schema for the POST Capella V4 API response that creates the database credential.
type CreateDatabaseCredentialResponse struct {
	Id       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

// GetDatabaseCredentialResponse represents the schema for the GET Capella V4 API request that fetches the database credential details.
type GetDatabaseCredentialResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Id A GUID4 identifier of the project.
	Id uuid.UUID `json:"id"`

	Name string `json:"name"`

	Access []Access `json:"access"`
}
