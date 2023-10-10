package api

import (
	"github.com/google/uuid"
)

// Access defines the level of access that the database credential will have across buckets and scopes.
// This access is currently defined for all buckets and all scopes in the cluster.
// todo: Support for granular access per bucket and per scope will be added in AV-62864
type Access struct {
	Privileges []string `json:"privileges"`
	// Resources is the level at which the above privileges are defined.
	// Ex: Access of read/write privilege can be defined at the bucket level or scope level resource.
	Resources *AccessibleResources `json:"resources,omitempty"`
}

// CreateDatabaseCredentialRequest represents the schema for the POST Capella V4 API request that creates the database credential.
// Password is an optional field, if not passed, the password for the database credential is auto-generated.
type CreateDatabaseCredentialRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Access   []Access `json:"access"`
}

// AccessibleResources is the level at which the above privileges are defined.
// Ex: Access of read/write privilege can be defined at the bucket level or scope level resource.
type AccessibleResources struct {
	// Buckets contains the details of all buckets with scope and collection level information to which the access applies.
	Buckets []Bucket `json:"buckets"`
}

// Bucket contains the details of a single bucket with scope and collection level information.
// Scopes can be a subset of all scopes inside the bucket, since this is defined only to govern the access.
type Bucket struct {
	Name string `json:"name"`
	// Scopes is the details of the scopes inside the bucket to which we want to apply access privileges.
	Scopes []Scope `json:"scopes,omitempty"`
}

// Scope is the details of a single scope inside the bucket, and it contains the collections details too.
// This collections can be a subset of all collections inside the scope, since this is defined only to govern the access.
type Scope struct {
	Name        string   `json:"name"`
	Collections []string `json:"collections,omitempty"`
}

// CreateDatabaseCredentialResponse represents the schema for the POST Capella V4 API response that creates the database credential.
type CreateDatabaseCredentialResponse struct {
	Id       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

// GetDatabaseCredentialResponse represents the schema for the GET Capella V4 API request that fetches the database credential details.
type GetDatabaseCredentialResponse struct {
	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// Id A GUID4 identifier of the project.
	Id uuid.UUID `json:"id"`

	Name string `json:"name"`

	Access []Access `json:"access"`
}

// PutDatabaseCredentialRequest represents the schema for the PUT Capella V4 API request that updates an existing database credential.
// Password is an optional field, if not passed, the existing password is not updated.
type PutDatabaseCredentialRequest struct {
	Password string   `json:"password,omitempty"`
	Access   []Access `json:"access"`
}
