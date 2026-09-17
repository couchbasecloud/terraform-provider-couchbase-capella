# Auditing and extending acceptance coverage

Read this in the coverage phase. The reason it exists: in this repository it is
common for tests to *look* like coverage for a change while guarding nothing at
all. Counting tests is misleading; the useful question is always the same.

> Would this test have failed if the change were reverted, or made wrong?

If the answer is no, the test is not coverage for this change - and reporting
that honestly is more valuable than a reassuring test count.

## Audit first

Find everything that exercises the changed resource or data source, including
tests that only use it as a fixture for something else:

```bash
grep -rln "<resource_type_name>" acceptance_tests/
grep -rn -A6 "resource \"couchbase-capella_<name>\"" acceptance_tests/
```

Read the configurations, not just the test names. Then check each against the
change. Two patterns account for most false coverage:

**The configuration never reaches the changed code.** Attributes that are always
omitted mean the changed branch is never executed. Before AV-139841, every
database-credential test used `access = [{ privileges = ["data_writer"] }]` and
no test anywhere set `resources`, so the entire bucket/scope/collection mapping
path - the thing the change was about - had no coverage.

Check the *examples* directory too. Deeper shapes often exist there
(`examples/*/terraform.template.tfvars`), which is easy to mistake for coverage,
but examples are not executed by CI.

**The assertion cannot distinguish the old behaviour from the new.** See below;
this one is subtle enough that experienced reviewers miss it.

## Set and List assertions look identical

`terraform-plugin-testing` flattens tuples, lists and sets alike by numeric
index - one branch handles all three in
`internal/configs/hcl2shim/flatmap.go`. So `access.0.privileges.0` is a valid
address whether the attribute is a `Set` or a `List`, and an assertion on it
passes under both.

That has a consequence worth internalising: **a positional assertion only proves
ordering if the configured order differs from the canonical order.** If a test
lists `["default", "metadata"]` or `["_default", "col_a", "col_b"]` - already
alphabetical - then `buckets.0.name` holds the same value under a Set as under a
List, and the test cannot detect a reversion.

So when writing ordering tests, choose values whose configured order is
deliberately *not* the sorted order, and say why in a comment so the next person
does not "tidy" them.

`randomStringWithPrefix` cannot help here: its suffix is random, so the relative
order of two generated names is unknown. Generate a shared random suffix with a
deterministic ordering infix instead:

```go
// Two names sharing a random suffix but sorting a-before-z, so callers can
// configure them in a known non-canonical order.
func orderedFixtureNames(prefix string) (a, z string) {
    suffix := randomString()
    return prefix + "a_" + suffix, prefix + "z_" + suffix
}
```

For fixture values you do not control - the shared bucket names, say - compare at
runtime and order them so the first sorts after the second, rather than
hard-coding an assumption about what they are called.

## Choose assertions that match the declared type

Positional assertions are right for a `List`. For an attribute that is genuinely
a `Set`, order is not meaningful, so assert membership instead - a positional
assertion there is a latent flake:

```go
resource.TestCheckResourceAttr(ref, "access.0.privileges.#", "2"),
resource.TestCheckTypeSetElemAttr(ref, "access.0.privileges.*", "data_reader"),
resource.TestCheckTypeSetElemAttr(ref, "access.0.privileges.*", "data_writer"),
```

Mixed nesting is common in this provider - a `Set` inside a `List` - so check the
schema for each attribute you assert on rather than applying one style
throughout.

## Import needs `ImportStateVerify`

`ImportState: true` alone only proves the import ID parses and the import does
not error. It does not prove the resource round-trips. An import that succeeds
and silently returns an empty attribute passes such a step.

Only `ImportStateVerify: true` compares state before and after import. Most
import steps in this repository do not set it, so "import silently loses data" is
invisible across much of the provider - worth flagging when the change touches an
attribute that import should carry.

```go
{
    ResourceName:            resourceReference,
    ImportStateIdFunc:       generate<Feature>ImportIdForResource(resourceReference),
    ImportState:             true,
    ImportStateVerify:       true,
    ImportStateVerifyIgnore: []string{"password"}, // never returned by GET
}
```

Read the resource's `ImportState` for the composite ID format; several resources
here take a comma-separated key such as
`id=...,cluster_id=...,project_id=...,organization_id=...`.

## Steps come with a free idempotency check

Each `TestStep` runs a plan after its apply and fails on a non-empty plan unless
`ExpectNonEmptyPlan` is set. Every step therefore already asserts idempotency -
so a perpetual diff is caught without writing anything extra, and a second step
repeating the same configuration adds little. Spend the extra step on a real
transition instead: reorder a list, cross a `nil`-to-empty boundary, narrow a
grant from bucket to scope to collection.

## Assert full state after an update

An update step that checks only the field that changed will not notice that
updating field A silently reset field B. Assert the whole shape after each
update; this is where mapping helpers with mirrored gaps get caught.

## Writing the tests

Use the `tf-acceptance-test-gen` skill so structure, naming and helper usage match
the repository. For tests attached to a specific bug, `tf-bugfix` sets the naming
convention `TestAcc<Feature>_AV_XXXXX`.

Verify compilation always - it is fast and catches most mistakes:

```bash
go test -c -o /dev/null ./acceptance_tests
```

Run against a live cluster only when asked. When you have not run them, list the
assumptions that remain unverified - whether the API accepts an empty nested
list, whether it tolerates duplicate entries, what shape a GET actually returns -
so nobody mistakes "compiles" for "passes".

## Documenting a defect you are not fixing

A test that asserts the correct behaviour and is skipped with a ticket reference
is one line away from becoming a regression guard, and it keeps the suite green
meanwhile. The convention here:

```go
t.Skip("AV-12345: <what is wrong today>; unskip once <what must change>")
```

Write the assertions for the *fixed* behaviour, not the buggy behaviour. A test
that pins a bug in place has to be rewritten when the bug is fixed, which is
exactly when nobody wants to be reinterpreting it.
