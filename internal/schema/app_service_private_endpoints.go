package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServicePrivateEndpoint represents a private endpoint connected to an App
// Service's private endpoint service.
type AppServicePrivateEndpoint struct {
	// EndpointId is the id of the private endpoint.
	EndpointId types.String `tfsdk:"endpoint_id"`

	// Status is the endpoint status. Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`

	// ClusterId is the ID of the cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AppServiceId is the ID of the App Service associated with the private endpoint.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// ServiceName is the name of the private endpoint service.
	ServiceName types.String `tfsdk:"service_name"`
}

// AppServicePrivateEndpoints defines a structure used by the LIST endpoint for
// private endpoints connected to an App Service's private endpoint service.
type AppServicePrivateEndpoints struct {
	ClusterId      types.String                    `tfsdk:"cluster_id"`
	ProjectId      types.String                    `tfsdk:"project_id"`
	OrganizationId types.String                    `tfsdk:"organization_id"`
	AppServiceId   types.String                    `tfsdk:"app_service_id"`
	Data           []AppServicePrivateEndpointData `tfsdk:"data"`
}

// AppServicePrivateEndpointData defines a single private endpoint connected to an
// App Service's private endpoint service.
type AppServicePrivateEndpointData struct {
	// Id is the endpoint id.
	Id types.String `tfsdk:"id"`
	// Status is the endpoint status. Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`
	// ServiceName is the name of the endpoint service.
	ServiceName types.String `tfsdk:"service_name"`
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
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
