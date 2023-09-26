package api_key

import (
	"github.com/google/uuid"
	"terraform-provider-capella/internal/api"
)

const (
	OrganizationMember ApiKeyOrganizationRoles = "organizationMember"
	OrganizationOwner  ApiKeyOrganizationRoles = "organizationOwner"
	ProjectCreator     ApiKeyOrganizationRoles = "projectCreator"
)

// Defines values for APIKeyResourcesItemsRoles.
const (
	ProjectDataReader       APIKeyResourcesItemsRoles = "projectDataReader"
	ProjectDataReaderWriter APIKeyResourcesItemsRoles = "projectDataReaderWriter"
	ProjectManager          APIKeyResourcesItemsRoles = "projectManager"
	ProjectOwner            APIKeyResourcesItemsRoles = "projectOwner"
	ProjectViewer           APIKeyResourcesItemsRoles = "projectViewer"
)

// APIKeyDescription Description for the API key.
type APIKeyDescription = string

// APIKeyExpiry Expiry of the API key in number of days.
// If set to -1, the token will not expire.
type APIKeyExpiry = float32

// APIKeyName Name of the API key.
type APIKeyName = string

// APIKeyOrganizationRoles Organization roles assigned to the API key.
//
// To learn more, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
type ApiKeyOrganizationRoles string

// APIKeyResources Resources are the resource level permissions associated with the API key.
//
// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
type APIKeyResources = []APIKeyResourcesItems

// APIKeyResourcesItems defines model for APIKeyResourcesItems.
type APIKeyResourcesItems struct {
	// Id ID of the project.
	Id uuid.UUID `json:"id"`

	// Roles Project Roles associated with the API key.
	//
	// To learn more about Project Roles, see [Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
	Roles []APIKeyResourcesItemsRoles `json:"roles"`

	// Type Type of the resource.
	Type *string `json:"type,omitempty"`
}

// APIKeyResourcesItemsRoles defines model for APIKeyResourcesItems.Roles.
type APIKeyResourcesItemsRoles string

// CreateApiKeyRequest defines model for CreateAPIKeyRequest.
type CreateApiKeyRequest struct {
	// AllowedCIDRs List of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs *[]string `json:"allowedCIDRs,omitempty"`

	// Description Description for the API key.
	Description *APIKeyDescription `json:"description,omitempty"`

	// Expiry Expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry *APIKeyExpiry `json:"expiry,omitempty"`

	// Name Name of the API key.
	Name              APIKeyName                `json:"name"`
	OrganizationRoles []ApiKeyOrganizationRoles `json:"organizationRoles"`

	// Resources Resources are the resource level permissions associated with the API key.
	//
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources *APIKeyResources `json:"resources,omitempty"`
}

// CreateAPIKeyResponse defines model for CreateAPIKeyResponse.
type CreateAPIKeyResponse struct {
	// Id The id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Token The Token is a confidential piece of information that is used to authorize requests made to v4 endpoints.
	Token string `json:"token"`
}

// GetApiKeyResponse defines model for GetAPIKey.
type GetApiKeyResponse struct {
	// AllowedCIDRs List of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs []string               `json:"allowedCIDRs"`
	Audit        api.CouchbaseAuditData `json:"audit"`

	// Description Description for the API key.
	Description APIKeyDescription `json:"description"`

	// Expiry Expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry APIKeyExpiry `json:"expiry"`

	// Id The id is a unique identifier for an apiKey.
	Id string `json:"id"`

	// Name Name of the API key.
	Name              APIKeyName                `json:"name"`
	OrganizationRoles []ApiKeyOrganizationRoles `json:"organizationRoles"`

	// Resources Resources are the resource level permissions associated with the API key.
	//
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources APIKeyResources `json:"resources"`
}

// GetAPIKeys defines model for GetAPIKeys.
type GetAPIKeys struct {
	Data []GetApiKeyResponse `json:"data"`
}
