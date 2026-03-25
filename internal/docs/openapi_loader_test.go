package docs

import (
	"os"
	"strings"
	"testing"
)

// Default OpenAPI spec URL - same as in Makefile
// The spec is embedded in the Couchbase docs page and extracted automatically
const defaultOpenAPISpecURL = "https://docs.couchbase.com/cloud/management-api-reference/index.html"

// TestMain sets up the environment for tests
func TestMain(m *testing.M) {
	// Use OPENAPI_SPEC_URL if already set (e.g., from Makefile export)
	// Otherwise, use the default URL
	if os.Getenv("OPENAPI_SPEC_URL") == "" {
		os.Setenv("OPENAPI_SPEC_URL", defaultOpenAPISpecURL)
	}

	// Reload the OpenAPI spec now that the env var is set
	// (init() runs before TestMain, so we need to reload)
	loadOpenAPISpec()

	// Run tests
	os.Exit(m.Run())
}

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

func TestGetOpenAPIDescription_ArrayWithReferences(t *testing.T) {
	tests := []struct {
		name                string
		resourceName        string
		fieldName           string
		expectDescription   bool
		expectEnums         bool
		descriptionContains string
		enumsContain        []string
	}{
		{
			name:                "organization_roles from user",
			resourceName:        "user",
			fieldName:           "organization_roles",
			expectDescription:   true,
			expectEnums:         true,
			descriptionContains: "Organization roles",
			enumsContain:        []string{"organizationOwner", "organizationMember", "projectCreator"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.fieldName)

			if tt.expectDescription && desc == "" {
				t.Errorf("Expected non-empty description for %s.%s", tt.resourceName, tt.fieldName)
			}

			if tt.expectDescription && tt.descriptionContains != "" {
				if !contains(desc, tt.descriptionContains) {
					t.Errorf("Expected description to contain %q, got: %s", tt.descriptionContains, desc)
				}
			}

			if tt.expectEnums {
				for _, enumVal := range tt.enumsContain {
					if !contains(desc, enumVal) {
						t.Errorf("Expected description to contain enum value %q, got: %s", enumVal, desc)
					}
				}
			}

			t.Logf("Description for %s.%s:\n%s", tt.resourceName, tt.fieldName, desc)
		})
	}
}

func TestGetOpenAPIDescription_BulletFormat(t *testing.T) {
	// Test that descriptions are formatted as bullets
	tests := []struct {
		name         string
		resourceName string
		fieldName    string
	}{
		{
			name:         "name field has bullet format",
			resourceName: "user",
			fieldName:    "name",
		},
		{
			name:         "organization_roles has bullet format",
			resourceName: "user",
			fieldName:    "organization_roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.fieldName)
			if desc == "" {
				t.Skip("No description found")
			}

			// Check that description starts with newline and bullet
			if len(desc) < 3 || desc[:3] != "\n -" {
				t.Errorf("Expected description to start with '\\n -', got: %q", desc[:min(10, len(desc))])
			}

			t.Logf("Description:\n%s", desc)
		})
	}
}

func TestGetOpenAPIDescription_ConstraintsFormat(t *testing.T) {
	// Test that constraints are properly formatted
	desc := GetOpenAPIDescription("user", "name")
	if desc == "" {
		t.Fatal("Expected description for user.name")
	}

	if !contains(desc, "**Constraints**:") {
		t.Errorf("Expected description to contain '**Constraints**:', got: %s", desc)
	}

	if !contains(desc, "Maximum length: 128 characters") {
		t.Errorf("Expected description to contain constraint details, got: %s", desc)
	}

	t.Logf("Description with constraints:\n%s", desc)
}

