package appservice_self_service_upgrade

import "time"

// The get and list self-service upgrade endpoints return the App Service maintenance shape, so
// their responses decode into appservice_maintenance.GetMaintenanceResponse.

// CreateRequest upgrades an App Service now, or schedules the upgrade at QueueAt.
type CreateRequest struct {
	// QueueAt is when the upgrade is queued. Nil omits the field, which upgrades the App Service now.
	QueueAt *time.Time `json:"queueAt,omitempty"`

	// VersionId is the ID of an available version to upgrade to.
	VersionId string `json:"versionId"`
}

// CreateResponse is returned with 202 Accepted, before the upgrade record exists.
type CreateResponse struct {
	// Id is the upgrade ID, used in the get, update and cancel paths.
	Id string `json:"id"`
}

// UpdateRequest replaces the version and queue time of a pending upgrade.
type UpdateRequest struct {
	QueueAt   time.Time `json:"queueAt"`
	VersionId string    `json:"versionId"`
}

// AvailableVersion is a version the App Service can be upgraded to.
type AvailableVersion struct {
	// Id is the version ID to pass as versionId.
	Id string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	// Priority is a free-form string, for example "critical".
	Priority string `json:"priority"`

	ReleaseNotes []string `json:"releaseNotes"`
}

// ListAvailableVersionsResponse is the available versions list. It is not paginated.
type ListAvailableVersionsResponse struct {
	Data []AvailableVersion `json:"data"`
}
