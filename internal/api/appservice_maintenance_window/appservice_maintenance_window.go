package appservice_maintenance_window

// CreateMaintenanceWindowRequest is the payload sent to the Capella V4 Public API to set
// the preferred weekly maintenance window of an App Service.
type CreateMaintenanceWindowRequest struct {
	// PreferredDay is the lowercase day of the week, for example "sunday".
	PreferredDay string `json:"preferredDay"`

	// PreferredHour is the hour of the day, from 0 to 23.
	PreferredHour int64 `json:"preferredHour"`

	// PreferredMinute is the minute of the hour, from 0 to 59.
	PreferredMinute int64 `json:"preferredMinute"`
}

// UpdateMaintenanceWindowRequest is the payload sent to the Capella V4 Public API to replace
// the preferred weekly maintenance window of an App Service. All fields are always sent.
type UpdateMaintenanceWindowRequest struct {
	PreferredDay    string `json:"preferredDay"`
	PreferredHour   int64  `json:"preferredHour"`
	PreferredMinute int64  `json:"preferredMinute"`
}

// GetMaintenanceWindowResponse is the response received from the Capella V4 Public API when
// reading the preferred weekly maintenance window of an App Service.
type GetMaintenanceWindowResponse struct {
	PreferredDay    string `json:"preferredDay"`
	PreferredHour   int64  `json:"preferredHour"`
	PreferredMinute int64  `json:"preferredMinute"`
}
