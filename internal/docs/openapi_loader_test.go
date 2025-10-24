package docs

import (
	"testing"
)

func TestGetOpenAPIDescription(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		tfFieldName  string
		wantEmpty    bool
	}{
		{
			name:         "project.name",
			resourceName: "project",
			tfFieldName:  "name",
			wantEmpty:    false,
		},
		{
			name:         "project.description",
			resourceName: "project",
			tfFieldName:  "description",
			wantEmpty:    false,
		},
		{
			name:         "project.id",
			resourceName: "project",
			tfFieldName:  "id",
			wantEmpty:    false,
		},
		{
			name:         "nonexistent.field",
			resourceName: "nonexistent",
			tfFieldName:  "field",
			wantEmpty:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.tfFieldName)
			if tt.wantEmpty && desc != "" {
				t.Errorf("Expected empty description, got: %s", desc)
			}
			if !tt.wantEmpty && desc == "" {
				t.Errorf("Expected non-empty description for %s.%s", tt.resourceName, tt.tfFieldName)
			}
			if !tt.wantEmpty {
				t.Logf("Description for %s.%s: %s", tt.resourceName, tt.tfFieldName, desc)
			}
		})
	}
}

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
			result := snakeToCamel(tt.input)
			if result != tt.expected {
				t.Errorf("snakeToCamel(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSnakeToCapitalizedCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"organization_id", "OrganizationId"},
		{"project_id", "ProjectId"},
		{"cluster_id", "ClusterId"},
		{"bucket_id", "BucketId"},
		{"app_service_id", "AppServiceId"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := snakeToCapitalizedCamel(tt.input)
			if result != tt.expected {
				t.Errorf("snakeToCapitalizedCamel(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetOpenAPIDescription_PathParameters(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		contains  string
	}{
		{
			name:      "organization_id",
			fieldName: "organization_id",
			contains:  "organization",
		},
		{
			name:      "project_id",
			fieldName: "project_id",
			contains:  "project",
		},
		{
			name:      "cluster_id",
			fieldName: "cluster_id",
			contains:  "cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resource name doesn't matter for path parameters
			desc := GetOpenAPIDescription("any_resource", tt.fieldName)
			if desc == "" {
				t.Errorf("Expected non-empty description for path parameter %s", tt.fieldName)
			}
			t.Logf("Path parameter %s: %s", tt.fieldName, desc)
		})
	}
}
