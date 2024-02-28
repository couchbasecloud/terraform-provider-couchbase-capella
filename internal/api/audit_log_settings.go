package api

// CreateClusterAuditSettingsRequest is the payload sent to the Capella V4 Public API when configuring audit log settings.
//
//	In order to access this endpoint, the provided API key must have at least one of the following roles:
//	- Organization Owner
//	- Project Owner
//	- Project Manager
//	To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateClusterAuditSettingsRequest struct {
	// AuditEnabled Determines whether audit logging is enabled or not on the cluster.
	AuditEnabled bool `json:"auditEnabled"`

	// DisabledUsers List of users whose filterable events will not be logged.
	DisabledUsers AuditSettingsDisabledUsers `json:"disabledUsers"`

	// EnabledEventIDs List of enabled filterable audit events for the cluster.
	EnabledEventIDs []int32 `json:"enabledEventIDs"`
}

// GetClusterAuditSettingsResponse is the response received from the Capella V4 Public API when retrieving cluster audit log settings.
type GetClusterAuditSettingsResponse struct {
	// AuditEnabled Determines whether audit logging is enabled or not on the cluster.
	AuditEnabled bool `json:"auditEnabled"`

	// DisabledUsers List of users whose filterable events will not be logged.
	DisabledUsers AuditSettingsDisabledUsers `json:"disabledUsers"`

	// EnabledEventIDs List of enabled filterable audit events for the cluster.
	EnabledEventIDs []int32 `json:"enabledEventIDs"`
}

// AuditSettingsDisabledUsers List of users whose filterable events will not be logged.
type AuditSettingsDisabledUsers = []AuditSettingsDisabledUser

type AuditSettingsDisabledUser struct {
	// Domain Specifies whether the user is local or external.
	Domain *string `json:"domain,omitempty"`

	// Name The user name.
	Name *string `json:"name,omitempty"`
}
