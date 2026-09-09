package api

// GetAppServicePrivateEndpointServiceStatusResponse is the response received from the
// Capella V4 Public API when getting the private endpoint service status for an App
// Service (Sync Gateway).
type GetAppServicePrivateEndpointServiceStatusResponse struct {
	// State is the lifecycle state of the private endpoint service derived from the
	// most recent enable/disable operation (for example "enabling" or "failed"). It is
	// optional and best-effort.
	State *string `json:"state,omitempty"`

	// TargetState is the desired end state of the private endpoint service (enabled or
	// disabled). When State is "failed", TargetState indicates which operation
	// (enable or disable) failed. It is optional and best-effort.
	TargetState *string `json:"targetState,omitempty"`
}
