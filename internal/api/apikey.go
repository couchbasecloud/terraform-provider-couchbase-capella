package api

import (
	"github.com/google/uuid"
)

// Resources  are the resource level permissions associated with the API key.
// To learn more about Organization Roles, see
// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
type Resources = []ResourcesItems

// ResourcesItems defines model for APIKeyResourcesItems.
type ResourcesItems struct {
	// Id is the id of the project.
	Id uuid.UUID `json:"id"`

	// Roles are the project roles associated with the API key.
	// To learn more about Project Roles, see
	//[Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
	Roles []string `json:"roles"`

	// Type is the type of the resource.
	Type *string `json:"type,omitempty"`
}

// CreateApiKeyRequest defines model for CreateAPIKeyRequest.
type CreateApiKeyRequest struct {
	// AllowedCIDRs is the list of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs *[]string `json:"allowedCIDRs,omitempty"`

	// Description is the description for the API key.
	Description *string `json:"description,omitempty"`

	// Expiry is the expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry *float32 `json:"expiry,omitempty"`

	// Name is the name of the API key.
	Name string `json:"name"`

	// OrganizationRoles are the organization level roles granted to the API key.
	OrganizationRoles []string `json:"organizationRoles"`

	// Resources are the resource level permissions associated with the API key.
	// To learn more about Organization Roles, see
	// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources *Resources `json:"resources,omitempty"`
}

// CreateApiKeyResponse defines model for CreateAPIKeyResponse.
type CreateApiKeyResponse struct {
	// Id The id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Token The Token is a confidential piece of information that is used to authorize requests made to v4 endpoints.
	Token string `json:"token"`
}

// GetApiKeyResponse defines model for GetAPIKey.
type GetApiKeyResponse struct {
	// AllowedCIDRs is the list of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs []string           `json:"allowedCIDRs"`
	Audit        CouchbaseAuditData `json:"audit"`

	// Description is the description for the API key.
	Description string `json:"description"`

	// Expiry is the expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry float32 `json:"expiry"`

	// Id is the id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Name is the name of the API key.
	Name string `json:"name"`

	// OrganizationRoles are the organization level roles granted to the API key.
	OrganizationRoles []string `json:"organizationRoles"`

	// Resources is the resources are the resource level permissions
	// associated with the API key. To learn more about Organization Roles, see
	// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources Resources `json:"resources"`
}

// RotateAPIKeyRequest defines model for RotateAPIKeyRequest.
type RotateAPIKeyRequest struct {
	// Secret represents the secret associated with an API key. One has to follow the secret key policy, such as allowed characters and a length of 64 characters.
	// If this field is left empty, a secret will be auto-generated.
	Secret *string `json:"secret,omitempty"`
}

// RotateAPIKeyResponse defines model for RotateAPIKeyResponse.
type RotateAPIKeyResponse struct {
	// SecretKey is a confidential token that is paired with the Access key.
	// The API key is made of an Access key and a Secret key.
	SecretKey string `json:"secretKey"`
}

// GetApiKeysResponse defines the model for a GetApiKeysResponse.
type GetApiKeysResponse struct {
	Data []GetApiKeyResponse `json:"data"`
}
