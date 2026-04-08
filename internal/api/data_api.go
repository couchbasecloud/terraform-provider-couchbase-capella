package api

// GetDataAPIStatusResponse is the API response received from the Capella V4 Public API
// while attempting to fetch the Data API status for a cluster.
type GetDataAPIStatusResponse struct {
	Enabled                  bool   `json:"enabled"`
	State                    string `json:"state"`
	EnabledForNetworkPeering bool   `json:"enabledForNetworkPeering"`
	StateForNetworkPeering   string `json:"stateForNetworkPeering"`
	ConnectionString         string `json:"connectionString"`
}

// UpdateDataAPIRequest is the request payload sent to the Capella V4 Public API
// in order to update the Data API configuration for a cluster.
type UpdateDataAPIRequest struct {
	EnableDataApi        bool `json:"enableDataApi"`
	EnableNetworkPeering bool `json:"enableNetworkPeering"`
}
