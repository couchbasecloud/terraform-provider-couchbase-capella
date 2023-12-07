package organization

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

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
	Description *string                `json:"description"`
	Preferences *Preferences           `json:"preferences"`
	Name        string                 `json:"name"`
	Audit       api.CouchbaseAuditData `json:"audit"`
	Id          uuid.UUID              `json:"id"`
}
