package app_endpoints

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/google/uuid"
)

// CreateAppEndpointRequest is the request payload sent to the Capella V4 Public API in order to create a new app endpoint.
// An App Endpoint provides a REST API for accessing data in a Couchbase Capella bucket through an App Service.
//
// To learn more about App Endpoints, see https://docs.couchbase.com/cloud/app-services/index.html
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateAppEndpointRequest struct {
	// Bucket is the name of the bucket associated with this App Endpoint.
	Bucket string `json:"bucket"`

	// Name is the name of the App Endpoint.
	Name string `json:"name"`

	// UserXattrKey is the user extended attribute key for the App Endpoint.
	UserXattrKey *string `json:"userXattrKey,omitempty"`

	// DeltaSyncEnabled enables or disables delta sync on this App Endpoint.
	DeltaSyncEnabled bool `json:"deltaSyncEnabled,omitempty"`

	// Scopes is the configuration for scopes and collections within the App Endpoint.
	Scopes ScopesConfig `json:"scopes,omitempty"`

	// Cors is the CORS configuration for the App Endpoint.
	Cors *AppEndpointCors `json:"cors,omitempty"`

	// Oidc is the list of OIDC providers for the App Endpoint.
	Oidc []AppEndpointOidc `json:"oidc,omitempty"`
}

// CreateAppEndpointResponse is the response received from the Capella V4 Public API when asked to create a new app endpoint.
type CreateAppEndpointResponse struct {
	// ID is the UUID of the app endpoint
	Id uuid.UUID `json:"id"`
}

// GetAppEndpointResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing app endpoint.
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetAppEndpointResponse struct {
	// Id is the UUID of the app endpoint
	Id uuid.UUID `json:"id"`

	// Bucket is the name of the bucket associated with this App Endpoint.
	Bucket string `json:"bucket"`

	// Name is the name of the App Endpoint.
	Name string `json:"name"`

	// UserXattrKey is the user extended attribute key for the App Endpoint.
	UserXattrKey *string `json:"userXattrKey,omitempty"`

	// DeltaSyncEnabled indicates whether delta sync is enabled for this App Endpoint.
	DeltaSyncEnabled bool `json:"deltaSyncEnabled"`

	// Scopes is the configuration for scopes within the App Endpoint.
	Scopes ScopesConfig `json:"scopes,omitempty"`

	// Cors is the CORS configuration for the App Endpoint.
	Cors *AppEndpointCors `json:"cors,omitempty"`

	// Oidc is the list of OIDC configurations for the App Endpoint.
	Oidc []AppEndpointOidc `json:"oidc,omitempty"`

	// RequireResync is the list of collections that require resync, keyed by scope.
	RequireResync map[string][]string `json:"requireResync,omitempty"`

	// AdminURL is the admin URL for the App Endpoint.
	AdminURL string `json:"adminURL"`

	// MetricsURL is the metrics URL for the App Endpoint.
	MetricsURL string `json:"metricsURL"`

	// PublicURL is the public URL for the App Endpoint.
	PublicURL string `json:"publicURL"`

	// Etag represents the version of the document
	Etag string

	// Audit contains all audit-related fields.
	Audit api.CouchbaseAuditData `json:"audit"`
}

// UpdateAppEndpointRequest is the request payload sent to the Capella V4 Public API in order to update an existing app endpoint.
type UpdateAppEndpointRequest struct {
	// Bucket is the name of the bucket associated with this App Endpoint.
	Bucket *string `json:"bucket,omitempty"`

	// Name is the name of the App Endpoint.
	Name *string `json:"name,omitempty"`

	// UserXattrKey is the user extended attribute key for the App Endpoint.
	UserXattrKey *string `json:"userXattrKey,omitempty"`

	// DeltaSyncEnabled enables or disables delta sync on this App Endpoint.
	DeltaSyncEnabled *bool `json:"deltaSyncEnabled,omitempty"`

	// Scopes is the configuration for scopes within the App Endpoint.
	Scopes ScopesConfig `json:"scopes,omitempty"`

	// Cors is the CORS configuration for the App Endpoint.
	Cors *AppEndpointCors `json:"cors,omitempty"`

	// Oidc is the list of OIDC providers for the App Endpoint.
	Oidc []AppEndpointOidc `json:"oidc,omitempty"`
}

// ScopesConfig maps scope name to a list of collection names
type (
	ScopesConfig map[string]ScopeConfig
	ScopeConfig  struct {
		Collections map[string]AppEndpointCollection `json:"collections,omitempty"` // Collection-specific config options.
	}
)

// AppEndpointCollection represents a collection configuration.
type AppEndpointCollection struct {
	AccessControlFunction *string `json:"accessControlFunction,omitempty"`
	ImportFilter          *string `json:"importFilter,omitempty"`
}

// AppEndpointCors represents the CORS configuration for an app endpoint.
type AppEndpointCors struct {
	Origin      []string `json:"origin,omitempty"`
	LoginOrigin []string `json:"loginOrigin,omitempty"`
	Headers     []string `json:"headers,omitempty"`
	MaxAge      *int64   `json:"maxAge,omitempty"`
	Disabled    *bool    `json:"disabled,omitempty"`
}

// AppEndpointOidc represents an OIDC configuration within an app endpoint.
type AppEndpointOidc struct {
	Issuer        string  `json:"issuer"`
	Register      *bool   `json:"register,omitempty"`
	ClientId      string  `json:"clientId"`
	UserPrefix    *string `json:"userPrefix,omitempty"`
	DiscoveryUrl  *string `json:"discoveryUrl,omitempty"`
	UsernameClaim *string `json:"usernameClaim,omitempty"`
	RolesClaim    *string `json:"rolesClaim,omitempty"`
	ProviderId    *string `json:"providerId,omitempty"`
	IsDefault     *bool   `json:"isDefault,omitempty"`
}

// ListAppEndpointsResponse is the response received from the Capella V4 Public API when listing app endpoints.
type ListAppEndpointsResponse struct {
	// Data contains the list of app endpoints
	Data []GetAppEndpointResponse `json:"data"`
}