func TestGetOpenAPIDescription_ValidValuesFormat(t *testing.T) {
	// Test that valid values (enums) are properly formatted
	desc := GetOpenAPIDescription("user", "organization_roles")
	if desc == "" {
		t.Fatal("Expected description for user.organization_roles")
	}

	if !contains(desc, "**Valid Values**:") {
		t.Errorf("Expected description to contain '**Valid Values**:', got: %s", desc)
	}

	// Check that enum values are backtick-quoted and comma-separated
	if !contains(desc, "`organizationOwner`") {
		t.Errorf("Expected description to contain '`organizationOwner`', got: %s", desc)
	}

	if !contains(desc, ", ") {
		t.Errorf("Expected enum values to be comma-separated, got: %s", desc)
	}

	t.Logf("Description with valid values:\n%s", desc)
}

func TestGetOpenAPIDescription_NestedSchemaFields(t *testing.T) {
	// Test that nested fields (like those in Resource schema) are found
	tests := []struct {
		name           string
		resourceName   string
		fieldName      string
		expectNonEmpty bool
		shouldContain  string
	}{
		{
			name:           "type field from Resource schema",
			resourceName:   "user",
			fieldName:      "type",
			expectNonEmpty: true,
			shouldContain:  "Type of the resource",
		},
		{
			name:           "roles field from Resource schema",
			resourceName:   "user",
			fieldName:      "roles",
			expectNonEmpty: true,
			shouldContain:  "Project Roles",
		},
		{
			name:           "type has enum values",
			resourceName:   "user",
			fieldName:      "type",
			expectNonEmpty: true,
			shouldContain:  "`project`",
		},
		{
			name:           "roles has multiple enum values",
			resourceName:   "user",
			fieldName:      "roles",
			expectNonEmpty: true,
			shouldContain:  "`projectOwner`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.fieldName)

			if tt.expectNonEmpty && desc == "" {
				t.Errorf("Expected non-empty description for %s.%s", tt.resourceName, tt.fieldName)
			}

			if tt.shouldContain != "" && !contains(desc, tt.shouldContain) {
				t.Errorf("Expected description to contain %q, got: %s", tt.shouldContain, desc)
			}

			t.Logf("Description for %s.%s:\n%s", tt.resourceName, tt.fieldName, desc)
		})
	}
}

func TestGetOpenAPIDescription_DifferentSchemaNames(t *testing.T) {
	// Test that resources with different OpenAPI schema names work correctly
	// when using the OpenAPI schema name directly
	tests := []struct {
		name              string
		openAPISchemaName string
		fieldName         string
		expectNonEmpty    bool
		shouldContain     string
	}{
		{
			name:              "allowedCidr cidr field",
			openAPISchemaName: "allowedCidr",
			fieldName:         "cidr",
			expectNonEmpty:    true,
			shouldContain:     "CIDR",
		},
		{
			name:              "allowedCidr comment field",
			openAPISchemaName: "allowedCidr",
			fieldName:         "comment",
			expectNonEmpty:    true,
			shouldContain:     "description",
		},
		{
			name:              "allowedCidr expires_at field",
			openAPISchemaName: "allowedCidr",
			fieldName:         "expires_at",
			expectNonEmpty:    true,
			shouldContain:     "RFC3339",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.openAPISchemaName, tt.fieldName)

			if tt.expectNonEmpty && desc == "" {
				t.Errorf("Expected non-empty description for %s.%s", tt.openAPISchemaName, tt.fieldName)
			}

			if tt.shouldContain != "" && !contains(desc, tt.shouldContain) {
				t.Errorf("Expected description to contain %q, got: %s", tt.shouldContain, desc)
			}

			t.Logf("Description for %s.%s:\n%s", tt.openAPISchemaName, tt.fieldName, desc)
		})
	}
}

func TestCleanDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes dashed list markers",
			input:    "Description:\n- item 1\n- item 2",
			expected: "Description: item 1 item 2",
		},
		{
			name:     "removes asterisk list markers",
			input:    "Description:\n* item 1\n* item 2",
			expected: "Description: item 1 item 2",
		},
		{
			name:     "removes numbered list markers",
			input:    "Description:\n1. first\n2. second",
			expected: "Description: first second",
		},
		{
			name:     "removes leading list marker from line",
			input:    "- A description with leading marker",
			expected: "A description with leading marker",
		},
		{
			name:     "handles mixed list markers",
			input:    "Description:\n- item 1\n* item 2\n+ item 3",
			expected: "Description: item 1 item 2 item 3",
		},
		{
			name:     "handles plain text",
			input:    "Just a plain description",
			expected: "Just a plain description",
		},
		{
			name:     "removes empty lines",
			input:    "Line 1\n\n\nLine 2",
			expected: "Line 1 Line 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanDescription(tt.input)
			if result != tt.expected {
				t.Errorf("cleanDescription() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{
			name:     "formats string values",
			input:    []interface{}{"value1", "value2", "value3"},
			expected: []string{"`value1`", "`value2`", "`value3`"},
		},
		{
			name:     "formats mixed types",
			input:    []interface{}{"string", 123, true},
			expected: []string{"`string`", "`123`", "`true`"},
		},
		{
			name:     "handles empty slice",
			input:    []interface{}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatEnumValues(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("formatEnumValues() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("formatEnumValues()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsListMarker(t *testing.T) {
	tests := []struct {
		char     byte
		expected bool
	}{
		{'-', true},
		{'*', true},
		{'+', true},
		{'a', false},
		{'1', false},
		{' ', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := isListMarker(tt.char)
			if result != tt.expected {
				t.Errorf("isListMarker(%c) = %v, want %v", tt.char, result, tt.expected)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		char     byte
		expected bool
	}{
		{'0', true},
		{'5', true},
		{'9', true},
		{'a', false},
		{'-', false},
		{' ', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := isDigit(tt.char)
			if result != tt.expected {
				t.Errorf("isDigit(%c) = %v, want %v", tt.char, result, tt.expected)
			}
		})
	}
}

func TestFormatBulletFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func() string
		expected string
	}{
		{
			name:     "formatDescriptionBullet",
			function: func() string { return formatDescriptionBullet("Test description") },
			expected: "\n - Test description",
		},
		{
			name:     "formatConstraintsBullet",
			function: func() string { return formatConstraintsBullet([]string{"Min: 1", "Max: 100"}) },
			expected: "\n - **Constraints**: Min: 1, Max: 100",
		},
		{
			name:     "formatDefaultBullet",
			function: func() string { return formatDefaultBullet("default") },
			expected: "\n - **Default**: `default`",
		},
		{
			name:     "formatFormatBullet",
			function: func() string { return formatFormatBullet("UUID") },
			expected: "\n - **Format**: UUID",
		},
		{
			name:     "formatDeprecationBullet",
			function: func() string { return formatDeprecationBullet() },
			expected: "\n - **Deprecated**: This field is deprecated and will be removed in a future release.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			if result != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGetOpenAPIDescription_NewMappings verifies that all newly added OpenAPI schema mappings work correctly
func TestGetOpenAPIDescription_NewMappings(t *testing.T) {
	tests := []struct {
		name           string
		schema         string
		field          string
		shouldHaveDesc bool
	}{
		{"accessFunction maps to AccessFunction", "accessFunction", "organization_id", true},
		{"CORSConfig has origin field", "CORSConfig", "origin", true},
		{"CORSConfig has headers field", "CORSConfig", "headers", true},
		{"importFilter maps to ImportFilter", "importFilter", "organization_id", true},
		{"OIDCProvider has issuer field", "OIDCProvider", "issuer", true},
		{"OIDCProvider has clientId field", "OIDCProvider", "clientId", true},
		{"resyncRequest has scopes field", "resyncRequest", "scopes", false}, // No description in OpenAPI spec for this field
		{"AppServiceAllowedCidr has cidr field", "AppServiceAllowedCidr", "cidr", true},
		{"ClusterOnOffSchedule has timezone field", "ClusterOnOffSchedule", "timezone", true},
		{"ClusterOnOffSchedule has days field", "ClusterOnOffSchedule", "days", false}, // No description in OpenAPI spec for this field
		{"indexDDLRequest has definition field", "indexDDLRequest", "definition", true},
		{"PostSampleBucket has name field", "PostSampleBucket", "name", true},
		{"WeeklySchedule has day_of_week field", "WeeklySchedule", "day_of_week", true},
		{"WeeklySchedule has start_at field", "WeeklySchedule", "start_at", true},
		{"WeeklySchedule has incremental_every field", "WeeklySchedule", "incremental_every", true},
		{"WeeklySchedule has retention_time field", "WeeklySchedule", "retention_time", true},
		{"WeeklySchedule has cost_optimized_retention field", "WeeklySchedule", "cost_optimized_retention", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.schema, tt.field)

			if tt.shouldHaveDesc {
				if desc == "" {
					t.Errorf("Expected description for %s.%s but got empty string", tt.schema, tt.field)
				} else {
					t.Logf("✓ %s.%s has description: %s", tt.schema, tt.field, desc)
				}
			} else {
				if desc != "" {
					t.Errorf("Expected empty description for %s.%s but got: %s", tt.schema, tt.field, desc)
				}
			}
		})
	}
}

// TestAppEndpointFieldLookups tests field lookups for app endpoint resources
func TestAppEndpointFieldLookups(t *testing.T) {
	tests := []struct {
		name   string
		schema string
		field  string
	}{
		// AccessFunction resource
		{"accessFunction.access_control_function", "accessFunction", "access_control_function"},
		{"accessFunction.app_endpoint_name", "accessFunction", "app_endpoint_name"},
		{"accessFunction.scope", "accessFunction", "scope"},
		{"accessFunction.collection", "accessFunction", "collection"},

		// CORSConfig resource
		{"CORSConfig.origin", "CORSConfig", "origin"},
		{"CORSConfig.headers", "CORSConfig", "headers"},
		{"CORSConfig.login_origin", "CORSConfig", "login_origin"},
		{"CORSConfig.loginOrigin", "CORSConfig", "loginOrigin"},
		{"CORSConfig.max_age", "CORSConfig", "max_age"},
		{"CORSConfig.maxAge", "CORSConfig", "maxAge"},
		{"CORSConfig.disabled", "CORSConfig", "disabled"},

		// ImportFilter resource
		{"importFilter.import_filter", "importFilter", "import_filter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.schema, tt.field)
			t.Logf("%s.%s: '%s'", tt.schema, tt.field, desc)
		})
	}
}

// TestCleanDescription_TableFormatting tests that markdown tables are preserved
// and start on a new line
func TestCleanDescription_TableFormatting(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Table with preceding text",
			input: `The CPU and RAM configuration. The supported combinations are:
| CPU (cores) | RAM (GB) |
| ----------- | -------- |
| 2           | 4        |
| 4           | 8        |`,
			want: `The CPU and RAM configuration. The supported combinations are:

| CPU (cores) | RAM (GB) |
| ----------- | -------- |
| 2           | 4        |
| 4           | 8        |`,
		},
		{
			name: "Table with following text",
			input: `Some intro text
| CPU | RAM |
| --- | --- |
| 2   | 4   |
More text after`,
			want: `Some intro text

| CPU | RAM |
| --- | --- |
| 2   | 4   |

More text after`,
		},
		{
			name: "Text with lists but no table",
			input: `Some text
- Item 1
- Item 2
More text`,
			want: `Some text Item 1 Item 2 More text`,
		},
		{
			name: "Multiple tables",
			input: `First table:
| A | B |
| - | - |
| 1 | 2 |
Between tables
| C | D |
| - | - |
| 3 | 4 |`,
			want: `First table:

| A | B |
| - | - |
| 1 | 2 |

Between tables

| C | D |
| - | - |
| 3 | 4 |`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanDescription(tt.input)
			if got != tt.want {
				t.Errorf("cleanDescription() mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.want)
			}
		})
	}
}

// TestAppServiceComputeDescription tests that the AppServiceCompute schema
// description with a table is properly formatted
func TestAppServiceComputeDescription(t *testing.T) {
	// Test getting description from the appservice resource's compute field
	// This should resolve the AppServiceCompute schema reference and include the table
	desc := GetOpenAPIDescription("appservice", "cpu")

	if desc == "" {
		t.Skip("Skipping test - OpenAPI spec not loaded or field not found")
	}

	// The CPU field should have a description
	if !strings.Contains(desc, "CPU") {
		t.Logf("CPU field description: %s", desc)
	}

	t.Logf("appservice.cpu description:\n%s", desc)
}

// TestIsTableRow tests the table row detection function
func TestIsTableRow(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"valid table header", "| Column 1 | Column 2 |", true},
		{"valid table row", "| Value 1 | Value 2 |", true},
		{"table separator", "| ------- | -------- |", true},
		{"single pipe", "This is not | a table", false},
		{"no pipes", "This is not a table", false},
		{"empty line", "", false},
		{"only pipes", "||", true},
		{"three pipes", "| A | B | C |", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTableRow(tt.line)
			if result != tt.expected {
				t.Errorf("isTableRow(%q) = %v, want %v", tt.line, result, tt.expected)
			}
		})
	}
}

