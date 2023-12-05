package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

func (p Project) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: p.OrganizationId,
		Id:             p.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// Projects defines the attributes for a list of projects in Capella.
type Projects struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

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
