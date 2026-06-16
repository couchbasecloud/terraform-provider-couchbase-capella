package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestMergeValues(t *testing.T) {
	cases := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{"no overlap", []string{"a", "b"}, []string{"c", "d"}, []string{"a", "b", "c", "d"}},
		{"full overlap", []string{"a", "b"}, []string{"a", "b"}, []string{"a", "b"}},
		{"partial overlap", []string{"a", "b"}, []string{"b", "c"}, []string{"a", "b", "c"}},
		{"empty a", []string{}, []string{"x"}, []string{"x"}},
		{"empty b", []string{"x"}, []string{}, []string{"x"}},
		{"both empty", []string{}, []string{}, []string{}},
		{"preserves first-seen order", []string{"b", "a"}, []string{"a", "c", "b"}, []string{"b", "a", "c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeValues(tc.a, tc.b)
			if !equalStrings(got, tc.want) {
				t.Errorf("mergeValues(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDedupByID(t *testing.T) {
	t.Run("no duplicates preserves order", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x"}},
			{ID: "B", Type: "string", Values: []string{"y"}},
		}
		got, err := dedupByID(sites)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 || got[0].ID != "A" || got[1].ID != "B" {
			t.Fatalf("want [A, B], got %v", idsOf(got))
		}
	})

	t.Run("merges duplicate IDs and unions values", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x", "y"}},
			{ID: "A", Type: "string", Values: []string{"y", "z"}},
		}
		got, err := dedupByID(sites)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Fatalf("want 1 site, got %d", len(got))
		}
		want := []string{"x", "y", "z"}
		if !equalStrings(got[0].Values, want) {
			t.Errorf("values = %v, want %v", got[0].Values, want)
		}
	})

	t.Run("conflicting types returns error", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x"}, SourcePath: "path1"},
			{ID: "A", Type: "integer", Values: []string{"1"}, SourcePath: "path2"},
		}
		_, err := dedupByID(sites)
		if err == nil {
			t.Fatal("expected error for conflicting types")
		}
	})

	t.Run("merges three branches in order", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"a"}},
			{ID: "A", Type: "string", Values: []string{"b"}},
			{ID: "A", Type: "string", Values: []string{"c"}},
		}
		got, err := dedupByID(sites)
		if err != nil {
			t.Fatal(err)
		}
		if !equalStrings(got[0].Values, []string{"a", "b", "c"}) {
			t.Errorf("values = %v, want [a b c]", got[0].Values)
		}
	})
}

func TestToPascal(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"foo", "Foo"},
		{"foo_bar", "FooBar"},
		{"foo_bar_baz", "FooBarBaz"},
		{"already", "Already"},
		{"_leading", "Leading"},
		{"trailing_", "Trailing"},
		{"double__underscore", "DoubleUnderscore"},
		{"a", "A"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			if got := toPascal(tc.in); got != tc.want {
				t.Errorf("toPascal(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	cases := []struct {
		base, part, want string
	}{
		{"", "field", "field"},
		{"parent", "child", "parent.child"},
		{"a.b", "c", "a.b.c"},
		{"", "", ""},
		{"a", "[]", "a.[]"},
	}
	for _, tc := range cases {
		t.Run(tc.base+"+"+tc.part, func(t *testing.T) {
			if got := joinPath(tc.base, tc.part); got != tc.want {
				t.Errorf("joinPath(%q, %q) = %q, want %q", tc.base, tc.part, got, tc.want)
			}
		})
	}
}

func TestEnumValues(t *testing.T) {
	cases := []struct {
		name string
		in   []any
		want []string
	}{
		{"strings", []any{"a", "b"}, []string{"a", "b"}},
		{"integers", []any{1, 2}, []string{"1", "2"}},
		{"mixed", []any{"a", 1, true}, []string{"a", "1", "true"}},
		{"nil value rendered as null", []any{nil, "x"}, []string{"null", "x"}},
		{"empty", []any{}, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := enumValues(tc.in)
			if !equalStrings(got, tc.want) {
				t.Errorf("enumValues(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestTypeOf(t *testing.T) {
	t.Run("nil type returns empty", func(t *testing.T) {
		if got := typeOf(&openapi3.Schema{}); got != "" {
			t.Errorf("typeOf(empty) = %q, want \"\"", got)
		}
	})

	t.Run("returns first type", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		if got := typeOf(s); got != "string" {
			t.Errorf("typeOf(string) = %q, want \"string\"", got)
		}
	})

	t.Run("union picks first", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"integer", "null"}}
		if got := typeOf(s); got != "integer" {
			t.Errorf("typeOf(union) = %q, want \"integer\"", got)
		}
	})
}

func TestWalker_PlainEnum(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Status": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"string"},
					Enum: []any{"active", "expired"},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 site, got %d: %v", len(sites), sites)
	}
	got := sites[0]
	if got.ID != "Status" || got.SchemaName != "Status" || got.FieldPath != "" || got.Type != "string" {
		t.Errorf("unexpected site: %+v", got)
	}
	if !equalStrings(got.Values, []string{"active", "expired"}) {
		t.Errorf("values = %v, want [active expired]", got.Values)
	}
	if got.Scope != scopeSchema {
		t.Errorf("scope = %q, want %q", got.Scope, scopeSchema)
	}
}

