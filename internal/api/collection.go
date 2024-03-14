package api

// CreateCollectionRequest is the payload passed to V4 Capella Public API to create a collection in a scope.
// Creates a new collection in a scope.
//
// To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateCollectionRequest struct {
	// Name The name of the collection. The name should adhere to the following rules:
	//
	// 1. The name must be between 1 and 251 characters in length.
	//
	// 2. The name can contain only the characters A-Z, a-z, 0-9, and the symbols _, -, and %.
	//
	// 3. The name cannot start with _ or %.
	//
	// Note that scope and collection names are case-sensitive.
	Name string `json:"name"`

	// MaxTTL Specify the time to live (TTL) value in seconds. Defines the duration (Seconds) for which the documents in a collection are kept before automatic removal from the database. -  For server versions < 7.6.0, this is a non-negative value. Set to 0 to use the bucket's maxTTL value. -  For server versions >= 7.6.0, this value should be >= -1. Set to -1 to disable expiry for that collection. Set to 0 to use the bucket's maxTTL value. -  The maximum value that can be set for maxTTL is 2147483647.
	MaxTTL *int64 `json:"maxTTL,omitempty"`
}

// GetCollectionResponse is the response received from Capella V4 Public API on requesting to information about an existing collection.
// Fetches the details of the given collection.
//
// To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
// - Project Viewer
// - Database Data Reader/Writer
// - Database Data Reader
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type GetCollectionResponse struct {
	// MaxTTL is the Max TTL of the collection.
	MaxTTL *int64 `json:"maxTTL,omitempty"`

	// Name is the Name of the collection.
	Name *string `json:"name,omitempty"`
}

// GetCollectionsResponse is the response received from the Capella V4 Public API when asked to list all collections.
// Lists all the collections in a scope.
//
// To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
// - Project Viewer
// - Database Data Reader/Writer
// - Database Data Reader
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type GetCollectionsResponse struct {
	Data []GetCollectionResponse `json:"data"`
}

// UpdateCollectionRequest is the payload sent to the Capella V4 Public API when asked to update an existing collection.
// Updates an existing collection.
//
// This operation is only allowed for a cluster with server version >= 7.6.0. A collection cannot be updated for the server versions lower than this.
//
// This allows to update the maxTTL of the given collection.
//
// To learn more about scopes and collections, see [Buckets, Scopes, and Collections](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html).
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type UpdateCollectionRequest struct {
	// MaxTTL Specify the new time to live (TTL) value in seconds.
	//  -  This value should be >= -1. Set to -1 to disable expiry for that collection.
	//  -  Set to 0 to use the bucket's maxTTL value.
	//  -  The maximum value that can be set for maxTTL is 2147483647.
	MaxTTL int64 `json:"maxTTL"`
}
