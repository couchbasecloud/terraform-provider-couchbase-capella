package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// User maps User resource schema data
type User struct {
	// Id is a GUID4 identifier of the user.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name represents the name of the user.
	Name types.String `tfsdk:"name"`

	// Name represents the email of the user.
	Email types.String `tfsdk:"email"`

	// OrganizationRoles is an array of strings representing the roles granted to the user
	OrganizationRoles []types.String `tfsdk:"organizationRoles"`

	// Resources is an array of objects representing the resources the user has access to
	Resources types.Object `tfsdk:"resources"`

	// ETag is a unique indentifier which the client uses to determine if the resource has changed.
	ETag types.String `tfsdk:"if_match"`

	// IfMatch is used to check if a request should be made. The request will only proceed if
	// the resources current ETag matches this value.
	IfMatch types.String `tfsdk:"if_match"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// OneAllowList maps user resource schema data; there is a separate response object to avoid conversion error for nested fields.
type OneUser struct {
	// Audit represents all audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Id is a GUID4 identifier of the user.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name represents the name of the user.
	Name types.String `tfsdk:"name"`

	// Name represents the email of the user.
	Email types.String `tfsdk:"email"`

	// OrganizationRoles is an array of strings representing the roles granted to the user
	OrganizationRoles []types.String `tfsdk:"organizationRoles"`

	// Resources is an array of objects representing the resources the user has access to
	Resources types.Object `tfsdk:"resources"`

	// ETag is a unique indentifier which the client uses to determine if the resource has changed.
	ETag types.String `tfsdk:"if_match"`

	// IfMatch is used to check if a request should be made. The request will only proceed if
	// the resources current ETag matches this value.
	IfMatch types.String `tfsdk:"if_match"`
}
