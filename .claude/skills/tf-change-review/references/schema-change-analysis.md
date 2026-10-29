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

The apply is not merely noisy. The resource is still written to state, and it is
written **tainted**, so the next `terraform apply` destroys and recreates
infrastructure the practitioner never asked to touch. Rate this class high
severity even when the mismatched value itself looks cosmetic, and always say in
the report whether recovery needs `terraform untaint`.

### Why it keeps happening in this provider

Models in `internal/schema/` come in two styles, and only one is safe by
construction:

| Model field style | how null and `[]` differ | risk |
|---|---|---|
| `types.List` / `types.Set` | distinct values, tested with `IsNull()` | explicit, hard to get wrong |
| `[]T` plain slice | nil slice versus empty slice | erased by ordinary Go idioms |

With a plain slice the distinction survives only as long as nobody writes
`len(x) == 0` and nobody lets `append` build the slice up from nil - both of
which are unremarkable Go that reads as correct. That is why this needs a
deliberate sweep rather than a careful read: the defect looks like good code.

### The three idioms that erase it

**A `len(...) == 0` guard**, which conflates "absent" with "present but empty":

```go
if b == nil || (len(b.Buckets) == 0 && len(b.Urls) == 0) {
    return nil, nil          // configuration said {}, state now says null
}
```

**`append` onto a nil accumulator.** When the loop runs zero times the slice is
still nil, and nil becomes a null list rather than `[]`:

```go
bindings := &Bindings{}                            // Urls is nil
for _, u := range b.Urls {                         // zero iterations
    bindings.Urls = append(bindings.Urls, ...)
}                                                  // still nil -> null in state
```

`make([]T, 0, len(src))` is the safe form: zero iterations still yields `[]`.
Both idioms frequently sit in the same file, which is a useful tell - if one
helper uses `make` and another builds from nil, the second is probably the bug.

### Where the nil ends up decides whether it matters

A nil slice is only a defect if it reaches state as null, and the framework
constructors disagree about that. This is the single thing that separates a real
finding from the dozens of harmless hits the sweep returns, so check it on every
hit rather than reasoning from the `append` alone:

| what consumes the nil slice | result | verdict |
|---|---|---|
| `types.ListValue` / `types.SetValue(t, nil)` | **empty** `[]` | safe |
| `types.ListValueFrom` / `types.SetValueFrom(ctx, t, nil)` | **null** | bug |
| a plain `[]T` field with a `tfsdk` tag | **null** | bug |

The `ValueFrom` constructors reflect over the Go value and preserve its nil-ness;
the direct constructors treat a nil element slice as zero elements. So `append`
onto nil is harmless when the result is handed to `types.SetValue` - which is why
`MorphAllowedCidrs` in `internal/schema/apikey.go` is fine - and a real defect
when it is handed to `SetValueFrom` or assigned straight to a model field.

Do not recall which is which; it is four lines and the answer is not intuitive:

```go
var nilStrs []types.String
lf, _ := types.ListValueFrom(ctx, types.StringType, nilStrs)  // IsNull=true
lv, _ := types.ListValue(types.StringType, nil)               // IsNull=false, len 0
```

**A nested guard that drops a level.** This shape looks careful but leaves the
container nil whenever the outer value is present and the inner one is absent:

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

### Sweep for it mechanically

The first two idioms are greppable, so sweep the changed files instead of hoping
to notice them. This lists every `append` accumulator in the diff and reports the
ones that were never `make`-initialised:

```bash
for f in $(git diff --cached main --name-only -- '*.go' | grep -v openapi.gen.go); do
  grep -oE '[A-Za-z_][A-Za-z0-9_.]*[[:space:]]*=[[:space:]]*append\(' "$f" |
    sed -E 's/[[:space:]]*=[[:space:]]*append\(//' | sort -u |
  while read -r v; do
    grep -qE "(^|[^A-Za-z0-9_])${v##*.}[[:space:]]*:?=[[:space:]]*make\(" "$f" ||
      echo "$f: $v is nil when the loop does not run"
  done
done
```

Then, over the same file list, `grep -n 'len(.*) == 0'` and `grep -n '!= nil'`
inside the mapping helpers.

Expect most hits to be noise; the sweep is wide on purpose. Three filters, in
this order, take a repository-wide run down to a handful:

1. **Does the nil reach state as null?** Apply the constructor table above. This
   discards the majority on its own.
2. **Can the practitioner write the empty value?** A `Computed`-only attribute
   has no configuration to contradict, so null versus `[]` cannot be an
   inconsistent result - it is at most a cosmetic wart. A `Required` attribute
   carrying `SizeAtLeast(1)` cannot reach the empty case either. The dangerous
   profile is `Optional` **without** `Computed`, where state must round-trip the
   configured value exactly.
3. **Is it on a state path at all?** An accumulator that feeds an API request has
   no plan to be inconsistent with.

Then read the validators before concluding, because they often cover part of the
surface and not the rest - and the part they miss is the finding. AV-145137 is
exactly that: `ValidateConfig` rejected an empty `cors.origin`, but returned early
when `cors.disabled` was true and never looked at `login_origin` or `headers`.

### Prove it offline, in seconds

These are ordinary Go functions, so you never need a cluster or an hour-long
acceptance run to settle whether the collapse is real. Call the helper directly
from a scratch test in its own package - unexported helpers are reachable there -
and print what comes back for the empty case:

```go
got, err := bindingsToSchema(&eventingapi.Bindings{
    Buckets: []eventingapi.BucketBinding{{Alias: "src", Bucket: "travel"}},
})
t.Logf("Urls nil=%t Constants nil=%t", got.Urls == nil, got.Constants == nil)
```

Anything that prints `nil=true` for an attribute the configuration set to `[]` is
a confirmed inconsistent-result bug, and you can say so in the ticket without
hedging. Delete the scratch file afterwards and prove you did:

```bash
git status --short
```

This is real: AV-145093 is exactly the above. `bindingsToSchema`
(`internal/schema/eventing_function.go:396`) returns `nil` when all three binding
lists are empty, and builds `Buckets`, `Urls` and `Constants` by appending onto
nil - so `bindings.urls = []` applies, errors, and taints the function, while
`newEventingFunctionBindingsObject` in the same file uses `make` and is correct.

### Check the mirrored directions

Fixing only the state side leaves two matching gaps, and both were live in
AV-145093. Look at all three directions before calling the class closed:

- **The request builder.** If state-mapping drops empty collections, the helper
  building the API payload usually drops them too, so the server never learns the
  practitioner cleared the field.
- **Change-detection helpers.** This provider hand-rolls `<thing>Changed()`
  predicates to decide what to PUT, and the usual shape is `if plan == nil {
  return false }`. That reads as "nothing to send" but means "clearing is not a
  change": the PUT omits the field, the server keeps the old value, and the
  read-back contradicts the plan. `eventingBindingsChanged`
  (`internal/resources/eventing_function.go:248`) does this. Grep the diff for
  `func [a-zA-Z]*Changed(` and test each predicate against a cleared plan.

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
