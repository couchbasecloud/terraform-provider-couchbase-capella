// Package enums exposes OpenAPI-spec enum metadata used by the schema
// builder to auto-attach OneOf validators. The generated file
// (enums.gen.go) provides the underlying enumTable; this file defines the
// types and the Lookup function that consumers call.
package enums

import "strings"

// SchemaBuilder is the subset of *capellaschema.SchemaBuilder needed by
// Lookup. Defined here so the schema package can import enums without
// creating an import cycle.
type SchemaBuilder interface {
	GetOpenAPISchemaName() string
	GetResourceName() string
}

// EnumDef is defined in the generated enums.gen.go file (kept there so
// the generated file is self-contained for tools that analyse it in
// isolation, e.g. golangci-lint's typecheck pass).

// Lookup returns the enum definition associated with a Terraform
// attribute on the given builder, walking the same schema-name pattern
// list used by docs.GetOpenAPIDescription. Returns nil when no enum is
// associated.
func Lookup(b SchemaBuilder, alternateSchemas []string, tfFieldName string) *EnumDef {
	field := snakeToCamel(tfFieldName)

	patterns := append([]string(nil), alternateSchemas...)

	// Patterns from the OpenAPI schema name. Capitalize the seed (same as
	// docs.GetOpenAPIDescription) so builders that pass lower-camel names
	// like "allowedCidr" still resolve to PascalCase enumTable keys.
	patterns = append(patterns, schemaPatterns(b.GetOpenAPISchemaName())...)

	// Patterns from the Terraform resource name, in case it differs from the
	// OpenAPI schema name (e.g. NewSchemaBuilder("allowlist", "allowedCidr")
	// where the enum lives under "Allowlist" rather than "AllowedCidr").
	if b.GetResourceName() != b.GetOpenAPISchemaName() {
		patterns = append(patterns, schemaPatterns(b.GetResourceName())...)
	}

	for _, p := range patterns {
		if p == "" {
			continue
		}
		if def, ok := enumTable[p][field]; ok {
			out := def
			return &out
		}
	}
	return nil
}

func schemaPatterns(name string) []string {
	if name == "" {
		return nil
	}
	capName := capitalize(name)
	return []string{
		capName,
		"Create" + capName + "Request",
		"Get" + capName + "Response",
		"Update" + capName + "Request",
		capName + "Request",
		capName + "Response",
		name,
	}
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += capitalize(p)
	}
	return out
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
