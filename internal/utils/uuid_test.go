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

func TestParseHierarchyUUIDs(t *testing.T) {
	newID := func() string { return uuid.New().String() }

	t.Run("all valid IDs parse successfully", func(t *testing.T) {
		orgID, projID, clusterID, appSvcID := newID(), newID(), newID(), newID()
		orgUUID, projUUID, clusterUUID, appSvcUUID, err := ParseHierarchyUUIDs(orgID, projID, clusterID, appSvcID)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if orgUUID.String() != orgID {
			t.Errorf("organization_id mismatch: want %s, got %s", orgID, orgUUID)
		}
		if projUUID.String() != projID {
			t.Errorf("project_id mismatch: want %s, got %s", projID, projUUID)
		}
		if clusterUUID.String() != clusterID {
			t.Errorf("cluster_id mismatch: want %s, got %s", clusterID, clusterUUID)
		}
		if appSvcUUID.String() != appSvcID {
			t.Errorf("app_service_id mismatch: want %s, got %s", appSvcID, appSvcUUID)
		}
	})

	validID := newID()
	invalidCases := []struct {
		name           string
		orgID          string
		projID         string
		clusterID      string
		appSvcID       string
		wantFieldInErr string
	}{
		{"invalid organization_id", "bad", validID, validID, validID, "organization_id"},
		{"invalid project_id", validID, "bad", validID, validID, "project_id"},
		{"invalid cluster_id", validID, validID, "bad", validID, "cluster_id"},
		{"invalid app_service_id", validID, validID, validID, "bad", "app_service_id"},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, err := ParseHierarchyUUIDs(tc.orgID, tc.projID, tc.clusterID, tc.appSvcID)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantFieldInErr) {
				t.Errorf("expected error to contain %q, got: %v", tc.wantFieldInErr, err)
			}
		})
	}
}
