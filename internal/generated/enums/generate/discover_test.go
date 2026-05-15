package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeValues(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want []string
	}{
		{"both empty", nil, nil, nil},
		{"a empty", nil, []string{"x"}, []string{"x"}},
		{"b empty", []string{"x"}, nil, []string{"x"}},
		{"no overlap", []string{"a", "b"}, []string{"c"}, []string{"a", "b", "c"}},
		{"full overlap", []string{"a", "b"}, []string{"a", "b"}, []string{"a", "b"}},
		{"partial overlap", []string{"a", "b"}, []string{"b", "c"}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeValues(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDedupByID(t *testing.T) {
	t.Run("no duplicates", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x"}},
			{ID: "B", Type: "string", Values: []string{"y"}},
		}
		got, err := dedupByID(sites)
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("merge values on same ID", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x", "y"}},
			{ID: "A", Type: "string", Values: []string{"y", "z"}},
		}
		got, err := dedupByID(sites)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, []string{"x", "y", "z"}, got[0].Values)
	})

	t.Run("conflicting types error", func(t *testing.T) {
		sites := []enumSite{
			{ID: "A", Type: "string", Values: []string{"x"}, SourcePath: "path1"},
			{ID: "A", Type: "integer", Values: []string{"1"}, SourcePath: "path2"},
		}
		_, err := dedupByID(sites)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicting enum types")
	})
}

func TestToPascal(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"foo", "Foo"},
		{"foo_bar", "FooBar"},
		{"foo_bar_baz", "FooBarBaz"},
		{"already_Good", "AlreadyGood"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, toPascal(tt.in))
		})
	}
}

func TestJoinPath(t *testing.T) {
	assert.Equal(t, "child", joinPath("", "child"))
	assert.Equal(t, "parent.child", joinPath("parent", "child"))
	assert.Equal(t, "a.b.c", joinPath("a.b", "c"))
}

func TestEnumValues(t *testing.T) {
	got := enumValues([]any{"a", "b", nil, 42})
	assert.Equal(t, []string{"a", "b", "null", "42"}, got)
}

func TestTypeOf(t *testing.T) {
	t.Run("nil type", func(t *testing.T) {
		s := &openapi3.Schema{}
		assert.Equal(t, "", typeOf(s))
	})

	t.Run("string type", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		assert.Equal(t, "string", typeOf(s))
	})

	t.Run("integer type", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"integer"}}
		assert.Equal(t, "integer", typeOf(s))
	})
}

func TestWalkerPlainEnum(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Status": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
						Enum: []any{"active", "inactive"},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	require.Len(t, w.sites, 1)
	assert.Equal(t, "Status", w.sites[0].ID)
	assert.Equal(t, scopeSchema, w.sites[0].Scope)
	assert.Equal(t, "string", w.sites[0].Type)
	assert.Equal(t, []string{"active", "inactive"}, w.sites[0].Values)
	assert.Equal(t, "", w.sites[0].FieldPath)
}

func TestWalkerNestedProperties(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Cluster": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Properties: openapi3.Schemas{
							"status": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"string"},
									Enum: []any{"healthy", "degraded"},
								},
							},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	require.Len(t, w.sites, 1)
	assert.Equal(t, "Cluster_Status", w.sites[0].ID)
	assert.Equal(t, "status", w.sites[0].FieldPath)
	assert.Equal(t, "Cluster", w.sites[0].SchemaName)
}

func TestWalkerComposition(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Combined": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						AllOf: openapi3.SchemaRefs{
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"a", "b"},
							}},
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"b", "c"},
							}},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	// Both branches share the same ID, so dedupByID merges them.
	sites, err := dedupByID(w.sites)
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, []string{"a", "b", "c"}, sites[0].Values)
}

func TestWalkerOneOf(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Choice": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						OneOf: openapi3.SchemaRefs{
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"x"},
							}},
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"y"},
							}},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	sites, err := dedupByID(w.sites)
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, []string{"x", "y"}, sites[0].Values)
}

func TestWalkerAnyOf(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Mixed": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						AnyOf: openapi3.SchemaRefs{
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"p"},
							}},
							{Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"q"},
							}},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	sites, err := dedupByID(w.sites)
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, []string{"p", "q"}, sites[0].Values)
}

func TestWalkerArrayItems(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"TagList": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"array"},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"tag1", "tag2"},
							},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	require.Len(t, w.sites, 1)
	assert.Equal(t, "TagList_Item", w.sites[0].ID)
	assert.Equal(t, "[]", w.sites[0].FieldPath)
}

func TestWalkerParameters(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Parameters: openapi3.ParametersMap{
				"sortBy": &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name: "sortBy",
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
								Enum: []any{"name", "date"},
							},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	require.Len(t, w.sites, 1)
	assert.Equal(t, scopeParameter, w.sites[0].Scope)
	assert.Equal(t, "Param_sortBy", w.sites[0].ID)
}

func TestWalkerSkipsRefProperties(t *testing.T) {
	doc := &openapi3.T{
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"Parent": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Properties: openapi3.Schemas{
							"child": &openapi3.SchemaRef{
								Ref: "#/components/schemas/Child",
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"string"},
									Enum: []any{"a"},
								},
							},
						},
					},
				},
			},
		},
	}
	w := &walker{doc: doc}
	w.run()

	assert.Empty(t, w.sites)
}

func TestWalkerNilComponents(t *testing.T) {
	doc := &openapi3.T{}
	w := &walker{doc: doc}
	w.run()
	assert.Empty(t, w.sites)
}
