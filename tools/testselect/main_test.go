package main

import "testing"

func TestSelectPattern(t *testing.T) {
	const all = ".*"
	feature := func(stems string) string {
		return "(?i)^Test(Acc)?(Datasource|Read|Cloud)*(" + stems + ")"
	}

	cases := []struct {
		name  string
		files []string
		want  string
	}{
		{"resource", []string{"internal/resources/cluster.go"}, feature("cluster")},
		{"datasource", []string{"internal/datasources/apikeys.go"}, feature("apikeys")},
		{"acceptance test file", []string{"acceptance_tests/backup_acceptance_test.go"}, feature("backup")},
		{"schema file maps to its resource", []string{"internal/resources/cluster_onoff_schema.go"}, feature("clusteronoff")},
		{"multiple features, sorted and deduped", []string{"internal/resources/user.go", "internal/datasources/users.go", "internal/resources/apikey.go", "internal/resources/apikey_schema.go"}, feature("apikey|user|users")},
		{"core package -> run all", []string{"internal/api/client.go"}, all},
		{"package subdirectory -> run all", []string{"internal/resources/custom_plan_modifiers/x.go"}, all},
		{"shared test helper -> run all", []string{"acceptance_tests/globals.go"}, all},
		{"shared test setup -> run all", []string{"acceptance_tests/setup_test.go"}, all},
		{"go.mod -> run all", []string{"go.mod"}, all},
		{"feature plus core -> run all", []string{"internal/resources/cluster.go", "go.mod"}, all},
		{"nothing test-relevant -> empty", []string{"docs/resources/cluster.md", "README.md"}, ""},
	}

	for _, c := range cases {
		if got := selectPattern(c.files); got != c.want {
			t.Errorf("%s:\n selectPattern(%v)\n = %q\n want %q", c.name, c.files, got, c.want)
		}
	}
}
