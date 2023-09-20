package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-capella/internal/api"
)

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

func (c CouchbaseAuditData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"version":     types.Int64Type,
	}
}

func NewCouchbaseAuditData(audit api.CouchbaseAuditData) CouchbaseAuditData {
	return CouchbaseAuditData{
		CreatedAt:  types.StringValue(audit.CreatedAt.String()),
		CreatedBy:  types.StringValue(audit.CreatedBy),
		ModifiedAt: types.StringValue(audit.ModifiedAt.String()),
		ModifiedBy: types.StringValue(audit.ModifiedBy),
		Version:    types.Int64Value(int64(audit.Version)),
	}
}
