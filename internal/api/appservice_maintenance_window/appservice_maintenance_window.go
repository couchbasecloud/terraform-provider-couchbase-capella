package appservice_maintenance_window

// MaintenanceWindow is the preferred weekly maintenance window of an App Service, in UTC. The
// Capella V4 Public API uses the same shape for the create and update payloads and the GET response.
type MaintenanceWindow struct {
	// PreferredDay is the lowercase day of the week, for example "sunday".
	PreferredDay string `json:"preferredDay"`

	// PreferredHour is the hour of the day, from 0 to 23.
	PreferredHour int64 `json:"preferredHour"`

	// PreferredMinute is the minute of the hour, from 0 to 59.
	PreferredMinute int64 `json:"preferredMinute"`
}
