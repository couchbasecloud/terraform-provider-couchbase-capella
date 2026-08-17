package cluster

import (
	"slices"
	"strings"
)

const (
	Degraded         State = "degraded"
	Deploying        State = "deploying"
	DeploymentFailed State = "deploymentFailed"
	DestroyFailed    State = "destroyFailed"
	Destroying       State = "destroying"
	Draft            State = "draft"
	Healthy          State = "healthy"
	Offline          State = "offline"
	Peering          State = "peering"
	PeeringFailed    State = "peeringFailed"
	RebalanceFailed  State = "rebalanceFailed"
	Rebalancing      State = "rebalancing"
	ScaleFailed      State = "scaleFailed"
	Scaling          State = "scaling"
	TurnedOff        State = "turnedOff"
	TurningOff       State = "turningOff"
	TurningOffFailed State = "turningOffFailed"
	TurningOn        State = "turningOn"
	TurningOnFailed  State = "turningOnFailed"
	UpgradeFailed    State = "upgradeFailed"
	Upgrading        State = "upgrading"
)

// State is the state that a cluster can have based on the fact if deployment of the cluster was successful or not.
type State string

func (s1 State) Equal(s2 State) bool {
	return strings.EqualFold(string(s1), string(s2)) // case in-sensitive equality
}

// IsTerminalState checks whether cluster is in a final non-healthy state.
func IsTerminalState(state State) bool {
	finalStates := []State{
		Degraded,
		DeploymentFailed,
		DestroyFailed,
		PeeringFailed,
		RebalanceFailed,
		ScaleFailed,
		UpgradeFailed,
	}
	return slices.Contains(finalStates, state)
}
