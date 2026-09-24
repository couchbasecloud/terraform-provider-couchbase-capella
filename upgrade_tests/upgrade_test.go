// Package upgrade_tests checks that state written by an older build of the
// provider still plans cleanly against the current one.
//
// Acceptance tests cannot cover this by construction: terraform-plugin-testing
// creates state from scratch with the provider under test on every run, so no
// state written by a previous version ever exists. That leaves a whole class of
// change invisible to CI - most importantly a schema type change, such as the
// Set -> List move in AV-139841, which alters how existing state is interpreted
// rather than what the API is asked to do.
//
// These tests need no Capella credentials and make no network calls. The
// provider's Configure only reads configuration and environment variables
// before constructing an HTTP client, so a placeholder token is enough to reach
// the plan phase, and `plan -refresh=false` skips the one step that would
// otherwise call the API. What is exercised is the part that matters here:
// Terraform decoding prior state with the current schema and diffing it against
// the configuration.
package upgrade_tests

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// mirrorVersion is the synthetic version the locally built provider is served
// under. It must be a version no real release will ever carry, so that a
// fixture's version constraint can only resolve to the binary built from the
// working tree.
const mirrorVersion = "99.0.0"

const providerSource = "registry.terraform.io/couchbasecloud/couchbase-capella"

type fixture struct {
	// dir is the testdata directory holding main.tf and terraform.tfstate.
	dir string

	// expectEmptyPlan is whether planning the current provider against this
	// fixture's prior state should report no changes. A fixture that expects a
	// diff is pinning a known, accepted break rather than asserting good
	// behaviour, and must say why in reason.
	expectEmptyPlan bool

	// reason explains what the fixture pins down, and is printed on failure.
	reason string
}

var fixtures = []fixture{
	// Both fixtures hold the same credential and the same configuration, and
	// differ only in the element order recorded in prior state. Under the Set
	// schema the provider ships today, order carries no meaning, so both must
	// plan clean - that is the property being asserted.
	//
	// Keeping both is what makes this a type guard. If these attributes ever
	// become Lists again, order becomes significant: the config-ordered fixture
	// would still plan clean, while the canonically-ordered one would diff, and
	// this suite would fail. That is the regression AV-139841 shipped in v1.11.0
	// and AV-142331 reverted in v1.11.1.
	{
		dir:             "dbcred_list_ordered_state",
		expectEmptyPlan: true,
		reason: "prior state recording elements in configuration order must plan clean; a diff means " +
			"a change altered how existing database credential state is read",
	},
	{
		dir:             "dbcred_set_ordered_state",
		expectEmptyPlan: true,
		reason: "prior state recording elements in canonical order - which is what a Set-schema " +
			"provider persists - must also plan clean, because a Set does not distinguish orderings. " +
			"This is the fixture that fails if these attributes return to List, since a List would " +
			"then diff against the configuration's own order",
	},
}

func TestProviderUpgradeFromPriorState(t *testing.T) {
	tfBin, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("terraform CLI not on PATH; skipping provider upgrade tests")
	}

	cliConfig := providerMirror(t)

	for _, f := range fixtures {
		t.Run(f.dir, func(t *testing.T) {
			t.Parallel()

			workDir := t.TempDir()
			copyDir(t, filepath.Join("testdata", f.dir), workDir)

			env := terraformEnv(cliConfig)

			if code, out := terraform(t, tfBin, workDir, env, "init", "-input=false", "-no-color"); code != 0 {
				t.Fatalf("terraform init failed (exit %d):\n%s", code, out)
			}

			// -detailed-exitcode: 0 = no changes, 1 = error, 2 = changes.
			code, out := terraform(t, tfBin, workDir, env,
				"plan", "-refresh=false", "-detailed-exitcode", "-input=false", "-no-color")

			switch code {
			case 0:
				if !f.expectEmptyPlan {
					t.Fatalf("expected a non-empty plan but got none.\n%s\n\nplan output:\n%s", f.reason, out)
				}
			case 2:
				if f.expectEmptyPlan {
					t.Fatalf("upgrading the provider against this prior state produces a diff.\n%s\n\nplan output:\n%s", f.reason, out)
				}
				t.Logf("known upgrade diff reproduced as expected.\n%s\n\nplan output:\n%s", f.reason, out)
			default:
				t.Fatalf("terraform plan failed (exit %d):\n%s", code, out)
			}
		})
	}
}

