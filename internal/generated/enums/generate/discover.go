package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type scope string

const (
	scopeParameter scope = "parameter"
	scopeSchema    scope = "schema"
)

type enumSite struct {
	ID         string
	Scope      scope
	SchemaName string
	FieldPath  string
	Type       string
	Values     []string
	SourcePath string
}

type compositionKind string

const (
	kindOneOf compositionKind = "oneOf"
	kindAnyOf compositionKind = "anyOf"
	kindAllOf compositionKind = "allOf"
)

type compositionSite struct {
	SchemaName string
	FieldPath  string
	Kind       compositionKind
	Branches   []string
	SourcePath string
}

type requiredSite struct {
	SchemaName string
	FieldPath  string
	SourcePath string
}

type constraintSite struct {
	SchemaName string
	FieldPath  string
	Minimum    *float64
	Maximum    *float64
	MinLength  *int64
	MaxLength  *int64
	MinItems   *int64
	MaxItems   *int64
	SourcePath string
}

func discover(specPath string) ([]enumSite, error) {
	enums, _, _, _, err := discoverAll(specPath)
	return enums, err
}

func discoverAll(specPath string) ([]enumSite, []compositionSite, []requiredSite, []constraintSite, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("load spec: %w", err)
	}

	w := &walker{doc: doc}
	w.run()

	sites, err := dedupByID(w.sites)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	sort.Slice(sites, func(i, j int) bool {
		a, b := sites[i], sites[j]
		if a.Scope != b.Scope {
			return a.Scope < b.Scope
		}
		if a.SchemaName != b.SchemaName {
			return a.SchemaName < b.SchemaName
		}
		return a.FieldPath < b.FieldPath
	})

	compSites := dedupComposition(w.compositionSites)
	sort.Slice(compSites, func(i, j int) bool {
		a, b := compSites[i], compSites[j]
		if a.SchemaName != b.SchemaName {
			return a.SchemaName < b.SchemaName
		}
		return a.FieldPath < b.FieldPath
	})

	reqSites := dedupRequired(w.requiredSites)
	sort.Slice(reqSites, func(i, j int) bool {
		a, b := reqSites[i], reqSites[j]
		if a.SchemaName != b.SchemaName {
			return a.SchemaName < b.SchemaName
		}
		return a.FieldPath < b.FieldPath
	})

	constrSites := dedupConstraints(w.constraintSites)
	sort.Slice(constrSites, func(i, j int) bool {
		a, b := constrSites[i], constrSites[j]
		if a.SchemaName != b.SchemaName {
			return a.SchemaName < b.SchemaName
		}
		return a.FieldPath < b.FieldPath
	})

	return sites, compSites, reqSites, constrSites, nil
}

func dedupComposition(sites []compositionSite) []compositionSite {
	type key struct {
		schema, field string
		kind          compositionKind
	}
	seen := make(map[key]int, len(sites))
	out := make([]compositionSite, 0, len(sites))
	for _, s := range sites {
		k := key{s.SchemaName, s.FieldPath, s.Kind}
		if idx, ok := seen[k]; ok {
			existing := &out[idx]
			existing.Branches = mergeBranches(existing.Branches, s.Branches)
			continue
		}
		seen[k] = len(out)
		out = append(out, s)
	}
	return out
}

