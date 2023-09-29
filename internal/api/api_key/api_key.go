package api_key

import (
	"github.com/google/uuid"
	"terraform-provider-capella/internal/api"
)

// ApiKeyResources  are the resource level permissions associated with the API key.
// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
type ApiKeyResources = []ApiKeyResourcesItems

// ApiKeyResourcesItems defines model for APIKeyResourcesItems.
type ApiKeyResourcesItems struct {
	// Id is the id of the project.
	Id uuid.UUID `json:"id"`

	// Roles are the project roles associated with the API key.
	// To learn more about Project Roles, see [Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
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
	Name              string   `json:"name"`
	OrganizationRoles []string `json:"organizationRoles"`

	// Resources are the resource level permissions associated with the API key.
	//
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources *ApiKeyResources `json:"resources,omitempty"`
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
	AllowedCIDRs []string               `json:"allowedCIDRs"`
	Audit        api.CouchbaseAuditData `json:"audit"`

	// Description is the description for the API key.
	Description string `json:"description"`

	// Expiry is the expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry float32 `json:"expiry"`

	// Id is the id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Name is the name of the API key.
	Name              string   `json:"name"`
	OrganizationRoles []string `json:"organizationRoles"`

	// Resources is the resources are the resource level permissions associated with the API key.
	//
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources ApiKeyResources `json:"resources"`
}
