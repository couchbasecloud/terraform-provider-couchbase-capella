package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DatabasePrivileges defines the top-level model for the list database privileges datasource.
type DatabasePrivileges struct {
	OrganizationId types.String            `tfsdk:"organization_id"`
	ProjectId      types.String            `tfsdk:"project_id"`
	ClusterId      types.String            `tfsdk:"cluster_id"`
	Data           []DatabasePrivilegeItem `tfsdk:"data"`
}

// DatabasePrivilegeItem represents a single Capella privilege in the list datasource response.
type DatabasePrivilegeItem struct {
	Name      types.String `tfsdk:"name"`
	Group     types.String `tfsdk:"group"`
	Resources *Resources   `tfsdk:"resources"`
}

// Validate checks that all required IDs are present.
func (d DatabasePrivileges) Validate() (clusterId, projectId, organizationId string, err error) {
	if d.OrganizationId.IsNull() {
		return "", "", "", errors.ErrOrganizationIdMissing
	}
	if d.ProjectId.IsNull() {
		return "", "", "", errors.ErrProjectIdMissing
	}
	if d.ClusterId.IsNull() {
		return "", "", "", errors.ErrClusterIdMissing
	}
	return d.ClusterId.ValueString(), d.ProjectId.ValueString(), d.OrganizationId.ValueString(), nil
}
