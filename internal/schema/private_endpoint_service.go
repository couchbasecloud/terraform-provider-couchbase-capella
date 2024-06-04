package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type PrivateEndpointService struct {
	OrganizationId types.String `tfsdk:"organization_id"`

	ProjectId types.String `tfsdk:"project_id"`

	ClusterId types.String `tfsdk:"cluster_id"`

	Enabled types.Bool `tfsdk:"enabled"`
}

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
