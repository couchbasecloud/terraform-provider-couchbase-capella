package api

import (
	"github.com/google/uuid"
)

// CreateAllowListRequest defines the model for CreateAllowListRequest.
type CreateAllowListRequest struct {
	// Cidr is the trusted CIDR to allow database connections from.
	Cidr string `json:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// CreateAllowListResponse defines the model for CreateAllowListResponse.
type CreateAllowListResponse struct {
	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}

// GetAllowListResponse defines the model for GetAllowListResponse.
type GetAllowListResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Cidr is the trusted CIDR to allow database connections from.
	Cidr string `json:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment *string `json:"comment"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt *string `json:"expiresAt"`

	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}
