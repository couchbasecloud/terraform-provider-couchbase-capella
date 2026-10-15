# What customer escalations actually are

Read this when reviewing any change to a resource's CRUD path, its wait loops,
or its request building - and read it before Phase 4, because it tells you which
cases the acceptance suite is most likely to be missing.

Every class below is drawn from a CBSE escalation or a release-blocking bug on
the Terraform provider, with the ticket numbers kept so you can read the
original. They repeat. A reviewer who checks these seven things will catch more
than one who reads the diff line by line, because each one is invisible in a
diff and obvious in a customer's plan output.

## 1. A zero value sent for a field the practitioner never configured

The single most expensive class. Terraform's null and unknown both arrive in Go
as a zero value, so an unguarded `.ValueInt64()` sends `0`, `.ValueBool()` sends
`false`, and `.ValueString()` sends `""`. The API cannot tell these from a
deliberate choice, so it either rejects the update or silently applies them.

* AV-138744 - bucket update sent empty values for unconfigured fields: the
  update failed on `durabilityLevel`, and would have zeroed `replicas` and TTL.
* AV-143040 - Azure clusters created without `autoexpansion` got auto-expansion
  switched OFF, overriding the server default.
* AV-143042 - the Azure Ultra disk guard checked `Storage.IsUnknown()` while
  assigning `IOPS`, so `iops: 0` was sent.
* AV-137441 - the inverse: no way to send `false` to disable cluster audit log.

The house idiom is to guard both states before reading:

```go
if !plan.Replicas.IsNull() && !plan.Replicas.IsUnknown() {
    body.Replicas = plan.Replicas.ValueInt64()
}
```

**How to find it.** AV-143042 was a *split pair* - one field checked for null,
a different one for unknown, on the same line. That is mechanically detectable,
and the check is quiet enough to run on every diff:

```bash
python3 - <<'PY'
import re, pathlib
isnull = re.compile(r'([A-Za-z_][\w.\[\]]*)\.IsNull\(\)')
isunk  = re.compile(r'([A-Za-z_][\w.\[\]]*)\.IsUnknown\(\)')
for f in pathlib.Path("internal").rglob("*.go"):
    if f.name.endswith("_test.go") or "openapi.gen" in str(f): continue
    for i, ln in enumerate(f.read_text().splitlines()):
        n = {m.group(1) for m in isnull.finditer(ln)}
        u = {m.group(1) for m in isunk.finditer(ln)}
        if n and u and n != u:
            print(f"{f}:{i+1}: IsNull{sorted(n)} != IsUnknown{sorted(u)}")
PY
```

Run against the commit before the AV-143042 fix this prints exactly the one
buggy line; against current `main` it prints two, both benign. Treat a hit as a
question: an asymmetry is deliberate when the code is *testing for absence*
(`!Secret.IsNull() && Secret.IsUnknown()` style rules in `ValidateConfig`), and a
bug when it guards an assignment into a request body.

Then check the other direction, which no grep finds: an `Optional` attribute the
practitioner can never set back to its zero value, because the request builder
skips it when it is `false` or `0`. That was AV-137441.

## 2. An apply that reports success when the operation failed

The customer's `terraform apply` goes green, the resource lands in state, and
the thing does not exist or is broken. Worse than an error, because state now
disagrees with reality and nobody knows.

* AV-138726 - App Service create and scale failures reported as successful
  applies: `IsFinalState` treated failure states as success.
* AV-132323 - cluster create succeeded even when the Capella deployment failed.
* AV-138608 - Data API resource left the cluster reconfigured but untracked when
  Create failed after the update had been accepted.
* AV-138606 - Data API resource wrote contradictory state when the status
  endpoint had not yet left its previous terminal state.

The trap is in the name. These predicates answer *"has it stopped moving?"*, not
*"did it work?"* - they deliberately include the failure states, because a wait
loop must stop on them:

```go
func IsFinalState(state State) bool {
    finalStates := []State{Complete, Failed}    // Failed is a final state
    return slices.Contains(finalStates, state)
}
```

Worse, the same name means opposite things across packages:
`appservice.IsFinalState` includes `Healthy`, while `cluster.IsTerminalState`
lists **only** failure states. Never reason about one from another.

**How to find it.** Every exit from a wait loop needs an explicit failure check
before it returns success:

```bash
grep -rn "IsFinalState(\|IsTerminalState(" internal/resources/ | grep -v _test
grep -rn "func IsFailureState" internal/api/
```

Compare the two lists. `internal/resources/appservice.go` shows the correct
shape - stop on `IsFinalState`, then `if appservice.IsFailureState(...) { return
err }`. Where no `IsFailureState` exists for that package, the call site almost
certainly treats "stopped" as "succeeded"; read it and say so.

