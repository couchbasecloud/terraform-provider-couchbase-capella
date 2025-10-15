# Terraform Provider Capella 

This is the repository for Couchbase's Terraform-Provider-Capella which forms a Terraform plugin for use with Couchbase Capella.

## Requirements

- [Git](https://git-scm.com/)
- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.2
- [Go](https://golang.org/doc/install) >= 1.21

### Environment

- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.

## Using the Provider

To use a released provider in your Terraform environment, run `terraform init` and Terraform will automatically install the provider.
Documentation about the provider specific configuration options can be found on the [provider's website](https://developer.hashicorp.com/terraform/language/providers).

## Contributing to the Provider
See [Contributing.md](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/CONTRIBUTING.md)

## Discovering New API features

Most of the new features of the provider are using [capella-public-apis](https://docs.couchbase.com/cloud/management-api-guide/management-api-intro.html)
Public APIs are updated automatically, tracking all new Capella features.

## Generated API client

This repository includes an OpenAPI-generated client in `internal/generated/api` (keeping the existing hand-written client in `internal/api` for backward compatibility).

Generate/update the client before working on a new resource or data source:

1) Ensure the generator is installed:

   `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`

2) Regenerate from the root `openapi.generated.yaml`:

   `make gen-api`

The command writes the client/types to `internal/generated/api/openapi.gen.go`.

Notes:
- Provider wiring makes both clients available:
  - `providerschema.Data.ClientV1`: legacy HTTP client (`internal/api`)
  - `providerschema.Data.ClientV2`: generated client with typed methods (`internal/generated/api`)
- When adding a new resource/data source, prefer calling `ClientV2` for new endpoints and migrate incrementally.

## Documentation

### Generating Documentation

Provider documentation is automatically generated from Go schema definitions and custom templates:

```bash
make build-docs
```

This command:
1. Generates documentation using `tfplugindocs`
2. Copies custom content from `templates/` to `docs/`
3. **Automatically reads** OpenAPI spec from `openapi.generated.yaml` at runtime

> **Note:** The OpenAPI spec is loaded from the filesystem at runtime, so documentation is always up-to-date with the latest spec changes - no rebuild needed!

### OpenAPI-Enhanced Descriptions

Field descriptions are automatically extracted from the OpenAPI spec using the `SchemaBuilder` and a generic function:

```go
// In internal/resources/project_schema.go
var projectBuilder = capellaschema.NewSchemaBuilder("project")

"name": capellaschema.WithOpenAPIDescription(
    projectBuilder,
    stringAttribute([]string{required}),
    "name",  // Automatically finds description in OpenAPI spec
),  // ✨ No type assertion needed - fully type-safe!
```

This provides rich documentation with:
- Detailed field descriptions
- Validation constraints (min/max length, patterns)
- Enum values and defaults
- Format specifications (UUID, date-time, etc.)

See `internal/docs/README.md` for implementation details.

### Custom Documentation

To preserve custom content (upgrade guides, custom pages):

1. **Add files to `templates/` directory:**
   ```
   templates/
     ├── guides/              # Upgrade guides
     ├── index.md.tmpl        # Custom homepage
     └── ...
   ```

2. **Run `make build-docs`** - your custom content will be preserved

#### Why `templates/` is at Root, Not Under `docs/`

The directory structure follows Terraform ecosystem conventions:

```
/
├── templates/          ← SOURCE (input files you edit)
├── docs/               ← OUTPUT (generated, DO NOT EDIT)
└── examples/
```

**Separation of Concerns:**
- `templates/` = What you **write** (source files)
- `docs/` = What `tfplugindocs` **generates** (output)

When `tfplugindocs` runs, it:
1. Clears/updates the `docs/` directory
2. Generates docs from Go schemas
3. Copies from `templates/` → `docs/`

If `templates/` were inside `docs/`, it would be deleted during generation! This pattern is used by all official Terraform providers (AWS, Google, Azure, etc.) and is the default behavior of `tfplugindocs`.

⚠️ **Important:** Never edit files in `docs/` directly - they will be overwritten!

- ✅ **DO** edit files in `templates/` for custom content
- ✅ **DO** edit Go schema files for field descriptions  
- ✅ **DO** run `make build-docs` after changes
- ❌ **DON'T** manually edit generated files in `docs/`

See `templates/README.md` for more details.
