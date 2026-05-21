package cluster

// UpdateDeletionProtectionRequest is the request payload sent to the Capella V4 Public API
// to enable or disable deletion protection on a cluster.
type UpdateDeletionProtectionRequest struct {
	DeletionProtection bool `json:"deletionProtection"`
}
