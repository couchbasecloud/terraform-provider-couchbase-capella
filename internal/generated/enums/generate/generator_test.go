package main

import (
	"strings"
	"testing"
)

func TestToExported(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", "Foo"},
		{"Foo", "Foo"},
		{"fooBar", "FooBar"},
		{"param_status", "Param_status"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := toExported(tc.in)
			if got != tc.want {
				t.Errorf("toExported(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFieldKey(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"roles.[]", "roles"},
		{"[]", ""},
		{"status", "status"},
		{"weeklySchedule.dayOfWeek", "weeklySchedule.dayOfWeek"},
		{"nested.roles.[]", "nested.roles"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := fieldKey(tc.in)
			if got != tc.want {
				t.Errorf("fieldKey(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSiteDecl(t *testing.T) {
	cases := []struct {
		name       string
		site       enumSite
		wantPrefix string
		wantErr    bool
	}{
		{
			name:       "string type",
			site:       enumSite{Type: "string", Values: []string{"a", "b"}, SourcePath: "x"},
			wantPrefix: `[]string{`,
		},
		{
			name:       "empty type treated as string",
			site:       enumSite{Type: "", Values: []string{"x"}, SourcePath: "x"},
			wantPrefix: `[]string{`,
		},
		{
			name:       "integer type",
			site:       enumSite{Type: "integer", Values: []string{"1", "2"}, SourcePath: "x"},
			wantPrefix: `[]int64{`,
		},
		{
			name:    "integer type with non-integer value",
			site:    enumSite{Type: "integer", Values: []string{"notanint"}, SourcePath: "x"},
			wantErr: true,
		},
		{
			name:    "unsupported type",
			site:    enumSite{Type: "boolean", Values: []string{"true"}, SourcePath: "x"},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := siteDecl(tc.site)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.HasPrefix(got, tc.wantPrefix) {
				t.Errorf("siteDecl() = %q, want prefix %q", got, tc.wantPrefix)
			}
		})
	}
}

func TestJoinQuoted(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"a", "b"}, `"a", "b"`},
		{[]string{"single"}, `"single"`},
		{[]string{}, ""},
		{[]string{`has"quote`}, `"has\"quote"`},
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.in, ","), func(t *testing.T) {
			got := joinQuoted(tc.in)
			if got != tc.want {
				t.Errorf("joinQuoted(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBuildLookup(t *testing.T) {
	sites := []enumSite{
		{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "status", Type: "string", Values: []string{"active", "inactive"}},
		{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "roles.[]", Type: "string", Values: []string{"owner", "member"}},
		{Scope: scopeSchema, SchemaName: "Foo", FieldPath: "", Type: "string", Values: []string{"top-level"}},
		{Scope: scopeParameter, SchemaName: "Param_Filter", FieldPath: "filter", Type: "string", Values: []string{"v"}},
		{Scope: scopeSchema, SchemaName: "Bar", FieldPath: "count", Type: "integer", Values: []string{"1", "2"}},
	}
	lookup := buildLookup(sites)

	t.Run("schema field included", func(t *testing.T) {
		vals, ok := lookup["Foo"]["status"]
		if !ok {
			t.Fatal("Foo.status not in lookup")
		}
		if len(vals) != 2 || vals[0] != "active" {
			t.Errorf("unexpected values: %v", vals)
		}
	})

	t.Run("array field key stripped", func(t *testing.T) {
		if _, ok := lookup["Foo"]["roles"]; !ok {
			t.Fatal("Foo.roles not in lookup after stripping []")
		}
		if _, ok := lookup["Foo"]["roles.[]"]; ok {
			t.Fatal("Foo.roles.[] should not be in lookup")
		}
	})

	t.Run("top-level enum excluded", func(t *testing.T) {
		if fields, ok := lookup["Foo"]; ok {
			if _, ok := fields[""]; ok {
				t.Fatal("top-level enum (empty FieldPath) should be excluded")
			}
		}
	})

	t.Run("parameter scope excluded", func(t *testing.T) {
		if _, ok := lookup["Param_Filter"]; ok {
			t.Fatal("parameter scope should be excluded from lookup")
		}
	})

	t.Run("integer type excluded", func(t *testing.T) {
		if _, ok := lookup["Bar"]; ok {
			t.Fatal("integer type should be excluded from lookup")
		}
	})
}

func TestGenerate(t *testing.T) {
	sites := []enumSite{
		{
			ID:         "MyEnum",
			Scope:      scopeSchema,
			SchemaName: "MyEnum",
			FieldPath:  "",
			Type:       "string",
			Values:     []string{"a", "b"},
			SourcePath: "components.schemas.MyEnum",
		},
		{
			ID:         "MySchema_Status",
			Scope:      scopeSchema,
			SchemaName: "MySchema",
			FieldPath:  "status",
			Type:       "string",
			Values:     []string{"active", "inactive"},
			SourcePath: "components.schemas.MySchema.properties.status",
		},
	}

	src, err := generate(sites)
	if err != nil {
		t.Fatalf("generate() error: %v", err)
	}

	out := string(src)

	t.Run("contains package declaration", func(t *testing.T) {
		if !strings.Contains(out, "package enums") {
			t.Error("output missing package declaration")
		}
	})

	t.Run("contains generated header", func(t *testing.T) {
		if !strings.Contains(out, "DO NOT EDIT") {
			t.Error("output missing generated header")
		}
	})

	t.Run("contains var declaration", func(t *testing.T) {
		if !strings.Contains(out, "var MyEnum") {
			t.Error("output missing var MyEnum")
		}
	})

	t.Run("contains lookup map", func(t *testing.T) {
		if !strings.Contains(out, "var Table") {
			t.Error("output missing Table map")
		}
	})

	t.Run("lookup contains schema field", func(t *testing.T) {
		if !strings.Contains(out, `"MySchema"`) {
			t.Error("output missing MySchema in Lookup")
		}
	})
}

func TestGenerateDuplicateVarError(t *testing.T) {
	sites := []enumSite{
		{ID: "Same", Type: "string", Values: []string{"a"}, SourcePath: "path1"},
		{ID: "Same", Type: "string", Values: []string{"b"}, SourcePath: "path2"},
	}
	_, err := generate(sites)
	if err == nil {
		t.Fatal("expected error for duplicate var name")
	}
	if !strings.Contains(err.Error(), "duplicate var") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	got := sortedKeys(m)
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q want %q", i, got[i], want[i])
		}
	}
}
