package schema

// CommonDescriptions contains standard descriptions for common fields used across resources.
// These are typically path parameters, headers, or standard Terraform metadata that don't
// appear in OpenAPI request/response bodies.
var CommonDescriptions = map[string]string{
	// ID fields (path parameters)
	"organization_id": "The GUID4 ID of the organization.",
	"project_id":      "The GUID4 ID of the project.",
	"cluster_id":      "The GUID4 ID of the cluster.",
	"bucket_id":       "The GUID4 ID of the bucket.",
	"app_service_id":  "The GUID4 ID of the app service.",
	"allowlist_id":    "The GUID4 ID of the allowlist entry.",
	"user_id":         "The GUID4 ID of the user.",
	"id":              "The GUID4 ID of this resource.",

	// HTTP headers
	"if_match": "A precondition header that specifies the entity tag of a resource. Used for optimistic concurrency control to prevent concurrent modifications.",
	"etag":     "The ETag header value returned by the server, used for optimistic concurrency control.",

	// Standard metadata
	"audit": "Audit metadata tracking when and by whom the resource was created and last modified.",
}