func TestWalker_NestedProperty(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Bucket": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"resolution": refOf(&openapi3.Schema{
							Type: &openapi3.Types{"string"},
							Enum: []any{"seqno", "lww"},
						}),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 site, got %d: %v", len(sites), sites)
	}
	got := sites[0]
	if got.SchemaName != "Bucket" || got.FieldPath != "resolution" {
		t.Errorf("unexpected site: %+v", got)
	}
	if got.ID != "Bucket_Resolution" {
		t.Errorf("ID = %q, want Bucket_Resolution", got.ID)
	}
	if got.SourcePath != "components.schemas.Bucket.properties.resolution" {
		t.Errorf("SourcePath = %q", got.SourcePath)
	}
}

func TestWalker_PropertyRefSkipped(t *testing.T) {
	// An unresolved $ref property (Value nil) is skipped; the referenced schema
	// is visited as its own top-level entry instead. (Resolved scalar-enum refs
	// are captured — see TestWalker_PropertyRefScalarEnumCaptured.)
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Status": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"string"},
					Enum: []any{"on", "off"},
				}),
				"Container": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"status": {Ref: "#/components/schemas/Status"},
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 site (Status only), got %d: %v", len(sites), sites)
	}
	if sites[0].SchemaName != "Status" {
		t.Errorf("unexpected site: %+v", sites[0])
	}
}

func TestWalker_ArrayItems(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Roles": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"roles": refOf(&openapi3.Schema{
							Type: &openapi3.Types{"array"},
							Items: refOf(&openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"admin", "viewer"},
							}),
						}),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 site, got %d: %v", len(sites), sites)
	}
	got := sites[0]
	if got.FieldPath != "roles.[]" {
		t.Errorf("FieldPath = %q, want roles.[]", got.FieldPath)
	}
	if got.ID != "Roles_Roles_Item" {
		t.Errorf("ID = %q, want Roles_Roles_Item", got.ID)
	}
}

func TestWalker_PropertyRefScalarEnumCaptured(t *testing.T) {
	// A property whose resolved $ref target is a scalar enum is captured under
	// (schema, field) — the case the generator previously dropped (the enum
	// reached only the field-less top-level site, which buildEnumTable discards).
	tz := &openapi3.Schema{Type: &openapi3.Types{"string"}, Enum: []any{"PT", "ET"}}
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"onOffTimezone": refOf(tz),
				"Schedule": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"timezone": resolvedRef("#/components/schemas/onOffTimezone", tz),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	got, ok := findSite(sites, "Schedule", "timezone")
	if !ok {
		t.Fatalf("want a Schedule.timezone site, got %v", idsOf(sites))
	}
	if got.ID != "Schedule_Timezone" || got.Type != "string" {
		t.Errorf("unexpected site: %+v", got)
	}
	if !equalStrings(got.Values, []string{"PT", "ET"}) {
		t.Errorf("values = %v, want [PT ET]", got.Values)
	}
}

func TestWalker_ArrayItemRefScalarEnumCaptured(t *testing.T) {
	// An array whose items resolve to a scalar-enum $ref is captured at "[]".
	role := &openapi3.Schema{Type: &openapi3.Types{"string"}, Enum: []any{"admin", "viewer"}}
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Role": refOf(role),
				"User": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"roles": refOf(&openapi3.Schema{
							Type:  &openapi3.Types{"array"},
							Items: resolvedRef("#/components/schemas/Role", role),
						}),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	got, ok := findSite(sites, "User", "roles.[]")
	if !ok {
		t.Fatalf("want a User.roles.[] site, got %v", idsOf(sites))
	}
	if !equalStrings(got.Values, []string{"admin", "viewer"}) {
		t.Errorf("values = %v, want [admin viewer]", got.Values)
	}
}

