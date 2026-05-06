package main

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file is .../internal/generated/enums/gen/main_test.go
	return filepath.Join(filepath.Dir(file), "../../../..")
}

func TestDiscoverAndGeneratePipeline(t *testing.T) {
	root := repoRoot(t)
	spec := filepath.Join(root, specPath)
	if _, err := os.Stat(spec); err != nil {
		t.Skipf("spec file not found at %s: %v", spec, err)
	}

	sites, err := discover(spec)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(sites) == 0 {
		t.Fatal("discover returned no sites")
	}

	src, err := generate(sites)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(src) == 0 {
		t.Fatal("generate returned empty output")
	}

	// generated source must be valid gofmt output
	formatted, err := format.Source(src)
	if err != nil {
		t.Fatalf("output is not valid Go: %v", err)
	}
	if !bytes.Equal(src, formatted) {
		t.Error("output is not gofmt-formatted")
	}
}

func TestDiscoverAndGenerateMatchesCommitted(t *testing.T) {
	root := repoRoot(t)
	spec := filepath.Join(root, specPath)
	committed := filepath.Join(root, "internal/generated/enums/enums.gen.go")

	for _, path := range []string{spec, committed} {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("required file not found (%s): %v", path, err)
		}
	}

	sites, err := discover(spec)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}

	src, err := generate(sites)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	want, err := os.ReadFile(committed)
	if err != nil {
		t.Fatalf("read committed file: %v", err)
	}

	if !bytes.Equal(src, want) {
		t.Error("generated output does not match committed enums.gen.go — run make gen-enums to regenerate")
	}
}
