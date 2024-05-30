package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type PrivateEndpointService struct {
	OrganizationId types.String `tfsdk:"organization_id"`

	ProjectId types.String `tfsdk:"project_id"`

	ClusterId types.String `tfsdk:"cluster_id"`

	Enabled types.Bool `tfsdk:"enabled"`
}

func (p *PrivateEndpointService) Validate() error {
	if p.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if p.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if p.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
