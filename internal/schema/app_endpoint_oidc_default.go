package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppEndpointDefaultOidcProvider describes the resource data model for default OIDC provider selection.
type AppEndpointDefaultOidcProvider struct {
	OrganizationId  types.String `tfsdk:"organization_id"`
	ProjectId       types.String `tfsdk:"project_id"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`
	ProviderId      types.String `tfsdk:"provider_id"`
}

// NewAppEndpointDefaultOidcProvider creates the default OIDC provider state from the IDs.
func NewAppEndpointDefaultOidcProvider(
	organizationId, projectId, clusterId, appServiceId, appEndpointName, providerId string,
) *AppEndpointDefaultOidcProvider {
	return &AppEndpointDefaultOidcProvider{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		ProviderId:      types.StringValue(providerId),
	}
}

// Validate validates base identifiers using shared helper.
func (a *AppEndpointDefaultOidcProvider) Validate() (map[Attr]string, error) {
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
