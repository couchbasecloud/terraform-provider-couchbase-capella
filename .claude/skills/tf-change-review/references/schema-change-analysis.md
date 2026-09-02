# Analysing a schema change

Read this when a diff touches a `*_schema.go` file, a model struct in
`internal/schema/`, or the mapping helpers between them. Schema changes are the
highest-risk category in this provider because they alter how *existing* state is
interpreted, and because several of their failure modes compile and vet cleanly.

## 1. Does the Go model still match the declared type?

Framework reflection is more permissive than people expect, which is why this is
worth checking explicitly rather than assuming.

| Model field | `SetNestedAttribute` | `ListNestedAttribute` |
|---|---|---|
| `[]Access` (plain slice) | works | works |
| `types.Set` | works | **breaks** |
| `types.List` | **breaks** | works |

So a Set to List change is invisible to the compiler when the model uses plain
slices - which is the common style in `internal/schema/`. A clean `go build` says
nothing about whether the change is safe.

## 2. Does every read of the attribute agree with the declaration?

This is the failure mode that hurts most, so hunt for it deliberately. If the
schema declares a `List` while some code path reads the same attribute into a
`types.Set` (or the reverse), the provider builds and vets cleanly and then fails
**every plan** with a framework `Value Conversion Error` - including for
configurations that never touch the changed nesting.

Search every path that reads the attribute, not just CRUD:

```bash
grep -rn "ValidateConfig\|ModifyPlan\|planmodifier\|types.Set\|types.List" \
  internal/resources/<resource>.go internal/schema/<model>.go
```

`ValidateConfig` and `ModifyPlan` are the usual offenders because they are easy
to forget - they are not part of the Create/Read/Update/Delete path a reviewer
naturally traces. This is real: AV-139978 was exactly this, an `access` attribute
declared `ListNestedAttribute` while `ValidateConfig` read it into a `types.Set`.

## 3. Are the resource and data source still consistent?

The same conceptual attribute is often declared twice - once in
`internal/resources/` and once in `internal/datasources/` - and they drift. Drift
is not always a bug, but it is always worth reporting, because practitioners
reasonably expect `privileges` to behave the same in both places and the
generated docs will contradict each other.

Compare them side by side and check the generated docs agree:

```bash
grep -n "Attributes List\|Attributes Set\|List of String\|Set of String" \
  docs/resources/<name>.md docs/data-sources/<name>s.md
```

## 4. Where does the attribute's state come from?

Trace this before reasoning about diffs - it determines which risks apply at all.

**Echoed from plan/state.** `Read` rebuilds the attribute from prior state rather
than the API response. Consequences: configuration order is preserved, so no
API-driven perpetual diff is possible - but `import` cannot populate the
attribute at all, because an imported resource has no prior state. It lands empty.

**Mapped from the API response.** Consequences: import works, but any server-side
reordering or normalisation shows up as a permanent diff. Once an attribute is a
`List`, order is significant, so this becomes a real risk where it was harmless
under a `Set`.

Changing an attribute from Set to List while `Read` maps from the API is the
dangerous combination. Verify the API actually preserves order before accepting
it; do not assume either way.

## 5. Is `nil` distinguished from empty at every level?

Terraform treats an absent list and an empty list as different values. A mapping
helper that collapses one into the other makes the applied state differ from the
plan, and Terraform aborts with:

```
Provider produced inconsistent result after apply
```

The bug hides in nested guards. This shape looks careful but drops a level:

```go
if acc.Resources != nil {
    if acc.Resources.Buckets != nil {          // Resources stays nil when
        access[i].Resources = &Resources{...}  // Buckets is nil
    }
}
```

Given `resources = {}` in configuration, the plan holds
`resources = {buckets = null}` while the applied state holds `resources = null`.
Walk each nesting level and ask what happens when the outer value is present and
the inner one is absent. Fix by assigning the container first, then filling it:

```go
if acc.Resources == nil {
    continue
}
access[i].Resources = &Resources{}   // present, even with no buckets
if acc.Resources.Buckets == nil {
    continue
}
```

Check the request-building helper too. The two directions usually have mirrored
gaps, and fixing only the state side leaves the API receiving nothing.

## 6. Ordering: what a Set/List change actually does

`cty` stores set elements in canonical order, not configuration order. A `List`
preserves configuration order. So the same configuration produces different
state depending on which type the attribute has - and existing state written by
the old provider disagrees with the new schema's reading of it.

Confirm this cheaply rather than asserting it. A throwaway test in the repo makes
it concrete in seconds:

```go
vals := []cty.Value{cty.StringVal("orders"), cty.StringVal("customers")}
setJSON, _ := ctyjson.Marshal(cty.SetVal(vals), cty.Set(cty.String))
listJSON, _ := ctyjson.Marshal(cty.ListVal(vals), cty.List(cty.String))
// set:  ["customers","orders"]   <- reordered
// list: ["orders","customers"]
```

Delete the scratch file afterwards.

Two consequences worth reporting separately:

- **Existing state disagrees with configuration** after the upgrade. Covered in
  `upgrade-and-breaking-changes.md`.
- **Duplicates are no longer collapsed.** A `Set` deduplicated identical elements
  silently; a `List` keeps them and sends them to the API. Usually more correct,
  occasionally a new API error.

## 7. Optionality and validator changes

`Required` becoming `Optional`, or a new `SizeAtLeast`/`ExactlyOneOf`, changes
which configurations are accepted. Loosening is safe; tightening rejects
configurations that used to apply, which is a breaking change even though nothing
about state changed. Check whether any shipped example or acceptance test uses a
configuration the new validators would now reject.

## Isolating a cause when behaviour is puzzling

Two techniques that save a lot of time:

**Compare against a structurally simpler resource.** If a hand-written fixture
behaves oddly, build the same fixture for a resource with no nesting (for example
`couchbase-capella_project`). If that one behaves correctly, the method is sound
and the cause is in the resource under review; if both misbehave, the fault is in
your harness. This separates two explanations that otherwise take hours to tease
apart.

**Do not null a Computed attribute to test whether it matters.** The framework
marks a null Computed attribute as unknown, which is itself a change. The
experiment then reports a diff caused by the experiment rather than by the thing
being tested. Set it to a concrete value in both configuration and state instead,
or bisect by toggling schema flags and rebuilding.