func TestWalker_ObjectRefPropertyNotInlined(t *testing.T) {
	// A property whose resolved $ref target is an object is NOT inlined; only
	// the object's own top-level visit records its inner enum.
	inner := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: openapi3.Schemas{
			"kind": refOf(&openapi3.Schema{Type: &openapi3.Types{"string"}, Enum: []any{"a", "b"}}),
		},
	}
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Inner": refOf(inner),
				"Container": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"inner": resolvedRef("#/components/schemas/Inner", inner),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if _, ok := findSite(sites, "Container", "inner.kind"); ok {
		t.Errorf("object $ref property must not be inlined; got Container.inner.kind in %v", idsOf(sites))
	}
	if _, ok := findSite(sites, "Inner", "kind"); !ok {
		t.Errorf("want Inner.kind from the top-level visit, got %v", idsOf(sites))
	}
}

func TestWalker_AllOfMerged(t *testing.T) {
	// Two allOf branches contribute values for the same logical field; the
	// walker emits both with the same ID and dedupByID merges them.
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Combined": refOf(&openapi3.Schema{
					AllOf: openapi3.SchemaRefs{
						refOf(&openapi3.Schema{
							Type: &openapi3.Types{"string"},
							Enum: []any{"a", "b"},
						}),
						refOf(&openapi3.Schema{
							Type: &openapi3.Types{"string"},
							Enum: []any{"b", "c"},
						}),
					},
				}),
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 merged site, got %d: %v", len(sites), sites)
	}
	if !equalStrings(sites[0].Values, []string{"a", "b", "c"}) {
		t.Errorf("values = %v, want [a b c]", sites[0].Values)
	}
}

func TestWalker_ParameterScope(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Parameters: openapi3.ParametersMap{
				"sortBy": &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Schema: refOf(&openapi3.Schema{
							Type: &openapi3.Types{"string"},
							Enum: []any{"asc", "desc"},
						}),
					},
				},
			},
		},
	}

	sites := runWalker(t, doc)
	if len(sites) != 1 {
		t.Fatalf("want 1 site, got %d", len(sites))
	}
	got := sites[0]
	if got.Scope != scopeParameter {
		t.Errorf("scope = %q, want %q", got.Scope, scopeParameter)
	}
	if got.ID != "Param_sortBy" || got.SchemaName != "sortBy" {
		t.Errorf("unexpected param site: %+v", got)
	}
}

func TestWalker_NoComponents(t *testing.T) {
	doc := &openapi3.T{}
	sites := runWalker(t, doc)
	if len(sites) != 0 {
		t.Errorf("want 0 sites for empty doc, got %d", len(sites))
	}
}

func TestWalker_SortedOutput(t *testing.T) {
	// discover() sorts by (Scope, SchemaName, FieldPath). Use input names that
	// would come out of a Go map in arbitrary order to exercise the sort.
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Zebra": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"string"},
					Enum: []any{"x"},
				}),
				"Alpha": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"string"},
					Enum: []any{"y"},
				}),
			},
			Parameters: openapi3.ParametersMap{
				"q": &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Schema: refOf(&openapi3.Schema{
							Type: &openapi3.Types{"string"},
							Enum: []any{"z"},
						}),
					},
				},
			},
		},
	}

	sites, err := discoverDoc(doc)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"q", "Alpha", "Zebra"} // parameter scope sorts before schema scope ("parameter" < "schema")
	got := make([]string, len(sites))
	for i, s := range sites {
		got[i] = s.SchemaName
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// helpers

func refOf(s *openapi3.Schema) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Value: s}
}

// resolvedRef mimics what the OpenAPI loader produces for a property/item $ref:
// both the Ref string and the resolved target Value are populated.
func resolvedRef(ref string, s *openapi3.Schema) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Ref: ref, Value: s}
}

func findSite(sites []enumSite, schemaName, fieldPath string) (enumSite, bool) {
	for _, s := range sites {
		if s.SchemaName == schemaName && s.FieldPath == fieldPath {
			return s, true
		}
	}
	return enumSite{}, false
}

