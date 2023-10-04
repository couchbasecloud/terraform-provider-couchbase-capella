package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// Organizations defines the model for GetOrganizations.
//type Organizations struct {
//	// OrganizationId is the organizationId of the capella.
//	OrganizationId types.String `tfsdk:"organization_id"`
//
//	Data []OneOrganization `tfsdk:"data"`
//}

type Organization struct {
	// Audit represents all audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Id is a GUID4 identifier of the project.
	//Id types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name represents the name of the organization
	Name types.String `tfsdk:"name"`

	// Description is a short description of the organization.
	Description types.String `tfsdk:"description"`

	Preferences Preferences `tfsdk:"preferences"`
}

type Preferences struct {
	SessionDuration types.Int64 `tfsdk:"session_duration"`
}
