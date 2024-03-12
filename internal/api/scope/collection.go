package scope

// Collection is the collection sent by the Capella V4 Public API for any existing scope or collection.
type Collection struct {
	// MaxTTL Max TTL of the collection.
	MaxTTL *int64 `json:"maxTTL,omitempty"`

	// Name is the Name of the collection.
	Name *string `json:"name,omitempty"`
}
