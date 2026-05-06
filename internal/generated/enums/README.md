# Enums Package

This package provides generated enum constants and a lookup map derived from the OpenAPI spec.

## Files

| File | Purpose |
|------|---------|
| `gen/main.go` | Generator entry point (`package main`). |
| `gen/discover.go` | Walks the OpenAPI document and collects enum sites. |
| `gen/generator.go` | Renders discovered sites as Go source. |
| `gen/*_test.go` | Unit tests for the generator. |
| `enums.gen.go` | Generated output — one typed `var` per enum site and a `Lookup` map. **Do not edit by hand.** |

## Regenerating

Run whenever `openapi.generated.yaml` changes:

```bash
make gen-enums
```

To run the generator unit tests:

```bash
go test ./internal/generated/enums/gen/
```

Commit the updated `enums.gen.go`.

## Contents of the generated file

### Per-site variables

One variable per discovered enum site, named after its location in the spec:

```go
// Source: components.schemas.BucketConflictResolution
var BucketConflictResolution = []string{"seqno", "lww"}

// Source: components.schemas.CreateScheduledBackupRequest.properties.weeklySchedule.properties.dayOfWeek
var CreateScheduledBackupRequest_WeeklySchedule_DayOfWeek = []string{"sunday", "monday", ...}
```

Integer enums use `[]int64`:

```go
// Source: components.schemas.CreateBucketRequest.properties.vbuckets
var CreateBucketRequest_Vbuckets = []int64{128, 1024}
```

### Lookup map

`Lookup` indexes string enum values by `(OpenAPI schemaName, dotPath)`. The path is
the dot-joined property chain from the schema root to the enum site, so nested
fields keep their full path:

```go
vals := enums.Lookup["AllowedCidr"]["status"]
// ["active", "expired"]

vals := enums.Lookup["CreateScheduledBackupRequest"]["weeklySchedule.dayOfWeek"]
// ["sunday", "monday", "tuesday", ...]
```

Array element paths use the property name only — the trailing `[]` is stripped
(e.g. `roles`, not `roles.[]`).

Top-level enum schemas (no parent property), parameter enums, and integer enums
are not included in `Lookup` — reference their named `var` directly.
