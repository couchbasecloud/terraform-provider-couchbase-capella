package cluster

// ClusterActivationRequest is the request payload sent to the Capella V4 Public API when asked to turn a cluster on or off.
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner
type ClusterActivationRequest struct {
	// TurnOnLinkedAppService if true will turn on the app service linked with the cluster,
	// and if false, will keep it turned off.
	// Default value for this is false.
	TurnOnLinkedAppService bool `json:"turnOnLinkedAppService"`
}
