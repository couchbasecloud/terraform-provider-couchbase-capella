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