// providerMirror builds the working tree's provider and lays it out as a
// filesystem mirror, returning the path of a CLI config file that points at it.
// The config excludes this provider from direct installation, so terraform init
// cannot silently fall back to the registry and test a released binary instead
// of the local one.
func providerMirror(t *testing.T) string {
	t.Helper()

	base := t.TempDir()
	mirror := filepath.Join(base, "mirror")
	platformDir := filepath.Join(mirror, "registry.terraform.io", "couchbasecloud",
		"couchbase-capella", mirrorVersion, runtime.GOOS+"_"+runtime.GOARCH)
	if err := os.MkdirAll(platformDir, 0o750); err != nil {
		t.Fatalf("creating mirror layout: %v", err)
	}

	binary := filepath.Join(platformDir, "terraform-provider-couchbase-capella_v"+mirrorVersion)
	if prebuilt := os.Getenv("CAPELLA_UPGRADE_PROVIDER_BINARY"); prebuilt != "" {
		copyFile(t, prebuilt, binary, 0o700)
	} else {
		// #nosec G204 -- the only variable is an output path this test just built.
		build := exec.Command("go", "build", "-o", binary, ".")
		build.Dir = repoRoot(t)
		if out, err := build.CombinedOutput(); err != nil {
			t.Fatalf("building provider: %v\n%s", err, out)
		}
	}

	cliConfig := filepath.Join(base, "cli.tfrc")
	contents := fmt.Sprintf(`provider_installation {
  filesystem_mirror {
    path    = %q
    include = [%q]
  }
  direct {
    exclude = [%q]
  }
}
`, mirror, providerSource, providerSource)
	if err := os.WriteFile(cliConfig, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing CLI config: %v", err)
	}

	return cliConfig
}

// terraformEnv builds the environment for the Terraform child process from an
// allow-list instead of inheriting the caller's.
//
// Inheriting would let a developer's or a runner's settings decide the result,
// which is the opposite of what a gate needs. CAPELLA_GLOBAL_API_REQUEST_TIMEOUT
// below the provider's 300s floor makes Configure fail outright; TF_CLI_ARGS can
// inject flags; TF_DATA_DIR and TF_WORKSPACE move where state is read from; and
// TF_LOG changes the output this test inspects. None of those should influence
// whether an upgrade regression is reported.
func terraformEnv(cliConfig string) []string {
	env := []string{
		"TF_CLI_CONFIG_FILE=" + cliConfig,
		"TF_IN_AUTOMATION=1",
		"TF_INPUT=0",
		"CHECKPOINT_DISABLE=1",
		// Configure never calls the API, so placeholders suffice. The timeout is
		// pinned at the provider's minimum so no inherited value can fail it.
		"CAPELLA_HOST=https://cloudapi.cloud.couchbase.com",
		"CAPELLA_AUTHENTICATION_TOKEN=upgrade-test-placeholder",
		"CAPELLA_GLOBAL_API_REQUEST_TIMEOUT=300",
	}

	// HOME and TMPDIR give Terraform its plugin and temporary directories, PATH
	// covers anything it shells out to, and the two Windows variables keep this
	// usable there. Everything else is deliberately dropped.
	for _, key := range []string{"HOME", "PATH", "TMPDIR", "SystemRoot", "USERPROFILE"} {
		if value, ok := os.LookupEnv(key); ok {
			env = append(env, key+"="+value)
		}
	}

	return env
}

// terraform runs the CLI and returns its exit code and combined output.
func terraform(t *testing.T, tfBin, dir string, env []string, args ...string) (int, string) {
	t.Helper()

	// #nosec G204 -- tfBin comes from exec.LookPath and args are literals
	// assembled in this file; driving the real CLI is the point of the test.
	cmd := exec.Command(tfBin, args...)
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		// A non-zero exit is a result here, not a failure: plan uses the exit
		// code to report whether there are changes.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), string(out)
		}
		t.Fatalf("running terraform %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return 0, string(out)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolving working directory: %v", err)
	}
	return filepath.Dir(wd)
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", src, err)
	}
	for _, e := range entries {
		// Directories are skipped, .terraform among them. A lock file is skipped
		// too: it pins a checksum of a provider binary that is rebuilt on every
		// run, so copying one makes init fail. A fixture should not carry either,
		// but skipping here means one that does still works.
		if e.IsDir() || e.Name() == ".terraform.lock.hcl" {
			continue
		}
		copyFile(t, filepath.Join(src, e.Name()), filepath.Join(dst, e.Name()), 0o600)
	}
}

// copyFile copies src to dst with the given mode. The mode is a parameter
// because the two callers differ: fixtures only need to be readable, while the
// provider binary has to keep its executable bit or terraform cannot run it.
func copyFile(t *testing.T, src, dst string, mode os.FileMode) {
	t.Helper()

	// #nosec G304 G703 -- src is either a path under testdata or one the operator
	// supplied via CAPELLA_UPGRADE_PROVIDER_BINARY.
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("reading %s: %v", src, err)
	}
	// #nosec G306 G703 -- 0o700 for the provider binary is deliberate, see above.
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("writing %s: %v", dst, err)
	}
}
