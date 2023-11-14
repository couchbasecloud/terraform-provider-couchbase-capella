package backup

import "slices"

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

// State is the state that a backup can have based on the fact if backup job has run successfully or not.
type State string

// IsFinalState checks whether backup job has successfully run.
func IsFinalState(state State) bool {
	// Returns True if the state is critical, False otherwise.
	finalStates := []State{
		Ready,
		Failed,
	}
	return slices.Contains(finalStates, state)
}
