package api

import (
	"github.com/google/uuid"
)

// CreateAllowListRequest  defines model for CreateAllowListRequest.
type CreateAllowListRequest struct {
	// The trusted CIDR to allow the database connections from.
	Cidr string `json:"cidr"`

	// A short description of the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	// An RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// CreateAllowListResponse defines model for CreateAllowListResponse.
type CreateAllowListResponse struct {
	// The ID of the cluster the AllowList was created for
	Id uuid.UUID `json:"id"`
}

// GetAllowListResponse defines model for GetAllowListResponse.
type GetAllowListResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// The trusted CIDR to allow the database connections from.
	Cidr string `json:"cidr"`

	// A short description of the allowed CIDR.
	Comment string `json:"comment"`

	// An RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt string `json:"expiresAt"`

	// The ID of the cluster the AllowList was created for
	Id uuid.UUID `json:"id"`
}
