package docs

import (
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var openAPIDoc *openapi3.T

func init() {
	// Get OpenAPI spec path from environment variable, or try multiple locations
	openAPIPath := os.Getenv("CAPELLA_OPENAPI_SPEC_PATH")
	if openAPIPath == "" {
		// Try multiple locations - first current dir, then parent dirs for tests
		possiblePaths := []string{
			"openapi.generated.yaml",          // From project root
			"../../openapi.generated.yaml",    // From internal/docs/
			"../../../openapi.generated.yaml", // From internal/docs/subdir if needed
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				openAPIPath = path
				break
			}
		}

		if openAPIPath == "" {
			openAPIPath = "openapi.generated.yaml" // Fallback
		}
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
// Priority order:
// 1. Path parameters (components.parameters) - for fields ending in _id
// 2. Header parameters (components.parameters) - for if_match
// 3. Response headers (components.headers) - for etag
// 4. Schema references (components.schemas) - for audit
// 5. Schema properties (CreateXRequest, GetXResponse, etc.)
// Converts snake_case field names to camelCase.
// Returns empty string if schema or field not found.
func GetOpenAPIDescription(resourceName, tfFieldName string) string {
	if openAPIDoc == nil || openAPIDoc.Components == nil {
		return ""
	}

	// Check for special fields first
	switch tfFieldName {
	case "if_match":
		// Header parameter
		if openAPIDoc.Components.Parameters != nil {
			if paramRef, ok := openAPIDoc.Components.Parameters["If-Match"]; ok && paramRef.Value != nil {
				if paramRef.Value.Description != "" {
					return strings.TrimSpace(paramRef.Value.Description)
				}
			}
		}

	case "etag":
		// Response header
		if openAPIDoc.Components.Headers != nil {
			if headerRef, ok := openAPIDoc.Components.Headers["ETag"]; ok && headerRef.Value != nil {
				if headerRef.Value.Description != "" {
					return strings.TrimSpace(headerRef.Value.Description)
				}
			}
		}

	case "audit":
		// Schema reference
		if openAPIDoc.Components.Schemas != nil {
			if schemaRef, ok := openAPIDoc.Components.Schemas["CouchbaseAuditData"]; ok && schemaRef.Value != nil {
				return "Couchbase audit data."
			}
		}

	default:
		// Check for fields that map to string-type schemas
		// These are fields where the entire field value IS the schema type
		stringSchemaMap := map[string]string{
			"access_control_function": "AccessFunction",
			"import_filter":           "ImportFilter",
		}

		if schemaName, ok := stringSchemaMap[tfFieldName]; ok {
			if schemaRef, ok := openAPIDoc.Components.Schemas[schemaName]; ok && schemaRef.Value != nil {
				if schemaRef.Value.Description != "" {
					return buildEnhancedDescription(schemaRef.Value, openAPIDoc)
				}
			}
		}

		// Check if this is a path parameter (e.g., organization_id, project_id, app_endpoint_name)
		if openAPIDoc.Components.Parameters != nil {
			// Try different parameter name patterns
			paramPatterns := []string{
				snakeToCapitalizedCamel(tfFieldName), // e.g., organization_id -> OrganizationId
			}

			// Special mappings for common parameters
			switch tfFieldName {
			case "app_endpoint_name":
				paramPatterns = append(paramPatterns, "appEndpointId", "appEndpointKeyspace")
			case "scope":
				paramPatterns = append(paramPatterns, "scopeName", "appEndpointKeyspace")
			case "collection":
				paramPatterns = append(paramPatterns, "collectionName", "appEndpointKeyspace")
			}

			for _, paramName := range paramPatterns {
				if paramRef, ok := openAPIDoc.Components.Parameters[paramName]; ok && paramRef.Value != nil {
					if paramRef.Value.Description != "" {
						return strings.TrimSpace(paramRef.Value.Description)
					}
				}
			}
		}
	}

	if openAPIDoc.Components.Schemas == nil {
		return ""
	}

	// Convert snake_case to camelCase for schema property lookup
	camelFieldName := snakeToCamel(tfFieldName)

	// Capitalize resource name for schema patterns
	capitalizedResource := capitalize(resourceName)

	// Try common schema name patterns
	schemaPatterns := []string{
		capitalizedResource, // Exact match (e.g., CORSConfig, AccessFunction)
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
			return buildEnhancedDescription(propRef.Value, openAPIDoc)
		}
	}

	// If not found in main schema, try common nested schemas
	// This handles fields like "type", "roles" that are inside nested Resource objects
	nestedSchemas := []string{
		"Resource",                       // For user resources, API key resources
		"ResourceBucket",                 // For bucket-specific resources
		capitalizedResource + "Resource", // e.g., UserResource if it exists
	}

	for _, schemaName := range nestedSchemas {
		schemaRef := openAPIDoc.Components.Schemas[schemaName]
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		propRef := schemaRef.Value.Properties[camelFieldName]
		if propRef != nil && propRef.Value != nil {
			return buildEnhancedDescription(propRef.Value, openAPIDoc)
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

// snakeToCapitalizedCamel converts snake_case to CapitalizedCamelCase (PascalCase)
// Used for OpenAPI parameter names: organization_id → OrganizationId
func snakeToCapitalizedCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	// Capitalize first letter of all parts
	result := ""
	for _, part := range parts {
		result += capitalize(part)
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

// buildEnhancedDescription creates a rich markdown description from OpenAPI property.
// It resolves references, extracts metadata, and formats the output as markdown bullets.
func buildEnhancedDescription(prop *openapi3.Schema, doc *openapi3.T) string {
	referencedSchema := resolveArrayReference(prop, doc)
	description := extractDescription(prop, referencedSchema)
	enumValues := collectEnumValues(prop, referencedSchema)

	var parts []string

	// Add main description or enum values if no description exists
	if description != "" {
		parts = append(parts, formatDescriptionBullet(description))
	} else if len(enumValues) > 0 {
		parts = append(parts, formatValidValuesBullet(enumValues))
	}

	// Add constraints
	if constraints := buildConstraints(prop); len(constraints) > 0 {
		parts = append(parts, formatConstraintsBullet(constraints))
	}

	// Add enum values if description exists (avoid duplication)
	if description != "" && len(enumValues) > 0 {
		parts = append(parts, formatValidValuesBullet(enumValues))
	}

	// Add default value
	if prop.Default != nil {
		parts = append(parts, formatDefaultBullet(prop.Default))
	}

	// Add format information
	if prop.Format != "" {
		if formatDesc := getFormatDescription(prop.Format); formatDesc != "" {
			parts = append(parts, formatFormatBullet(formatDesc))
		}
	}

	// Add deprecation warning
	if prop.Deprecated {
		parts = append(parts, formatDeprecationBullet())
	}

	return strings.Join(parts, "")
}

// resolveArrayReference extracts the referenced schema for array types with items.$ref.
// Returns nil if the property is not an array or has no reference.
func resolveArrayReference(prop *openapi3.Schema, doc *openapi3.T) *openapi3.Schema {
	if doc == nil || prop.Type == nil || !prop.Type.Is("array") {
		return nil
	}

	if prop.Items == nil || prop.Items.Ref == "" {
		return nil
	}

	// Extract schema name from $ref (e.g., "#/components/schemas/OrganizationRoles" -> "OrganizationRoles")
	refParts := strings.Split(prop.Items.Ref, "/")
	if len(refParts) == 0 {
		return nil
	}

	schemaName := refParts[len(refParts)-1]
	if schemaRef := doc.Components.Schemas[schemaName]; schemaRef != nil && schemaRef.Value != nil {
		return schemaRef.Value
	}

	return nil
}

// extractDescription gets the description from the property or falls back to referenced schema.
// It cleans the description by removing markdown list formatting.
func extractDescription(prop *openapi3.Schema, referencedSchema *openapi3.Schema) string {
	if prop.Description != "" {
		return cleanDescription(prop.Description)
	}

	if referencedSchema != nil && referencedSchema.Description != "" {
		return cleanDescription(referencedSchema.Description)
	}

	return ""
}

// cleanDescription removes markdown list formatting from description text while
// preserving important markdown structures like tables.
//
// Behavior:
//   - Removes list markers (-, *, +) and numbered lists (1., 2., etc.)
//   - Joins non-table text into paragraphs
//   - Preserves markdown tables with proper spacing
//   - Ensures tables start on a new line
func cleanDescription(desc string) string {
	if desc = strings.TrimSpace(desc); desc == "" {
		return ""
	}

	lines := strings.Split(desc, "\n")
	cleanedLines := make([]string, 0, len(lines))
	textBuffer := make([]string, 0)
	inTable := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		isTable := isTableRow(line)

		switch {
		case isTable && !inTable:
			// Entering table: flush text buffer, add spacing, start table
			cleanedLines = appendTextBuffer(cleanedLines, textBuffer)
			textBuffer = textBuffer[:0] // Reset buffer
			cleanedLines = append(cleanedLines, "", line)
			inTable = true

		case !isTable && inTable:
			// Exiting table: add spacing if next line has content
			inTable = false
			if line != "" {
				cleanedLines = append(cleanedLines, "")
			}
			// Fall through to process the line as text

		case inTable:
			// Inside table: preserve line as-is
			cleanedLines = append(cleanedLines, line)
			continue
		}

		// Process non-table text
		if !inTable && line != "" {
			cleaned := removeListMarkers(line)
			if cleaned != "" {
				textBuffer = append(textBuffer, cleaned)
			}
		}
	}

	// Flush any remaining text buffer
	cleanedLines = appendTextBuffer(cleanedLines, textBuffer)

	return strings.Join(cleanedLines, "\n")
}

// isTableRow determines if a line is part of a markdown table.
// Tables are identified by having pipe characters (|) and at least 2 of them.
// This simple heuristic works for most OpenAPI descriptions.
func isTableRow(line string) bool {
	// Count pipe characters - tables typically have at least 2 pipes per row
	return strings.Count(line, "|") >= 2
}

// removeListMarkers removes markdown list markers from the beginning of a line.
// Handles:
//   - Unordered lists: -, *, +
//   - Ordered lists: 1., 2., 3., 42., etc.
func removeListMarkers(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	// Remove unordered list markers (-, *, +)
	if len(line) > 0 && isListMarker(line[0]) {
		line = strings.TrimSpace(line[1:])
		return line
	}

	// Remove ordered list markers (1., 2., 42., etc.)
	// Find where the number ends and the period begins
	if len(line) > 0 && isDigit(line[0]) {
		dotIndex := strings.Index(line, ".")
		if dotIndex > 0 && dotIndex < len(line)-1 {
			// Check if all characters before the dot are digits
			allDigits := true
			for i := 0; i < dotIndex; i++ {
				if !isDigit(line[i]) {
					allDigits = false
					break
				}
			}
			if allDigits {
				line = strings.TrimSpace(line[dotIndex+1:])
			}
		}
	}

	return line
}

// appendTextBuffer flushes the text buffer by joining it into a paragraph
// and appending to the result. Returns the result with the flushed text.
func appendTextBuffer(result []string, buffer []string) []string {
	if len(buffer) > 0 {
		result = append(result, strings.Join(buffer, " "))
	}
	return result
}

// collectEnumValues gets enum values from the property or referenced schema.
// Prefers the property's own enums over the referenced schema.
func collectEnumValues(prop *openapi3.Schema, referencedSchema *openapi3.Schema) []interface{} {
	if len(prop.Enum) > 0 {
		return prop.Enum
	}

	if referencedSchema != nil && len(referencedSchema.Enum) > 0 {
		return referencedSchema.Enum
	}

	return nil
}

// formatEnumValues converts enum values to backtick-quoted strings.
func formatEnumValues(enumValues []interface{}) []string {
	formatted := make([]string, len(enumValues))
	for i, val := range enumValues {
		formatted[i] = fmt.Sprintf("`%v`", val)
	}
	return formatted
}

// Formatting helper functions for consistent markdown output

func formatDescriptionBullet(desc string) string {
	return fmt.Sprintf("\n - %s", desc)
}

func formatConstraintsBullet(constraints []string) string {
	return fmt.Sprintf("\n - **Constraints**: %s", strings.Join(constraints, ", "))
}

func formatValidValuesBullet(enumValues []interface{}) string {
	formatted := formatEnumValues(enumValues)
	return fmt.Sprintf("\n - **Valid Values**: %s", strings.Join(formatted, ", "))
}

func formatDefaultBullet(defaultVal interface{}) string {
	return fmt.Sprintf("\n - **Default**: `%v`", defaultVal)
}

func formatFormatBullet(formatDesc string) string {
	return fmt.Sprintf("\n - **Format**: %s", formatDesc)
}

func formatDeprecationBullet() string {
	return "\n - **Deprecated**: This field is deprecated and will be removed in a future release."
}

// Helper functions for character checking

func isListMarker(c byte) bool {
	return c == '-' || c == '*' || c == '+'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
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
