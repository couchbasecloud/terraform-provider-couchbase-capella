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
	// A $ref'd property should not be inlined; the referenced schema is
	// visited as its own top-level entry instead.
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
