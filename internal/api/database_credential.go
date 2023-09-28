package api

import (
	"github.com/google/uuid"
)

type Access struct {
	Privileges []string `json:"privileges"`
}

type CreateDatabaseCredentialRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Access   []Access `json:"access"`
}

type CreateDatabaseCredentialResponse struct {
	Id       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type GetDatabaseCredentialResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Id A GUID4 identifier of the project.
	Id uuid.UUID `json:"id"`

	Etag string

	Access []Access `json:"access"`

	Name string `json:"name"`
}
