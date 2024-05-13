package api

import "time"

// CreateClusterAuditLogExportRequest is the payload sent to the Capella V4 Public API when creating export job request.
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateClusterAuditLogExportRequest struct {
	// End Specifies the audit log's end date and time.
	End time.Time `json:"end"`

	// Start Specifies the audit log's start date and time.
	Start time.Time `json:"start"`
}

// CreateClusterAuditLogExportResponse is the response received from the Capella V4 Public API after creating export job request.
type CreateClusterAuditLogExportResponse struct {
	// ExportId The export ID of the export job.
	ExportId string `json:"exportId"`
}

// GetClusterAuditLogExportResponse is the response received from the Capella V4 Public API when asked to fetch details of an export job.
type GetClusterAuditLogExportResponse struct {
	// AuditLogDownloadURL Pre-signed URL to download cluster audit logs.
	AuditLogDownloadURL *string `json:"auditLogDownloadURL,omitempty"`

	// AuditLogExportId The export ID of the audit log export job.
	AuditLogExportId string `json:"auditLogExportId"`

	// CreatedAt The timestamp when the audit logs were exported.
	CreatedAt time.Time `json:"createdAt"`

	// End The timestamp of when audit logs should end.
	End time.Time `json:"end"`

	// Expiration The timestamp when the download link expires.
	Expiration *time.Time `json:"expiration,omitempty"`

	// Start The timestamp of when audit logs should start.
	Start time.Time `json:"start"`

	// Status Indicates status of audit log creation. When status is complete, the compressed file can be manually downloaded.
	Status string `json:"status"`
}
