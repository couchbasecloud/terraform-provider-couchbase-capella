package docs

import (
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var openAPIDoc *openapi3.T

func init() {
	// Get OpenAPI spec path from environment variable, or assume project root
	openAPIPath := os.Getenv("CAPELLA_OPENAPI_SPEC_PATH")
	if openAPIPath == "" {
		openAPIPath = "openapi.generated.yaml"
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
		// Check if this is a path parameter (e.g., organization_id, project_id)
		if strings.HasSuffix(tfFieldName, "_id") && openAPIDoc.Components.Parameters != nil {
			paramName := snakeToCapitalizedCamel(tfFieldName)
			if paramRef, ok := openAPIDoc.Components.Parameters[paramName]; ok && paramRef.Value != nil {
				if paramRef.Value.Description != "" {
					return strings.TrimSpace(paramRef.Value.Description)
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
