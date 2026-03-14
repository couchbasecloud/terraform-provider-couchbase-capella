## Project Overview

Terraform Provider for Couchbase Capella. Written in Go using the HashiCorp Terraform Plugin Framework (`terraform-plugin-framework`). Manages Couchbase Capella cloud resources (clusters, buckets, projects, users, app services, etc.) via the Capella Public API.

## Agent Guidelines

Behavioral rules to reduce common AI coding mistakes. Biased toward caution over speed.

### Think Before Coding

- State assumptions explicitly. If uncertain, ask.
- If multiple approaches exist, present the tradeoffs — don't pick silently.
- If something is unclear, stop and ask. Don't guess at Capella API behavior or Terraform lifecycle semantics.

### Simplicity First

- No features beyond what was asked.
- No abstractions for single-use code.
- No speculative "flexibility" or "configurability."
- If you write 200 lines and it could be 50, rewrite it.
- Ask: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

### Surgical Changes

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- If you notice unrelated issues (dead code, linting problems), mention them — don't fix them.
- Remove imports/variables/functions that YOUR changes made unused. Don't remove pre-existing dead code.
- Every changed line should trace directly to the user's request.

## Repository Structure

```
main.go                          # Provider entrypoint
internal/
  provider/provider.go           # Provider configuration, registers all resources & data sources
  resources/                     # Terraform resource implementations (CRUD lifecycle)
    <resource>.go                # Resource logic (Create/Read/Update/Delete)
    <resource>_schema.go         # Resource schema definition
  datasources/                   # Terraform data source implementations (read-only)
    <datasource>.go              # Data source logic
    <datasource>_schema.go       # Data source schema definition
  schema/                        # Shared schema types, builder, validation helpers
  api/                           # API client layer (v1 hand-written client, sub-packages per resource)
  generated/api/                 # Auto-generated OpenAPI client (oapi-codegen) - DO NOT EDIT openapi.gen.go
  errors/errors.go               # Centralized error constants
  docs/                          # OpenAPI spec loader for auto-generated field descriptions
acceptance_tests/                # Acceptance tests (create real Capella resources)
examples/                        # HCL example configurations per resource
templates/                       # Doc templates; `make build-docs` copies to docs/
docs/                            # Generated documentation - do not edit directly; edit templates/ instead
openapi.generated.yaml           # OpenAPI spec - source of truth for generated client and descriptions
version/version.go               # Provider version constant
```

## Key Patterns

### Adding a New Resource

Every resource has two files:
- `internal/resources/<name>.go` — implements `resource.Resource` with CRUD methods
- `internal/resources/<name>_schema.go` — defines the Terraform schema using `SchemaBuilder` and `AddAttr` from `internal/schema`

Corresponding data source follows the same pattern in `internal/datasources/`.

Register new resources/data sources in `internal/provider/provider.go` in the `Resources()` or `DataSources()` method.

### API Client

Two API client layers coexist:
- **v1** (`internal/api/`): Hand-written HTTP client with per-resource sub-packages (e.g., `api/cluster/`, `api/bucket/`)
- **v2** (`internal/generated/api/`): Auto-generated via `oapi-codegen` from `openapi.generated.yaml`. Regenerate with `make gen-api`. **Never manually edit `openapi.gen.go`.**


## Build & Test Commands

```bash
make setup        # Install dev tools (golangci-lint, goimports, misspell)
make build        # Format code + build binary to ./bin/
make fmt          # gofmt + goimports
make lint-fix     # golangci-lint --fix
make test         # Unit tests (excludes acceptance_tests/)
make testacc      # Acceptance tests (requires TF_VAR_auth_token, TF_VAR_host, TF_VAR_organization_id)
make check        # Full quality gate: fmt + tffmt + docs-lint + lint-fix + test
make build-docs   # Generate provider docs from templates/ to docs/
make gen-api      # Regenerate OpenAPI client from openapi.generated.yaml
```

## Linting & Formatting

- Go linting: `golangci-lint` v1.64.8 (config in `.golangci.yml`); enabled linters include `gosec`, `errcheck`, `staticcheck`, `govet`, `misspell`, and others
- Go formatting: `gofmt -s` + `goimports` (with local module prefix `github.com/couchbasecloud/terraform-provider-couchbase-capella`)
- Terraform formatting: `terraform fmt -recursive`
- Doc linting: `misspell` on `docs/`

## Testing

- **Unit tests**: `make test` — runs with `-short -cover -race`, sets `CAPELLA_OPENAPI_SPEC_PATH` to the project root spec file
- **Acceptance tests**: `make testacc` — creates real cloud resources; requires env vars `TF_VAR_auth_token`, `TF_VAR_host`, `TF_VAR_organization_id`. Uses `TF_ACC=1` flag, 120m timeout
- Test framework: `github.com/stretchr/testify` for assertions, `terraform-plugin-testing` for acceptance tests

## CI (GitHub Actions)

- **PR**: unit tests (`make vet` + `make test`), golangci-lint (new issues only), terraform fmt check
- **PR + push to main**: acceptance tests (secrets from GitHub environment)
- **Tag push**: GoReleaser builds + signs + publishes to Terraform Registry

## Documentation

- Edit templates in `templates/`, NOT `docs/` (docs/ is generated)
- Run `make build-docs` to regenerate docs from templates + schema definitions
- Upgrade guides go in `templates/guides/<version>-upgrade-guide.md`

## Conventions

- Commit messages: `[AV-XXXXX] Description` (Jira ticket prefix)
- PR labels: `enhancement`, `bug`, `breaking-change`, `documentation`, `no-changelog-needed`
- Vendored dependencies (`go.mod` with `-mod=vendor` flag in `GOFLAGS`)
- Module path: `github.com/couchbasecloud/terraform-provider-couchbase-capella`
- Go 1.24+, Terraform >= 1.5.2
