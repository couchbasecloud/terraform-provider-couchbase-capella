package snapshot_backup

import "slices"

type State string

const (
	Pending  State = "pending"
	Complete State = "complete"
	Failed   State = "failed"
)

// IsFinalState checks whether snapshotbackup job has successfully run.
func IsFinalState(state State) bool {
	// Returns True if the state is critical, False otherwise.
	finalStates := []State{
		Complete,
		Failed,
	}
	return slices.Contains(finalStates, state)
}
