// Command testselect prints a `go test -run` regex selecting the acceptance
// tests for the resources or data sources touched by a set of changed files.
//
// A changed .go file directly under internal/resources or internal/datasources,
// or an acceptance_tests/*_test.go file, maps to a feature "stem" — its filename
// with the .go extension, the _test/_acceptance/_datasource/_schema suffixes and
// the underscores removed. The printed regex matches, case-insensitively, every
// test whose name *starts* with that stem (after the conventional TestAcc /
// Datasource / Read / Cloud prefixes), so a feature's tests are matched wherever
// they live — e.g. the bucket tests bundled in data_management_acceptance_test.go
// — without matching tests that merely mention the stem in a validation case
// (TestAccBucketResourceInvalidCluster).
//
// When a change can't be pinned to a specific resource or data source it prints
// ".*" to run the whole suite: a shared/core file (internal/api, internal/provider,
// internal/schema, a package subdirectory, ...), a shared acceptance-test helper
// or setup_test.go, or a go.mod/go.sum change. A change that touches nothing
// test-relevant (docs, examples, ...) prints nothing.
//
// Usage:
//
//	testselect <file>...            # explicit changed-file list
//	testselect -base origin/main    # changed files from git diff
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

// runAll matches every test, so the caller runs the full suite.
const runAll = ".*"

func main() {
	base := flag.String("base", "", "git ref to diff against (git diff --name-only <base>...HEAD)")
	flag.Parse()

	files, err := changedFiles(flag.Args(), *base)
	if err != nil {
		fmt.Fprintln(os.Stderr, "testselect:", err)
		os.Exit(1)
	}

	if pattern := selectPattern(files); pattern != "" {
		fmt.Println(pattern)
	}
}

// changedFiles returns the changed paths from explicit args, or from a git diff
// against base when no args are given.
func changedFiles(args []string, base string) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}
	if base == "" {
		return nil, fmt.Errorf("provide changed files as args or a -base ref")
	}
	out, err := exec.Command("git", "diff", "--name-only", base+"...HEAD").Output()
	if err != nil {
		return nil, fmt.Errorf("git diff: %w", err)
	}
	return strings.Fields(string(out)), nil
}

// selectPattern returns runAll to run everything, an anchored -run regex to run
// a subset, or "" when nothing test-relevant changed.
func selectPattern(files []string) string {
	stems := map[string]bool{}
	for _, f := range files {
		p := path.Clean(strings.TrimSpace(f))
		base := path.Base(p)
		dir := path.Dir(p)

		switch {
		case p == "go.mod" || p == "go.sum":
			return runAll
		case dir == "internal/resources" || dir == "internal/datasources":
			if strings.HasSuffix(base, ".go") {
				stems[stem(base)] = true
			}
		case dir == "acceptance_tests":
			switch {
			case base == "setup_test.go":
				return runAll // shared TestMain setup
			case strings.HasSuffix(base, "_test.go"):
				stems[stem(base)] = true
			case strings.HasSuffix(base, ".go"):
				return runAll // shared test helper
			}
		case strings.HasPrefix(p, "internal/") && strings.HasSuffix(base, ".go"):
			return runAll // shared/core code we can't pin to a resource
		}
	}

	if len(stems) == 0 {
		return ""
	}
	names := make([]string, 0, len(stems))
	for s := range stems {
		names = append(names, s)
	}
	sort.Strings(names)
	return "(?i)^Test(Acc)?(Datasource|Read|Cloud)*(" + strings.Join(names, "|") + ")"
}

// stem reduces a filename to its feature name, e.g.
// "cluster_onoff_acceptance_test.go" -> "clusteronoff", so it matches the
// squashed CamelCase Test* function names case-insensitively.
func stem(base string) string {
	s := strings.TrimSuffix(base, ".go")
	for {
		trimmed := s
		for _, suf := range []string{"_test", "_acceptance", "_datasource", "_schema"} {
			trimmed = strings.TrimSuffix(trimmed, suf)
		}
		if trimmed == s {
			break
		}
		s = trimmed
	}
	return strings.ReplaceAll(s, "_", "")
}
