package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestParseUUID(t *testing.T) {
	t.Run("valid UUID is parsed successfully", func(t *testing.T) {
		id := uuid.New()
		got, err := ParseUUID("organization_id", id.String())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if got != id {
			t.Errorf("expected %v, got %v", id, got)
		}
	})

	invalidCases := []struct {
		name      string
		fieldName string
		value     string
	}{
		{"invalid UUID returns descriptive error", "cluster_id", "not-a-uuid"},
		{"empty string returns descriptive error", "project_id", ""},
		{"field name is included in error message", "app_service_id", "bad-value"},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertInvalidUUID(t, tc.fieldName, tc.value)
		})
	}
}

// assertInvalidUUID is a test helper that verifies ParseUUID returns an error
// containing the given fieldName for an invalid value.
func assertInvalidUUID(t *testing.T, fieldName, value string) {
	t.Helper()
	_, err := ParseUUID(fieldName, value)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid "+fieldName) {
		t.Errorf("expected error to contain %q, got: %v", "invalid "+fieldName, err)
	}
}

func TestParseUUIDs(t *testing.T) {
	newID := func() string { return uuid.New().String() }

	// Verify that any number of fields can be parsed and that the output slice
	// preserves insertion order regardless of the field names used.
	orderedCases := []struct {
		name   string
		fields []IDField
	}{
		{
			"single field",
			[]IDField{{"field_a", newID()}},
		},
		{
			"three fields preserve order",
			[]IDField{{"field_a", newID()}, {"field_b", newID()}, {"field_c", newID()}},
		},
		{
			"four fields preserve order",
			[]IDField{{"field_a", newID()}, {"field_b", newID()}, {"field_c", newID()}, {"field_d", newID()}},
		},
		{
			"five fields preserve order",
			[]IDField{{"field_a", newID()}, {"field_b", newID()}, {"field_c", newID()}, {"field_d", newID()}, {"field_e", newID()}},
		},
	}
	for _, tc := range orderedCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			uuids, err := ParseUUIDs(tc.fields...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if len(uuids) != len(tc.fields) {
				t.Fatalf("expected %d UUIDs, got %d", len(tc.fields), len(uuids))
			}
			for i, f := range tc.fields {
				if uuids[i].String() != f.Value {
					t.Errorf("index %d (%s) mismatch: want %s, got %s", i, f.Name, f.Value, uuids[i])
				}
			}
		})
	}

	t.Run("zero fields returns empty slice without error", func(t *testing.T) {
		uuids, err := ParseUUIDs()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(uuids) != 0 {
			t.Errorf("expected empty slice, got %v", uuids)
		}
	})

	// Verify that a parse failure at any position is reported with the correct
	// field name and that processing stops at the first bad value.
	validID := newID()
	invalidCases := []struct {
		name           string
		fields         []IDField
		wantFieldInErr string
	}{
		{
			"invalid first field stops early",
			[]IDField{{"field_a", "bad"}, {"field_b", validID}, {"field_c", validID}},
			"field_a",
		},
		{
			"invalid middle field returns correct field name",
			[]IDField{{"field_a", validID}, {"field_b", "bad"}, {"field_c", validID}},
			"field_b",
		},
		{
			"invalid last field returns correct field name",
			[]IDField{{"field_a", validID}, {"field_b", validID}, {"field_c", "bad"}},
			"field_c",
		},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseUUIDs(tc.fields...)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantFieldInErr) {
				t.Errorf("expected error to contain %q, got: %v", tc.wantFieldInErr, err)
			}
		})
	}
}

// TestParseUUIDs_AnyArity verifies that ParseUUIDs correctly handles an arbitrary
// number of fields with arbitrary names, without assuming any fixed hierarchy.
func TestParseUUIDs_AnyArity(t *testing.T) {
	newID := func() string { return uuid.New().String() }

	arities := []int{1, 2, 3, 4, 5, 10}
	for _, n := range arities {
		n := n
		t.Run(fmt.Sprintf("%d fields all valid", n), func(t *testing.T) {
			fields := make([]IDField, n)
			for i := range fields {
				fields[i] = IDField{fmt.Sprintf("field_%d", i), newID()}
			}

			uuids, err := ParseUUIDs(fields...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if len(uuids) != n {
				t.Fatalf("expected %d UUIDs, got %d", n, len(uuids))
			}
			for i, f := range fields {
				if uuids[i].String() != f.Value {
					t.Errorf("index %d (%s) mismatch: want %s, got %s", i, f.Name, f.Value, uuids[i])
				}
			}
		})

		t.Run(fmt.Sprintf("%d fields invalid at last position", n), func(t *testing.T) {
			fields := make([]IDField, n)
			for i := range fields {
				if i == n-1 {
					fields[i] = IDField{fmt.Sprintf("field_%d", i), "not-a-uuid"}
				} else {
					fields[i] = IDField{fmt.Sprintf("field_%d", i), newID()}
				}
			}

			_, err := ParseUUIDs(fields...)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			wantField := fmt.Sprintf("field_%d", n-1)
			if !strings.Contains(err.Error(), wantField) {
				t.Errorf("expected error to contain %q, got: %v", wantField, err)
			}
		})
	}
}
