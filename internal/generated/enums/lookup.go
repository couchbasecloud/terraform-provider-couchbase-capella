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
	candidates := fieldCandidates(tfFieldName)

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

	// Resolve nested object fields keyed by OpenAPI path, e.g.
	// "weeklySchedule.dayOfWeek" under CreateScheduledBackupRequest. Dotted
	// candidates are tried before the bare name so a hinted nested field never
	// resolves to a same-named top-level field in the same schema.
	candidates = append(nestedFieldCandidates(alternateSchemas, candidates), candidates...)

	for _, p := range patterns {
		if p == "" {
			continue
		}
		for _, field := range candidates {
			if def, ok := enumTable[p][field]; ok {
				out := def
				return &out
			}
		}
	}
	return nil
}

func schemaPatterns(name string) []string {
	if name == "" {
		return nil
	}
	capName := capitalize(name)
	acroName := capitalizeAcronym(name)

	patterns := []string{
		capName,
		"Create" + capName + "Request",
		"Get" + capName + "Response",
		"Update" + capName + "Request",
		capName + "Request",
		capName + "Response",
		name,
	}
	if acroName != capName {
		patterns = append(patterns,
			acroName,
			"Create"+acroName+"Request",
			"Get"+acroName+"Response",
			"Update"+acroName+"Request",
			acroName+"Request",
			acroName+"Response",
		)
	}
	return patterns
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

// snakeToCamelAcronym is like snakeToCamel but uppercases segments that match
// known acronyms (id → ID, url → URL, etc.). The OpenAPI spec uses both styles
// inconsistently (e.g. `accountId` vs `vpcNetworkID`), so callers should try
// both this and snakeToCamel when looking a field up.
func snakeToCamelAcronym(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	out := parts[0]
	for _, p := range parts[1:] {
		if upper, ok := acronyms[p]; ok {
			out += upper
		} else {
			out += capitalize(p)
		}
	}
	return out
}

// acronyms maps lowercase segments to the uppercase form the OpenAPI spec
// uses when that segment is treated as an acronym. Used both for field names
// (e.g. `vpc_network_id` → `vpcNetworkID`) and for the leading word of schema
// builder names (e.g. `gcpPrivateEndpointCommand` → `GCPPrivateEndpointCommand`).
var acronyms = map[string]string{
	"api":  "API",
	"arn":  "ARN",
	"aws":  "AWS",
	"cidr": "CIDR",
	"cmek": "CMEK",
	"cors": "CORS",
	"gcp":  "GCP",
	"guid": "GUID",
	"id":   "ID",
	"ids":  "IDs",
	"ip":   "IP",
	"ips":  "IPs",
	"json": "JSON",
	"oidc": "OIDC",
	"sso":  "SSO",
	"tls":  "TLS",
	"ttl":  "TTL",
	"uri":  "URI",
	"url":  "URL",
	"uuid": "UUID",
	"vpc":  "VPC",
}

// capitalizeAcronym returns name with its leading lowercase prefix uppercased
// when that prefix matches a known acronym. Falls back to capitalize() otherwise.
// Example: "gcpPrivateEndpointCommand" → "GCPPrivateEndpointCommand";
// "projectFoo" → "ProjectFoo".
func capitalizeAcronym(name string) string {
	if name == "" {
		return name
	}
	i := 0
	for i < len(name) && name[i] >= 'a' && name[i] <= 'z' {
		i++
	}
	if i > 0 {
		if upper, ok := acronyms[name[:i]]; ok {
			return upper + name[i:]
		}
	}
	return capitalize(name)
}

// fieldCandidates returns possible camelCase variants of a snake_case TF field
// name. The first is the canonical snakeToCamel result; the second (only when
// it differs) is a variant where trailing acronym segments are uppercased.
func fieldCandidates(tfFieldName string) []string {
	base := snakeToCamel(tfFieldName)
	upper := snakeToCamelAcronym(tfFieldName)
	if upper == base {
		return []string{base}
	}
	return []string{base, upper}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func lowerCamel(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// nestedFieldCandidates prefixes each base field with the lower-camel of each
// alternate schema (the flattened nested object's name), yielding dotted keys
// like "weeklySchedule.dayOfWeek" that match the generated table.
func nestedFieldCandidates(alternateSchemas, base []string) []string {
	var out []string
	for _, s := range alternateSchemas {
		prefix := lowerCamel(s)
		if prefix == "" {
			continue
		}
		for _, f := range base {
			out = append(out, prefix+"."+f)
		}
	}
	return out
}

// CompositionLookup returns the composition definition (oneOf/anyOf/allOf)
// associated with a Terraform attribute on the given builder, walking the
// same schema-name pattern list used by Lookup and docs.GetOpenAPIDescription.
// Returns nil when no composition is associated with the attribute.
func CompositionLookup(b SchemaBuilder, alternateSchemas []string, tfFieldName string) *CompositionDef {
	candidates := fieldCandidates(tfFieldName)

	patterns := append([]string(nil), alternateSchemas...)

	// Patterns from the OpenAPI schema name
	patterns = append(patterns, schemaPatterns(b.GetOpenAPISchemaName())...)

	// Patterns from the Terraform resource name, in case it differs
	if b.GetResourceName() != b.GetOpenAPISchemaName() {
		patterns = append(patterns, schemaPatterns(b.GetResourceName())...)
	}

	for _, p := range patterns {
		if p == "" {
			continue
		}
		for _, field := range candidates {
			if def, ok := compositionTable[p][field]; ok {
				out := def
				return &out
			}
		}
	}
	return nil
}

// RequiredLookup returns true if the given Terraform attribute is marked
// as required in the OpenAPI spec, walking the same schema-name pattern
// list used by Lookup and docs.GetOpenAPIDescription.
func RequiredLookup(b SchemaBuilder, alternateSchemas []string, tfFieldName string) bool {
	candidates := fieldCandidates(tfFieldName)

	patterns := append([]string(nil), alternateSchemas...)

	// Patterns from the OpenAPI schema name
	patterns = append(patterns, schemaPatterns(b.GetOpenAPISchemaName())...)

	// Patterns from the Terraform resource name, in case it differs
	if b.GetResourceName() != b.GetOpenAPISchemaName() {
		patterns = append(patterns, schemaPatterns(b.GetResourceName())...)
	}

	for _, p := range patterns {
		if p == "" {
			continue
		}
		for _, field := range candidates {
			if _, ok := requiredTable[p][field]; ok {
				return true
			}
		}
	}
	return false
}

// ConstraintLookup returns the constraint definition (min/max values, lengths,
// item counts) associated with a Terraform attribute on the given builder,
// walking the same schema-name pattern list used by Lookup.
// Returns nil when no constraints are associated with the attribute.
func ConstraintLookup(b SchemaBuilder, alternateSchemas []string, tfFieldName string) *ConstraintDef {
	candidates := fieldCandidates(tfFieldName)

	patterns := append([]string(nil), alternateSchemas...)

	// Patterns from the OpenAPI schema name
	patterns = append(patterns, schemaPatterns(b.GetOpenAPISchemaName())...)

	// Patterns from the Terraform resource name, in case it differs
	if b.GetResourceName() != b.GetOpenAPISchemaName() {
		patterns = append(patterns, schemaPatterns(b.GetResourceName())...)
	}

	for _, p := range patterns {
		if p == "" {
			continue
		}
		for _, field := range candidates {
			if def, ok := constraintTable[p][field]; ok {
				out := def
				return &out
			}
		}
	}
	return nil
}
