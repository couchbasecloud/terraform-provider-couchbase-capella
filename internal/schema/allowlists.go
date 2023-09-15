package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AllowList maps AllowList resource schema data
type AllowList struct {
	// The trusted CIDR to allow the database connections from.
	Cidr types.String `tfsdk:"cidr"`

	// A short description of the allowed CIDR.
	Comment types.String `tfsdk:"comment"`

	// An RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expiresAt"`

	// Id A GUID4 identifier of the allowlist.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	IfMatch types.String `tfsdk:"if_match"`

	// Audit All audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// AllowLists defines model for GetAllowLists.
type AllowLists struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id""`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data It contains the list of resources.
	Data []OneProject `tfsdk:"data"`
}

// OneAllowList maps allowlist resource schema data; there is a separate response object to avoid conversion error for nested fields.
type OneAllowList struct {
	// Audit All audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// The trusted CIDR to allow the database connections from.
	Cidr types.String `tfsdk:"cidr"`

	// A short description of the allowed CIDR.
	Comment types.String `tfsdk:"comment"`

	// An RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expiresAt"`

	// Id A GUID4 identifier of the project.
	Id types.String `tfsdk:"id"`

	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	IfMatch types.String `tfsdk:"if_match"`
}
