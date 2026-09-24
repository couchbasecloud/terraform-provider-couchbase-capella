# Provider upgrade tests

Checks that configuration and state created by an **older** provider version
still behave when the provider is upgraded underneath them.

Acceptance tests cannot cover this. `terraform-plugin-testing` creates state
from scratch with the provider under test on every run, so state written by a
previous version never exists, and a change that alters how existing state is
*interpreted* — rather than what the API is asked to do — is invisible to CI.
The Set → List move in AV-139841 is exactly that kind of change.

## Two tiers

| | [`upgrade_test.go`](upgrade_test.go) | [`upgrade-check.sh`](upgrade-check.sh) |
|---|---|---|
| Runs | every PR | before a release, or when a change touches how state is read |
| Needs credentials | no | yes, a live Capella organisation |
| Needs network | no | yes, the registry and the V4 API |
| Runtime | ~3s plus one `go build` | a few minutes |
| Covers | state decode and plan merge | the whole lifecycle: install, apply, upgrade, re-plan, re-apply |

Tier 1 is cheap enough to gate every change and catches the common case. Tier 2
is the only one that can see an API response-shape change or a server-side
normalisation, because it is the only one that actually refreshes and applies.

---

## Tier 1 — offline plan check

```bash
go test ./upgrade_tests/
```

No credentials and no network. The provider's `Configure` only reads
configuration and environment variables before constructing an HTTP client, so a
placeholder token is enough to reach the plan phase, and `plan -refresh=false`
skips the one step that would call the API.

The tests skip themselves if the `terraform` CLI is not on `PATH`. Set
`CAPELLA_UPGRADE_PROVIDER_BINARY` to an already-built provider to skip the build.

### How it works

1. The working tree is built and laid out as a Terraform **filesystem mirror**
   under a synthetic version, `99.0.0`. The generated CLI config excludes this
   provider from `direct` installation, so `terraform init` cannot quietly fall
   back to the registry and test a released binary instead of the local one.
2. Each fixture in `testdata/` supplies a `main.tf` and a `terraform.tfstate`
   representing state as some earlier provider wrote it.
3. `terraform plan -refresh=false -detailed-exitcode` runs against it.
   Exit `0` means no changes, `2` means a diff.

### Adding a fixture

Add a directory under `testdata/` with `main.tf` and `terraform.tfstate`, then
register it in the `fixtures` table in `upgrade_test.go` with `expectEmptyPlan`
and a `reason`. A fixture that expects a diff is pinning a known, accepted break
rather than asserting good behaviour, so the `reason` has to say why.

The easiest way to get a realistic state file is to capture one: apply the
configuration with the released provider version and copy the resulting
`terraform.tfstate`. Then edit `terraform_version` down to something old — a
state claiming a version newer than the CLI running the test is rejected
outright.

**Put every list in non-canonical order.** This is the whole point of the
fixtures and the easiest thing to get wrong. `cty` stores set elements in
canonical order, so a config listing `["orders", "customers"]` round-trips
through a Set as `["customers", "orders"]` but through a List as
`["orders", "customers"]`. If a fixture's lists happen to already be sorted,
its positional assertions pass under **both** types and it guards nothing.

### Two things that will bite you when hand-writing state

**`sensitive_attributes` is not optional.** Terraform records each
schema-sensitive attribute there:

```json
"sensitive_attributes": [[{ "type": "get_attr", "value": "password" }]]
```

Leave it as `[]` and the prior value carries no sensitive mark while the planned
value does. The values are identical, but the marks are not, so Terraform plans
an in-place update and renders it with *no changed attributes* — a confusing
phantom diff that looks like a provider bug. Every resource with a `Sensitive`
attribute needs this.

**Do not commit `.terraform/` or `.terraform.lock.hcl` into a fixture.** The lock
file records a checksum of the provider binary, and the binary is rebuilt on
every run, so a committed lock makes `init` fail with "the cached package ...
does not match any of the checksums recorded in the dependency lock file".

The driver copies each fixture into a fresh temp directory and runs `init`
there, skipping directories and any `.terraform.lock.hcl` it finds, so a fixture
that carries one still works. Do not rely on that: a lock file in `testdata` is
a mistake either way, and only these two names are filtered.

Note that `.gitignore` excludes `terraform.tfstate` repo-wide; there is a
negation for `upgrade_tests/testdata/**/terraform.tfstate` so these fixtures are
tracked. A new fixture outside that path will be silently untracked.

---

## Tier 2 — end-to-end upgrade check

```bash
FROM_VERSION=1.10.0 \
TF_VAR_host=https://cloudapi.cloud.couchbase.com \
TF_VAR_auth_token=... \
TF_VAR_organization_id=... \
TF_VAR_project_id=... \
TF_VAR_cluster_id=... \
TF_VAR_bucket_name=default \
upgrade_tests/upgrade-check.sh
```

`FROM_VERSION` is the released version to upgrade *from*. Reusing an existing
project and cluster keeps the run short; the script creates and destroys only a
database credential.

### Phases

1. Install `FROM_VERSION` from the registry and apply a configuration whose
   nested lists are all in non-canonical order.
2. Build the working tree and serve it as `99.0.0` through a filesystem mirror.
3. Re-plan with the upgraded provider against the **unchanged** configuration.
   A non-empty plan here is what a practitioner sees after upgrading.
4. On a diff, apply once and re-plan, to report whether the diff is
   **self-healing** or **permanent**. That distinction is what belongs in the
   release notes.

Exit codes: `0` clean upgrade, `1` setup or apply failure, `2` upgrade
regression. Test resources are destroyed by an `EXIT` trap.

---

## A note on `dev_overrides`

Both tiers set `TF_CLI_CONFIG_FILE` for **every** Terraform invocation, which
replaces the user's `~/.terraformrc` rather than merging with it. That is
deliberate. A `dev_overrides` block for this provider — the normal way to
develop against a local build, and present in most contributors' `~/.terraformrc`
— makes Terraform run the local binary no matter which version it just
downloaded. In an upgrade check that turns the released-version phase into
another run of the working tree, and the comparison passes for the wrong reason.

`terraform version` does not help you notice: it reports the version selected in
the lock file even while `dev_overrides` is substituting a different binary. The
only reliable signal is the `Provider development overrides are in effect`
warning, which `upgrade-check.sh` greps for and treats as fatal rather than
reporting a meaningless pass.

## Known result: AV-139841

Confirmed with both tiers. Upgrading from v1.10.0 or earlier, where `access`,
`buckets`, `scopes` and `collections` were Sets, produces a **non-empty first
plan** for any database credential whose configuration lists those elements in
an order other than the canonical one. v1.10.0's own create plan reorders the
blocks, which is what leaves its state disagreeing with the configuration.

The diff is **self-healing**: `Update` sends the same effective permissions and
rewrites state in configuration order, so it clears after a single apply and
destroys nothing.

A state upgrader cannot fix this. An upgrader receives only prior state, and the
canonical order recorded there carries no information about the order the
practitioner wrote. The upgrade is therefore something to document, not
something to migrate.
