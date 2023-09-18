package utils

import (
	"terraform-provider-capella/internal/api"
)

// IsFinalState checks whether cluster is successfully deployed/updated or not while creation/updation
//TODO: Degraded, draft, peeringFailed, turningOffFailed, and turningOnFailed are not known when it occurs and What happens if rebalancing fails? Will it retry?"

func IsFinalState(state api.CurrentState) bool {
	//"""Returns True if the state is critical, False otherwise."""
	finalStates := []api.CurrentState{
		api.Healthy,
		api.Degraded,
		api.DeploymentFailed,
		api.DestroyFailed,
		api.PeeringFailed,
		api.RebalanceFailed,
		api.ScaleFailed,
		api.UpgradeFailed,
	}
	return Contains(finalStates, state)
}

// Contains checks whether passed element presents in array or not
func Contains[T comparable](s []T, e T) bool {
	for _, r := range s {
		if r == e {
			return true
		}
	}
	return false
}

// AreEqual returns true if the two arrays contain the same elements, without any extra values, False otherwise.
func AreEqual[T comparable](array1 []T, array2 []T) bool {
	set1 := make(map[T]bool)
	for _, element := range array1 {
		set1[element] = true
	}

	for _, element := range array2 {
		if !set1[element] {
			return false
		}
	}

	return len(set1) == len(array1)
}