## 3. Unknown values surviving an error path into state

When Create fails partway, every computed attribute that was never resolved is
still unknown. Returning early without settling them writes unknowns into state,
which Terraform reports as an inconsistent result or as a permanent diff.

* AV-145075 - `app_endpoint_oidc_provider` left `discovery_url` and
  `username_claim` unknown in state when the post-create refresh failed.
* AV-138146 - `app_service_log_streaming` Create returned unknown computed
  `config_state`/`streaming_state` when the enable wait failed.

**How to find it.** Trace every `return` between the API call and the final
`resp.State.Set`. Each one must either abort without saving state, or null the
computed attributes first - this repository has
`set<Resource>ComputedAttributesToNull` helpers for exactly that. Check that
each call site appends the returned diagnostics; discarding them is its own bug
(AV-145093).

## 4. Import or refresh that drops an attribute and forces replacement

The highest blast radius in the catalogue: the next `terraform apply` destroys
and recreates real infrastructure. Any attribute carrying `RequiresReplace` that
Read populates imperfectly is a candidate.

* AV-134985 - the `INCLUDE MISSING` modifier was dropped from `index_keys` on
  import, causing an unnecessary forced replacement.
* `network_peer` import/read could erase `provider_type`.
* AV-139985 - refresh re-derived `credential_type` from the GET response, so an
  advanced credential whose `userRoles` were absent was destroyed and recreated.
* AV-139981 - a new `credential_type` attribute with a default and
  `RequiresReplace` forced destroy-and-recreate of every existing credential on
  upgrade.

**How to find it.** Cross-reference the two lists - any attribute in both is
worth reading closely:

```bash
grep -rn "RequiresReplace" internal/resources/*_schema.go
```

For each, ask whether Read re-derives it rather than echoing what is stored, and
whether the API returns it losslessly. A normalising API (dropping a modifier,
reordering, lower-casing) plus `RequiresReplace` equals a destroy. When the
attribute cannot be recovered from the API at all, say so plainly: no state
upgrader can invent information the old schema never recorded.

## 5. Validation deferred to apply

A configuration error that surfaces at apply instead of plan costs the customer
a half-finished change. The 1.9.1 cycle alone logged roughly twenty-five of
these, which is why generated validators exist - and why it matters that they
sometimes do not attach.

* AV-132142 / AV-132143 - `terraform validate` accepted an invalid cluster disk
  type and an empty service list, despite the generated OpenAPI enum and
  `minItems` metadata.
* Generated validators did not attach to nested backup schedule fields keyed by
  OpenAPI property paths.
* AV-132147 - free tier cluster name length validation was skipped because a
  manual validator already existed.
* AV-138335 - App Service version-change rejection was enforced only at apply
  time, and bypassed entirely on replacement.

**How to find it.** When a change adds an attribute, check the constraint
actually reached the schema rather than trusting the builder:

```bash
grep -rn "Validators:" internal/resources/<resource>_schema.go
```

The two failure modes are a nested attribute whose OpenAPI path does not match
what the builder looks up, and `builder.go` skipping a generated validator
because a manual one is already present. Both produce a schema that looks
validated and is not. A plan-time test is the cheap proof: a `terraform
validate` acceptance step with `ExpectError` costs nothing and never needs a
cluster.

## 6. Out-of-band deletion and same-name rebinding

Customers delete things in the UI. Read must treat a 404 as "remove from state"
rather than erroring, and recovery logic must not bind state to the wrong
object.

* CBSE-22525 - App Endpoint deleted from underneath Terraform.
* Create-error recovery could bind state to an *older* same-name failed network
  peer.

**How to find it.** In Read, confirm a 404 calls `resp.State.RemoveResource`
rather than raising. Where recovery looks an object up by name, check it
disambiguates by id or creation time - names are not unique over time, and the
stale match is usually the failed predecessor.

## 7. A fix that reintroduces an earlier escalation

Escalation fixes cluster in the same files, so the fix for this month's CBSE
frequently breaks last month's.

* AV-137114 was shipped, then reverted - "No Go" - because it regressed the
  corner case behind CBSE-22798 / AV-134634.
* AV-139978 - the fine-grained RBAC branch conflicted with AV-139841, and
  rebasing onto main broke every `database_credential` plan with a Value
  Conversion Error.
* CBSE-22207 was reverted from a release and had to be re-landed later.

**How to find it.** Before approving a change, ask what else has been fixed in
these files, and check the change does not undo it:

```bash
git log --oneline -15 -- <changed file>
git log --all --oneline --grep="CBSE" -- <changed file>
```

Read the commit message of anything touching the same attribute. When a previous
fix added a guard and this change removes or reorders it, that is the finding -
and it is worth naming the old ticket in the review, because the person who
filed it will be asked whether their case still works.
