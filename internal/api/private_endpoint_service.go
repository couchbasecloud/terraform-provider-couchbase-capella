package api

// GetPrivateEndpointServiceStatusResponse is the response received from the Capella V4 Public API
// when getting private endpoint service status.
type GetPrivateEndpointServiceStatusResponse struct {
	Enabled bool `json:"enabled"`
}
