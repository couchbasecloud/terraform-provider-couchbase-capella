package api

// UpdateClusterAuditSettingsRequest is the payload sent to the Capella V4 Public API when configuring audit log settings.
//
//	In order to access this endpoint, the provided API key must have at least one of the following roles:
//	- Organization Owner
//	- Project Owner
//	- Project Manager
//	To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type UpdateClusterAuditSettingsRequest struct {
	DisabledUsers   AuditSettingsDisabledUsers `json:"disabledUsers"`
	EnabledEventIDs []int32                    `json:"enabledEventIDs"`
	AuditEnabled    bool                       `json:"auditEnabled"`
}

// GetClusterAuditSettingsResponse is the response received from the Capella V4 Public API when retrieving cluster audit log settings.
type GetClusterAuditSettingsResponse struct {
	DisabledUsers   AuditSettingsDisabledUsers `json:"disabledUsers"`
	EnabledEventIDs []int32                    `json:"enabledEventIDs"`
	AuditEnabled    bool                       `json:"auditEnabled"`
}

// AuditSettingsDisabledUsers List of users whose filterable events will not be logged.
type AuditSettingsDisabledUsers = []AuditSettingsDisabledUser

type AuditSettingsDisabledUser struct {
	// Domain Specifies whether the user is local or external.
	Domain *string `json:"domain,omitempty"`

	// Name The user name.
	Name *string `json:"name,omitempty"`
}
