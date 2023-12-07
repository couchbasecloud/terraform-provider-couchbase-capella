package api

import (
	"github.com/google/uuid"
)

// Access defines the level of access that the database credential will have across buckets and scopes.
// This access is currently defined for all buckets and all scopes in the cluster.
type Access struct {
	// Resources is the level at which the above privileges are defined.
	// Ex: Access of read/write privilege can be defined at the bucket level or scope level resource.
	// Leaving this empty will grant access to all buckets.
	Resources *AccessibleResources `json:"resources,omitempty"`

	// Privileges defines the privileges field in this API represents the privilege level for users.
	// It accepts one of the following values:
	// data_reader
	// data_writer
	// read: Equivalent to data_reader.
	// write: Equivalent to data_writer.
	Privileges []string `json:"privileges"`
}

// CreateDatabaseCredentialRequest represents the schema for the POST Capella V4 API request that creates the database credential.
// Password is an optional field, if not passed, the password for the database credential is auto-generated.
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateDatabaseCredentialRequest struct {
	// Name is Username for the database credential. The name should be according to the following rules:
	// The name must be between 2 & 128 characters.
	// The name cannot contain spaces.
	// The name cannot contain the following characters - ) ( > < , ; : " \ / ] [ ? = } {
	// The name cannot begin with @ character.
	Name string `json:"name"`

	// Password is the password associated with the database credential.
	// If this field is left empty, a password will be auto-generated.
	// The password should be according to the following rules:
	// Password should have at least 8 or more characters.
	// Characters used for the password should contain at least one uppercase (A-Z), one lowercase (a-z), one numerical (0-9), and one special character.
	// Forbidden special characters for the password are: < > ; . * & | £
	Password string `json:"password,omitempty"`

	// Access describes the access information of the database credential.
	Access []Access `json:"access"`
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
	// Name is the name of the bucket.
	Name string `json:"name"`
	// Scopes is the details of the scopes inside the bucket to which we want to apply access privileges.
	Scopes []Scope `json:"scopes,omitempty"`
}

// Scope is the details of a single scope inside the bucket, and it contains the collections details too.
// This collections can be a subset of all collections inside the scope, since this is defined only to govern the access.
type Scope struct {
	// Name is the name of the scope.
	Name string `json:"name"`
	// Collections is the set of collections that are present in the provided scope to which the database credential should have access.
	Collections []string `json:"collections,omitempty"`
}

// CreateDatabaseCredentialResponse represents the schema for the POST Capella V4 API response that creates the database credential.
type CreateDatabaseCredentialResponse struct {
	Password string    `json:"password"`
	Id       uuid.UUID `json:"id"`
}

// GetDatabaseCredentialResponse represents the schema for the GET Capella V4 API request that fetches the database credential details.
type GetDatabaseCredentialResponse struct {
	Name           string             `json:"name"`
	Password       string             `json:"password"`
	OrganizationId string             `json:"organizationId"`
	ProjectId      string             `json:"projectId"`
	ClusterId      string             `json:"clusterId"`
	Audit          CouchbaseAuditData `json:"audit"`
	Access         []Access           `json:"access"`
	Id             uuid.UUID          `json:"id"`
}

// PutDatabaseCredentialRequest represents the schema for the PUT Capella V4 API request that updates an existing database credential.
type PutDatabaseCredentialRequest struct {
	// Password is an optional field, if not passed, the existing password is not updated.
	Password string `json:"password,omitempty"`
	// Access describes the access information of the database credential.
	Access []Access `json:"access"`
}
