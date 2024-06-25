package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// PrivateEndpoint represents a private endpoint resource.
type PrivateEndpoint struct {
	// EndpointId is the id of the private endpoint.
	EndpointId types.String `tfsdk:"endpoint_id"`

	// Status is the endpoint status.  Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`

	// ClusterId is the ID of the cluster associated with the private endpoint.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// PrivateEndpoints defines a structure used by the LIST endpoint for private endpoints.
type PrivateEndpoints struct {
	ClusterId      types.String          `tfsdk:"cluster_id"`
	ProjectId      types.String          `tfsdk:"project_id"`
	OrganizationId types.String          `tfsdk:"organization_id"`
	Data           []PrivateEndpointData `tfsdk:"data"`
}

// PrivateEndpointData defines a single private endpoint.
type PrivateEndpointData struct {
	// Id is the endpoint id.
	Id types.String `tfsdk:"id"`
	// Status is the endpoint status.  Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`
}

func (p *PrivateEndpoint) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: p.OrganizationId,
		ProjectId:      p.ProjectId,
		ClusterId:      p.ClusterId,
		EndpointId:     p.EndpointId,
	}

	IDs, err := validateSchemaState(state, EndpointId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
