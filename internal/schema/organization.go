package schema

import (
	"terraform-provider-capella/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Organizations defines the model for GetOrganizations.
//type Organizations struct {
//	// OrganizationId is the organizationId of the capella.
//	OrganizationId types.String `tfsdk:"organization_id"`
//
//	Data []OneOrganization `tfsdk:"data"`
//}

type Organization struct {
	// Audit represents all audit-related fields.
	Audit types.Object `tfsdk:"audit"`

	// Id is a GUID4 identifier of the project.
	//Id types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`

	// Name represents the name of the organization
	Name types.String `tfsdk:"name"`

	// Description is a short description of the organization.
	Description types.String `tfsdk:"description"`

	Preferences types.Object `tfsdk:"preferences"`
}

type Preferences struct {
	SessionDuration types.Int64 `tfsdk:"session_duration"`
}

func (p Preferences) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"session_duration": types.Int64Type,
	}
}

func NewPreferences(preference api.Preferences) Preferences {
	var sessionDuration int64
	if preference.SessionDuration != nil {
		sessionDuration = int64(*preference.SessionDuration)
	}
	return Preferences{
		SessionDuration: types.Int64Value(sessionDuration),
	}
}
