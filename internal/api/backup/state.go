package backup

const (
	// Pending communicates that a backup record has been created but the
	// backup job has not yet run. This happens when dp-backup first
	// communicates the intent to start a backup to the internal api, which
	// results in a backup record created in pending status.
	Pending State = "pending"

	// Ready represents a backup record that has had its backup job run successfully.
	Ready State = "ready"

	// Failed represents a backup record that has had its backup job fail.
	Failed State = "failed"
)

type State string

// IsFinalState checks whether backup job has successfully run.
func IsFinalState(state State) bool {
	//"""Returns True if the state is critical, False otherwise."""
	finalStates := []State{
		Ready,
		Failed,
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
