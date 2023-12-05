package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/organization"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Organization struct {
	// Audit represents all audit-related fields.
	Audit types.Object `tfsdk:"audit"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name represents the name of the organization
	Name types.String `tfsdk:"name"`

	// Description is a short description of the organization.
	Description types.String `tfsdk:"description"`

	// Preferences stores preferences for the tenant.
	Preferences types.Object `tfsdk:"preferences"`
}

type Preferences struct {
	// SessionDuration: Maximum allowed time in seconds inside the tenant for a user.
	SessionDuration types.Int64 `tfsdk:"session_duration"`
}

func (p Preferences) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"session_duration": types.Int64Type,
	}
}

// NewPreferences create new preferences object
func NewPreferences(preference organization.Preferences) Preferences {
	var sessionDuration int64
	if preference.SessionDuration != nil {
		sessionDuration = int64(*preference.SessionDuration)
	}
	return Preferences{
		SessionDuration: types.Int64Value(sessionDuration),
	}
}
