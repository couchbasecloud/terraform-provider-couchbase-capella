package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationOptions configures the schema pattern validation.
type ValidationOptions struct {
	// LegacyFiles is a map of filenames that should be skipped from validation
	LegacyFiles map[string]bool
	// AllowLegacyAttributes allows specific attribute patterns that haven't been migrated yet
	AllowLegacyAttributes []string
}

// ValidationResult holds the results of schema validation.
type ValidationResult struct {
	Failures []string
	Files    []string
}

// ValidateSchemaPatterns checks that all schema files follow the AddAttr pattern.
// It enforces:
// 1. No direct assignments to attrs map (must use capellaschema.AddAttr)
// 2. No MarkdownDescription inside attribute definitions (only on top-level schema)
// 3. No WithDescription calls (deprecated)
// 4. Required capellaschema import in all schema files
func ValidateSchemaPatterns(schemaDir string, opts ValidationOptions) (*ValidationResult, error) {
	pattern := filepath.Join(schemaDir, "*_schema.go")
	schemaFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob schema files: %w", err)
	}

	if len(schemaFiles) == 0 {
		return nil, fmt.Errorf("no schema files found in %s - pattern may be incorrect", schemaDir)
	}

	result := &ValidationResult{
		Files: schemaFiles,
	}

	for _, file := range schemaFiles {
		filename := filepath.Base(file)

		// Skip files in legacy list
		if opts.LegacyFiles[filename] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			result.Failures = append(result.Failures, formatError(filename, 0, "Failed to read file: "+err.Error(), ""))
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
			if strings.Contains(trimmed, `attrs["`) && strings.Contains(trimmed, "] =") {
				// Check if it's using AddAttr (which is OK)
				if !strings.Contains(line, "capellaschema.AddAttr(attrs,") {
					// Check if it's a known legacy attribute that's allowed
					isAllowed := false
					for _, allowedPattern := range opts.AllowLegacyAttributes {
						if strings.Contains(trimmed, allowedPattern) {
							isAllowed = true
							break
						}
					}
					if !isAllowed {
						result.Failures = append(result.Failures, formatError(filename, lineNum, "Direct attrs assignment - use capellaschema.AddAttr instead", line))
					}
				}
			}

			// Rule 2: No MarkdownDescription inside attribute definitions
			if strings.Contains(trimmed, "MarkdownDescription:") {
				if !strings.Contains(trimmed, "schema.Schema{") && !strings.Contains(trimmed, "return schema.Schema{") {
					// Check if the previous few lines have 'return schema.Schema'
					isTopLevel := false
					for j := maxInt(0, i-5); j < i; j++ {
						if strings.Contains(lines[j], "return schema.Schema{") {
							isTopLevel = true
							break
						}
					}
					// Check if it's in a filter block (some datasources have this)
					isFilterBlock := false
					for j := maxInt(0, i-10); j < i; j++ {
						if strings.Contains(lines[j], `filterAttrs["`) || strings.Contains(lines[j], `filterAttrs :=`) {
							isFilterBlock = true
							break
						}
					}
					if !isTopLevel && !isFilterBlock {
						result.Failures = append(result.Failures, formatError(filename, lineNum, "MarkdownDescription inside attribute - remove it, AddAttr handles this", line))
					}
				}
			}

			// Rule 3: No WithDescription calls (deprecated)
			if strings.Contains(trimmed, "WithDescription(") && !strings.Contains(trimmed, "WithOpenAPIDescription(") {
				result.Failures = append(result.Failures, formatError(filename, lineNum, "WithDescription is deprecated - use capellaschema.AddAttr instead", line))
			}
		}

		// Rule 4: File should import capellaschema if it defines schemas
		if strings.Contains(fileContent, "func ") && strings.Contains(fileContent, "Schema() schema.Schema") {
			if !strings.Contains(fileContent, `capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"`) {
				result.Failures = append(result.Failures, formatError(filename, 0, "Missing capellaschema import", ""))
			}
		}
	}

	return result, nil
}

// ValidateAttributesFile ensures attributes.go contains only helper functions, not schema definitions.
func ValidateAttributesFile(attributesPath string) error {
	content, err := os.ReadFile(attributesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, that's OK
		}
		return fmt.Errorf("failed to read attributes.go: %w", err)
	}

	// attributes.go should not define any Schema() functions
	if strings.Contains(string(content), "func ") && strings.Contains(string(content), "Schema() schema.Schema") {
		return fmt.Errorf("attributes.go should not define schema functions - move them to separate *_schema.go files")
	}

	return nil
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
