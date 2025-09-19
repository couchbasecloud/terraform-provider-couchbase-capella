package api

// AppEndpointOIDCDefaultProviderRequest represents the payload to set the default OIDC provider.
type AppEndpointOIDCDefaultProviderRequest struct {
	ProviderID string `json:"providerId"`
}

// AppEndpointOIDCProviderRequest represents the payload to create/update an OIDC provider for an App Endpoint.
type AppEndpointOIDCProviderRequest struct {
	Issuer        string  `json:"issuer"`
	ClientID      string  `json:"clientId"`
	DiscoveryURL  *string `json:"discoveryUrl,omitempty"`
	Register      *bool   `json:"register,omitempty"`
	RolesClaim    *string `json:"rolesClaim,omitempty"`
	UserPrefix    *string `json:"userPrefix,omitempty"`
	UsernameClaim *string `json:"usernameClaim,omitempty"`
}

// AppEndpointOIDCProviderResponse models the OIDC provider as returned by the API.
type AppEndpointOIDCProviderResponse struct {
	ProviderID    string `json:"providerId,omitempty"`
	Issuer        string `json:"issuer,omitempty"`
	ClientID      string `json:"clientId,omitempty"`
	DiscoveryURL  string `json:"discoveryUrl,omitempty"`
	Register      bool   `json:"register,omitempty"`
	RolesClaim    string `json:"rolesClaim,omitempty"`
	UserPrefix    string `json:"userPrefix,omitempty"`
	UsernameClaim string `json:"usernameClaim,omitempty"`
	IsDefault     bool   `json:"isDefault,omitempty"`
}

// AppEndpointOIDCProviderListResponse models a list of OIDC providers.
type AppEndpointOIDCProviderListResponse struct {
	Data []AppEndpointOIDCProviderResponse `json:"data"`
}
