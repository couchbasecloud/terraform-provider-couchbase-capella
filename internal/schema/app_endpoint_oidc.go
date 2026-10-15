package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// AppEndpointOidcProvider describes the resource data model.
type AppEndpointOidcProvider struct {
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	ClientId        types.String `tfsdk:"client_id"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	DiscoveryUrl    types.String `tfsdk:"discovery_url"`
	Issuer          types.String `tfsdk:"issuer"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	ProjectId       types.String `tfsdk:"project_id"`
	ProviderId      types.String `tfsdk:"provider_id"`
	Register        types.Bool   `tfsdk:"register"`
	RolesClaim      types.String `tfsdk:"roles_claim"`
	UserPrefix      types.String `tfsdk:"user_prefix"`
	UsernameClaim   types.String `tfsdk:"username_claim"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
}

// NewAppEndpointOidcProvider creates the OIDC provider state from the IDs and the get OIDC provider response.
func NewAppEndpointOidcProvider(
	organizationId, projectId, clusterId, appServiceId, appEndpointName string,
	resp api.AppEndpointOIDCProviderResponse,
) *AppEndpointOidcProvider {
	return &AppEndpointOidcProvider{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		ProviderId:      types.StringValue(resp.ProviderID),
		Issuer:          types.StringValue(resp.Issuer),
		ClientId:        types.StringValue(resp.ClientID),
		DiscoveryUrl:    types.StringValue(resp.DiscoveryURL),
		Register:        types.BoolValue(resp.Register),
		RolesClaim:      types.StringValue(resp.RolesClaim),
		UserPrefix:      types.StringValue(resp.UserPrefix),
		UsernameClaim:   types.StringValue(resp.UsernameClaim),
		IsDefault:       types.BoolValue(resp.IsDefault),
	}
}

// Validate validates the AppEndpointActivationStatus resource for import.
func (a *AppEndpointOidcProvider) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  a.OrganizationId,
		ProjectId:       a.ProjectId,
		ClusterId:       a.ClusterId,
		AppServiceId:    a.AppServiceId,
		AppEndpointName: a.AppEndpointName,
		ProviderId:      a.ProviderId,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
