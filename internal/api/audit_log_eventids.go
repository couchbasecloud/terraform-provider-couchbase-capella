package api

// GetAuditLogEventIDsResponse is the response received from the Capella V4 Public API when retrieving audit log event ids.
//
// Retrieves a list of audit event IDs. The list of filterable event IDs can be specified while configuring audit log for cluster.
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

type GetAuditLogEventIDsResponse struct {
	Events []AuditFilterableEvent `json:"events"`
}

type AuditFilterableEvent struct {
	Description string `json:"description"`
	Id          int32  `json:"id"`
	Module      string `json:"module"`
	Name        string `json:"name"`
}
