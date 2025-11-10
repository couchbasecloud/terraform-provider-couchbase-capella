package datasources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAllSchemasUseAddAttrPattern validates that all datasource schema files follow the AddAttr pattern.
// This test enforces the following rules:
// 1. No direct assignments to attrs map (e.g., attrs["field"] = ...)
// 2. No MarkdownDescription inside attribute definitions (only allowed on top-level schema)
// 3. No WithDescription calls (deprecated in favor of AddAttr)
func TestAllSchemasUseAddAttrPattern(t *testing.T) {
	// Get all *_schema.go files in this directory
	schemaFiles, err := filepath.Glob("*_schema.go")
	if err != nil {
		t.Fatalf("Failed to glob schema files: %v", err)
	}

	if len(schemaFiles) == 0 {
		t.Fatal("No schema files found - this test may be running from wrong directory")
	}

	// Files that are being migrated or have legacy patterns
	// Remove files from this list as they are committed in PRs
	legacyFiles := map[string]bool{
		// All datasource schemas have been migrated

		// Helper file, not a schema
		"attributes.go": true,
	}

	var failures []string

	for _, file := range schemaFiles {
		// Skip files in legacy list
		if legacyFiles[file] {
			t.Logf("Skipping legacy file: %s", file)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", file, err)
			continue
		}

		fileContent := string(content)
		lines := strings.Split(fileContent, "\n")

		// Check for antipatterns
		for i, line := range lines {
			lineNum := i + 1
			trimmed := strings.TrimSpace(line)

			// Skip comments
			if strings.HasPrefix(trimmed, "//") {
				continue
			}

			// Rule 1: No direct attrs assignments
			// Exception: Legacy attributes like computedEventAttributes and computedCursorAttribute
			// that are defined in attributes.go with MarkdownDescription are allowed
			if strings.Contains(trimmed, `attrs["`) && strings.Contains(trimmed, "] =") {
				// Check if it's using AddAttr (which is OK)
				if !strings.Contains(line, "capellaschema.AddAttr(attrs,") {
					// Check if it's a known legacy attribute (computedEventAttributes, computedCursorAttribute)
					if !strings.Contains(trimmed, "computedEventAttributes") &&
						!strings.Contains(trimmed, "computedCursorAttribute") &&
						!strings.Contains(trimmed, "computedAuditAttribute") {
						failures = append(failures, formatError(file, lineNum, "Direct attrs assignment - use capellaschema.AddAttr instead", line))
					}
				}
			}

			// Rule 2: No MarkdownDescription inside attribute definitions
			// Allow: MarkdownDescription: "..." on schema.Schema{}
			// Disallow: MarkdownDescription: "..." inside &schema.StringAttribute{}, etc.
			if strings.Contains(trimmed, "MarkdownDescription:") {
				// Check if this is inside an attribute definition (has '&schema.' nearby)
				// We need to look at context - if we see 'return schema.Schema{' before the next attribute, it's OK
				// For simplicity, flag any MarkdownDescription not on a line with 'schema.Schema'
				if !strings.Contains(trimmed, "schema.Schema{") && !strings.Contains(trimmed, "return schema.Schema{") {
					// Check if the previous few lines have 'return schema.Schema'
					isTopLevel := false
					for j := max(0, i-5); j < i; j++ {
						if strings.Contains(lines[j], "return schema.Schema{") {
							isTopLevel = true
							break
						}
					}
					// Also check if it's in the filter block (app_endpoints_schema.go has this)
					isFilterBlock := false
					for j := max(0, i-10); j < i; j++ {
						if strings.Contains(lines[j], `filterAttrs["`) || strings.Contains(lines[j], `filterAttrs :=`) {
							isFilterBlock = true
							break
						}
					}
					if !isTopLevel && !isFilterBlock {
						failures = append(failures, formatError(file, lineNum, "MarkdownDescription inside attribute - remove it, AddAttr handles this", line))
					}
				}
			}

			// Rule 3: No WithDescription calls (deprecated)
			if strings.Contains(trimmed, "WithDescription(") && !strings.Contains(trimmed, "WithOpenAPIDescription(") {
				failures = append(failures, formatError(file, lineNum, "WithDescription is deprecated - use capellaschema.AddAttr instead", line))
			}
		}

		// Rule 4: File should import capellaschema if it defines schemas
		if strings.Contains(fileContent, "func ") && strings.Contains(fileContent, "Schema() schema.Schema") {
			if !strings.Contains(fileContent, `capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"`) {
				failures = append(failures, formatError(file, 0, "Missing capellaschema import", ""))
			}
		}
	}

	if len(failures) > 0 {
		t.Errorf("Found %d AddAttr pattern violations:\n\n%s", len(failures), strings.Join(failures, "\n"))
	}
}

// TestAttributesFileDoesNotDefineSchemas ensures attributes.go is only helpers
func TestAttributesFileDoesNotDefineSchemas(t *testing.T) {
	content, err := os.ReadFile("attributes.go")
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("attributes.go does not exist")
		}
		t.Fatalf("Failed to read attributes.go: %v", err)
	}

	// attributes.go should not define any Schema() functions
	if strings.Contains(string(content), "func ") && strings.Contains(string(content), "Schema() schema.Schema") {
		t.Error("attributes.go should not define schema functions - move them to separate *_schema.go files")
	}
}

func formatError(file string, line int, message string, code string) string {
	if line > 0 && code != "" {
		return fmt.Sprintf("  %s:%d: %s\n    %s", file, line, message, strings.TrimSpace(code))
	}
	if line > 0 {
		return fmt.Sprintf("  %s:%d: %s", file, line, message)
	}
	return fmt.Sprintf("  %s: %s", file, message)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
