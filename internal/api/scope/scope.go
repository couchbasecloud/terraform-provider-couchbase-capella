package scope

// CreateScopeRequest is the payload passed to V4 Capella Public API to create a scope in a bucket.
// Creates a new scope in a bucket.
//
// To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateScopeRequest struct {
	// Name The name of the scope. The name should adhere to the following rules:
	//
	// 1. The name must be between 1 and 251 characters in length.
	// 2. The name can contain only the characters A-Z, a-z, 0-9, and the symbols _, -, and %.
	// 3. The name cannot start with _ or %.
	// Note that scope and collection names are case-sensitive.
	Name string `json:"name"`
}

// GetScopeResponse is the response received from Capella V4 Public API on requesting to information about an existing scope.
//Fetches the details of the given scope.
//
//To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
//In order to access this endpoint, the provided API key must have at least one of the following roles:
//- Organization Owner
//- Project Owner
//- Project Manager
//- Project Viewer
//- Database Data Reader/Writer
//- Database Data Reader
//
//To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).

type GetScopeResponse struct {
	Collections *[]Collection `json:"collections,omitempty"`

	// Name is the name of the scope.
	Name *string `json:"name,omitempty"`
}

// GetScopesResponse is the response received from the Capella V4 Public API when asked to list all scopes.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetScopesResponse struct {
	Scopes []GetScopeResponse `json:"scopes"`
}
