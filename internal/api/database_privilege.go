package api

// GetDatabasePrivilegeResponse represents a single privilege returned by the
// GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/privileges endpoint.
type GetDatabasePrivilegeResponse struct {
	// Name is the name of the Capella privilege (e.g. dataRead, queryIndex).
	Name string `json:"name"`

	// Group is the category the privilege belongs to (e.g. Data, Query, Analytics, FTS, Eventing, Global).
	Group string `json:"group"`

	// Resources is the RBAC template indicating which resource levels (bucket/scope/collection) are configurable.
	// A nil value means the privilege is global.
	Resources *AccessibleResources `json:"resources,omitempty"`
}
