package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type AppEndpoints struct {
	// OrganizationId is the ID of the organization to which the App Endpoints belong.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Endpoints belong.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Endpoints belong.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the App Endpoints belong.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Data is a list of App Endpoints.
	Data []OneAppEndpoint `tfsdk:"data"`
}

// AppEndpoint represents the Terraform schema for an App Endpoint configuration.
type AppEndpoint struct {
	// OrganizationId is the ID of the organization to which the App Endpoint belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Endpoint belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Endpoint belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the App Endpoint belongs.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Bucket The Capella Cluster backing bucket for the App Endpoint.
	Bucket types.String `tfsdk:"bucket"`

	// Name is the name of the App Endpoint.
	Name types.String `tfsdk:"name"`

	// UserXattrKey is the key used for user extended attributes in the App Endpoint.
	UserXattrKey types.String `tfsdk:"user_xattr_key"`

	// DeltaSyncEnabled Indicates whether Delta Sync is enabled for the App Endpoint.
	DeltaSyncEnabled types.Bool `tfsdk:"delta_sync_enabled"`

	// Scope is a map of scope names to scope configurations.
	Scopes types.Map `tfsdk:"scopes"`

	// Cors configures cross origin resource sharing (CORS) for the App Endpoint.
	Cors *AppEndpointCors `tfsdk:"cors"`

	// Oidc is a list of OIDC provider configurations for the App Endpoint.
	Oidc types.Set `tfsdk:"oidc"`

	// RequireResync is a map of scopes to a list of collection names that require resync.
	RequireResync types.Map `tfsdk:"require_resync"`

	// State is the current state of the App Endpoint including online, offline, resyncing, etc.
	State types.String `tfsdk:"state"`

	// AdminURL A URL for the admin API used for the administration of App Endpoints. For more information, read the [Capella App Services Admin API Reference](https://docs.couchbase.com/cloud/app-services/references/rest-api-introduction.html#:~:text=Capella%20App%20Services%20Admin%20API%20Reference)
	AdminURL types.String `tfsdk:"admin_url"`

	// MetricsURL A URL for the metrics API used for monitoring App Services performance metrics. For more information, read the [Capella App Services Metrics API Reference](https://docs.couchbase.com/cloud/app-services/references/rest_api_metric.html)
	MetricsURL types.String `tfsdk:"metrics_url"`

	// PublicURL A URL for the public API used for access to functions for data access and manipulation. For more information, read the [Capella App Services Public API Reference](https://docs.couchbase.com/cloud/app-services/references/rest_api_public.html)
	PublicURL types.String `tfsdk:"public_url"`
}

