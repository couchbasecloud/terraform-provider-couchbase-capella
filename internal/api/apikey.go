package api

import (
	"github.com/google/uuid"
)

// Resources  are the resource level permissions associated with the API key.
// To learn more about Organization Roles, see
// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
type Resources = []ResourcesItems

// ResourcesItems the individual item that is part of Resources.
// These items define the set of roles or access that can be had on a single type of resource.
type ResourcesItems struct {
	Type  *string   `json:"type,omitempty"`
	Roles []string  `json:"roles"`
	Id    uuid.UUID `json:"id"`
}

// CreateApiKeyRequest is the payload sent to the Capella V4 Public API when asked to create an API key in the organization.
// Organization Owners can create Organization and Project scoped API keys.
//
// Project Owner and Project Creator can create project scoped keys.
//
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateApiKeyRequest struct {
	AllowedCIDRs      *[]string  `json:"allowedCIDRs,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Expiry            *float32   `json:"expiry,omitempty"`
	Resources         *Resources `json:"resources,omitempty"`
	Name              string     `json:"name"`
	OrganizationRoles []string   `json:"organizationRoles"`
}

// CreateApiKeyResponse is the response received from the Capella V4 Public API when asked to create an API key in the organization.
type CreateApiKeyResponse struct {
	// Id The id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Token The Token is a confidential piece of information that is used to authorize requests made to v4 endpoints.
	Token string `json:"token"`
}

// GetApiKeyResponse is the response received from the Capella V4 Public API when asked to fetch an API key in an organization.
//
// Organization Owners can get any API key inside the Organization.
// Organization Members and Project Creator can get any Project scoped API key for which they are Project Owner.
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetApiKeyResponse struct {
	Description       string             `json:"description"`
	Id                string             `json:"id"`
	Name              string             `json:"name"`
	Audit             CouchbaseAuditData `json:"audit"`
	AllowedCIDRs      []string           `json:"allowedCIDRs"`
	OrganizationRoles []string           `json:"organizationRoles"`
	Resources         Resources          `json:"resources"`
	Expiry            float32            `json:"expiry"`
}

// RotateApiKeyRequest is the payload sent to the Capella V4 Public API when asked to rotate an API key.
//
// Organization Owners can rotate any API key inside the Organization.
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type RotateApiKeyRequest struct {
	// Secret represents the secret associated with an API key.
	// One has to follow the secret key policy, such as allowed characters and a length of 64 characters.
	// If this field is left empty, a secret will be auto-generated.
	Secret *string `json:"secret,omitempty"`
}

// RotateApiKeyResponse is the response received from the Capella V4 Public API when asked to rotate an API key.
type RotateApiKeyResponse struct {
	// SecretKey is a confidential token that is paired with the Access key.
	// The API key is made of an Access key and a Secret key.
	SecretKey string `json:"secretKey"`
}
