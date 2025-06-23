package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppEndpoint represents the Terraform schema for an app endpoint configuration
type AppEndpoint struct {
	Bucket           types.String    `tfsdk:"bucket"`
	Name             types.String    `tfsdk:"name"`
	UserXattrKey     types.String    `tfsdk:"user_xattr_key"`
	DeltaSyncEnabled types.Bool      `tfsdk:"delta_sync_enabled"`
	Scopes           types.Object    `tfsdk:"scopes"`
	Cors             AppEndpointCors `tfsdk:"cors"`
	Oidc             AppEndpointOidc `tfsdk:"oidc"`
}

// AppEndpointScopes represents the scopes configuration within an app endpoint
type AppEndpointScopes struct {
	Default types.Object `tfsdk:"_default"`
}

// AppEndpointDefaultScope represents the default scope configuration
type AppEndpointDefaultScope struct {
	Collections types.Object `tfsdk:"collections"`
}

// AppEndpointCollections represents the collections configuration within a scope
type AppEndpointCollections struct {
	Default types.Object `tfsdk:"_default"`
}

// AppEndpointDefaultCollection represents the default collection configuration
type AppEndpointDefaultCollection struct {
	AccessControlFunction types.String `tfsdk:"access_control_function"`
	ImportFilter          types.String `tfsdk:"import_filter"`
}

// AppEndpointCors represents the CORS configuration for an app endpoint
type AppEndpointCors struct {
	Origin      types.List `tfsdk:"origin"`
	LoginOrigin types.List `tfsdk:"login_origin"`
	Headers     types.List `tfsdk:"headers"`
	Disabled    types.Bool `tfsdk:"disabled"`
}

// AppEndpointOidc represents an OIDC configuration within an app endpoint
type AppEndpointOidc struct {
	Issuer        types.String `tfsdk:"issuer"`
	Register      types.Bool   `tfsdk:"register"`
	ClientId      types.String `tfsdk:"client_id"`
	UserPrefix    types.String `tfsdk:"user_prefix"`
	DiscoveryUrl  types.String `tfsdk:"discovery_url"`
	UsernameClaim types.String `tfsdk:"username_claim"`
	RolesClaim    types.String `tfsdk:"roles_claim"`
}
