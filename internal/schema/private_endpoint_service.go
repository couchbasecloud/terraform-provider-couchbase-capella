package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// PrivateEndpointService represents the status of private endpoint service on a cluster.
type PrivateEndpointService struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster associated with the private endpoint service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Enabled indicates if private endpoint service is enabled/disabled on cluster.
	Enabled types.Bool `tfsdk:"enabled"`
}

// Validate is used to verify that IDs have been properly imported.
func (p *PrivateEndpointService) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: p.OrganizationId,
		ProjectId:      p.ProjectId,
		ClusterId:      p.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