// Validate verifies required identifiers for the App Endpoint are present and returns them.
func (a *AppEndpoint) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		EndpointName:   a.Name,
	}

	IDs, err := validateSchemaState(state, EndpointName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

// OneAppEndpoint represents the Terraform schema for an App Endpoint configuration.
type OneAppEndpoint struct {
	// Bucket The Capella Cluster backing bucket for the App Endpoint.
	Bucket types.String `tfsdk:"bucket"`

	// Name is the name of the App Endpoint.
	Name types.String `tfsdk:"name"`

	// UserXattrKey is the key used for user extended attributes in the App Endpoint.
	UserXattrKey types.String `tfsdk:"user_xattr_key"`

	// DeltaSyncEnabled Indicates whether Delta Sync is enabled for the App Endpoint.
	DeltaSyncEnabled types.Bool `tfsdk:"delta_sync_enabled"`

	// Scopes is a map of scope names to scope configurations.
	Scopes types.Map `tfsdk:"scopes"`

	// Cors configures cross origin resource sharing (CORS) for the App Endpoint.
	Cors *AppEndpointCors `tfsdk:"cors"`

	// Oidc is a list of OIDC provider configurations for the App Endpoint.
	Oidc types.Set `tfsdk:"oidc"`

	// RequireResync is a map of scopes to a list of collection names that require resync.
	RequireResync types.Map `tfsdk:"require_resync"`

	// AdminURL A URL for the admin API used for the administration of App Endpoints.
	AdminURL types.String `tfsdk:"admin_url"`

	// MetricsURL A URL for the metrics API used for monitoring App Services performance metrics.
	MetricsURL types.String `tfsdk:"metrics_url"`

	// PublicURL A URL for the public API used for access to functions for data access and manipulation.
	PublicURL types.String `tfsdk:"public_url"`

	// State is the current state of the App Endpoint including online, offline, resyncing, etc.
	State types.String `tfsdk:"state"`
}

// AppEndpointScope represents a scope configuration for app endpoints with tfsdk tags
type AppEndpointScope struct {
	Collections types.Map `tfsdk:"collections"`
}

// AppEndpointCollection represents a collection configuration.
type AppEndpointCollection struct {
	// AccessControlFunction The Javascript function that is used to specify the access control policies to be applied to documents in this collection.
	// Every document update is processed by this function. The default access control function is 'function(doc){channel(doc.channels);}'
	// for the default collection and 'function(doc){channel(collectionName);}' for named collections.
	AccessControlFunction types.String `tfsdk:"access_control_function"`

	// ImportFilter The Javascript function used to specify the documents in this collection that are to be imported by the App Endpoint. By default, all documents in corresponding collection are imported.
	ImportFilter types.String `tfsdk:"import_filter"`
}

// AppEndpointCors represents the CORS configuration for an App Endpoint.
type AppEndpointCors struct {
	// Origin List of allowed origins, use ['*'] to allow access from everywhere
	Origin types.Set `tfsdk:"origin"`

	// LoginOrigin List of allowed login origins
	LoginOrigin types.Set `tfsdk:"login_origin"`

	// Headers List of allowed headers
	Headers types.Set `tfsdk:"headers"`

	// MaxAge Specifies the duration (in seconds) for which the results of a preflight request can be cached.
	MaxAge types.Int64 `tfsdk:"max_age"`

	// Disabled indicated whether CORS is disabled.
	Disabled types.Bool `tfsdk:"disabled"`
}

// AppEndpointOidc represents an OIDC configuration within an App Endpoint.
type AppEndpointOidc struct {
	// Issuer The URL for the OpenID Connect issuer.
	Issuer types.String `tfsdk:"issuer"`

	// Register Indicates whether to register a new App Service user account when a user logs in using OpenID Connect.
	Register types.Bool `tfsdk:"register"`

	// ClientId The OpenID Connect provider client ID.
	ClientId types.String `tfsdk:"client_id"`

	// UserPrefix Username prefix for all users created for this provider
	UserPrefix types.String `tfsdk:"user_prefix"`

	// DiscoveryUrl The URL for the non-standard discovery endpoint.
	DiscoveryUrl types.String `tfsdk:"discovery_url"`

	// UsernameClaim Allows a different OpenID Connect field to be specified instead of the Subject (sub).
	UsernameClaim types.String `tfsdk:"username_claim"`

	// RolesClaim If set, the value(s) of the given OpenID Connect authentication token claim will be added to the user's roles.
	// The value of this claim in the OIDC token must be either a string or an array of strings, any other type will result in an error.
	RolesClaim types.String `tfsdk:"roles_claim"`

	// ProviderId UUID of the provider.
	ProviderId types.String `tfsdk:"provider_id"`

	// IsDefault Indicates whether this is the default OpenID Connect provider.
	IsDefault types.Bool `tfsdk:"is_default"`
}

func (o AppEndpointOidc) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"issuer":         types.StringType,
		"register":       types.BoolType,
		"client_id":      types.StringType,
		"user_prefix":    types.StringType,
		"discovery_url":  types.StringType,
		"username_claim": types.StringType,
		"roles_claim":    types.StringType,
		"provider_id":    types.StringType,
		"is_default":     types.BoolType,
	}
}

func (c AppEndpointCors) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"origin":       types.SetType{ElemType: types.StringType},
		"login_origin": types.SetType{ElemType: types.StringType},
		"headers":      types.SetType{ElemType: types.StringType},
		"max_age":      types.Int64Type,
		"disabled":     types.BoolType,
	}
}

func (s AppEndpointScope) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"collections": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: AppEndpointCollection{}.AttributeTypes()},
		},
	}
}

func (c AppEndpointCollection) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_control_function": types.StringType,
		"import_filter":           types.StringType,
	}
}
