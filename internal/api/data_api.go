package api

// UpdateDataAPIRequest is the request body for the PUT /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/dataAPI endpoint.
type UpdateDataAPIRequest struct {
	EnableDataApi        bool `json:"enableDataApi"`
	EnableNetworkPeering bool `json:"enableNetworkPeering"`
}

// GetDataAPIStatusResponse is the response from the GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/dataAPI endpoint.
type GetDataAPIStatusResponse struct {
	Enabled                  bool   `json:"enabled"`
	State                    string `json:"state"`
	EnabledForNetworkPeering bool   `json:"enabledForNetworkPeering"`
	StateForNetworkPeering   string `json:"stateForNetworkPeering"`
	ConnectionString         string `json:"connectionString"`
}

// IsDataAPIFinalState returns true if the Data API state is a final (non-transitional) state.
func IsDataAPIFinalState(state string) bool {
	switch state {
	case "enabled", "disabled":
		return true
	default:
		return false
	}
}
