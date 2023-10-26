package api

import (
	"github.com/google/uuid"
)

// CreateAllowListRequest is the payload sent to the Capella V4 Public API when asked to add an allowlist to gain access to a Capella cluster.
// Couchbase Capella only allows trusted IP addresses to connect to databases.
// Each database has a configurable Allowed IP list that can include up to 75 entries.
// Each entry can be a single IP address or an IP address space.
// Any IP address you add to this list can have a user-specified expiration time for temporary access, or be permanent.
// Capella automatically denies any connection attempts to and from an IP not in the allowed IP list.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateAllowListRequest struct {
	// Cidr is the trusted CIDR to allow the database connections from.
	// To add a single IP address, use a subnet mask of 32.
	Cidr string `json:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	// If this field is empty/omitted then the allowed CIDR is permanent and will never automatically expire.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// CreateAllowListResponse is the response received from the Capella V4 Public API when asked to add an allowlist to gain access to a Capella cluster.
type CreateAllowListResponse struct {
	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}

// GetAllowListResponse is the response received from the Capella V4 Public API when asked to fetch details of a particular allowlist.
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
type GetAllowListResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Cidr is the trusted CIDR to allow the database connections from.
	// To add a single IP address, use a subnet mask of 32.
	Cidr string `json:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment *string `json:"comment"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt *string `json:"expiresAt"`

	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}
