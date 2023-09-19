package schema

import (
	"strings"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Project maps project resource schema data
type Project struct {
	// Description The description of a particular project.
	Description types.String `tfsdk:"description"`

	// Id A GUID4 identifier of the project.
	Id types.String `tfsdk:"id"`

	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name The name of the project.
	Name types.String `tfsdk:"name"`

	Etag types.String `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`

	// Audit All audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

func (p Project) Validate() (projectId string, organizationId string, err error) {
	organizationId = p.OrganizationId.ValueString()
	projectId = p.Id.ValueString()
	var found bool

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if p.OrganizationId.IsNull() {
		strs := strings.Split(p.Id.ValueString(), ",")
		if len(strs) > 1 {
			_, projectId, found = strings.Cut(strs[0], "id=")
			if !found || projectId == "" {
				return "", "", errors.ErrProjectIdMissing
			}

			_, organizationId, found = strings.Cut(strs[1], "organization_id=")
			if !found || organizationId == "" {
				return "", "", errors.ErrOrganizationIdMissing
			}
		} else {
			return "", "", errors.ErrOrganizationIdMissing
		}
	}

	if projectId == "" {
		return "", "", errors.ErrProjectIdCannotBeEmpty
	}

	if organizationId == "" {
		return "", "", errors.ErrOrganizationIdCannotBeEmpty
	}

	return projectId, organizationId, nil
}

// Projects defines model for GetProjectsResponse.
type Projects struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id""`

	// Data It contains the list of resources.
	Data []OneProject `tfsdk:"data"`
}

// OneProject maps project resource schema data; there is a separate response object to avoid conversion error for nested fields.
type OneProject struct {
	// Audit All audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Description The description of a particular project.
	Description types.String `tfsdk:"description"`

	// Id A GUID4 identifier of the project.
	Id types.String `tfsdk:"id"`

	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name The name of the project.
	Name types.String `tfsdk:"name"`

	Etag types.String `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`
}
