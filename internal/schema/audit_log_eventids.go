package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type AuditLogEventIDs struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	Data []AuditLogEventID `tfsdk:"data"`
}

type AuditLogEventID struct {
	Description types.String `tfsdk:"description"`
	Id          types.Int64  `tfsdk:"id"`
	Module      types.String `tfsdk:"module"`
	Name        types.String `tfsdk:"name"`
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (a *AuditLogEventIDs) Validate() error {
	if a.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if a.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if a.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
