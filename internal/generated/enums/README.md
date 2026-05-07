# Enums Package

Generated enum constants derived from the OpenAPI spec, indexed by
schema name and field path.

## Files

| File | Purpose |
|------|---------|
| `generate/main.go` | Generator entry point (`package main`). |
| `generate/discover.go` | Walks the OpenAPI document and collects enum sites. |
| `generate/generator.go` | Renders discovered sites as Go source. |
| `enums.gen.go` | Generated output. **Do not edit by hand.** |

## Regenerating

Run whenever `openapi.generated.yaml` changes:

```bash
make gen-enums
```

Commit the updated `enums.gen.go`.

## Generated contents

The generated file declares an `EnumDef` type and a package-private
`enumTable` indexed by `(schemaName, fieldPath)`:

```go
type EnumDef struct {
    Type    string   // "string" or "integer"
    Values  []string // integer values are stringified
    IsArray bool     // true when the property is an array; values constrain each element
}

var enumTable = map[string]map[string]EnumDef{
    "AllowedCidr": {
        "status": {Type: "string", Values: []string{"active", "expired"}},
        "type":   {Type: "string", Values: []string{"temporary", "permanent"}},
    },
    "APIKeyResourcesItems": {
        "roles": {Type: "string", IsArray: true, Values: []string{"projectOwner", ...}},
    },
    "CreateBucketRequest": {
        "vbuckets": {Type: "integer", Values: []string{"128", "1024"}},
    },
    ...
}
```

Field paths are dot-joined from the schema root; array element paths
drop the trailing `[]` (e.g. `roles`, not `roles.[]`) and set
`IsArray: true`.

Excluded from the table:

- Top-level enum schemas (no parent property).
- Parameter enums (path/query parameter scope).
- Non-string, non-integer enums.
