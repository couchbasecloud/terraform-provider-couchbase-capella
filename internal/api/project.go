package api

import (
	"github.com/google/uuid"
)

// CreateProjectRequest is the payload sent to the Capella V4 Public API when asked to create a project in the organization.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Creator
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateProjectRequest struct {
	// Description A short description about the project.
	Description string `json:"description,omitempty"`

	// Name The name of the project.
	Name string `json:"name"`
}

// CreateProjectResponse is the response received from the Capella V4 Public API when asked to create a project in the organization.
type CreateProjectResponse struct {
	// Id The ID of the project created.
	Id uuid.UUID `json:"id"`
}

// GetProjectResponse is the response received from the Capella V4 Public API when asked to get project details in an organization.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetProjectResponse struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Etag        string
	Audit       CouchbaseAuditData `json:"audit"`
	Id          uuid.UUID          `json:"id"`
}

// PutProjectRequest is the payload sent to the Capella V4 Public API when asked to update a project in an organization.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type PutProjectRequest struct {
	// Description represents a short description of the project.
	Description string `json:"description,omitempty"`

	// Name is the name of the project.
	Name string `json:"name"`
}
