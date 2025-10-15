package schema

// CommonDescriptions contains standard descriptions for common fields used across resources.
// These are for fields that don't appear in the OpenAPI spec:
// - HTTP headers (if_match, etag)
// - Special nested attributes (audit)
//
// Note: Path parameter IDs (organization_id, project_id, etc.) are now
// automatically loaded from the OpenAPI spec's components.parameters section.
var CommonDescriptions = map[string]string{
	// HTTP headers
	"if_match": "A precondition header that specifies the entity tag of a resource. Used for optimistic concurrency control to prevent concurrent modifications.",
	"etag":     "The ETag header value returned by the server, used for optimistic concurrency control.",

	// Standard metadata (nested attributes)
	"audit": "Couchbase audit data.",
}
