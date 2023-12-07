package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Certificate struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	Data []OneCertificate `tfsdk:"data"`
}

type OneCertificate struct {
	Certificate types.String `tfsdk:"certificate"`
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (c *Certificate) Validate() error {
	if c.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if c.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if c.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
