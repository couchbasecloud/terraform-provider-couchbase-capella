---
name: tf-change-review
description: Review a change to the Couchbase Capella Terraform provider - a PR number, branch, commit range, or the working-tree diff - for regressions, newly introduced bugs, pre-existing bugs the change exposes, and breaking changes that support and SEs need to relay to customers. Then audit acceptance coverage and add what is missing, draft and file AV bug tickets, and fix what is fixable. Use this whenever the user asks to review a PR or a change, asks whether something causes regressions or breaks anything, asks whether a change is covered by tests, or asks about upgrade impact on existing state - including when they only paste a PR URL or say "check this diff".
---

# Terraform provider change review

You are reviewing a change to the Couchbase Capella Terraform provider as someone
who knows both Capella and the Terraform plugin framework well. The value you add
is not restating the diff: it is working out what the change does to practitioners
who already have state, configuration and pipelines built on the current provider.

Work through the phases below in order. Each phase feeds the next - you cannot
judge acceptance coverage before you know what the change actually risks.

## Ground rule: ask rather than assume

The whole point of this review is to be trusted, and a confident wrong answer
costs more than a question. When something material is genuinely ambiguous, stop
and ask. Things worth asking about, rather than guessing:

- Which change to review, if the user's reference is ambiguous (a branch that has
  moved, a PR with several revisions, "my change" with a dirty working tree).
- The intended semantics when the diff is consistent with more than one intent.
  A `Set` becoming a `List` might be about ordering, about plan output, or about
  working around a framework bug - the right review differs.
- Whether a pre-existing bug you uncover is in scope to fix (see Phase 5).
- The target release or version number for anything version-stamped.
- Whether to run anything against a live Capella organisation.

Do not ask about things you can determine yourself from the repository. Read the
code first; ask about intent and scope.

## Phase 1 - Establish exactly what you are reviewing

Getting this wrong invalidates everything downstream, so spend real effort here.

`gh` is often not installed. Fetch a PR by ref instead:

```bash
git fetch origin "refs/pull/<N>/head:refs/remotes/pr/<N>"
git merge-base main pr/<N>
git diff --stat $(git merge-base main pr/<N>) pr/<N>
```

Then establish three facts that reviewers routinely get wrong:

**Is the change merged, and is it released?** These are different questions with
different consequences. A change can be on `main` yet absent from the newest tag,
which decides whether a breaking-change note belongs in an existing upgrade guide
or a new one:

```bash
git describe --tags --abbrev=0
git show <tag>:<changed-file> | head -40      # what practitioners actually have
git log --oneline <tag>..main -- <changed-file>
```

**Is there related work in flight?** Other branches may already fix, extend or
conflict with the change, and duplicating or contradicting a colleague's work
wastes everyone's time:

```bash
git fetch origin -q
git branch -a --contains <commit>
git branch -r | grep -i "<ticket-id>"
git log --all --oneline --grep="<ticket-id>"
```

If you find a related branch, read its version of the changed code. A refactor
in flight may carry the same defect, which changes your recommendation from
"fix here" to "fix here and port to that branch".

**What is in the diff?** Follow the repository's own review workflow in
`CLAUDE.md` - stage the change, list changed Go files, and never read
`internal/generated/api/openapi.gen.go`; it will consume the context window for
no benefit.

## Phase 2 - Analyse for regressions and bugs

Read every changed Go file against its `main` counterpart, then reason about
consequences rather than syntax. Three questions carry most of the signal:

1. **Where does state come from?** Trace `Create`, `Read` and `Update`. An
   attribute whose state is echoed back from the plan behaves completely
   differently from one mapped out of the API response - the first cannot drift
   from configuration but breaks import, the second can produce a perpetual diff
   if the API reorders or normalises anything. This single distinction explains
   most database-credential behaviour in this provider.

2. **Does the declared schema still match every place the code reads it?** A
   type declared one way and read another compiles and vets cleanly, then fails
   at runtime. This is the highest-severity class of defect in this codebase.

3. **Is `nil` distinguished from empty at every nesting level?** Terraform treats
   an absent list and an empty list as different values, and a mapping helper
   that collapses one into the other produces "Provider produced inconsistent
   result after apply" - an error that aborts the apply.

When the diff touches a schema file, read
`references/schema-change-analysis.md`. It covers the specific failure modes of
type changes (Set/List, nesting, optionality), how to verify model compatibility,
resource/datasource parity, and how to reproduce an ordering effect cheaply.

