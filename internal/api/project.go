package api

import (
	"github.com/google/uuid"
)

// CreateProjectRequest defines model for CreateProjectRequest.
type CreateProjectRequest struct {
	// Description A short description about the project.
	Description string `json:"description,omitempty"`

	// Name The name of the project.
	Name string `json:"name"`
}

// CreateProjectResponse defines model for CreateProjectResponse.
type CreateProjectResponse struct {
	// Id The ID of the project created.
	Id uuid.UUID `json:"id"`
}

// GetProjectResponse defines model for GetProjectResponse.
type GetProjectResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Description The description of a particular project.
	Description string `json:"description"`

	// Id A GUID4 identifier of the project.
	Id uuid.UUID `json:"id"`

	// Name The name of the project.
	Name string `json:"name"`

	Etag string
}