func runWalker(t *testing.T, doc *openapi3.T) []enumSite {
	t.Helper()
	w := &walker{doc: doc}
	w.run()
	sites, err := dedupByID(w.sites)
	if err != nil {
		t.Fatalf("dedupByID: %v", err)
	}
	return sites
}

// discoverDoc mirrors discover() but operates on an in-memory document so
// tests don't need a temp YAML file. The sort criteria match discover().
func discoverDoc(doc *openapi3.T) ([]enumSite, error) {
	w := &walker{doc: doc}
	w.run()
	sites, err := dedupByID(w.sites)
	if err != nil {
		return nil, err
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
	return sites, nil
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func idsOf(sites []enumSite) []string {
	out := make([]string, len(sites))
	for i, s := range sites {
		out[i] = s.ID
	}
	return out
}

// Composition discovery tests

func TestExtractSchemaName(t *testing.T) {
	cases := []struct {
		ref, want string
	}{
		{"#/components/schemas/AWS", "AWS"},
		{"#/components/schemas/GCPConfig", "GCPConfig"},
		{"#/components/schemas/AzureConfigData", "AzureConfigData"},
		{"#/components/parameters/foo", ""},
		{"", ""},
		{"AWS", ""},
	}
	for _, tc := range cases {
		t.Run(tc.ref, func(t *testing.T) {
			got := extractSchemaName(tc.ref)
			if got != tc.want {
				t.Errorf("extractSchemaName(%q) = %q, want %q", tc.ref, got, tc.want)
			}
		})
	}
}

func TestMergeBranches(t *testing.T) {
	cases := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{"no overlap", []string{"A", "B"}, []string{"C", "D"}, []string{"A", "B", "C", "D"}},
		{"full overlap", []string{"A", "B"}, []string{"A", "B"}, []string{"A", "B"}},
		{"partial overlap", []string{"A", "B"}, []string{"B", "C"}, []string{"A", "B", "C"}},
		{"empty a", []string{}, []string{"X"}, []string{"X"}},
		{"empty b", []string{"X"}, []string{}, []string{"X"}},
		{"both empty", []string{}, []string{}, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeBranches(tc.a, tc.b)
			if !equalStrings(got, tc.want) {
				t.Errorf("mergeBranches(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestWalker_OneOfComposition(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"AWSConfig": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
				}),
				"GCPConfig": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
				}),
				"CMEKRequest": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"config": refOf(&openapi3.Schema{
							OneOf: openapi3.SchemaRefs{
								{Ref: "#/components/schemas/AWSConfig"},
								{Ref: "#/components/schemas/GCPConfig"},
							},
						}),
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	if len(w.compositionSites) != 1 {
		t.Fatalf("want 1 composition site, got %d: %v", len(w.compositionSites), w.compositionSites)
	}
	got := w.compositionSites[0]
	if got.SchemaName != "CMEKRequest" {
		t.Errorf("SchemaName = %q, want CMEKRequest", got.SchemaName)
	}
	if got.FieldPath != "config" {
		t.Errorf("FieldPath = %q, want config", got.FieldPath)
	}
	if got.Kind != kindOneOf {
		t.Errorf("Kind = %q, want %q", got.Kind, kindOneOf)
	}
	if !equalStrings(got.Branches, []string{"AWSConfig", "GCPConfig"}) {
		t.Errorf("Branches = %v, want [AWSConfig GCPConfig]", got.Branches)
	}
}

func TestWalker_AnyOfComposition(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"AWS":   refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"GCP":   refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"Azure": refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"NetworkPeer": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"providerConfig": refOf(&openapi3.Schema{
							AnyOf: openapi3.SchemaRefs{
								{Ref: "#/components/schemas/AWS"},
								{Ref: "#/components/schemas/GCP"},
								{Ref: "#/components/schemas/Azure"},
							},
						}),
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	if len(w.compositionSites) != 1 {
		t.Fatalf("want 1 composition site, got %d", len(w.compositionSites))
	}
	got := w.compositionSites[0]
	if got.Kind != kindAnyOf {
		t.Errorf("Kind = %q, want %q", got.Kind, kindAnyOf)
	}
	if !equalStrings(got.Branches, []string{"AWS", "GCP", "Azure"}) {
		t.Errorf("Branches = %v, want [AWS GCP Azure]", got.Branches)
	}
}

func TestWalker_AllOfComposition(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"BaseConfig":     refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"ExtendedConfig": refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"MergedConfig": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"settings": refOf(&openapi3.Schema{
							AllOf: openapi3.SchemaRefs{
								{Ref: "#/components/schemas/BaseConfig"},
								{Ref: "#/components/schemas/ExtendedConfig"},
							},
						}),
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	if len(w.compositionSites) != 1 {
		t.Fatalf("want 1 composition site, got %d", len(w.compositionSites))
	}
	got := w.compositionSites[0]
	if got.Kind != kindAllOf {
		t.Errorf("Kind = %q, want %q", got.Kind, kindAllOf)
	}
}

func TestWalker_InlineSchemaSkipsComposition(t *testing.T) {
	// When composition branches are inline (not $ref), no composition site is recorded.
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"InlineUnion": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"value": refOf(&openapi3.Schema{
							OneOf: openapi3.SchemaRefs{
								refOf(&openapi3.Schema{Type: &openapi3.Types{"string"}}),
								refOf(&openapi3.Schema{Type: &openapi3.Types{"integer"}}),
							},
						}),
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	if len(w.compositionSites) != 0 {
		t.Errorf("want 0 composition sites for inline branches, got %d: %v",
			len(w.compositionSites), w.compositionSites)
	}
}

func TestWalker_MixedCompositionSkipped(t *testing.T) {
	// Mixed compositions (some $ref, some inline) should be skipped entirely
	// to avoid incomplete branch metadata.
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"AWSConfig": refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"MixedUnion": refOf(&openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"config": refOf(&openapi3.Schema{
							OneOf: openapi3.SchemaRefs{
								{Ref: "#/components/schemas/AWSConfig"},
								refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}), // inline
							},
						}),
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	if len(w.compositionSites) != 0 {
		t.Errorf("want 0 composition sites for mixed $ref/inline branches, got %d: %v",
			len(w.compositionSites), w.compositionSites)
	}
}

