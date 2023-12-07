package organization

// Preferences is the setting for session duration for the entire Capella Organization.
type Preferences struct {
	// SessionDuration: Maximum allowed time in seconds inside the tenant for a user.
	SessionDuration *int32 `json:"sessionDuration"`
}
