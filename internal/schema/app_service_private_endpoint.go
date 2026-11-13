package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServicePrivateEndpoint represents a private endpoint accepted on an App Service.
type AppServicePrivateEndpoint struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Service is linked.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// EndpointId is the ID of the private endpoint.
	EndpointId types.String `tfsdk:"endpoint_id"`

	// ServiceName is the service name associated with the private endpoint.
	ServiceName types.String `tfsdk:"service_name"`

	// Status is the endpoint status. Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`
}

// Validate is used to verify that IDs have been properly imported.
func (a *AppServicePrivateEndpoint) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		EndpointId:     a.EndpointId,
	}

	IDs, err := validateSchemaState(state, EndpointId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
