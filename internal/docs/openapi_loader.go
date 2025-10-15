package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var openAPIDoc *openapi3.T

// findProjectRoot walks up the directory tree to find go.mod
func findProjectRoot() (string, error) {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up until we find go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func init() {
	var openAPIPath string

	// Try environment variable first (for build/test scenarios)
	if envPath := os.Getenv("CAPELLA_OPENAPI_SPEC_PATH"); envPath != "" {
		openAPIPath = envPath
	} else {
		// Find project root by locating go.mod
		projectRoot, err := findProjectRoot()
		if err != nil {
			// Gracefully degrade - descriptions will be empty but provider still works
			fmt.Fprintf(os.Stderr, "Warning: Could not locate project root: %v\n", err)
			fmt.Fprintf(os.Stderr, "Field descriptions will not be enhanced with OpenAPI metadata.\n")
			return
		}
		openAPIPath = filepath.Join(projectRoot, "openapi.generated.yaml")
	}

	data, err := os.ReadFile(openAPIPath)
	if err != nil {
		// Gracefully degrade - descriptions will be empty but provider still works
		fmt.Fprintf(os.Stderr, "Warning: Could not load OpenAPI spec at %s: %v\n", openAPIPath, err)
		fmt.Fprintf(os.Stderr, "Field descriptions will not be enhanced with OpenAPI metadata.\n")
		return
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to parse OpenAPI spec: %v\n", err)
		fmt.Fprintf(os.Stderr, "Field descriptions will not be enhanced with OpenAPI metadata.\n")
		return
	}

	openAPIDoc = doc
}

// GetOpenAPIDescription retrieves an enhanced description for a field from the OpenAPI spec.
// Automatically tries common schema patterns (CreateXRequest, GetXResponse, UpdateXRequest).
// Converts snake_case field names to camelCase.
// Returns empty string if schema or field not found.
func GetOpenAPIDescription(resourceName, tfFieldName string) string {
	if openAPIDoc == nil || openAPIDoc.Components == nil || openAPIDoc.Components.Schemas == nil {
		return ""
	}

	// Convert snake_case to camelCase
	camelFieldName := snakeToCamel(tfFieldName)

	// Capitalize resource name for schema patterns
	capitalizedResource := capitalize(resourceName)

	// Try common schema name patterns
	schemaPatterns := []string{
		"Create" + capitalizedResource + "Request",
		"Get" + capitalizedResource + "Response",
		"Update" + capitalizedResource + "Request",
		capitalizedResource + "Request",
		capitalizedResource + "Response",
	}

	// Try each schema pattern until we find the field
	for _, schemaName := range schemaPatterns {
		schemaRef := openAPIDoc.Components.Schemas[schemaName]
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		propRef := schemaRef.Value.Properties[camelFieldName]
		if propRef != nil && propRef.Value != nil {
			return buildEnhancedDescription(propRef.Value)
		}
	}

	return ""
}

// snakeToCamel converts snake_case to camelCase
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	// First part stays lowercase
	result := parts[0]

	// Capitalize first letter of remaining parts
	for i := 1; i < len(parts); i++ {
		result += capitalize(parts[i])
	}

	return result
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// buildEnhancedDescription creates a rich markdown description from OpenAPI property
func buildEnhancedDescription(prop *openapi3.Schema) string {
	var parts []string

	// Base description
	if prop.Description != "" {
		parts = append(parts, strings.TrimSpace(prop.Description))
	}

	// Add constraints section if any exist
	constraints := buildConstraints(prop)
	if len(constraints) > 0 {
		parts = append(parts, "\n\n**Constraints:**\n")
		for _, constraint := range constraints {
			parts = append(parts, fmt.Sprintf("  - %s\n", constraint))
		}
	}

	// Add enum values
	if len(prop.Enum) > 0 {
		parts = append(parts, "\n**Valid Values:**\n")
		for _, val := range prop.Enum {
			parts = append(parts, fmt.Sprintf("  - `%v`\n", val))
		}
	}

	// Add default value
	if prop.Default != nil {
		parts = append(parts, fmt.Sprintf("\n**Default:** `%v`\n", prop.Default))
	}

	// Add format information
	if prop.Format != "" {
		formatDesc := getFormatDescription(prop.Format)
		if formatDesc != "" {
			parts = append(parts, fmt.Sprintf("\n**Format:** %s\n", formatDesc))
		}
	}

	// Add deprecation warning
	if prop.Deprecated {
		parts = append(parts, "\n **Deprecated**: This field is deprecated and will be removed in a future release.\n")
	}

	return strings.Join(parts, "")
}

// buildConstraints extracts constraint information from schema
func buildConstraints(prop *openapi3.Schema) []string {
	var constraints []string

	if prop.MinLength > 0 {
		constraints = append(constraints, fmt.Sprintf("Minimum length: %d characters", prop.MinLength))
	}

	if prop.MaxLength != nil && *prop.MaxLength > 0 {
		constraints = append(constraints, fmt.Sprintf("Maximum length: %d characters", *prop.MaxLength))
	}

	if prop.Min != nil {
		constraints = append(constraints, fmt.Sprintf("Minimum: %v", *prop.Min))
	}

	if prop.Max != nil {
		constraints = append(constraints, fmt.Sprintf("Maximum: %v", *prop.Max))
	}

	if prop.Pattern != "" {
		constraints = append(constraints, fmt.Sprintf("Pattern: `%s`", prop.Pattern))
	}

	return constraints
}

// getFormatDescription returns a human-readable description for common formats
func getFormatDescription(format string) string {
	formats := map[string]string{
		"uuid":      "UUID (GUID4)",
		"date":      "Date in RFC3339 format",
		"date-time": "Date-time in RFC3339 format",
		"email":     "Email address",
		"uri":       "URI",
		"hostname":  "Hostname",
		"ipv4":      "IPv4 address",
		"ipv6":      "IPv6 address",
	}

	if desc, ok := formats[format]; ok {
		return desc
	}
	return ""
}
