package api

// AppEndpointOIDCDefaultProviderRequest represents the payload to set the default OIDC provider.
type AppEndpointOIDCDefaultProviderRequest struct {
	ProviderID string `json:"providerId"`
}

// AppEndpointOIDCProviderRequest represents the payload to create/update an OIDC provider for an App Endpoint.
type AppEndpointOIDCProviderRequest struct {
	// Issuer is the OIDC issuer URL.
	Issuer string `json:"issuer"`
	// ClientID is the client ID registered with the OIDC provider.
	ClientID string `json:"clientId"`
	// DiscoveryURL is the URL to fetch the OIDC discovery document.
	DiscoveryURL *string `json:"discoveryUrl,omitempty"`
	// Register indicates if users can self-register.
	Register *bool `json:"register,omitempty"`
	// RolesClaim is the claim in the OIDC token that contains user roles.
	RolesClaim *string `json:"rolesClaim,omitempty"`
	// UserPrefix is the prefix added to usernames from this provider.
	UserPrefix *string `json:"userPrefix,omitempty"`
	// UsernameClaim is the claim in the OIDC token that contains the username.
	UsernameClaim *string `json:"usernameClaim,omitempty"`
}

// AppEndpointOIDCProviderResponse models the OIDC provider as returned by the API.
type AppEndpointOIDCProviderResponse struct {
	// ProviderID is the unique identifier for the OIDC provider.
	ProviderID string `json:"providerId,omitempty"`
	// Issuer is the OIDC issuer URL.
	Issuer string `json:"issuer,omitempty"`
	// ClientID is the client ID registered with the OIDC provider.
	ClientID string `json:"clientId,omitempty"`
	// DiscoveryURL is the URL to fetch the OIDC discovery document.
	DiscoveryURL string `json:"discoveryUrl,omitempty"`
	// Register indicates if users can self-register.
	Register bool `json:"register,omitempty"`
	// RolesClaim is the claim in the OIDC token that contains user roles.
	RolesClaim string `json:"rolesClaim,omitempty"`
	// UserPrefix is the prefix added to usernames from this provider.
	UserPrefix string `json:"userPrefix,omitempty"`
	// UsernameClaim is the claim in the OIDC token that contains the username.
	UsernameClaim string `json:"usernameClaim,omitempty"`
	// IsDefault indicates if this provider is the default for the App Endpoint.
	IsDefault bool `json:"isDefault,omitempty"`
}

// AppEndpointOIDCProviderListResponse models a list of OIDC providers.
type AppEndpointOIDCProviderListResponse struct {
	Data []AppEndpointOIDCProviderResponse `json:"data"`
}
