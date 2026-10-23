package cluster

// ExpressScalingState is the express scaling status reported on a cluster read.
type ExpressScalingState string

const (
	ExpressScalingEnabling  ExpressScalingState = "enabling"
	ExpressScalingEnabled   ExpressScalingState = "enabled"
	ExpressScalingDisabling ExpressScalingState = "disabling"
	ExpressScalingDisabled  ExpressScalingState = "disabled"
)

// IsEnabled reports whether the state is enabled or heading towards enabled.
// A cluster that reports no state has never had express scaling enabled.
func (s ExpressScalingState) IsEnabled() bool {
	return s == ExpressScalingEnabling || s == ExpressScalingEnabled
}

// IsSettled reports whether the state is a final enabled or disabled state for the
// requested setting. A cluster that reports no state counts as disabled.
func (s ExpressScalingState) IsSettled(enabled bool) bool {
	if enabled {
		return s == ExpressScalingEnabled
	}
	return s == ExpressScalingDisabled || s == ""
}

// IsMovingTo reports whether a change towards the requested setting is in progress.
func (s ExpressScalingState) IsMovingTo(enabled bool) bool {
	if enabled {
		return s == ExpressScalingEnabling
	}
	return s == ExpressScalingDisabling
}

// UpdateExpressScalingRequest is the request payload sent to the Capella V4 Public API
// to enable or disable express scaling on a cluster. The API rejects a missing or null
// value, so the field is always sent.
type UpdateExpressScalingRequest struct {
	ExpressScaling bool `json:"expressScaling"`
}
