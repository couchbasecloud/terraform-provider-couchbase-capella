package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
}

// Validate validates the AppEndpointActivationStatus resource for import.
func (a *AppEndpointOidcProvider) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  a.OrganizationId,
		ProjectId:       a.ProjectId,
		ClusterId:       a.ClusterId,
		AppServiceId:    a.AppServiceId,
		AppEndpointName: a.AppEndpointName,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
