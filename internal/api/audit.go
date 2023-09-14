package api

import "time"

// CouchbaseAuditData contains all audit-related fields.
type CouchbaseAuditData struct {
	// CreatedAt The RFC3339 timestamp associated with when the resource was initially
	// created.
	CreatedAt time.Time `json:"createdAt"`

	// CreatedBy The user who created the resource; this will be a UUID4 ID for standard
	// users and will be a string such as "internal-support" for internal
	// Couchbase support users.
	CreatedBy string `json:"createdBy"`

	// ModifiedAt The RFC3339 timestamp associated with when the resource was last modified.
	ModifiedAt time.Time `json:"modifiedAt"`

	// ModifiedBy The user who last modified the resource; this will be a UUID4 ID for
	// standard users and wilmal be a string such as "internal-support" for
	// internal Couchbase support users.
	ModifiedBy string `json:"modifiedBy"`

	// Version The version of the document. This value is incremented each time the
	// resource is modified.
	Version int `json:"version"`
}
