# Upgrade impact and breaking changes

Read this when a change touches schema types, state mapping, defaults, validators
or anything else a practitioner has already persisted.

The population at risk is people who upgrade the provider and change *nothing
else*. They are invisible to acceptance tests, because
`terraform-plugin-testing` creates state from scratch with the provider under
test on every run - state written by an older version never exists in that
harness. `upgrade_tests/` exists to close exactly this gap; read its README for
the full mechanics.

## Two tiers

**`upgrade_tests/upgrade_test.go`** - offline, a few seconds, no credentials and
no network. Builds the working tree, serves it through a Terraform filesystem
mirror as a synthetic version, and runs `terraform plan -refresh=false
-detailed-exitcode` against a checked-in state fixture. Exit `0` means no
changes, `2` means a diff. Cheap enough to gate every change, and it is what
proves an ordering or state-interpretation break.

This works without credentials because the provider's `Configure` only reads
configuration and environment variables before constructing an HTTP client, so a
placeholder token reaches the plan phase, and `-refresh=false` skips the only
step that would call the API.

**`upgrade_tests/upgrade-check.sh`** - end-to-end against a live organisation.
Installs a released version from the registry, applies, swaps in the working
tree, re-plans, and on a diff applies once more to report whether the diff is
self-healing or permanent. This is the only tier that can see an API
response-shape change or server-side normalisation. Recommend it when the finding
depends on real API behaviour; do not run it without asking.

```bash
FROM_VERSION=<last release with the OLD behaviour> \
TF_VAR_host=... TF_VAR_auth_token=... TF_VAR_organization_id=... \
TF_VAR_project_id=... TF_VAR_cluster_id=... TF_VAR_bucket_name=... \
upgrade_tests/upgrade-check.sh
```

`FROM_VERSION` must be a release that actually has the old behaviour. Confirm
with `git show <tag>:<file>` rather than assuming the newest tag predates the
change - a change on `main` is often not in the newest tag.

## Adding a fixture

Add a directory under `upgrade_tests/testdata/` with `main.tf` and
`terraform.tfstate`, and register it in the `fixtures` table with
`expectEmptyPlan` and a `reason`. A fixture expecting a diff records a known,
accepted break; its `reason` has to say why, so the next reader can tell a
deliberate decision from a bug.

Add both kinds when you can. A fixture in the *new* shape expecting a clean plan
is the forward-looking guard; a fixture in the *old* shape recording the break is
the documentation. Together they also prove the harness can distinguish the two -
a harness that only ever reports diffs is worthless.

The most reliable fixture is a captured one: apply with the released provider,
copy the resulting `terraform.tfstate`, then edit `terraform_version` down to
something old, because a state claiming a version newer than the running CLI is
rejected outright.

### Put every list in non-canonical order

This is the whole point of the fixture and the easiest thing to get wrong.
Because `cty` stores sets canonically, a fixture whose lists happen to be sorted
already will produce identical state under both Set and List - so its assertions
pass either way and it guards nothing. Order elements so configuration order and
sorted order differ, and say so in a comment.

### `sensitive_attributes` is not optional

Terraform records each schema-sensitive attribute in the instance:

```json
"sensitive_attributes": [[{ "type": "get_attr", "value": "password" }]]
```

Omit it and the prior value carries no sensitive mark while the planned value
does. The values are identical but the marks are not, so Terraform plans an
in-place update and renders it with **no changed attributes** - a phantom diff
that looks exactly like a provider bug and will send you hunting in the wrong
place. Any resource with a `Sensitive` attribute needs this line.

### Do not commit `.terraform/` or `.terraform.lock.hcl`

The lock file records a checksum of the provider binary, and the binary is
rebuilt on every run, so a committed lock makes `init` fail with "the cached
package ... does not match any of the checksums recorded in the dependency lock
file". The driver copies only regular files into a fresh temp directory, which
keeps this out of the fixtures.

Note also that `.gitignore` excludes `terraform.tfstate` repo-wide, with a
negation for `upgrade_tests/testdata/**/terraform.tfstate`. A fixture placed
outside that path is silently untracked - and an untracked fixture means a
harness that passes locally and is broken for everyone else.

## The `dev_overrides` trap

Most contributors have a `dev_overrides` block in `~/.terraformrc` pointing at a
local build - it is the normal way to develop against the working tree. It makes
Terraform run that binary regardless of which version it just downloaded, and
Terraform reports "Installed v1.9.1" while running your local code.

In an upgrade check this silently turns the released-version phase into a second
run of the working tree, and the comparison passes for the wrong reason. A false
pass is worse than a failure because nobody investigates it.

Both tiers set `TF_CLI_CONFIG_FILE` on every invocation, which *replaces* the
user's config rather than merging with it. If you write any new Terraform-driving
check, do the same.

`terraform version` will not warn you: it reports the version selected in the
lock file even while `dev_overrides` substitutes a different binary. The only
reliable signal is the `Provider development overrides are in effect` warning.
Treat it as fatal.

## Characterising the break

Four things, because a diff alone is not enough for someone writing customer
communications:

- **Who is affected**, in configuration terms. "Any credential with more than one
  access block, or more than one collection listed in a non-alphabetical order"
  tells a support engineer whether a given customer is at risk. "Users of
  database credentials" does not.
- **What they see** - the real plan output, quoted.
- **Self-healing or permanent.** Apply once and re-plan. If it clears, this is a
  release note; if it persists, it is an emergency.
- **Whether it can be migrated away.** A state upgrader receives only prior
  state. Information the old schema never recorded - such as the order the
  practitioner wrote elements in, when canonical order overwrote it - cannot be
  recovered, so no upgrader can fix it. Say that plainly instead of leaving the
  impression that a migration is merely unimplemented.

## Customer-facing note

Write it for a support engineer, not a provider developer. No framework
vocabulary - no `cty`, no `SetNestedAttribute`, no "plan merge".

```markdown
### <Resource>: first plan after upgrading to <version> may show changes

**Who is affected:** <concrete configuration shapes>

**What you see:** on the first `terraform plan` after upgrading, <resource> shows
an in-place update even though nothing in the configuration changed. <Quote the
shape of the diff.>

**Why:** <one or two sentences, no framework jargon>

**What to do:** <the action - typically: review the plan, confirm it only
reorders, and apply once. State that the diff clears after a single apply and
that no Capella-side change occurs.>

**Is anything at risk?** <Say plainly whether permissions, data or availability
are affected. If nothing is, say nothing is.>
```

Version-stamped artefacts need the target release number, which you generally
cannot infer - ask. Note that upgrade guides in `docs/guides/` appear to be
produced by `scripts/generate-upgrade-guide.sh`, so check that before
hand-writing one.
