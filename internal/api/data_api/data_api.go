package data_api

// UpdateDataApiRequest is the payload sent to the Capella V4 Public API when enabling or disabling
// the Data API and network peering on a cluster.
type UpdateDataApiRequest struct {
	// EnableDataApi enables or disables the Data API for the cluster.
	EnableDataApi bool `json:"enableDataApi"`

	// EnableNetworkPeering enables or disables network peering when the Data API is enabled.
	EnableNetworkPeering *bool `json:"enableNetworkPeering"`
}

// GetDataApiStatusResponse is the response received from the Capella V4 Public API when retrieving
// the Data API and network peering status of a cluster.
type GetDataApiStatusResponse struct {
	// Enabled indicates whether the Data API is enabled or disabled on the cluster.
	Enabled bool `json:"enabled"`

	// State is the current status of the Data API.
	State State `json:"state"`

	// EnabledForNetworkPeering indicates whether network peering was enabled or disabled for the Data API.
	EnabledForNetworkPeering bool `json:"enabledForNetworkPeering"`

	// StateForNetworkPeering is the current status of network peering for the Data API.
	StateForNetworkPeering State `json:"stateForNetworkPeering"`

	// ConnectionString is the connection string for the Data API service.
	// If the Data API is not enabled, this is an empty string.
	ConnectionString string `json:"connectionString"`
}
