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
