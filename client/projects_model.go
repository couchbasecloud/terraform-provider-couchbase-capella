package client

import (
	"github.com/google/uuid"
	"time"
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

// PutProjectRequest defines model for CreateProjectRequest.
type PutProjectRequest struct {
	// Description A short description about the project.
	Description string `json:"description,omitempty"`

	// Name The name of the project.
	Name string `json:"name"`
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

// CouchbaseAuditData contains all audit-related fields.
type CouchbaseAuditData struct {
	// CreatedAt The RFC3339 timestamp associated with when the resource was initially
	// created.
	CreatedAt time.Time `json:"createdAt"`

	// CreatedBy The user who created the resource; this will be a UUID4 ID for standard
	// users and will be a string such as "internal-support" for internal
	// Couchbase support users.
	CreatedBy string `json:"createdBy"`

	// ModifiedAt The RFC3339 timestamp associated with when the resource was last modified.
	ModifiedAt time.Time `json:"modifiedAt"`

	// ModifiedBy The user who last modified the resource; this will be a UUID4 ID for
	// standard users and wilmal be a string such as "internal-support" for
	// internal Couchbase support users.
	ModifiedBy string `json:"modifiedBy"`

	// Version The version of the document. This value is incremented each time the
	// resource is modified.
	Version int `json:"version"`
}

// GetProjectsResponse defines model for GetProjectsResponse.
type GetProjectsResponse struct {
	Data []GetProjectResponse `json:"data"`
}

type PutProjectParams struct {
	// IfMatch A precondition header that specifies the entity tag of a resource.
	IfMatch *string `json:"If-Match,omitempty"`
}
