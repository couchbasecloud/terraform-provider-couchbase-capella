package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppEndpointActivationStatus provides the means to turn the given app endpoint to on or off state.
type AppEndpointActivationStatus struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the appServiceId of the capella tenant.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpointId is the appEndpointId of the capella tenant.
	AppEndpointId types.String `tfsdk:"app_endpoint_id"`

	// State is the state to which the app endpoint needs to be turned to i.e. on or off.
	State types.String `tfsdk:"state"`
}

func (a *AppEndpointActivationStatus) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		EndpointId:     a.AppEndpointId,
	}

	IDs, err := validateSchemaState(state, EndpointId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
