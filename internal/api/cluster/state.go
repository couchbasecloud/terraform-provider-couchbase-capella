package cluster

// Defines values for State.
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

// State defines model for State.
type State string

// IsFinalState checks whether cluster is successfully deployed/updated or not while creation/updation
//TODO: Degraded, draft, peeringFailed, turningOffFailed, and turningOnFailed are not known when it occurs and What happens if rebalancing fails? Will it retry?"

func IsFinalState(state State) bool {
	//"""Returns True if the state is critical, False otherwise."""
	finalStates := []State{
		Healthy,
		Degraded,
		DeploymentFailed,
		DestroyFailed,
		PeeringFailed,
		RebalanceFailed,
		ScaleFailed,
		UpgradeFailed,
	}
	return Contains(finalStates, state)
}