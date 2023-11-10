package appservice

import "slices"

const (
	// Pending communicates that the sync gateway is waiting to be
	// provisioned.
	Pending State = "pending"

	// Deploying communicates that the sync gateway is deploying. This
	// status is driven by control plane activities and is not related to
	// the act of scaling the node count up or down.
	Deploying        State = "deploying"
	DeploymentFailed State = "deploymentFailed"

	// Destroying communicates that the sync gateway is being torn down.
	// This state can only occur if a user has chosen to remove the sync
	// gateway.
	Destroying    State = "destroying"
	DestroyFailed State = "destroyFailed"

	// Healthy communicates that the sync gateway was successfully deployed
	// and is operational.
	Healthy  State = "healthy"
	Degraded State = "degraded"

	Scaling     State = "scaling"
	ScaleFailed State = "scaleFailed"

	Upgrading     State = "upgrading"
	UpgradeFailed State = "upgradeFailed"

	TurnedOff     State = "turnedOff"
	TurningOff    State = "turningOff"
	TurnOffFailed State = "turnOffFailed"
	TurningOn     State = "turningOn"
	TurnOnFailed  State = "turnOnFailed"
)

type State string

func IsFinalState(state State) bool {
	// Returns True if the state is an end state, False otherwise.
	finalStates := []State{
		Healthy,
		Degraded,
		DeploymentFailed,
		DestroyFailed,
		TurnedOff,
		TurnOffFailed,
		TurnOnFailed,
		ScaleFailed,
		UpgradeFailed,
	}
	return slices.Contains(finalStates, state)
}
