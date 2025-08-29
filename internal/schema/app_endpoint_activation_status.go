package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppEndpointActivationStatus provides the means to activate or deactivate an App Endpoint.
type AppEndpointActivationStatus struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the appServiceId of the capella tenant.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpointName is the name of the app endpoint.
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`

	// Online indicates whether the app endpoint should be online (true) or offline (false).
	Online types.Bool `tfsdk:"online"`
}

// Validate validates the AppEndpointActivationStatus resource for import.
func (a *AppEndpointActivationStatus) Validate() (map[Attr]string, error) {
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
