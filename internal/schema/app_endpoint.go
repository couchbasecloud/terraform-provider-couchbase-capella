package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpoint represents the Terraform schema for an app endpoint configuration.
type AppEndpoint struct {
	OrganizationId   types.String                    `tfsdk:"organization_id"`
	ProjectId        types.String                    `tfsdk:"project_id"`
	ClusterId        types.String                    `tfsdk:"cluster_id"`
	AppServiceId     types.String                    `tfsdk:"app_service_id"`
	Bucket           types.String                    `tfsdk:"bucket"`
	Name             types.String                    `tfsdk:"name"`
	UserXattrKey     types.String                    `tfsdk:"userXattrKey"`
	DeltaSyncEnabled types.Bool                      `tfsdk:"deltaSyncEnabled"`
	Scopes           AppEndpointScopes               `tfsdk:"scopes"`
	Cors             AppEndpointCors                 `tfsdk:"cors"`
	Oidc             []AppEndpointOidc               `tfsdk:"oidc"`
	RequireResync    map[types.String][]types.String `tfsdk:"requireResync"`
	AdminURL         types.String                    `tfsdk:"adminURL"`
	MetricsURL       types.String                    `tfsdk:"metricsURL"`
	PublicURL        types.String                    `tfsdk:"publicURL"`
}

// ScopesConfig maps scope name to a list of collection names
type (
	AppEndpointScopes      map[types.String]AppEndpointScopeConfig
	AppEndpointScopeConfig struct {
		Collections map[types.String]AppEndpointCollection `json:"collections,omitempty"` // Collection-specific config options.
	}
)

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
