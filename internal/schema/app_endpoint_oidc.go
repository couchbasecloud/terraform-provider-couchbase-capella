package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppEndpointOidcProvider describes the resource data model.
type AppEndpointOidcProvider struct {
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`
	AppServiceID    types.String `tfsdk:"app_service_id"`
	ClientID        types.String `tfsdk:"client_id"`
	ClusterID       types.String `tfsdk:"cluster_id"`
	DiscoveryURL    types.String `tfsdk:"discovery_url"`
	Issuer          types.String `tfsdk:"issuer"`
	OrganizationID  types.String `tfsdk:"organization_id"`
	ProjectID       types.String `tfsdk:"project_id"`
	ProviderID      types.String `tfsdk:"provider_id"`
	Register        types.Bool   `tfsdk:"register"`
	RolesClaim      types.String `tfsdk:"roles_claim"`
	UserPrefix      types.String `tfsdk:"user_prefix"`
	UsernameClaim   types.String `tfsdk:"username_claim"`
}
