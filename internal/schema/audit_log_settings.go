package schema

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ClusterAuditSettings struct {
	// AuditEnabled Determines whether audit logging is enabled or not on the cluster.
	AuditEnabled types.Bool `tfsdk:"audit_enabled"`

	// DisabledUsers List of users whose filterable events will not be logged.
	DisabledUsers []AuditSettingsDisabledUser `tfsdk:"disabled_users"`

	// EnabledEventIDs List of enabled filterable audit events for the cluster.
	EnabledEventIDs []types.Int64 `tfsdk:"enabled_event_ids"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`
}

type AuditSettingsDisabledUser struct {
	// Domain Specifies whether the user is local or external.
	Domain types.String `tfsdk:"domain"`

	// Name The user name.
	Name types.String `tfsdk:"name"`
}

// Validate is used to verify that IDs have been properly imported.
func (c *ClusterAuditSettings) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		Id:             c.ClusterId,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
