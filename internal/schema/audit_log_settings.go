package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

type ClusterAuditSettings struct {
	OrganizationId  types.String                `tfsdk:"organization_id"`
	ProjectId       types.String                `tfsdk:"project_id"`
	ClusterId       types.String                `tfsdk:"cluster_id"`
	DisabledUsers   []AuditSettingsDisabledUser `tfsdk:"disabled_users"`
	EnabledEventIDs []types.Int64               `tfsdk:"enabled_event_ids"`
	AuditEnabled    types.Bool                  `tfsdk:"audit_enabled"`
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
