package utils

import (
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

	t.Run("invalid UUID returns descriptive error", func(t *testing.T) {
		_, err := ParseUUID("cluster_id", "not-a-uuid")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid cluster_id") {
			t.Errorf("expected error to contain 'invalid cluster_id', got: %v", err)
		}
	})

	t.Run("empty string returns descriptive error", func(t *testing.T) {
		_, err := ParseUUID("project_id", "")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid project_id") {
			t.Errorf("expected error to contain 'invalid project_id', got: %v", err)
		}
	})

	t.Run("field name is included in error message", func(t *testing.T) {
		fieldName := "app_service_id"
		_, err := ParseUUID(fieldName, "bad-value")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), fieldName) {
			t.Errorf("expected error to contain field name %q, got: %v", fieldName, err)
		}
	})
}
