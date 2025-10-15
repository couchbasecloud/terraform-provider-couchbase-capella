# OpenAPI Documentation Package

This package provides automatic documentation enhancement by extracting field descriptions from the OpenAPI specification.

## Usage

### Using the SchemaBuilder Interface

The `SchemaBuilder` interface is defined in `internal/schema/builder.go` and can be used by both **resources** and **data sources**.

Each resource or data source should create a `SchemaBuilder` instance at the package level:

```go
// At the top of your schema file (e.g., project_schema.go)
import (
    capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectBuilder = capellaschema.NewSchemaBuilder("project")

func ProjectSchema() schema.Schema {
    return schema.Schema{
        Attributes: map[string]schema.Attribute{
            "name": projectBuilder.WithOpenAPIDescription(
                stringAttribute([]string{required}),
                "name",  // Only the field name!
            ).(*schema.StringAttribute),
        },
    }
}
```

**Key Benefits:**
- ✅ Shared between resources and data sources
- ✅ Resource name defined once per file
- ✅ Interface enforces contract across all schemas
- ✅ Type-safe at compile time
- ✅ Clean, readable code

The `SchemaBuilder` automatically:
- Converts `field_name` from snake_case to camelCase (`fieldName`)
- Tries common OpenAPI schema patterns:
  - `CreateResourceRequest`
  - `GetResourceResponse`
  - `UpdateResourceRequest`
  - etc.
- Returns the first matching field description found

### For Terraform-specific fields

Continue using the existing `WithDescription()` function:

```go
"if_match": WithDescription(
    stringAttribute([]string{optional}),
    "Custom description for Terraform-specific field"
)
```

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
```go
import (
    capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// At the top of the file
var projectBuilder = capellaschema.NewSchemaBuilder("project")

// In ProjectSchema()
"name": projectBuilder.WithOpenAPIDescription(
    stringAttribute([]string{required}),
    "name",
).(*schema.StringAttribute),

"description": projectBuilder.WithOpenAPIDescription(
    stringAttribute([]string{optional}),
    "description",
).(*schema.StringAttribute),
```

### Bucket Resource (in internal/resources/bucket_schema.go)
```go
import (
    capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// At the top of the file
var bucketBuilder = capellaschema.NewSchemaBuilder("bucket")

// In BucketSchema()
"memory_allocation_in_mb": bucketBuilder.WithOpenAPIDescription(
    int64Attribute([]string{optional}),
    "memory_allocation_in_mb",  // Automatically converts to memoryAllocationInMb
).(*schema.Int64Attribute),

"storage_backend": bucketBuilder.WithOpenAPIDescription(
    stringAttribute([]string{optional}),
    "storage_backend",  // Automatically converts to storageBackend
).(*schema.StringAttribute),
```

### Data Source Example (in internal/datasources/users_schema.go)
```go
import (
    capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// At the top of the file
var usersBuilder = capellaschema.NewSchemaBuilder("users")

// In UsersSchema()
"email": usersBuilder.WithOpenAPIDescription(
    stringAttribute([]string{optional}),
    "email",
).(*schema.StringAttribute),
```

## Maintaining the OpenAPI Spec

The `openapi.generated.json` file in this directory is embedded at compile time. 

To update it:
```bash
cp ../../openapi.generated.json .
```

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
- Ensure `openapi.generated.json` exists in this directory
- Run `go mod tidy` if dependencies are missing

