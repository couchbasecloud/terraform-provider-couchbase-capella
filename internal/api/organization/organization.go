package organization

import (
	"github.com/google/uuid"
	"terraform-provider-capella/internal/api"
)

// GetOrganizationResponse defines the model for GetOrganizationResponse.
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
