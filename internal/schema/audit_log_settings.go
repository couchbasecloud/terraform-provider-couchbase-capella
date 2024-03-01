package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ClusterAuditSettings struct {
	// AuditEnabled Determines whether audit logging is enabled or not on the cluster.
	AuditEnabled types.Bool `tfsdk:"auditenabled"`

	// DisabledUsers List of users whose filterable events will not be logged.
	DisabledUsers []AuditSettingsDisabledUser `tfsdk:"disabledusers"`

	// EnabledEventIDs List of enabled filterable audit events for the cluster.
	EnabledEventIDs []types.Int64 `tfsdk:"enabledeventids"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`
}

type AuditSettingsDisabledUser struct {
	// Domain Specifies whether the user is local or external.
	Domain types.String `tfsdk:"domain,omitempty"`

	// Name The user name.
	Name types.String `tfsdk:"name,omitempty"`
}

// Validate is used to verify that IDs have been properly imported.
func (c *ClusterAuditSettings) Validate() error {
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
