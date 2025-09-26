package app_endpoints

// AppEndpointRequest is the request payload sent to the Capella V4 Public API in order to create a new App Endpoint.
// An App Endpoint provides a REST API for accessing data in a Couchbase Capella bucket through an App Service.
//
// To learn more about App Endpoints see https://docs.couchbase.com/cloud/app-services/connect/connect-apps-to-endpoint.html
type AppEndpointRequest struct {
	// Bucket is the name of the bucket associated with this App Endpoint.
	Bucket string `json:"bucket"`

	// Name is the name of the App Endpoint.
	Name string `json:"name"`

	// UserXattrKey is the user extended attribute key for the App Endpoint.
	UserXattrKey string `json:"userXattrKey,omitempty"`

	// DeltaSyncEnabled enables or disables delta sync on this App Endpoint.
	DeltaSyncEnabled bool `json:"deltaSyncEnabled,omitempty"`

	// Scopes is the configuration for scopes and collections within the App Endpoint.
	Scopes Scopes `json:"scopes,omitempty"`

	// Cors is the CORS configuration for the App Endpoint.
	Cors *AppEndpointCors `json:"cors,omitempty"`

	// Oidc is the list of OIDC providers for the App Endpoint.
	Oidc []AppEndpointOidc `json:"oidc,omitempty"`
}

// GetAppEndpointResponse is the response received from the Capella V4 Public API when asked to fetch details of
// an existing app endpoint.
type GetAppEndpointResponse struct {
	// Bucket is the name of the bucket associated with this App Endpoint.
	Bucket string `json:"bucket"`

	// Name is the name of the App Endpoint.
	Name string `json:"name"`

	// UserXattrKey is the user extended attribute key for the App Endpoint.
	UserXattrKey string `json:"userXattrKey"`

	// DeltaSyncEnabled indicates whether delta sync is enabled for this App Endpoint.
	DeltaSyncEnabled bool `json:"deltaSyncEnabled"`

	// Scopes is the configuration for scopes within the App Endpoint.
	Scopes Scopes `json:"scopes"`

	// Cors is the CORS configuration for the App Endpoint.
	Cors *AppEndpointCors `json:"cors"`

	// Oidc is the list of OIDC configurations for the App Endpoint.
	Oidc []AppEndpointOidc `json:"oidc"`

	// RequireResync is the list of collections that require resync, keyed by scope.
	RequireResync map[string][]string `json:"requireResync"`

	// AdminURL is the admin URL for the App Endpoint.
	AdminURL string `json:"adminURL"`

	// MetricsURL is the metrics URL for the App Endpoint.
	MetricsURL string `json:"metricsURL"`

	// PublicURL is the public URL for the App Endpoint.
	PublicURL string `json:"publicURL"`

	// State is the current state of the App Endpoint, such as online, offline, resyncing, etc.
	State string `json:"state"`
}

type Scopes map[string]Scope

type Scope struct {
	Collections map[string]Collection `json:"collections"`
}

type Collection struct {
	AccessControlFunction string `json:"accessControlFunction,omitempty"`
	ImportFilter          string `json:"importFilter,omitempty"`
}

// AppEndpointCors represents the CORS configuration for an app endpoint.
type AppEndpointCors struct {
	Origin      []string `json:"origin,omitempty"`
	LoginOrigin []string `json:"loginOrigin,omitempty"`
	Headers     []string `json:"headers,omitempty"`
	MaxAge      int64    `json:"maxAge,omitempty"`
	Disabled    bool     `json:"disabled,omitempty"`
}

// AppEndpointOidc represents an OIDC configuration within an app endpoint.
type AppEndpointOidc struct {
	Issuer        string `json:"issuer"`
	Register      bool   `json:"register,omitempty"`
	ClientId      string `json:"clientId"`
	UserPrefix    string `json:"userPrefix,omitempty"`
	DiscoveryUrl  string `json:"discoveryUrl,omitempty"`
	UsernameClaim string `json:"usernameClaim,omitempty"`
	RolesClaim    string `json:"rolesClaim,omitempty"`
	ProviderId    string `json:"providerId,omitempty"`
	IsDefault     bool   `json:"isDefault,omitempty"`
}

// ListAppEndpointsResponse is the response received from the Capella V4 Public API when listing app endpoints.
type ListAppEndpointsResponse struct {
	// Data contains the list of app endpoints
	Data []GetAppEndpointResponse `json:"data"`
}