func mergeBranches(a, b []string) []string {
	seen := make(map[string]struct{}, len(a))
	for _, v := range a {
		seen[v] = struct{}{}
	}
	out := append([]string(nil), a...)
	for _, v := range b {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func dedupRequired(sites []requiredSite) []requiredSite {
	type key struct {
		schema, field string
	}
	seen := make(map[key]bool, len(sites))
	out := make([]requiredSite, 0, len(sites))
	for _, s := range sites {
		k := key{s.SchemaName, s.FieldPath}
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, s)
	}
	return out
}

func dedupConstraints(sites []constraintSite) []constraintSite {
	type key struct {
		schema, field string
	}
	seen := make(map[key]int, len(sites))
	out := make([]constraintSite, 0, len(sites))
	for _, s := range sites {
		k := key{s.SchemaName, s.FieldPath}
		idx, ok := seen[k]
		if !ok {
			seen[k] = len(out)
			out = append(out, s)
			continue
		}
		mergeConstraintSite(&out[idx], s)
	}
	return out
}

// mergeConstraintSite copies non-nil constraint fields from src onto dst,
// preserving any value dst already has. Used to combine multiple discoveries
// of the same (schema, field) pair into a single site.
func mergeConstraintSite(dst *constraintSite, src constraintSite) {
	if dst.Minimum == nil {
		dst.Minimum = src.Minimum
	}
	if dst.Maximum == nil {
		dst.Maximum = src.Maximum
	}
	if dst.MinLength == nil {
		dst.MinLength = src.MinLength
	}
	if dst.MaxLength == nil {
		dst.MaxLength = src.MaxLength
	}
	if dst.MinItems == nil {
		dst.MinItems = src.MinItems
	}
	if dst.MaxItems == nil {
		dst.MaxItems = src.MaxItems
	}
}

// dedupByID merges sites that share an ID — emitted by composition keyword
// branches (allOf/oneOf/anyOf) that converge on the same logical field. Values
// are unioned in first-seen order. Conflicting Types are a spec error.
func dedupByID(sites []enumSite) ([]enumSite, error) {
	byID := make(map[string]int, len(sites))
	out := make([]enumSite, 0, len(sites))
	for _, s := range sites {
		idx, ok := byID[s.ID]
		if !ok {
			byID[s.ID] = len(out)
			out = append(out, s)
			continue
		}
		existing := &out[idx]
		if existing.Type != s.Type {
			return nil, fmt.Errorf("conflicting enum types for %s: %q at %s vs %q at %s",
				s.ID, existing.Type, existing.SourcePath, s.Type, s.SourcePath)
		}
		existing.Values = mergeValues(existing.Values, s.Values)
	}
	return out, nil
}

func mergeValues(a, b []string) []string {
	seen := make(map[string]struct{}, len(a))
	for _, v := range a {
		seen[v] = struct{}{}
	}
	out := append([]string(nil), a...)
	for _, v := range b {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

type walker struct {
	doc              *openapi3.T
	sites            []enumSite
	compositionSites []compositionSite
	requiredSites    []requiredSite
	constraintSites  []constraintSite
}

func (w *walker) run() {
	if w.doc.Components == nil {
		return
	}
	for name, ref := range w.doc.Components.Schemas {
		if ref == nil || ref.Value == nil {
			continue
		}
		w.schema(name, ref.Value, name, "", "components.schemas."+name, scopeSchema)
	}
	for name, ref := range w.doc.Components.Parameters {
		if ref == nil || ref.Value == nil || ref.Value.Schema == nil {
			continue
		}
		s := ref.Value.Schema
		if s.Ref != "" || s.Value == nil {
			continue
		}
		w.schema("Param_"+name, s.Value, name, "", "components.parameters."+name+".schema", scopeParameter)
	}
}

// schema visits a single schema node. Property $refs are skipped — those schemas
// are visited independently as top-level entries in run. Composition keywords
// (allOf/oneOf/anyOf) are traversed transparently: branches contribute to the
// same logical field, and dedupByID merges any colliding sites at the end.
func (w *walker) schema(id string, s *openapi3.Schema, schemaName, fieldPath, sourcePath string, sc scope) {
	if s == nil {
		return
	}
	if len(s.Enum) > 0 {
		w.sites = append(w.sites, enumSite{
			ID:         id,
			Scope:      sc,
			SchemaName: schemaName,
			FieldPath:  fieldPath,
			Type:       typeOf(s),
			Values:     enumValues(s.Enum),
			SourcePath: sourcePath,
		})
		return
	}
	w.composition(id, s.AllOf, schemaName, fieldPath, sourcePath, sc, "allOf")
	w.composition(id, s.OneOf, schemaName, fieldPath, sourcePath, sc, "oneOf")
	w.composition(id, s.AnyOf, schemaName, fieldPath, sourcePath, sc, "anyOf")
	// Only record composition sites for schema-scope (not parameters)
	if sc == scopeSchema {
		w.recordComposition(s.OneOf, schemaName, fieldPath, sourcePath, kindOneOf)
		w.recordComposition(s.AnyOf, schemaName, fieldPath, sourcePath, kindAnyOf)
		w.recordComposition(s.AllOf, schemaName, fieldPath, sourcePath, kindAllOf)
	}
	// Record required fields from the schema's required array
	if sc == scopeSchema && len(s.Required) > 0 {
		w.recordRequired(s.Required, schemaName, fieldPath, sourcePath)
	}
	// Record min/max constraints for this schema node
	if sc == scopeSchema {
		w.recordConstraints(s, schemaName, fieldPath, sourcePath)
	}
	for fieldName, propRef := range s.Properties {
		if propRef == nil || propRef.Ref != "" || propRef.Value == nil {
			continue
		}
		w.schema(
			id+"_"+toPascal(fieldName), propRef.Value,
			schemaName, joinPath(fieldPath, fieldName),
			sourcePath+".properties."+fieldName, sc,
		)
	}
	w.items(id, s, schemaName, fieldPath, sourcePath, sc)
}

func (w *walker) composition(id string, refs openapi3.SchemaRefs, schemaName, fieldPath, sourcePath string, sc scope, kw string) {
	for i, ref := range refs {
		if ref == nil || ref.Ref != "" || ref.Value == nil {
			continue
		}
		w.schema(id, ref.Value, schemaName, fieldPath,
			fmt.Sprintf("%s.%s[%d]", sourcePath, kw, i), sc)
	}
}

// recordComposition records a composition site only when ALL branches are $refs
// to named schemas. Mixed compositions (some inline, some $ref) are skipped
// entirely to avoid incomplete branch metadata.
func (w *walker) recordComposition(refs openapi3.SchemaRefs, schemaName, fieldPath, sourcePath string, kind compositionKind) {
	if len(refs) == 0 {
		return
	}

	branches := make([]string, 0, len(refs))
	for _, ref := range refs {
		// Skip nil refs
		if ref == nil {
			continue
		}
		// If any branch is not a $ref (inline schema), skip the entire composition
		if ref.Ref == "" {
			return
		}
		name := extractSchemaName(ref.Ref)
		// If $ref doesn't point to a component schema, skip the entire composition
		if name == "" {
			return
		}
		branches = append(branches, name)
	}

	if len(branches) == 0 {
		return
	}

	w.compositionSites = append(w.compositionSites, compositionSite{
		SchemaName: schemaName,
		FieldPath:  fieldPath,
		Kind:       kind,
		Branches:   branches,
		SourcePath: sourcePath + "." + string(kind),
	})
}

func extractSchemaName(ref string) string {
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return ""
	}
	return strings.TrimPrefix(ref, prefix)
}

// recordRequired captures required field names from a schema's required array.
// Each required field is stored as a requiredSite for later code generation.
// Only top-level required fields are recorded (parentFieldPath must be empty)
// because RequiredLookup matches by single field name, not dot-paths.
func (w *walker) recordRequired(required []string, schemaName, parentFieldPath, sourcePath string) {
	// Skip nested required fields - they can't be matched by RequiredLookup
	// which only looks up single Terraform field keys.
	if parentFieldPath != "" {
		return
	}
	for _, fieldName := range required {
		w.requiredSites = append(w.requiredSites, requiredSite{
			SchemaName: schemaName,
			FieldPath:  fieldName,
			SourcePath: sourcePath + ".required[" + fieldName + "]",
		})
	}
}

// recordConstraints captures min/max constraints (minimum, maximum, minLength,
// maxLength, minItems, maxItems) from a schema node. Only records when at least
// one constraint is present and fieldPath is not empty (top-level schema-only
// constraints are skipped since they apply to the schema itself, not a field).
func (w *walker) recordConstraints(s *openapi3.Schema, schemaName, fieldPath, sourcePath string) {
	// Skip top-level schemas - constraints only make sense for fields
	if fieldPath == "" {
		return
	}
	// Only record if at least one constraint is present
	if s.Min == nil && s.Max == nil && s.MinLength == 0 && s.MaxLength == nil && s.MinItems == 0 && s.MaxItems == nil {
		return
	}
	site := constraintSite{
		SchemaName: schemaName,
		FieldPath:  fieldPath,
		SourcePath: sourcePath,
	}
	if s.Min != nil {
		site.Minimum = s.Min
	}
	if s.Max != nil {
		site.Maximum = s.Max
	}
	if s.MinLength > 0 && s.MinLength <= math.MaxInt64 {
		v := int64(s.MinLength)
		site.MinLength = &v
	}
	if s.MaxLength != nil && *s.MaxLength <= math.MaxInt64 {
		v := int64(*s.MaxLength)
		site.MaxLength = &v
	}
	if s.MinItems > 0 && s.MinItems <= math.MaxInt64 {
		v := int64(s.MinItems)
		site.MinItems = &v
	}
	if s.MaxItems != nil && *s.MaxItems <= math.MaxInt64 {
		v := int64(*s.MaxItems)
		site.MaxItems = &v
	}
	w.constraintSites = append(w.constraintSites, site)
}

func (w *walker) items(id string, s *openapi3.Schema, schemaName, fieldPath, sourcePath string, sc scope) {
	if s.Items == nil || s.Items.Ref != "" || s.Items.Value == nil {
		return
	}
	w.schema(id+"_Item", s.Items.Value, schemaName, joinPath(fieldPath, "[]"), sourcePath+".items", sc)
}

func typeOf(s *openapi3.Schema) string {
	if s.Type == nil || len(s.Type.Slice()) == 0 {
		return ""
	}
	return s.Type.Slice()[0]
}

func enumValues(enum []any) []string {
	result := make([]string, len(enum))
	for i, v := range enum {
		if v == nil {
			result[i] = "null"
			continue
		}
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}

func toPascal(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	for _, part := range strings.Split(s, "_") {
		if part != "" {
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String()
}

func joinPath(base, part string) string {
	if base == "" {
		return part
	}
	return base + "." + part
}
