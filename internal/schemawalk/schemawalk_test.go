package schemawalk

import "testing"

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "name"},
		{"description", "description"},
		{"memory_allocation_in_mb", "memoryAllocationInMb"},
		{"storage_backend", "storageBackend"},
		{"bucket_conflict_resolution", "bucketConflictResolution"},
		{"time_to_live_in_seconds", "timeToLiveInSeconds"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := SnakeToCamel(tt.input); got != tt.expected {
				t.Errorf("SnakeToCamel(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"foo", "Foo"},
		{"Foo", "Foo"},
		{"fOO", "FOO"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Capitalize(tt.input); got != tt.expected {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSchemaNameCandidates(t *testing.T) {
	t.Run("lowercase name produces capitalized patterns", func(t *testing.T) {
		got := SchemaNameCandidates("replication")
		want := []string{
			"Replication",
			"CreateReplicationRequest",
			"GetReplicationResponse",
			"UpdateReplicationRequest",
			"ReplicationRequest",
			"ReplicationResponse",
			"replication",
		}
		if len(got) != len(want) {
			t.Fatalf("len=%d want %d: %v", len(got), len(want), got)
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("[%d] got %q want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("already-capitalized name deduplicates last entry", func(t *testing.T) {
		got := SchemaNameCandidates("AllowedCidr")
		if len(got) != 6 {
			t.Fatalf("expected 6 candidates (no duplicate), got %d: %v", len(got), got)
		}
		if got[0] != "AllowedCidr" {
			t.Errorf("first candidate should be AllowedCidr, got %q", got[0])
		}
	})
}

func TestEnumValues(t *testing.T) {
	t.Run("finds enum by openAPISchemaName direct match", func(t *testing.T) {
		vals := EnumValues("AllowedCidr", "allowlist", nil, "status")
		if len(vals) == 0 {
			t.Fatal("expected enum values for AllowedCidr.status")
		}
		if vals[0] != "active" {
			t.Errorf("expected first value 'active', got %q", vals[0])
		}
	})

	t.Run("finds enum via Create...Request pattern", func(t *testing.T) {
		vals := EnumValues("replication", "replication", nil, "direction")
		if len(vals) == 0 {
			t.Fatal("expected enum values for CreateReplicationRequest.direction")
		}
	})

	t.Run("finds enum via alternateSchemas", func(t *testing.T) {
		vals := EnumValues("unknown", "unknown", []string{"Days"}, "day")
		if len(vals) == 0 {
			t.Fatal("expected enum values for Days.day via alternateSchemas")
		}
		if vals[0] != "monday" {
			t.Errorf("expected first value 'monday', got %q", vals[0])
		}
	})

	t.Run("alternateSchema takes priority over pattern match", func(t *testing.T) {
		vals := EnumValues("SomeOtherSchema", "other", []string{"Days"}, "state")
		if len(vals) == 0 {
			t.Fatal("expected enum values for Days.state via alternateSchemas")
		}
		if vals[0] != "on" {
			t.Errorf("expected first value 'on', got %q", vals[0])
		}
	})

	t.Run("returns nil for unknown schema/field", func(t *testing.T) {
		if vals := EnumValues("NonExistentSchema", "nonExistent", nil, "nonExistentField"); vals != nil {
			t.Errorf("expected nil, got %v", vals)
		}
	})

	t.Run("returns nil for empty inputs", func(t *testing.T) {
		if vals := EnumValues("", "", nil, "status"); vals != nil {
			t.Errorf("expected nil for empty schema names, got %v", vals)
		}
	})
}
