package main

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEnumTable_ScopeFiltering(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "S", FieldPath: "field", Type: "string", Values: []string{"a"}},
		{Scope: scopeParameter, SchemaName: "P", FieldPath: "field", Type: "string", Values: []string{"b"}},
	}
	table := buildEnumTable(sites)

	assert.Contains(t, table, "S")
	assert.NotContains(t, table, "P")
}

func TestBuildEnumTable_TypeFiltering(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "A", FieldPath: "f", Type: "string", Values: []string{"x"}},
		{Scope: scopeSchema, SchemaName: "B", FieldPath: "f", Type: "integer", Values: []string{"1"}},
		{Scope: scopeSchema, SchemaName: "C", FieldPath: "f", Type: "boolean", Values: []string{"true"}},
		{Scope: scopeSchema, SchemaName: "D", FieldPath: "f", Type: "", Values: []string{"?"}},
	}
	table := buildEnumTable(sites)

	assert.Contains(t, table, "A")
	assert.Contains(t, table, "B")
	assert.NotContains(t, table, "C")
	assert.NotContains(t, table, "D")
}

func TestBuildEnumTable_TopLevelKeySkip(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "Status", FieldPath: "", Type: "string", Values: []string{"active"}},
	}
	table := buildEnumTable(sites)
	assert.Empty(t, table)
}

func TestBuildEnumTable_ArraySuffix(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "S", FieldPath: "tags.[]", Type: "string", Values: []string{"a"}},
		{Scope: scopeSchema, SchemaName: "S", FieldPath: "items[]", Type: "string", Values: []string{"b"}},
	}
	table := buildEnumTable(sites)

	require.Contains(t, table, "S")
	tagsEntry, ok := table["S"]["tags"]
	require.True(t, ok)
	assert.True(t, tagsEntry.IsArray)

	itemsEntry, ok := table["S"]["items"]
	require.True(t, ok)
	assert.True(t, itemsEntry.IsArray)
}

func TestBuildEnumTable_NonArrayField(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "S", FieldPath: "status", Type: "string", Values: []string{"on"}},
	}
	table := buildEnumTable(sites)
	require.Contains(t, table, "S")
	assert.False(t, table["S"]["status"].IsArray)
}

func TestIsArrayLiteral(t *testing.T) {
	assert.Equal(t, "IsArray: true, ", isArrayLiteral(true))
	assert.Equal(t, "", isArrayLiteral(false))
}

func TestJoinQuoted(t *testing.T) {
	assert.Equal(t, `"a", "b", "c"`, joinQuoted([]string{"a", "b", "c"}))
	assert.Equal(t, `"x"`, joinQuoted([]string{"x"}))
	assert.Equal(t, "", joinQuoted(nil))
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	assert.Equal(t, []string{"a", "b", "c"}, sortedKeys(m))
}

func TestSortedKeysEmpty(t *testing.T) {
	m := map[string]int{}
	assert.Empty(t, sortedKeys(m))
}

func TestGenerateEndToEnd(t *testing.T) {
	sites := []enumSite{
		{
			ID:         "Cluster_Status",
			Scope:      scopeSchema,
			SchemaName: "Cluster",
			FieldPath:  "status",
			Type:       "string",
			Values:     []string{"healthy", "degraded"},
		},
		{
			ID:         "Cluster_Priority",
			Scope:      scopeSchema,
			SchemaName: "Cluster",
			FieldPath:  "priority",
			Type:       "integer",
			Values:     []string{"1", "2", "3"},
		},
		{
			ID:         "Bucket_Tags_Item",
			Scope:      scopeSchema,
			SchemaName: "Bucket",
			FieldPath:  "tags.[]",
			Type:       "string",
			Values:     []string{"dev", "prod"},
		},
		// Should be excluded: parameter scope
		{
			ID:         "Param_sort",
			Scope:      scopeParameter,
			SchemaName: "sort",
			FieldPath:  "sort",
			Type:       "string",
			Values:     []string{"name"},
		},
		// Should be excluded: top-level enum
		{
			ID:         "Status",
			Scope:      scopeSchema,
			SchemaName: "Status",
			FieldPath:  "",
			Type:       "string",
			Values:     []string{"on", "off"},
		},
	}

	src, err := generate(sites)
	require.NoError(t, err)

	output := string(src)

	// Verify it parses as valid Go
	fset := token.NewFileSet()
	_, parseErr := parser.ParseFile(fset, "enums.gen.go", src, parser.AllErrors)
	require.NoError(t, parseErr, "generated code must be valid Go")

	// Verify expected map entries
	assert.Contains(t, output, `"Cluster"`)
	assert.Contains(t, output, `"status"`)
	assert.Contains(t, output, `"healthy"`)
	assert.Contains(t, output, `"degraded"`)
	assert.Contains(t, output, `"priority"`)
	assert.Contains(t, output, `"Bucket"`)
	assert.Contains(t, output, `"tags"`)
	assert.Contains(t, output, "IsArray: true")

	// Verify excluded entries
	assert.NotContains(t, output, `"Param_sort"`)
	assert.NotContains(t, output, `"sort": {`)

	// Verify the "Status" top-level enum with empty field is excluded
	// The schema "Status" should not appear as a top-level map key
	// because its FieldPath is empty
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `"Status"`) && strings.Contains(trimmed, "{") {
			t.Error("top-level enum 'Status' with empty FieldPath should be excluded from the table")
		}
	}
}
