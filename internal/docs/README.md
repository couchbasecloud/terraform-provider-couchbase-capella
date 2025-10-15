# OpenAPI Documentation Package

This package provides automatic documentation enhancement by extracting field descriptions from the OpenAPI specification.

## Usage

### Using the SchemaBuilder

The `SchemaBuilder` type is defined in `internal/schema/builder.go` and can be used by both **resources** and **data sources**.

Each resource or data source should create a `SchemaBuilder` instance at the package level:

```go
// At the top of your schema file (e.g., project_schema.go)
import (
    capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectBuilder = capellaschema.NewSchemaBuilder("project")

func ProjectSchema() schema.Schema {
    attrs := make(map[string]schema.Attribute)
    
    // AddAttr automatically finds descriptions from OpenAPI OR common registry
    capellaschema.AddAttr(attrs, "name", projectBuilder, stringAttribute([]string{required}))
    capellaschema.AddAttr(attrs, "organization_id", projectBuilder, stringAttribute([]string{required}))
    capellaschema.AddAttr(attrs, "if_match", projectBuilder, stringAttribute([]string{optional}))
    
    // Special attributes that are pre-built
    attrs["audit"] = computedAuditAttribute()
    
    return schema.Schema{
        Attributes: attrs,
    }
}
```

**Key Benefits:**
- Shared between resources and data sources
- Resource name defined once per file
- **No field name duplication** - `AddAttr` takes the field name once
- **Uniform syntax** - all fields use the same pattern
- Type-safe at compile time with generics
- Clean, readable code

### Description Resolution Priority

`AddAttr` automatically finds descriptions using a three-tier fallback:

**1. OpenAPI Path Parameters (highest priority for *_id fields)**
- For fields ending in `_id` (e.g., `organization_id`, `project_id`, `cluster_id`)
- Looks up in `components.parameters` section of OpenAPI spec
- Converts snake_case → CapitalizedCamelCase (e.g., `organization_id` → `OrganizationId`)
- Returns the path parameter description

**2. OpenAPI Schema Properties (primary for other fields)**
- Converts `field_name` from snake_case to camelCase (`fieldName`)
- Tries common OpenAPI schema patterns:
  - `CreateResourceRequest`
  - `GetResourceResponse`
  - `UpdateResourceRequest`
  - `ResourceRequest`, `ResourceResponse`
- Extracts rich constraints (maxLength, pattern, enum, etc.)
- Formats as readable markdown with examples

**3. Common Descriptions Registry (fallback for special fields)**
- For fields not in OpenAPI spec:
  - **HTTP headers**: `if_match`, `etag`
  - **Special nested attributes**: `audit`
- Defined in `internal/schema/common_descriptions.go`
- Minimal set - most descriptions come from OpenAPI now

**4. Empty** (if not found anywhere)

This means **ALL fields can use `AddAttr`** with consistent syntax, and most descriptions come directly from the OpenAPI spec!

## What Gets Enhanced

The OpenAPI loader automatically extracts and formats:

- **Base description** from OpenAPI field description
- **Constraints:**
  - Maximum/minimum length
  - Maximum/minimum values
  - Pattern validation
- **Valid Values:** Enum options
- **Default:** Default values
- **Format:** UUID, date-time, email, etc.
- **Deprecation warnings**

## Example Output

**Input (OpenAPI spec):**
```json
{
  "name": {
    "type": "string",
    "description": "The name of the project (up to 128 characters).",
    "maxLength": 128
  }
}
```

**Output (Generated docs):**
```markdown
- name (String) The name of the project (up to 128 characters).

**Constraints:**
  - Maximum length: 128 characters
```

## Examples

### Project Resource (in internal/resources/project_schema.go)

**Before** (mixed patterns):
```go
Attributes: map[string]schema.Attribute{
    "name": WithOpenAPIDescription(..., "name"),
    "organization_id": WithDescription(..., "The GUID4 ID..."),
    "if_match": WithDescription(..., "A precondition header..."),
}
```

**After** (uniform `AddAttr` for everything):
```go
var projectBuilder = capellaschema.NewSchemaBuilder("project")

func ProjectSchema() schema.Schema {
    attrs := make(map[string]schema.Attribute)
    
    capellaschema.AddAttr(attrs, "name", projectBuilder, stringAttribute([]string{required}))
    capellaschema.AddAttr(attrs, "organization_id", projectBuilder, stringAttribute([]string{required}))
    capellaschema.AddAttr(attrs, "if_match", projectBuilder, stringAttribute([]string{optional}))
    capellaschema.AddAttr(attrs, "etag", projectBuilder, stringAttribute([]string{computed}))
    
    // Only special pre-built attributes are different
    attrs["audit"] = computedAuditAttribute()
    
    return schema.Schema{Attributes: attrs}
}
```

**All fields use the same pattern!** Descriptions are automatically found from OpenAPI or the common registry.

### Bucket Resource (in internal/resources/bucket_schema.go)
```go
var bucketBuilder = capellaschema.NewSchemaBuilder("bucket")

func BucketSchema() schema.Schema {
    attrs := make(map[string]schema.Attribute)
    
    // Automatically converts memory_allocation_in_mb → memoryAllocationInMb
    capellaschema.AddAttr(attrs, "memory_allocation_in_mb", bucketBuilder,
        int64Attribute([]string{optional}))
    
    // Automatically converts storage_backend → storageBackend
    capellaschema.AddAttr(attrs, "storage_backend", bucketBuilder,
        stringAttribute([]string{optional}))
    
    return schema.Schema{Attributes: attrs}
}
```

### Data Source Example (in internal/datasources/users_schema.go)
```go
var usersBuilder = capellaschema.NewSchemaBuilder("users")

func UsersSchema() schema.Schema {
    attrs := make(map[string]schema.Attribute)
    
    capellaschema.AddAttr(attrs, "email", usersBuilder,
        stringAttribute([]string{optional}))
    
    return schema.Schema{Attributes: attrs}
}
```

## How OpenAPI Spec Loading Works

The OpenAPI spec is **loaded from the filesystem at runtime**, not embedded in the binary. This means:
- Always uses the latest `openapi.generated.yaml` from the project root
- No need to copy or sync files
- No need to rebuild when the spec changes
- Simpler architecture

### Finding the Spec

The loader automatically finds `openapi.generated.yaml` by:
1. **Checking environment variable**: `CAPELLA_OPENAPI_SPEC_PATH` (if set)
2. **Walking up directories**: Looking for `go.mod` to find project root
3. **Reading from project root**: `<project_root>/openapi.generated.yaml`

If the spec can't be found, the provider gracefully degrades - it still works, but field descriptions won't be enhanced with OpenAPI metadata.

### Environment Variable (for special cases)

If running from a non-standard location, set:
```bash
export CAPELLA_OPENAPI_SPEC_PATH="/path/to/openapi.generated.yaml"
```

The `make build-docs` target automatically sets this for you.

**Note:** We use YAML (not JSON) as the single source of truth. The `kin-openapi` library parses YAML directly, so no conversion is needed.

## Testing

Run tests to verify the OpenAPI loader:
```bash
cd internal/docs
go test -v
```

## Troubleshooting

**Empty descriptions in generated docs?**
- Check that the schema name matches the OpenAPI spec exactly
- Check that the field name matches (use camelCase as in OpenAPI)
- Run tests to verify the loader can find the field

**Build errors?**
- Ensure `openapi.generated.yaml` exists in this directory
- Run `go mod tidy` if dependencies are missing
- If you see "cannot embed irregular file", ensure the YAML file is a regular file, not a symlink

