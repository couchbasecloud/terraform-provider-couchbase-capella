package network_peer

// PeeringStatus communicates the state of the VPC peering relationship. It is the state and reasoning for VPC peer.
type PeeringStatus struct {
	Reasoning *string `json:"reasoning,omitempty"`
	State     *string `json:"state,omitempty"`
}
