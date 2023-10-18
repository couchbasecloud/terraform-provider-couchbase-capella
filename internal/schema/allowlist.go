package schema

import (
	"fmt"
	"strings"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

	IfMatch types.String `tfsdk:"if_match"`

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
func (a *AllowList) Validate() (map[string]string, error) {
	const idDelimiter = ","
	var found bool

	organizationId := a.OrganizationId.ValueString()
	projectId := a.ProjectId.ValueString()
	clusterId := a.ClusterId.ValueString()
	allowListId := a.Id.ValueString()

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if a.OrganizationId.IsNull() {
		strs := strings.Split(a.Id.ValueString(), idDelimiter)
		if len(strs) != 4 {
			return nil, errors.ErrIdMissing
		}

		_, allowListId, found = strings.Cut(strs[0], "id=")
		if !found {
			return nil, errors.ErrAllowListIdMissing
		}

		_, organizationId, found = strings.Cut(strs[1], "organization_id=")
		if !found {
			return nil, errors.ErrOrganizationIdMissing
		}

		_, projectId, found = strings.Cut(strs[2], "project_id=")
		if !found {
			return nil, errors.ErrProjectIdMissing
		}

		_, clusterId, found = strings.Cut(strs[3], "cluster_id=")
		if !found {
			return nil, errors.ErrClusterIdMissing
		}
	}

	resourceIDs := a.generateResourceIdMap(organizationId, projectId, clusterId, allowListId)

	err := a.checkEmpty(resourceIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errors.ErrUnableToImportResource, err)
	}

	return resourceIDs, nil
}

// generateResourceIdmap is used to populate a map with selected IDs
func (a *AllowList) generateResourceIdMap(organizationId, projectId, clusterId, allowListId string) map[string]string {
	return map[string]string{
		"organizationId": organizationId,
		"projectId":      projectId,
		"clusterId":      clusterId,
		"allowListId":    allowListId,
	}
}

// checkEmpty is used to verify that a supplied resourceId map has been populated
func (a *AllowList) checkEmpty(resourceIdMap map[string]string) error {
	if resourceIdMap["allowListId"] == "" {
		return errors.ErrAllowListIdCannotBeEmpty
	}

	if resourceIdMap["clusterId"] == "" {
		return errors.ErrClusterIdCannotBeEmpty
	}

	if resourceIdMap["projectId"] == "" {
		return errors.ErrProjectIdCannotBeEmpty
	}

	if resourceIdMap["organizationId"] == "" {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	return nil
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

	IfMatch types.String `tfsdk:"if_match"`
}