func TestWalker_TopLevelCompositionSkipped(t *testing.T) {
	// Top-level oneOf (empty FieldPath) is skipped when building the table
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"A": refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"B": refOf(&openapi3.Schema{Type: &openapi3.Types{"object"}}),
				"TopLevelOneOf": refOf(&openapi3.Schema{
					OneOf: openapi3.SchemaRefs{
						{Ref: "#/components/schemas/A"},
						{Ref: "#/components/schemas/B"},
					},
				}),
			},
		},
	}

	w := &walker{doc: doc}
	w.run()

	// The site is recorded but has empty FieldPath
	if len(w.compositionSites) != 1 {
		t.Fatalf("want 1 composition site, got %d", len(w.compositionSites))
	}
	if w.compositionSites[0].FieldPath != "" {
		t.Errorf("FieldPath should be empty for top-level composition")
	}

	// buildCompositionTable should exclude it
	table := buildCompositionTable(w.compositionSites)
	if len(table) != 0 {
		t.Errorf("top-level composition should be excluded from table, got %v", table)
	}
}

func TestDedupComposition(t *testing.T) {
	t.Run("merges same kind", func(t *testing.T) {
		sites := []compositionSite{
			{SchemaName: "A", FieldPath: "config", Kind: kindOneOf, Branches: []string{"X", "Y"}},
			{SchemaName: "A", FieldPath: "config", Kind: kindOneOf, Branches: []string{"Y", "Z"}},
		}
		got := dedupComposition(sites)
		if len(got) != 1 {
			t.Fatalf("want 1 site after dedup, got %d", len(got))
		}
		if !equalStrings(got[0].Branches, []string{"X", "Y", "Z"}) {
			t.Errorf("Branches = %v, want [X Y Z]", got[0].Branches)
		}
	})

	t.Run("keeps different kinds separate", func(t *testing.T) {
		sites := []compositionSite{
			{SchemaName: "A", FieldPath: "config", Kind: kindOneOf, Branches: []string{"X", "Y"}},
			{SchemaName: "A", FieldPath: "config", Kind: kindAnyOf, Branches: []string{"Y", "Z"}},
		}
		got := dedupComposition(sites)
		if len(got) != 2 {
			t.Fatalf("want 2 sites (different kinds), got %d", len(got))
		}
		// Should have both oneOf and anyOf
		kinds := []compositionKind{got[0].Kind, got[1].Kind}
		if (kinds[0] != kindOneOf || kinds[1] != kindAnyOf) && (kinds[0] != kindAnyOf || kinds[1] != kindOneOf) {
			t.Errorf("expected oneOf and anyOf, got %v", kinds)
		}
	})
}

