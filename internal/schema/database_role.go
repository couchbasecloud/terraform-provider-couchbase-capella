package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DatabaseRole maps the schema for the database role resource in Capella.
type DatabaseRole struct {
	// Id is the UUID of the created database role.
	Id types.String `tfsdk:"id"`

	// Name is the name of the database role.
	Name types.String `tfsdk:"name"`

	// Description is an optional description of the database role.
	Description types.String `tfsdk:"description"`

	// OrganizationId is the ID of the organization to which the cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster under which the role is created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Audit contains all audit-related fields.
	Audit types.Object `tfsdk:"audit"`

	// Access is a list of access entries defining privileges and resource scopes.
	Access []Access `tfsdk:"access"`
}

// NewDatabaseRole constructs a DatabaseRole from the GET API response fields.
func NewDatabaseRole(
	id, name, description types.String,
	organizationId, projectId, clusterId types.String,
	auditObject basetypes.ObjectValue,
) *DatabaseRole {
	return &DatabaseRole{
		Id:             id,
		Name:           name,
		Description:    description,
		OrganizationId: organizationId,
		ProjectId:      projectId,
		ClusterId:      clusterId,
		Audit:          auditObject,
	}
}

// Validate validates the required IDs for the database role resource state.
func (d DatabaseRole) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: d.OrganizationId,
		ProjectId:      d.ProjectId,
		ClusterId:      d.ClusterId,
		Id:             d.Id,
	}
	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

// DatabaseRoles defines the model for the list database roles datasource response.
type DatabaseRoles struct {
	OrganizationId types.String       `tfsdk:"organization_id"`
	ProjectId      types.String       `tfsdk:"project_id"`
	ClusterId      types.String       `tfsdk:"cluster_id"`
	Data           []DatabaseRoleItem `tfsdk:"data"`
}

// Validate checks that all required IDs are present.
func (d DatabaseRoles) Validate() (clusterId, projectId, organizationId string, err error) {
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

// DatabaseRoleItem represents a single database role in the list datasource response.
type DatabaseRoleItem struct {
	Id             types.String       `tfsdk:"id"`
	Name           types.String       `tfsdk:"name"`
	Description    types.String       `tfsdk:"description"`
	OrganizationId types.String       `tfsdk:"organization_id"`
	ProjectId      types.String       `tfsdk:"project_id"`
	ClusterId      types.String       `tfsdk:"cluster_id"`
	Access         []Access           `tfsdk:"access"`
	Audit          CouchbaseAuditData `tfsdk:"audit"`
}
