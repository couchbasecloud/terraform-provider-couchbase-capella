package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServicePrivateEndpointService represents the status of the private endpoint
// service on an App Service.
type AppServicePrivateEndpointService struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Enabled indicates if the private endpoint service is enabled/disabled on the App Service.
	Enabled types.Bool `tfsdk:"enabled"`

	// State is the lifecycle state of the private endpoint service derived from the
	// most recent enable/disable operation. Terminal states are enabled, disabled, and
	// failed; transient states are enabling and disabling. It may be empty when the
	// control plane does not report a state.
	State types.String `tfsdk:"state"`

	// TargetState is the desired end state of the private endpoint service (enabled or
	// disabled). It may be empty when the control plane does not report a target state.
	TargetState types.String `tfsdk:"target_state"`
}

// Validate is used to verify that IDs have been properly imported.
func (a *AppServicePrivateEndpointService) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
