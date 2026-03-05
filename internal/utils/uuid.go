package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID parses a string into a UUID, returning a descriptive error if parsing fails.
// The fieldName parameter is used to construct a human-readable error message indicating
// which field contained the invalid UUID value.
func ParseUUID(fieldName, value string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return parsed, nil
}

// IDField pairs a human-readable field name with its raw string ID value,
// used as input to ParseUUIDs.
type IDField struct {
	Name  string
	Value string
}

// ParseUUIDs parses an ordered list of IDField pairs into UUID values.
// Results are returned in the same order as the inputs.
// It returns the first parse error encountered, with the field name included
// in the error message for easier debugging.
func ParseUUIDs(fields ...IDField) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(fields))
	for i, f := range fields {
		parsed, err := ParseUUID(f.Name, f.Value)
		if err != nil {
			return nil, err
		}
		uuids[i] = parsed
	}
	return uuids, nil
}
