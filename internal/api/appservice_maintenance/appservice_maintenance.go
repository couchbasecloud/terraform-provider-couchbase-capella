package appservice_maintenance

import (
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// GetMaintenanceResponse is a customer-visible App Service maintenance record, as returned by the
// Capella V4 Public API list and get maintenance endpoints.
type GetMaintenanceResponse struct {
	// QueueAt is when the maintenance is queued to start.
	QueueAt *time.Time `json:"queueAt"`

	// RescheduleBy is the latest time the maintenance can be queued at: the end of its scheduling window.
	RescheduleBy *time.Time `json:"rescheduleBy"`

	// CompletedAt is when the maintenance finished; nil until then.
	CompletedAt *time.Time `json:"completedAt"`

	// VersionId is the release the maintenance moves the App Service to; nil when no release
	// contains the record's image.
	VersionId *string `json:"versionId"`

	// Id is the maintenance ID, used in the get and reschedule paths.
	Id string `json:"id"`

	// Type is the kind of maintenance, for example "upgrade".
	Type string `json:"type"`

	// InitiatedBy is who created the record, for example "couchbase" or "selfService".
	InitiatedBy string `json:"initiatedBy"`

	// CurrentState is the state of the maintenance, for example "pending".
	CurrentState string `json:"currentState"`

	Audit api.CouchbaseAuditData `json:"audit"`
}

// RescheduleMaintenanceRequest moves a pending maintenance to a new queue time.
type RescheduleMaintenanceRequest struct {
	QueueAt time.Time `json:"queueAt"`
}