func TestDedupConstraints_IntersectsBounds(t *testing.T) {
	// When the same (schema, field) is discovered twice with different bounds,
	// the merge must take the most restrictive of each pair (allOf semantics):
	// max of lower bounds, min of upper bounds.
	t.Run("integer bounds intersect", func(t *testing.T) {
		min2 := int64(2)
		min10 := int64(10)
		max100 := int64(100)
		max50 := int64(50)
		sites := []constraintSite{
			{SchemaName: "S", FieldPath: "f", MinLength: &min2, MaxLength: &max100},
			{SchemaName: "S", FieldPath: "f", MinLength: &min10, MaxLength: &max50},
		}
		got := dedupConstraints(sites)
		if len(got) != 1 {
			t.Fatalf("want 1 merged site, got %d", len(got))
		}
		if got[0].MinLength == nil || *got[0].MinLength != 10 {
			t.Errorf("MinLength = %v, want 10 (max of 2, 10)", got[0].MinLength)
		}
		if got[0].MaxLength == nil || *got[0].MaxLength != 50 {
			t.Errorf("MaxLength = %v, want 50 (min of 100, 50)", got[0].MaxLength)
		}
	})

	t.Run("float bounds intersect", func(t *testing.T) {
		min1 := 1.0
		min5 := 5.0
		max32 := 32.0
		max16 := 16.0
		sites := []constraintSite{
			{SchemaName: "S", FieldPath: "n", Minimum: &min1, Maximum: &max32},
			{SchemaName: "S", FieldPath: "n", Minimum: &min5, Maximum: &max16},
		}
		got := dedupConstraints(sites)
		if got[0].Minimum == nil || *got[0].Minimum != 5 {
			t.Errorf("Minimum = %v, want 5 (max of 1, 5)", got[0].Minimum)
		}
		if got[0].Maximum == nil || *got[0].Maximum != 16 {
			t.Errorf("Maximum = %v, want 16 (min of 32, 16)", got[0].Maximum)
		}
	})

	t.Run("nil bounds take the other side", func(t *testing.T) {
		min10 := int64(10)
		max100 := int64(100)
		sites := []constraintSite{
			{SchemaName: "S", FieldPath: "f", MinLength: &min10},  // MaxLength nil
			{SchemaName: "S", FieldPath: "f", MaxLength: &max100}, // MinLength nil
		}
		got := dedupConstraints(sites)
		if got[0].MinLength == nil || *got[0].MinLength != 10 {
			t.Errorf("MinLength = %v, want 10 (only one side set)", got[0].MinLength)
		}
		if got[0].MaxLength == nil || *got[0].MaxLength != 100 {
			t.Errorf("MaxLength = %v, want 100 (only one side set)", got[0].MaxLength)
		}
	})

	t.Run("item-count bounds intersect", func(t *testing.T) {
		min1 := int64(1)
		min3 := int64(3)
		max30 := int64(30)
		max10 := int64(10)
		sites := []constraintSite{
			{SchemaName: "S", FieldPath: "list", MinItems: &min1, MaxItems: &max30},
			{SchemaName: "S", FieldPath: "list", MinItems: &min3, MaxItems: &max10},
		}
		got := dedupConstraints(sites)
		if got[0].MinItems == nil || *got[0].MinItems != 3 {
			t.Errorf("MinItems = %v, want 3", got[0].MinItems)
		}
		if got[0].MaxItems == nil || *got[0].MaxItems != 10 {
			t.Errorf("MaxItems = %v, want 10", got[0].MaxItems)
		}
	})

	t.Run("keeps different fields separate", func(t *testing.T) {
		min1 := int64(1)
		sites := []constraintSite{
			{SchemaName: "S", FieldPath: "a", MinLength: &min1},
			{SchemaName: "S", FieldPath: "b", MinLength: &min1},
		}
		got := dedupConstraints(sites)
		if len(got) != 2 {
			t.Errorf("want 2 sites (different fields), got %d", len(got))
		}
	})
}
