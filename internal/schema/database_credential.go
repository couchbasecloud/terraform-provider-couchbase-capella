package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DatabaseCredential maps database credential resource schema data
type DatabaseCredential struct {
	Id types.String `tfsdk:"id"`

	// Name The name of the project.
	Name types.String `tfsdk:"name"`

	Password types.String `tfsdk:"password"`

	OrganizationId types.String `tfsdk:"organization_id"`

	ProjectId types.String `tfsdk:"project_id"`

	ClusterId types.String `tfsdk:"cluster_id"`

	Etag types.String `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`

	// Audit All audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`

	Access []Access `tfsdk:"access"`
}

type Access struct {
	Privileges []types.String `tfsdk:"privileges"`
}

type OneDatabaseCredential struct {
	// Audit All audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Id A GUID4 identifier of the project.
	Id types.String `tfsdk:"id"`

	// Name The name of the project.
	Name types.String `tfsdk:"name"`

	Password types.String `tfsdk:"password"`

	OrganizationId types.String `tfsdk:"organization_id"`

	ProjectId types.String `tfsdk:"project_id"`

	ClusterId types.String `tfsdk:"cluster_id"`

	Etag types.String `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`

	Access []Access `tfsdk:"access"`
}
