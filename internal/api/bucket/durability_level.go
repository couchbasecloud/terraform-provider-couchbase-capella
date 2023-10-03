package bucket

const (
	// NoneDurability represents no durability option
	NoneDurability DurabilityLevel = "none"

	// Majority is the level in which the mutation must be replicated to a
	// majority of the Data Service nodes.
	Majority DurabilityLevel = "majority"

	// MajorityPersistActive is the level in which mutation must be replicated
	// to a majority of the Data Service nodes. It also must be persisted on the
	// node hosting the active vBucket for the data.
	MajorityPersistActive DurabilityLevel = "majorityAndPersistActive"

	// PersistMajority is the level in which the mutation must be persisted to
	// a majority of the Data Service Nodes.
	PersistMajority DurabilityLevel = "persistToMajority"
)

type DurabilityLevel string
