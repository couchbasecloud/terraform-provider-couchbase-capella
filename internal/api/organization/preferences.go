package organization

type Preferences struct {
	// SessionDuration: Maximum allowed time in seconds inside the tenant for a user.
	SessionDuration *int32 `json:"sessionDuration"`
}
