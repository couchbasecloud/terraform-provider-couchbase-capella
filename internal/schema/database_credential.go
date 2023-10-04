package schema

import (
	"strings"

	"terraform-provider-capella/internal/errors"

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

	// Audit contains all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
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
	// Resources is the level at which the above privileges are defined.
	// Ex: Access of read/write privilege can be defined at the bucket level or scope level resource.
	Resources *Resources `tfsdk:"resources"`
}

// Resources is the level at which the above privileges are defined.
// Ex: Access of read/write privilege can be defined at the bucket level or scope level resource.
type Resources struct {
	// Buckets contains the details of all buckets with scope and collection level information to which the access applies.
	Buckets []Bucket `tfsdk:"buckets"`
}

// Bucket contains the details of a single bucket with scope and collection level information.
// Scopes can be a subset of all scopes inside the bucket, since this is defined only to govern the access.
type Bucket struct {
	Name types.String `tfsdk:"name"`
	// Scopes is the details of the scopes inside the bucket to which we want to apply access privileges.
	Scopes []Scope `tfsdk:"scopes"`
}

// Scope is the details of a single scope inside the bucket, and it contains the collections details too.
// This collections can be a subset of all collections inside the scope, since this is defined only to govern the access.
type Scope struct {
	Name        types.String   `tfsdk:"name"`
	Collections []types.String `tfsdk:"collections"`
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

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import capella_database_credential.new_database_credential id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>`
func (c DatabaseCredential) Validate() (databaseCredentialId, clusterId, projectId, organizationId string, err error) {
	const (
		idDelimiter             = ","
		organizationIdSep       = "organization_id="
		projectIdSep            = "project_id="
		clusterIdSep            = "cluster_id="
		databaseCredentialIdSep = "id="
	)

	organizationId = c.OrganizationId.ValueString()
	projectId = c.ProjectId.ValueString()
	clusterId = c.ClusterId.ValueString()
	databaseCredentialId = c.Id.ValueString()
	var found bool

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if c.OrganizationId.IsNull() {
		strs := strings.Split(c.Id.ValueString(), idDelimiter)
		if len(strs) != 4 {
			err = errors.ErrIdMissing
			return
		}
		_, databaseCredentialId, found = strings.Cut(strs[0], databaseCredentialIdSep)
		if !found {
			err = errors.ErrDatabaseCredentialIdMissing
			return
		}

		_, clusterId, found = strings.Cut(strs[1], clusterIdSep)
		if !found {
			err = errors.ErrClusterIdMissing
			return
		}

		_, projectId, found = strings.Cut(strs[2], projectIdSep)
		if !found {
			err = errors.ErrProjectIdMissing
			return
		}

		_, organizationId, found = strings.Cut(strs[3], organizationIdSep)
		if !found {
			err = errors.ErrOrganizationIdMissing
			return
		}
	}

	if databaseCredentialId == "" {
		err = errors.ErrDatabaseCredentialIdCannotBeEmpty
		return
	}

	if clusterId == "" {
		err = errors.ErrClusterIdCannotBeEmpty
		return
	}

	if projectId == "" {
		err = errors.ErrProjectIdCannotBeEmpty
		return
	}

	if organizationId == "" {
		err = errors.ErrOrganizationIdCannotBeEmpty
		return
	}

	return databaseCredentialId, clusterId, projectId, organizationId, nil
}
