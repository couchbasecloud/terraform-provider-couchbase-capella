package data_api

type State string

const (
	Enabled     State = "enabled"
	Disabled    State = "disabled"
	Enabling    State = "enabling"
	Disabling   State = "disabling"
	Configuring State = "configuring"
)

// IsFinalState checks whether the Data API or its network peering has finished transitioning and reached a final state
func IsFinalState(state State) bool {
	switch state {
	case Enabled, Disabled:
		return true
	default:
		return false
	}
}
