package main

import (
	"testing"
)

func TestMergeValues(t *testing.T) {
	cases := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "no overlap",
			a:    []string{"a", "b"},
			b:    []string{"c", "d"},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "full overlap",
			a:    []string{"a", "b"},
			b:    []string{"a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "partial overlap",
			a:    []string{"a", "b"},
			b:    []string{"b", "c"},
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty a",
			a:    []string{},
			b:    []string{"x"},
			want: []string{"x"},
		},
		{
			name: "empty b",
			a:    []string{"x"},
			b:    []string{},
			want: []string{"x"},
		},
		{
			name: "both empty",
			a:    []string{},
			b:    []string{},
			want: []string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeValues(tc.a, tc.b)
			if len(got) != len(tc.want) {
				t.Fatalf("len=%d want %d: %v", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("[%d] got %q want %q", i, got[i], tc.want[i])
				}
			}
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
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("want 2 sites, got %d", len(got))
		}
	})

	t.Run("merges duplicate IDs", func(t *testing.T) {
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
		if len(got[0].Values) != len(want) {
			t.Fatalf("values len=%d want %d: %v", len(got[0].Values), len(want), got[0].Values)
		}
		for i, v := range want {
			if got[0].Values[i] != v {
				t.Errorf("[%d] got %q want %q", i, got[0].Values[i], v)
			}
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
}

func TestToPascal(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", "Foo"},
		{"foo_bar", "FooBar"},
		{"foo_bar_baz", "FooBarBaz"},
		{"already", "Already"},
		{"_leading", "Leading"},
		{"double__underscore", "DoubleUnderscore"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := toPascal(tc.in)
			if got != tc.want {
				t.Errorf("toPascal(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	cases := []struct {
		base string
		part string
		want string
	}{
		{"", "field", "field"},
		{"parent", "child", "parent.child"},
		{"a.b", "c", "a.b.c"},
	}
	for _, tc := range cases {
		t.Run(tc.base+"+"+tc.part, func(t *testing.T) {
			got := joinPath(tc.base, tc.part)
			if got != tc.want {
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
		{"nil value", []any{nil, "x"}, []string{"null", "x"}},
		{"empty", []any{}, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := enumValues(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("len=%d want %d", len(got), len(tc.want))
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("[%d] got %q want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}
