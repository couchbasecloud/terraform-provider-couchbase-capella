package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// APIKeyResourcesItems defines model for APIKeyResourcesItems.
type APIKeyResourcesItems struct {
	// Id ID of the project.
	Id types.String `tfsdk:"id"`

	// Roles Project Roles associated with the API key.
	//
	// To learn more about Project Roles, see [Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
	Roles []types.String `tfsdk:"roles"`

	// Type Type of the resource.
	Type types.String `tfsdk:"type"`
}

type ApiKey struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AllowedCIDRs List of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs []types.String `tfsdk:"allowed_cidrs"`
	Audit        types.Object   `tfsdk:"audit"`

	// Description Description for the API key.
	Description types.String `tfsdk:"description"`

	// Expiry Expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry types.Float64 `tfsdk:"expiry"`

	// Id The id is a unique identifier for an apiKey.
	Id types.String `tfsdk:"id"`

	// Name Name of the API key.
	Name              types.String   `tfsdk:"name"`
	OrganizationRoles []types.String `tfsdk:"organization_roles"`

	// Resources Resources are the resource level permissions associated with the API key.
	//
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources []APIKeyResourcesItems `tfsdk:"resources"`
}
