package resources

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema/validator"
)

// TestAllSchemasUseAddAttrPattern validates that all resource schema files follow the AddAttr pattern.
func TestAllSchemasUseAddAttrPattern(t *testing.T) {
	opts := validator.ValidationOptions{
		LegacyFiles: map[string]bool{
			// Helper file, not a schema
			"attributes.go": true,
		},
		// No legacy attributes allowed for resources - all migrated
		AllowLegacyAttributes: []string{},
	}

	result, err := validator.ValidateSchemaPatterns(".", opts)
	if err != nil {
		t.Fatalf("Failed to validate schemas: %v", err)
	}

	if len(result.Files) == 0 {
		t.Fatal("No schema files found - this test may be running from wrong directory")
	}

	if len(result.Failures) > 0 {
		t.Errorf("Found %d AddAttr pattern violations:\n\n%s", len(result.Failures), strings.Join(result.Failures, "\n"))
	}
}

// TestAttributesFileDoesNotDefineSchemas ensures attributes.go is only helpers
func TestAttributesFileDoesNotDefineSchemas(t *testing.T) {
	err := validator.ValidateAttributesFile(filepath.Join(".", "attributes.go"))
	if err != nil {
		t.Error(err)
	}
}
