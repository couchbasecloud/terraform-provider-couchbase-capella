package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AllowList maps AllowList resource schema data to the response received from V4 Capella Public API.
type AllowList struct {
	// Cidr represents the trusted CIDR to allow the database connections from.
	Cidr types.String `tfsdk:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment types.String `tfsdk:"comment"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expires_at"`

	// Id is a GUID4 identifier of the allowlist.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// AllowLists defines the attributes as received from the V4 Capella Public API when asked to list allowlists.
type AllowLists struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data contains the list of resources.
	Data []OneAllowList `tfsdk:"data"`
}

// Validate is used to verify that IDs have been properly imported
func (a *AllowList) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}

// OneAllowList maps allowlist resource schema data; there is a separate response object to avoid conversion error for nested fields.
type OneAllowList struct {
	// Audit represents all audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Cidr is the trusted CIDR to allow the database connections from.
	Cidr types.String `tfsdk:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment types.String `tfsdk:"comment"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expires_at"`

	// Id is a GUID4 identifier of the project.
	Id types.String `tfsdk:"id"`

	// OrganizationId is he organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`
}
