package api

// GetPrivateEndpointServiceStatusResponse is the response received from the Capella V4 Public API
// when getting private endpoint service status.
type GetPrivateEndpointServiceStatusResponse struct {
	Enabled bool `json:"enabled"`

	// Status is the lifecycle state of the private endpoint service derived from
	// the most recent enable/disable/update operation (for example "enableFailed"
	// or "enabling"). It is optional and best-effort: older control planes omit it,
	// in which case callers fall back to the Enabled boolean.
	Status *string `json:"status,omitempty"`

	PrivateDns string `json:"privateDns"`
}