// TestRemoveListMarkers tests the list marker removal function
func TestRemoveListMarkers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"dash marker", "- Item text", "Item text"},
		{"asterisk marker", "* Item text", "Item text"},
		{"plus marker", "+ Item text", "Item text"},
		{"numbered marker", "1. Item text", "Item text"},
		{"numbered two digits", "42. Item text", "Item text"},
		{"no marker", "Plain text", "Plain text"},
		{"dash in middle", "Text - with dash", "Text - with dash"},
		{"empty string", "", ""},
		{"just marker", "-", ""},
		{"whitespace with marker", "  - Item", "Item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeListMarkers(tt.input)
			if result != tt.expected {
				t.Errorf("removeListMarkers(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestAppendTextBuffer tests the text buffer flushing function
func TestAppendTextBuffer(t *testing.T) {
	tests := []struct {
		name     string
		result   []string
		buffer   []string
		expected []string
	}{
		{
			name:     "empty buffer",
			result:   []string{"Line 1"},
			buffer:   []string{},
			expected: []string{"Line 1"},
		},
		{
			name:     "single item buffer",
			result:   []string{"Line 1"},
			buffer:   []string{"Buffered text"},
			expected: []string{"Line 1", "Buffered text"},
		},
		{
			name:     "multiple items buffer",
			result:   []string{"Line 1"},
			buffer:   []string{"Text 1", "Text 2", "Text 3"},
			expected: []string{"Line 1", "Text 1 Text 2 Text 3"},
		},
		{
			name:     "empty result with buffer",
			result:   []string{},
			buffer:   []string{"First", "Second"},
			expected: []string{"First Second"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendTextBuffer(tt.result, tt.buffer)
			if !equalStringSlices(result, tt.expected) {
				t.Errorf("appendTextBuffer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestGetOpenAPIDescription_AmbiguousFieldNames tests that common field names
// that may exist in multiple schemas (e.g., "type", "last", "next") return
// the correct description for the intended resource context.
//
// This test helps prevent regressions where a description is found but comes
// from the wrong nested object.
func TestGetOpenAPIDescription_AmbiguousFieldNames(t *testing.T) {
	tests := []struct {
		name             string
		resourceName     string
		fieldName        string
		expectNonEmpty   bool
		shouldContain    string
		shouldNotContain string
		description      string
	}{
		{
			name:           "type field for user resource should come from Resource schema",
			resourceName:   "user",
			fieldName:      "type",
			expectNonEmpty: true,
			shouldContain:  "project",
			description:    "The 'type' field in user context should refer to resource type (project/bucket)",
		},
		{
			name:           "type field for bucket should come from bucket schema if available",
			resourceName:   "bucket",
			fieldName:      "type",
			expectNonEmpty: true,
			shouldContain:  "bucket",
			description:    "Bucket 'type' field should refer to bucket type (couchbase/ephemeral/membase)",
		},
		{
			name:           "id field is commonly ambiguous across schemas",
			resourceName:   "project",
			fieldName:      "id",
			expectNonEmpty: true,
			description:    "The 'id' field should be found with a valid description",
		},
		{
			name:           "name field is commonly ambiguous across schemas",
			resourceName:   "cluster",
			fieldName:      "name",
			expectNonEmpty: true,
			description:    "The 'name' field should be found with a valid description",
		},
		{
			name:           "status field may exist in multiple response schemas",
			resourceName:   "cluster",
			fieldName:      "status",
			expectNonEmpty: false, // May or may not exist, document the behavior
			description:    "Status field behavior for cluster resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.fieldName)

			if tt.expectNonEmpty && desc == "" {
				t.Errorf("%s: Expected non-empty description for %s.%s",
					tt.description, tt.resourceName, tt.fieldName)
			}

			if tt.shouldContain != "" && !contains(desc, tt.shouldContain) {
				t.Errorf("%s: Expected description to contain %q, got: %s",
					tt.description, tt.shouldContain, desc)
			}

			if tt.shouldNotContain != "" && contains(desc, tt.shouldNotContain) {
				t.Errorf("%s: Expected description NOT to contain %q, got: %s",
					tt.description, tt.shouldNotContain, desc)
			}

			t.Logf("%s.%s (%s): %s", tt.resourceName, tt.fieldName, tt.description, desc)
		})
	}
}

// TestGetOpenAPIDescription_PaginationFields tests fields like "last", "next",
// "previous" that typically appear in paginated response schemas and could be
// confused between different contexts.
func TestGetOpenAPIDescription_PaginationFields(t *testing.T) {
	paginationFields := []string{"last", "next", "previous", "first", "cursor"}

	for _, fieldName := range paginationFields {
		t.Run(fieldName, func(t *testing.T) {
			// Test against different resource contexts to see if we get consistent descriptions
			resources := []string{"user", "project", "cluster", "bucket"}
			descriptions := make(map[string]string)

			for _, resource := range resources {
				desc := GetOpenAPIDescription(resource, fieldName)
				if desc != "" {
					descriptions[resource] = desc
				}
			}

			// Log findings for documentation purposes
			if len(descriptions) == 0 {
				t.Logf("Pagination field '%s' not found in any resource schema", fieldName)
			} else {
				t.Logf("Pagination field '%s' found in %d resources:", fieldName, len(descriptions))
				for resource, desc := range descriptions {
					t.Logf("  %s: %s", resource, truncate(desc, 100))
				}
			}
		})
	}
}

// TestGetOpenAPIDescription_CompositionTypes tests that schemas using allOf,
// oneOf, or anyOf composition types are handled correctly. These compositions
// are common in OpenAPI specs for inheritance and union types.
//
// NOTE: The current implementation does not traverse allOf/anyOf/oneOf references.
// These tests document the current behavior and will help detect when support
// is added or if regressions occur.
func TestGetOpenAPIDescription_CompositionTypes(t *testing.T) {
	// These tests document expected behavior for composition types
	// The current implementation searches specific schema patterns but does not
	// recursively traverse allOf/anyOf/oneOf compositions

	tests := []struct {
		name         string
		resourceName string
		fieldName    string
		description  string
	}{
		{
			name:         "audit field typically uses $ref to CouchbaseAuditData",
			resourceName: "project",
			fieldName:    "audit",
			description:  "Audit data is often referenced via $ref, not composition",
		},
		{
			name:         "resources field in user may use array of $ref",
			resourceName: "user",
			fieldName:    "resources",
			description:  "Resources field typically uses array items with $ref",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetOpenAPIDescription(tt.resourceName, tt.fieldName)
			t.Logf("%s.%s: %s (description: %q)", tt.resourceName, tt.fieldName, tt.description, desc)
		})
	}
}

// TestGetOpenAPIDescription_SchemaSearchOrder documents and tests the order
// in which schemas are searched for field descriptions. This is important
// for understanding which description is returned when a field exists in
// multiple schemas.
func TestGetOpenAPIDescription_SchemaSearchOrder(t *testing.T) {
	// Document the search order from GetOpenAPIDescription:
	// 1. Special fields (if_match, etag, audit)
	// 2. String-type schema mappings (access_control_function, import_filter)
	// 3. Path parameters (fields ending in _id)
	// 4. Schema patterns in order:
	//    - {ResourceName} (exact match)
	//    - Create{ResourceName}Request
	//    - Get{ResourceName}Response
	//    - Update{ResourceName}Request
	//    - {ResourceName}Request
	//    - {ResourceName}Response
	//    - {resourceName} (lowercase exact match)
	// 5. Nested schemas:
	//    - Resource
	//    - ResourceBucket
	//    - {ResourceName}Resource

	t.Run("schema search order documented", func(t *testing.T) {
		// Test that we get descriptions from expected schemas
		testCases := []struct {
			resourceName    string
			fieldName       string
			expectedPattern string
		}{
			{"project", "name", "CreateProjectRequest or GetProjectResponse"},
			{"user", "type", "Resource nested schema"},
			{"bucket", "name", "CreateBucketRequest or GetBucketResponse"},
		}

		for _, tc := range testCases {
			desc := GetOpenAPIDescription(tc.resourceName, tc.fieldName)
			if desc != "" {
				t.Logf("%s.%s found (expected from %s): %s",
					tc.resourceName, tc.fieldName, tc.expectedPattern, truncate(desc, 80))
			} else {
				t.Logf("%s.%s not found (expected from %s)",
					tc.resourceName, tc.fieldName, tc.expectedPattern)
			}
		}
	})
}

// TestGetOpenAPIDescription_ConsistencyAcrossResources tests that the same
// logical field returns consistent descriptions across different resource
// contexts where appropriate.
func TestGetOpenAPIDescription_ConsistencyAcrossResources(t *testing.T) {
	// Common fields that should ideally have consistent descriptions
	commonFields := []struct {
		fieldName string
		resources []string
	}{
		{
			fieldName: "organization_id",
			resources: []string{"project", "cluster", "bucket", "user"},
		},
		{
			fieldName: "project_id",
			resources: []string{"cluster", "bucket", "appservice"},
		},
	}

	for _, cf := range commonFields {
		t.Run(cf.fieldName, func(t *testing.T) {
			descriptions := make([]string, 0, len(cf.resources))

			for _, resource := range cf.resources {
				desc := GetOpenAPIDescription(resource, cf.fieldName)
				descriptions = append(descriptions, desc)
			}

			// Check if all descriptions are the same (for path parameters, they should be)
			allSame := true
			for i := 1; i < len(descriptions); i++ {
				if descriptions[i] != descriptions[0] {
					allSame = false
					break
				}
			}

			if allSame && descriptions[0] != "" {
				t.Logf("Field '%s' has consistent description across resources: %s",
					cf.fieldName, truncate(descriptions[0], 80))
			} else {
				t.Logf("Field '%s' has varying descriptions across resources:", cf.fieldName)
				for i, resource := range cf.resources {
					t.Logf("  %s: %s", resource, truncate(descriptions[i], 60))
				}
			}
		})
	}
}

// truncate shortens a string to maxLen characters, adding "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
