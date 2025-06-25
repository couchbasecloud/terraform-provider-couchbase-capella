package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpoint represents the Terraform schema for an app endpoint configuration.
type AppEndpoint struct {
	Bucket           types.String                             `tfsdk:"bucket"`
	Name             types.String                             `tfsdk:"name"`
	UserXattrKey     types.String                             `tfsdk:"userXattrKey"`
	DeltaSyncEnabled types.Bool                               `tfsdk:"deltaSyncEnabled"`
	Scopes           AppEndpointScopes                        `tfsdk:"scopes"`
	Cors             AppEndpointCors                          `tfsdk:"cors"`
	Oidc             []AppEndpointOidc                        `tfsdk:"oidc"`
	RequireResync    map[string]AppEndpointRequireResyncScope `tfsdk:"requireResync"`
	AdminURL         types.String                             `tfsdk:"adminURL"`
	MetricsURL       types.String                             `tfsdk:"metricsURL"`
	PublicURL        types.String                             `tfsdk:"publicURL"`
}

// AppEndpointScopes represents the scopes configuration within an app endpoint.
type AppEndpointScopes struct {
	Default AppEndpointScope `tfsdk:"_default"`
}

// AppEndpointScope represents a scope configuration.
type AppEndpointScope struct {
	Collections AppEndpointCollections `tfsdk:"collections"`
}

// AppEndpointCollections represents the collections configuration within a scope.
type AppEndpointCollections struct {
	Default AppEndpointCollection `tfsdk:"_default"`
}

// AppEndpointCollection represents a collection configuration.
type AppEndpointCollection struct {
	AccessControlFunction types.String `tfsdk:"accessControlFunction"`
	ImportFilter          types.String `tfsdk:"importFilter"`
}

// AppEndpointCors represents the CORS configuration for an app endpoint.
type AppEndpointCors struct {
	Origin      []types.String `tfsdk:"origin"`
	LoginOrigin []types.String `tfsdk:"loginOrigin"`
	Headers     []types.String `tfsdk:"headers"`
	MaxAge      types.Int64    `tfsdk:"maxAge"`
	Disabled    types.Bool     `tfsdk:"disabled"`
}

// AppEndpointOidc represents an OIDC configuration within an app endpoint.
type AppEndpointOidc struct {
	Issuer        types.String `tfsdk:"issuer"`
	Register      types.Bool   `tfsdk:"register"`
	ClientId      types.String `tfsdk:"clientId"`
	UserPrefix    types.String `tfsdk:"userPrefix"`
	DiscoveryUrl  types.String `tfsdk:"discoveryUrl"`
	UsernameClaim types.String `tfsdk:"usernameClaim"`
	RolesClaim    types.String `tfsdk:"rolesClaim"`
	ProviderId    types.String `tfsdk:"providerId"`
	IsDefault     types.Bool   `tfsdk:"isDefault"`
}

// AppEndpointRequireResyncDefault represents the default require resync configuration.
type AppEndpointRequireResyncScope struct {
	Items []types.String `tfsdk:"items"`
}
