package appservice_maintenance

import (
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// InitiatedBySelfService marks a maintenance record created by a self-service upgrade. Its queue
// time is owned by the self-service upgrade resource, not by the maintenance schedule resource.
const InitiatedBySelfService = "selfService"

// GetMaintenanceResponse is a customer-visible App Service maintenance record, as returned by the
// Capella V4 Public API list and get maintenance endpoints.
type GetMaintenanceResponse struct {
	// QueueAt is when the maintenance is queued to start.
	QueueAt *time.Time `json:"queueAt"`

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