Then run the repository's verification loop from `CLAUDE.md`: `goimports -w
-local github.com/couchbasecloud/terraform-provider-couchbase-capella`, `go vet`,
and the version-stamped `go build`. Report failures with output rather than
paraphrasing them.

Distinguish clearly, throughout, between three categories, because the right
response differs for each:

- **Introduced** - the change causes it. Fix it in this change.
- **Exposed** - pre-existing, but the change makes it reachable or visible.
- **Adjacent** - pre-existing and untouched, found while reading. Worth a ticket;
  usually not worth expanding this change to fix.

## Phase 3 - Breaking changes and upgrade impact

This is the phase most reviews skip, and the one support and SEs actually need.

A change that alters how *existing state* is interpreted breaks practitioners who
change nothing at all. Acceptance tests are structurally blind to it, because
`terraform-plugin-testing` builds state from scratch with the provider under test
on every run, so state written by an older version never exists.

Read `references/upgrade-and-breaking-changes.md` when the change touches schema
types, state mapping, defaults, or anything a practitioner has already persisted.
It covers the two-tier upgrade harness in `upgrade_tests/`, how to author a
fixture that actually proves something, and the traps that produce false passes
and phantom diffs.

For any breaking change, establish and state four things - a diff alone is not
enough for someone writing customer comms:

- **Who is affected**, concretely: which resources, and which configuration
  shapes. "Any credential with more than one access block" is useful;
  "users of database credentials" is not.
- **What they see**: the actual plan output, quoted.
- **Whether it is self-healing** - clears after one apply - or permanent. This is
  the difference between a release note and an emergency.
- **Whether it can be migrated away.** A state upgrader receives only prior
  state, so information the old schema never recorded cannot be recovered. Say
  so plainly when that is the case, rather than implying a fix is possible.

Then write the customer-facing note using the template in that reference. Keep it
free of framework vocabulary: the reader is a support engineer, not a provider
developer.

## Phase 4 - Acceptance coverage

Audit before you write. Read `references/acceptance-coverage.md` - it explains
how to tell real coverage from tests that only look like coverage, which is the
common case here and the reason this phase exists.

Enumerate every existing test that touches the changed resource or data source,
then judge each against the change: would this test have failed if the change
were reverted or made wrong? A test that passes either way is not coverage, and
saying so is more useful than counting tests.

Write the missing tests with the `tf-acceptance-test-gen` skill so they match the
repository's conventions. Always verify they compile:

```bash
go test -c -o /dev/null ./acceptance_tests
```

Run them against a live cluster only when the user asks. They provision real
Capella resources and can take an hour, so surprise runs are unwelcome; say
plainly which assumptions remain unverified until someone does run them.

For a defect you are documenting but not fixing, a test that asserts the correct
behaviour and is skipped with a ticket reference is far more valuable than no
test - it is one line from becoming a regression guard:

```go
t.Skip("AV-12345: <what is wrong>; unskip once <what must change>")
```

## Phase 5 - Bug tickets

Draft a ticket for each confirmed finding, show them all to the user, and create
them in the **AV** project via the Atlassian MCP only after the user approves.
Filing is outward-facing and lands in someone's triage queue, so a misjudged
finding has a real cost.

Each draft needs: a summary naming the symptom; affected versions (use the
merged-versus-released facts from Phase 1); a minimal reproducing configuration;
the root cause with `file:line`; expected versus actual behaviour, quoting real
output where you have it; and whether it is introduced, exposed or adjacent.

Verify a finding before drafting a ticket for it. Reproducing it - a scratch
test, a plan against a fixture, a live run - is worth the effort, both because
unverified tickets erode trust and because the reproduction usually sharpens the
root cause. Say explicitly when something is reasoned but unreproduced.

## Phase 6 - Fixes

Present each bug with its proposed fix and let the user decide, one at a time.
Regressions the change introduced are usually clear-cut; pre-existing bugs expand
the scope of someone's PR, which is their call and not yours.

Use the `tf-bugfix` skill for approved fixes so root-cause analysis and test
naming follow the repository's conventions, including the
`TestAcc<Feature>_AV_XXXXX` test name format.

Keep fixes minimal and structural. When a mapping helper mishandles one branch,
fix that branch rather than rewriting the helper - a reviewer comparing against
the reported bug should see the correspondence immediately.

If you experiment on the source to isolate a cause, restore it and prove you did:

```bash
git diff --quiet <file> && echo "restored"
```

## Report

Lead with the verdict, then detail. Someone should be able to read the first two
sections and know whether the change is safe to merge.

```markdown
## Verdict
<Ship / ship with notes / do not ship, and the single reason why>

## What the change does
<2-4 sentences: mechanism, not a diff restatement>

## Regressions and bugs
<Per finding: introduced | exposed | adjacent; severity; file:line;
 what breaks and for whom; verified how>

## Breaking changes - for customer, support and SE comms
<The customer-facing note. Omit this section entirely if there are none.>

## Acceptance coverage
<What exists, what it actually guards, what does not guard anything and why,
 what was added>

## Tickets
<Drafted / created, with keys>

## Fixes
<Applied, with files; and what was deliberately left alone>

## Unverified
<Anything reasoned but not reproduced, and what would settle it>
```

An empty findings section is a good outcome, not a failure - say so directly
rather than inflating minor observations to fill the template. Equally, do not
soften a real breaking change: someone downstream is going to have to explain it
to a customer, and they need it stated plainly.
