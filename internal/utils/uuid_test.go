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

func TestParseUUIDs(t *testing.T) {
	newID := func() string { return uuid.New().String() }

	t.Run("three IDs (no app_service_id) parse successfully", func(t *testing.T) {
		orgID, projID, clusterID := newID(), newID(), newID()
		uuids, err := ParseUUIDs(
			IDField{"organization_id", orgID},
			IDField{"project_id", projID},
			IDField{"cluster_id", clusterID},
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(uuids) != 3 {
			t.Fatalf("expected 3 UUIDs, got %d", len(uuids))
		}
		if uuids[0].String() != orgID {
			t.Errorf("organization_id mismatch: want %s, got %s", orgID, uuids[0])
		}
		if uuids[1].String() != projID {
			t.Errorf("project_id mismatch: want %s, got %s", projID, uuids[1])
		}
		if uuids[2].String() != clusterID {
			t.Errorf("cluster_id mismatch: want %s, got %s", clusterID, uuids[2])
		}
	})

	t.Run("four IDs (with app_service_id) parse successfully", func(t *testing.T) {
		orgID, projID, clusterID, appSvcID := newID(), newID(), newID(), newID()
		uuids, err := ParseUUIDs(
			IDField{"organization_id", orgID},
			IDField{"project_id", projID},
			IDField{"cluster_id", clusterID},
			IDField{"app_service_id", appSvcID},
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(uuids) != 4 {
			t.Fatalf("expected 4 UUIDs, got %d", len(uuids))
		}
		if uuids[3].String() != appSvcID {
			t.Errorf("app_service_id mismatch: want %s, got %s", appSvcID, uuids[3])
		}
	})

	t.Run("zero fields returns empty slice without error", func(t *testing.T) {
		uuids, err := ParseUUIDs()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(uuids) != 0 {
			t.Errorf("expected empty slice, got %v", uuids)
		}
	})

	validID := newID()
	invalidCases := []struct {
		name           string
		fields         []IDField
		wantFieldInErr string
	}{
		{
			"invalid organization_id stops early",
			[]IDField{{"organization_id", "bad"}, {"project_id", validID}, {"cluster_id", validID}},
			"organization_id",
		},
		{
			"invalid second field returns correct field name",
			[]IDField{{"organization_id", validID}, {"project_id", "bad"}, {"cluster_id", validID}},
			"project_id",
		},
		{
			"invalid last field returns correct field name",
			[]IDField{{"organization_id", validID}, {"project_id", validID}, {"cluster_id", "bad"}},
			"cluster_id",
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

func TestParseUUIDs_HierarchyCoverage(t *testing.T) {
	newID := func() string { return uuid.New().String() }

	t.Run("all valid hierarchy IDs parse successfully", func(t *testing.T) {
		orgID, projID, clusterID, appSvcID := newID(), newID(), newID(), newID()
		uuids, err := ParseUUIDs(
			IDField{"organization_id", orgID},
			IDField{"project_id", projID},
			IDField{"cluster_id", clusterID},
			IDField{"app_service_id", appSvcID},
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if uuids[0].String() != orgID {
			t.Errorf("organization_id mismatch: want %s, got %s", orgID, uuids[0])
		}
		if uuids[1].String() != projID {
			t.Errorf("project_id mismatch: want %s, got %s", projID, uuids[1])
		}
		if uuids[2].String() != clusterID {
			t.Errorf("cluster_id mismatch: want %s, got %s", clusterID, uuids[2])
		}
		if uuids[3].String() != appSvcID {
			t.Errorf("app_service_id mismatch: want %s, got %s", appSvcID, uuids[3])
		}
	})

	validID := newID()
	invalidCases := []struct {
		name           string
		fields         []IDField
		wantFieldInErr string
	}{
		{
			"invalid organization_id",
			[]IDField{{"organization_id", "bad"}, {"project_id", validID}, {"cluster_id", validID}, {"app_service_id", validID}},
			"organization_id",
		},
		{
			"invalid project_id",
			[]IDField{{"organization_id", validID}, {"project_id", "bad"}, {"cluster_id", validID}, {"app_service_id", validID}},
			"project_id",
		},
		{
			"invalid cluster_id",
			[]IDField{{"organization_id", validID}, {"project_id", validID}, {"cluster_id", "bad"}, {"app_service_id", validID}},
			"cluster_id",
		},
		{
			"invalid app_service_id",
			[]IDField{{"organization_id", validID}, {"project_id", validID}, {"cluster_id", validID}, {"app_service_id", "bad"}},
			"app_service_id",
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
