package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DatabaseCredential maps the schema for the resource - database credential in Capella.
// A database credential is created on a cluster resource to gain read/write access to the cluster's data.
// This credential can have a fixed password supplied during creation or the password can be auto-generated.
//
// A database credential is simply a user in the couchbase server with some roles attached to it
// based on the Access field supplied during creation.
type DatabaseCredential struct {
	// Id is the id of the created database credential.
	Id types.String `tfsdk:"id"`

	// Name is the name of the database credential, the name of the database credential should follow this naming criteria:
	// A database credential name should have at least 2 characters and up to 256 characters and should not contain spaces.
	Name types.String `tfsdk:"name"`

	// Password is the password that you may want to use to create this database credential.
	// This password can later be used to authenticate connections to the underlying couchbase server.
	// The password should contain 8+ characters, at least 1 lower, 1 upper, 1 numerical and 1 special character.
	Password types.String `tfsdk:"password"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Audit All audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`

	// Access is a list of access which can be narrowed to the scope level of every bucket in the Capella cluster.
	// Access can be "read", "write" or both.
	Access []Access `tfsdk:"access"`
}

// Access is a list of privileges or permissions which can be narrowed to the scope level of every bucket in the Capella cluster.
type Access struct {
	// Privileges is a list of permissions that the database credential will have over the data in the given bucket or scope.
	// Privileges can be "read", "write" or both.
	Privileges []types.String `tfsdk:"privileges"`
}

// OneDatabaseCredential is used to retrieve the new state of a database credential after it is created by Terraform.
// This struct is separate from the DatabaseCredential struct because of the change in data type of its attributes after retrieval.
type OneDatabaseCredential struct {
	// Audit All audit-related fields.
	Audit CouchbaseAuditData `tfsdk:"audit"`

	// Id A GUID4 identifier of the created database credential.
	Id types.String `tfsdk:"id"`

	// Name is the name of the database credential, the name of the database credential should follow this naming criteria:
	// A database credential name should have at least 2 characters and up to 256 characters and should not contain spaces.
	Name types.String `tfsdk:"name"`

	// Password is the password that you may want to use to create this database credential.
	// This password can later be used to authenticate connections to the underlying couchbase server.
	// The password should contain 8+ characters, at least 1 lower, 1 upper, 1 numerical and 1 special character.
	Password types.String `tfsdk:"password"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Access is a list of access which can be narrowed to the scope level of every bucket in the Capella cluster.
	// Access can be "read", "write" or both.
	Access []Access `tfsdk:"access"`
}
