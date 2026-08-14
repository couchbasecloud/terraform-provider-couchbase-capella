package data_api

type State string

const (
	Enabled     State = "enabled"
	Disabled    State = "disabled"
	Enabling    State = "enabling"
	Disabling   State = "disabling"
	Configuring State = "configuring"
)

// DesiredState maps the requested enable/disable flag to the final state the resource should reach.
func DesiredState(enabled bool) State {
	if enabled {
		return Enabled
	}
	return Disabled
}
