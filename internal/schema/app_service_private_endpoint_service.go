package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServicePrivateEndpointService represents the private endpoint service on an App Service.
type AppServicePrivateEndpointService struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Service is linked.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service on which the private endpoint service is enabled.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// State is the status of the private endpoint service: enabled, enabling, disabled, disabling or failed.
	State types.String `tfsdk:"state"`
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
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
