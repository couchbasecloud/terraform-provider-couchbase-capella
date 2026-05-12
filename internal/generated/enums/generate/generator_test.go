package main

import (
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"
)

func TestBuildEnumTable(t *testing.T) {
	t.Run("indexes by schema and field", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "AllowedCidr", FieldPath: "status", Type: "string", Values: []string{"active", "expired"}},
			{Scope: scopeSchema, SchemaName: "AllowedCidr", FieldPath: "type", Type: "string", Values: []string{"temporary", "permanent"}},
		}
		got := buildEnumTable(sites)
		want := map[string]map[string]EnumDef{
			"AllowedCidr": {
				"status": {Type: "string", Values: []string{"active", "expired"}},
				"type":   {Type: "string", Values: []string{"temporary", "permanent"}},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("strips trailing []", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "APIKey", FieldPath: "roles.[]", Type: "string", Values: []string{"admin"}},
		}
		got := buildEnumTable(sites)
		def, ok := got["APIKey"]["roles"]
		if !ok {
			t.Fatalf("missing key 'roles'; got %#v", got)
		}
		if !def.IsArray {
			t.Errorf("IsArray should be true for %q", "roles.[]")
		}
	})

	t.Run("strips bare [] suffix", func(t *testing.T) {
		// Top-level array element path uses bare "[]" rather than ".[]".
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "bar[]", Type: "string", Values: []string{"a"}},
		}
		got := buildEnumTable(sites)
		if _, ok := got["Foo"]["bar"]; !ok {
			t.Errorf("expected key 'bar' (stripped of '[]'); got %#v", got["Foo"])
		}
	})

	t.Run("excludes parameter scope", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeParameter, SchemaName: "sortBy", FieldPath: "x", Type: "string", Values: []string{"asc"}},
		}
		got := buildEnumTable(sites)
		if len(got) != 0 {
			t.Errorf("parameter scope sites must be excluded, got %#v", got)
		}
	})

	t.Run("excludes non-string/integer types", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "flag", Type: "boolean", Values: []string{"true"}},
			{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "obj", Type: "object", Values: []string{"{}"}},
		}
		got := buildEnumTable(sites)
		if len(got) != 0 {
			t.Errorf("non-scalar types must be excluded, got %#v", got)
		}
	})

	t.Run("includes integer enums", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "CloudConfig", FieldPath: "compute.cpu", Type: "integer", Values: []string{"4", "32"}},
		}
		got := buildEnumTable(sites)
		def := got["CloudConfig"]["compute.cpu"]
		if def.Type != "integer" {
			t.Errorf("Type = %q, want \"integer\"", def.Type)
		}
	})

	t.Run("excludes top-level enums (empty FieldPath)", func(t *testing.T) {
		sites := []enumSite{
			{Scope: scopeSchema, SchemaName: "Status", FieldPath: "", Type: "string", Values: []string{"on", "off"}},
		}
		got := buildEnumTable(sites)
		if len(got) != 0 {
			t.Errorf("top-level enum sites must be excluded, got %#v", got)
		}
	})
}

func TestIsArrayLiteral(t *testing.T) {
	if got := isArrayLiteral(true); got != "IsArray: true, " {
		t.Errorf("isArrayLiteral(true) = %q", got)
	}
	if got := isArrayLiteral(false); got != "" {
		t.Errorf("isArrayLiteral(false) = %q, want empty", got)
	}
}

func TestJoinQuoted(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"a"}, `"a"`},
		{"multiple", []string{"a", "b", "c"}, `"a", "b", "c"`},
		{"escapes quotes", []string{`a"b`}, `"a\"b"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := joinQuoted(tc.in); got != tc.want {
				t.Errorf("joinQuoted(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	in := map[string]int{"zebra": 1, "apple": 2, "mango": 3}
	got := sortedKeys(in)
	want := []string{"apple", "mango", "zebra"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortedKeys = %v, want %v", got, want)
	}
}

func TestSortedKeys_Empty(t *testing.T) {
	got := sortedKeys(map[string]struct{}{})
	if len(got) != 0 {
		t.Errorf("sortedKeys(empty) = %v, want []", got)
	}
}

func TestGenerate_ProducesValidGo(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "AllowedCidr", FieldPath: "status", Type: "string", Values: []string{"active", "expired"}},
		{Scope: scopeSchema, SchemaName: "APIKey", FieldPath: "roles.[]", Type: "string", Values: []string{"admin", "viewer"}},
		{Scope: scopeSchema, SchemaName: "CloudConfig", FieldPath: "cpu", Type: "integer", Values: []string{"4", "32"}},
		{Scope: scopeParameter, SchemaName: "sortBy", FieldPath: "x", Type: "string", Values: []string{"asc"}},
	}

	src, err := generate(sites)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	// Output must parse as Go.
	if _, err := parser.ParseFile(token.NewFileSet(), "enums.gen.go", src, parser.AllErrors); err != nil {
		t.Fatalf("generated source does not parse: %v\n--- src ---\n%s", err, src)
	}

	got := string(src)
	checks := []string{
		"// Code generated by enums generator. DO NOT EDIT.",
		"package enums",
		"type EnumDef struct {",
		"var enumTable = map[string]map[string]EnumDef{",
		`"AllowedCidr"`,
		`"status": {Type: "string", Values: []string{"active", "expired"}}`,
		`"APIKey"`,
		`"roles": {Type: "string", IsArray: true, Values: []string{"admin", "viewer"}}`,
		`"CloudConfig"`,
		`"cpu": {Type: "integer", Values: []string{"4", "32"}}`,
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("generated output missing %q\n--- src ---\n%s", want, got)
		}
	}

	// Parameter-scope sites must not appear in the table.
	if strings.Contains(got, `"sortBy"`) {
		t.Errorf("parameter-scope site leaked into generated table:\n%s", got)
	}

	// Schema-name keys must be sorted alphabetically.
	idxAPI := strings.Index(got, `"APIKey"`)
	idxAllowed := strings.Index(got, `"AllowedCidr"`)
	idxCloud := strings.Index(got, `"CloudConfig"`)
	if idxAPI >= idxAllowed || idxAllowed >= idxCloud {
		t.Errorf("schema keys not sorted: APIKey=%d AllowedCidr=%d CloudConfig=%d", idxAPI, idxAllowed, idxCloud)
	}
}

func TestGenerate_EmptyInput(t *testing.T) {
	src, err := generate(nil)
	if err != nil {
		t.Fatalf("generate(nil): %v", err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "enums.gen.go", src, parser.AllErrors); err != nil {
		t.Fatalf("empty generated source does not parse: %v\n%s", err, src)
	}
	if !strings.Contains(string(src), "var enumTable = map[string]map[string]EnumDef{") {
		t.Errorf("empty output missing enumTable declaration:\n%s", src)
	}
}
