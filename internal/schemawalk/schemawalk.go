// Package schemawalk provides helpers for resolving Terraform attribute names
// to OpenAPI schema locations. Both the description loader (internal/docs) and
// the validator auto-attacher (internal/schema) walk the same candidate schema
// names; this package is the single home for that shared traversal logic.
package schemawalk

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/enums"
)

// SnakeToCamel converts a snake_case Terraform field name to the camelCase
// property name used in OpenAPI schemas.
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += Capitalize(parts[i])
	}
	return result
}

// SnakeToPascal converts snake_case to PascalCase (CapitalizedCamelCase).
// Used to map Terraform parameter names to OpenAPI parameter names
// (e.g. organization_id → OrganizationId).
func SnakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	result := ""
	for _, part := range parts {
		result += Capitalize(part)
	}
	return result
}

// Capitalize uppercases the first letter of s.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// SchemaNameCandidates returns the ordered list of OpenAPI schema names to try
// for a given resource name. It tries the capitalized form first, then the
// standard Create/Get/Update request/response patterns, with the original name
// as a final fallback for lowercase schema names (e.g. "datadog").
func SchemaNameCandidates(name string) []string {
	cap := Capitalize(name)
	seen := make(map[string]bool, 7)
	candidates := make([]string, 0, 7)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			candidates = append(candidates, s)
		}
	}
	add(cap)
	add("Create" + cap + "Request")
	add("Get" + cap + "Response")
	add("Update" + cap + "Request")
	add(cap + "Request")
	add(cap + "Response")
	add(name)
	return candidates
}

// RefName extracts the schema name from a $ref string like
// "#/components/schemas/Foo" → "Foo". Returns "" for unrecognised formats.
func RefName(ref string) string {
	const prefix = "#/components/schemas/"
	if strings.HasPrefix(ref, prefix) {
		return ref[len(prefix):]
	}
	return ""
}

// ResolveArrayRef returns the referenced schema for an array property whose
// items use a $ref. Returns nil when the property is not an array or has no ref.
func ResolveArrayRef(schemas openapi3.Schemas, prop *openapi3.Schema) *openapi3.Schema {
	if prop.Type == nil || !prop.Type.Is("array") {
		return nil
	}
	if prop.Items == nil || prop.Items.Ref == "" {
		return nil
	}
	name := RefName(prop.Items.Ref)
	if name == "" {
		return nil
	}
	if ref := schemas[name]; ref != nil && ref.Value != nil {
		return ref.Value
	}
	return nil
}

// FindProperty searches schema for a property named fieldName, following $ref
// links and allOf/anyOf/oneOf compositions. schemas is the components.schemas
// map used to resolve $ref references.
//
// To avoid nondeterminism from map iteration, it collects all matches and only
// returns a result when exactly one unique schema pointer is found — this
// prevents attaching the wrong description when the same field name appears in
// multiple nested schemas.
//
// Array items are not searched to avoid false positives from unrelated objects.
func FindProperty(schemas openapi3.Schemas, schema *openapi3.Schema, fieldName string) *openapi3.Schema {
	return findProperty(schemas, schema, fieldName, make(map[string]bool))
}

func findProperty(schemas openapi3.Schemas, schema *openapi3.Schema, fieldName string, visited map[string]bool) *openapi3.Schema {
	if schema == nil {
		return nil
	}

	// Direct property match — deterministic, checked first.
	if propRef, ok := schema.Properties[fieldName]; ok && propRef != nil && propRef.Value != nil {
		return propRef.Value
	}

	var matches []*openapi3.Schema

	for _, propRef := range schema.Properties {
		if found := searchRef(schemas, propRef, fieldName, visited); found != nil {
			matches = append(matches, found)
		}
	}

	for _, list := range [][]*openapi3.SchemaRef{schema.AllOf, schema.AnyOf, schema.OneOf} {
		for _, ref := range list {
			if found := searchRef(schemas, ref, fieldName, visited); found != nil {
				matches = append(matches, found)
			}
		}
	}

	if len(matches) == 1 {
		return matches[0]
	}
	if len(matches) > 1 {
		first := matches[0]
		for _, m := range matches[1:] {
			if m != first {
				return nil
			}
		}
		return first
	}
	return nil
}

func searchRef(schemas openapi3.Schemas, ref *openapi3.SchemaRef, fieldName string, visited map[string]bool) *openapi3.Schema {
	if ref == nil {
		return nil
	}
	if ref.Ref != "" {
		name := RefName(ref.Ref)
		if name == "" || visited[name] {
			return nil
		}
		visited[name] = true
		refSchema := schemas[name]
		if refSchema == nil || refSchema.Value == nil {
			return nil
		}
		return findProperty(schemas, refSchema.Value, fieldName, visited)
	}
	if ref.Value != nil && ref.Value.Type != nil && ref.Value.Type.Is("object") {
		return findProperty(schemas, ref.Value, fieldName, visited)
	}
	return nil
}

// EnumValues returns the string enum values for a Terraform attribute by walking
// the same schema-name candidates as SchemaNameCandidates. alternateSchemas are
// tried first as exact names before the pattern expansion of openAPISchemaName
// and resourceName.
//
// Returns nil if no enum is found for the attribute.
func EnumValues(openAPISchemaName, resourceName string, alternateSchemas []string, tfFieldName string) []string {
	camelField := SnakeToCamel(tfFieldName)
	seen := make(map[string]bool)

	try := func(name string) []string {
		if name == "" || seen[name] {
			return nil
		}
		seen[name] = true
		if inner, ok := enums.Table[name]; ok {
			if vals, ok := inner[camelField]; ok {
				return vals
			}
		}
		return nil
	}

	for _, alt := range alternateSchemas {
		if vals := try(alt); vals != nil {
			return vals
		}
	}

	for _, base := range []string{openAPISchemaName, resourceName} {
		for _, candidate := range SchemaNameCandidates(base) {
			if vals := try(candidate); vals != nil {
				return vals
			}
		}
	}

	return nil
}
