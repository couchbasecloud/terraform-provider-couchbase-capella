package capellaschema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectResourceModel maps project resource schema data
type ProjectResourceModel struct {
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

// ProjectResponse maps project resource schema data; there is a separate response object to avoid conversion errors for nested fields.
type ProjectResponse struct {
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

// CouchbaseAuditData contains all audit-related fields.
type CouchbaseAuditData struct {
	// CreatedAt The RFC3339 timestamp associated with when the resource was initially
	// created.
	CreatedAt types.String `tfsdk:"created_at"`

	// CreatedBy The user who created the resource; this will be a UUID4 ID for standard
	// users and will be a string such as "internal-support" for internal
	// Couchbase support users.
	CreatedBy types.String `tfsdk:"created_by"`

	// ModifiedAt The RFC3339 timestamp associated with when the resource was last modified.
	ModifiedAt types.String `tfsdk:"modified_at"`

	// ModifiedBy The user who last modified the resource; this will be a UUID4 ID for
	// standard users and wilmal be a string such asas "internal-support" for
	// internal Couchbase support users.
	ModifiedBy types.String `tfsdk:"modified_by"`

	// Version The version of the document. This value is incremented each time the
	// resource is modified.
	Version types.Int64 `tfsdk:"version"`
}

// ProjectsResourceModel defines model for GetProjectsResponse.
type ProjectsResourceModel struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id""`

	// Data It contains the list of projects.
	Data []ProjectResponse `tfsdk:"data"`
}
