package organization

import (
	"terraform-provider-capella/internal/api"

	"github.com/google/uuid"
)

// GetOrganizationResponse is the response received from the Capella V4 Public API when asked to fetch organization details.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Creator
// Organization Member
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-user-roles.html
type GetOrganizationResponse struct {
	// Audit contains all audit-related fields.
	Audit api.CouchbaseAuditData `json:"audit"`

	// Name represents the organization name.
	Name string `json:"name"`

	// Description is a short description of the organization.
	Description *string `json:"description"`

	// Preferences stores preferences for the tenant.
	Preferences *Preferences `json:"preferences"`

	// ID is the ID of the Organization
	Id uuid.UUID `json:"id"`
}
