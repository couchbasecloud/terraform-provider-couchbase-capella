package utils

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUUID(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		fieldName   string
		value       string
		expected    *uuid.UUID
		expectError bool
	}{
		{
			name:      "valid UUID is parsed successfully",
			fieldName: "organization_id",
			value:     validID.String(),
			expected:  &validID,
		},
		{
			name:        "invalid UUID returns descriptive error",
			fieldName:   "cluster_id",
			value:       "not-a-uuid",
			expectError: true,
		},
		{
			name:        "empty string returns descriptive error",
			fieldName:   "project_id",
			value:       "",
			expectError: true,
		},
		{
			name:        "field name is included in error message",
			fieldName:   "app_service_id",
			value:       "bad-value",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseUUID(test.fieldName, test.value)
			if test.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, "invalid "+test.fieldName)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, *test.expected, got)
		})
	}
}

// generateIDFields creates n IDField entries with random valid UUIDs
func generateIDFields(n int) []IDField {
	fields := make([]IDField, n)
	for i := range fields {
		fields[i] = IDField{Name: fmt.Sprintf("field_%d", i), Value: uuid.New().String()}
	}
	return fields
}

func TestParseUUIDs(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name           string
		fields         []IDField
		expectError    bool
		wantFieldInErr string
	}{
		{
			name:   "zero fields returns empty slice",
			fields: []IDField{},
		},
		{
			name:   "single field",
			fields: generateIDFields(1),
		},
		{
			name:   "three fields preserve order",
			fields: generateIDFields(3),
		},
		{
			name:   "five fields preserve order",
			fields: generateIDFields(5),
		},
		{
			name:           "invalid first field stops early",
			fields:         []IDField{{"field_a", "bad"}, {"field_b", validID}, {"field_c", validID}},
			expectError:    true,
			wantFieldInErr: "field_a",
		},
		{
			name:           "invalid middle field returns correct field name",
			fields:         []IDField{{"field_a", validID}, {"field_b", "bad"}, {"field_c", validID}},
			expectError:    true,
			wantFieldInErr: "field_b",
		},
		{
			name:           "invalid last field returns correct field name",
			fields:         []IDField{{"field_a", validID}, {"field_b", validID}, {"field_c", "bad"}},
			expectError:    true,
			wantFieldInErr: "field_c",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uuids, err := ParseUUIDs(test.fields...)
			if test.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantFieldInErr)
				return
			}

			require.NoError(t, err)
			require.Len(t, uuids, len(test.fields))
			for i, f := range test.fields {
				assert.Equal(t, f.Value, uuids[i].String())
			}
		})
	}
}
