# Droid Runbook: Generating Terraform Data Sources

A practical guide for using AI droids (coding agents) to generate Terraform data sources for the Couchbase Capella provider.

## Overview

This provider uses a **skill-based droid system** (located in `.factory/skills/`) to automate the creation of Terraform data sources from the OpenAPI spec. Instead of hand-writing boilerplate, you describe what you need and the droid generates the data source code, schema, tests, and provider registration.

## Requirements

- [Git](https://git-scm.com/)
- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.2
- A coding agent / AI droid with access to the workspace
- The OpenAPI spec (`openapi.generated.yaml`) at the project root

## Quick Start

### 1. Identify the Feature

Determine which Capella API resource you want to expose as a Terraform data source. Check the OpenAPI spec for available endpoints:

```bash
# Search for your feature in the spec
grep -i "your-feature" openapi.generated.yaml
```

You need to identify:
- **GET endpoint** — for fetching a single resource (e.g., `GET /v4/organizations/{id}/projects/{id}/clusters/{id}/buckets/{id}`)
- **LIST endpoint** — for fetching all resources (e.g., `GET /v4/organizations/{id}/projects/{id}/clusters/{id}/buckets`)

If only one endpoint exists, only that data source will be generated.

### 2. Invoke the Droid

Point the droid at the skill file and provide the feature context:

> "Generate Terraform data sources for **[Feature Name]** using the `tf-datasource-gen` skill and the OpenAPI spec."

The droid will follow the instructions in `.factory/skills/tf-datasource-gen/SKILL.md` to generate:

| Artifact | Location | Purpose |
|---|---|---|
| `feature.go` | `internal/datasources/` | Single-resource data source (Read) |
| `features.go` | `internal/datasources/` | List data source (Read) |
| `feature_schema.go` | `internal/datasources/` | Schema definition |
| `features_schema.go` | `internal/datasources/` | Schema definition (list) |
| API structs | `internal/api/` | Request/response types |
| Provider registration | `internal/provider/provider.go` | Wires up the data source |
| Acceptance tests | `acceptance_tests/` | End-to-end tests |

### 3. Implementation Steps (What the Droid Generates)

The droid follows the steps defined in `.factory/skills/tf-datasource-gen/SKILL.md`. Below is the full breakdown of each implementation step so you can verify — or manually replicate — the output.

#### Step 1 — Create Data Source Files

Place data source code in `internal/datasources/`. Two files are created depending on available endpoints:

| Endpoint available | File to create | Purpose |
|---|---|---|
| GET (single resource) | `feature.go` (e.g., `bucket.go`) | Fetch one resource by ID |
| LIST (all resources) | `features.go` (e.g., `buckets.go`) | Fetch all resources |

If only one endpoint exists in the spec, only that file is created.

#### Step 2 — Create Schema Files

Each data source gets its own schema file using the naming pattern `<feature>_schema.go` / `<features>_schema.go`.

**Add validation** for `organization_id`, `project_id`, and `cluster_id` (if present) using `requiredStringWithValidator()`:

```go
capellaschema.AddAttr(attrs, "organization_id", snapshotBackupBuilder, requiredStringWithValidator())

func requiredStringWithValidator() *schema.StringAttribute {
    return &schema.StringAttribute{
        Required:   true,
        Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
    }
}
```

#### Step 3 — Define the Data Source Struct

Create a struct named after the feature that embeds the shared `Data` struct:

```go
type Buckets struct {
    *providerschema.Data
}
```

#### Step 4 — Add the `New` Constructor Function

Every data source needs a `New` function that returns a `datasource.DataSource`:

```go
func NewBuckets() datasource.DataSource {
    return &Buckets{}
}
```

#### Step 5 — Assert Interface Compliance

Use compile-time nil assertions to guarantee the struct implements the required interfaces:

```go
var (
    _ datasource.DataSource              = (*Buckets)(nil)
    _ datasource.DataSourceWithConfigure = (*Buckets)(nil)
)
```

#### Step 6 — Implement `Metadata`

Return the Terraform type name (provider name + feature suffix):

```go
func (d *Buckets) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_buckets"
}
```

#### Step 7 — Implement `Configure`

Extract the shared provider data and store it on the struct:

```go
func (d *Buckets) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    data, ok := req.ProviderData.(*providerschema.Data)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.Data = data
}
```

#### Step 8 — Generate API Structs

Create request/response structs in `internal/api/` to model the JSON payloads returned by the Capella API. Derive field names and types from the OpenAPI spec.

#### Step 9 — Use `ClientV1` for API Calls

All new data sources **must** use `ClientV1` with retry logic:

```go
response, err := s.ClientV1.ExecuteWithRetry(...)
```

> **Important:** If the droid generates code using an older client, update it to `ClientV1`.

#### Step 10 — Register in the Provider

Add the new `New` function(s) to `internal/provider/provider.go`:

```go
func (p *capellaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        // ...existing data sources...
        datasources.NewBucket,
        datasources.NewBuckets,
    }
}
```

#### Step 11 — Create Acceptance Tests

Add acceptance tests in `acceptance_tests/` with the naming pattern `<feature>_acceptance_test.go` (e.g., `buckets_acceptance_test.go`).

#### Step 12 — Use Parallel Tests

All acceptance tests must use `resource.ParallelTest()` for parallel execution:

```go
func TestAccBucketsDataSource(t *testing.T) {
    resource.ParallelTest(t, resource.TestCase{
        // ...
    })
}
```

### 4. Review the Generated Code

After the droid finishes, verify the output matches the steps above and project conventions. Key things to check:

- [ ] **Struct** embeds `*providerschema.Data`
- [ ] **Interface assertions** exist for `datasource.DataSource` and `datasource.DataSourceWithConfigure`
- [ ] **`Metadata`** returns `req.ProviderTypeName + "_feature"`
- [ ] **`Configure`** extracts `*providerschema.Data` with proper error handling
- [ ] **Schema validators** are present for `organization_id`, `project_id`, `cluster_id`
- [ ] **API structs** in `internal/api/` match the OpenAPI spec
- [ ] **`ClientV1`** is used for all API calls (not an older client)
- [ ] **Provider registration** in `provider.go` includes both `New` functions
- [ ] **Acceptance tests** exist and use `resource.ParallelTest()`

### 5. Run the Code Review Checklist

Follow `AGENTS.md` to validate the generated code:

```bash
# Step 1 — Check changed files
git diff main --name-only -- '*.go' ':!internal/generated/api/openapi.gen.go'

# Step 2 — Format imports
goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella internal/datasources/your_feature.go

# Step 3 — Vet
go vet ./internal/datasources/...

# Step 4 — Build
VERSION=$(git describe --tags --abbrev=0)
go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella
```

Repeat steps 2–4 until clean. If errors persist after 5 retries, report them.

### 6. Run Acceptance Tests

> ⚠️ **Acceptance tests create real resources and may cost money.**

```bash
make testacc
```

Tests should use `resource.ParallelTest()` for parallel execution.

## File Naming Conventions

| Type | Naming Pattern | Example |
|---|---|---|
| Single resource datasource | `feature.go` | `snapshot_backup.go` |
| List datasource | `features.go` (plural) | `snapshot_backups.go` |
| Schema (single) | `feature_schema.go` | `snapshot_backup_schema.go` |
| Schema (list) | `features_schema.go` | `snapshot_backups_schema.go` |
| Acceptance test | `feature_acceptance_test.go` | `snapshot_backup_acceptance_test.go` |

## Project Structure Reference

```
internal/
├── api/                  # API structs and client helpers
├── datasources/          # All data source implementations + schemas
├── provider/             # Provider registration (provider.go)
├── resources/            # Resource implementations (not covered here)
└── schema/               # Shared schema helpers (Data struct, builders)

acceptance_tests/         # Acceptance tests for all resources & data sources

.factory/
├── droids/               # Droid configurations
└── skills/
    └── tf-datasource-gen/
        └── SKILL.md      # The skill definition the droid follows
```

## Troubleshooting

| Problem | Solution |
|---|---|
| Droid uses old API client | Tell it to use `ClientV1` explicitly |
| Missing schema validators | Add `requiredStringWithValidator()` for org/project/cluster IDs |
| Build fails after generation | Run `goimports` then `go vet`, fix, repeat up to 5 times |
| Data source not appearing in Terraform | Check it's registered in `provider.go` `DataSources()` |
| `openapi.gen.go` consuming context | Never read that file — it's excluded per `AGENTS.md` |

## Tips

- **One feature at a time.** Don't ask the droid to generate multiple unrelated data sources in one pass.
- **Provide the endpoint paths.** The more specific you are about which OpenAPI endpoints to use, the better the output.
- **Check the diff.** Always review `git diff` before committing generated code.
- **Iterate.** If the first pass isn't perfect, point the droid at the specific issue and ask it to fix just that.

