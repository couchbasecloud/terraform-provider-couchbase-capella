# Custom Droids for terraform-provider-capella

This directory contains custom droids configured for the terraform-provider-capella project.

## Available Droids

### capella-terraform-provider-generator

Automates adding new Couchbase Capella API resources and datasources to the provider.

#### What It Does

Given a YAML configuration file with API specifications, this droid generates:

1. **Unit tests** (`internal/resources/<resource>_test.go`) - Fast tests with mocked HTTP client (TDD-first development)
2. **Acceptance tests** (`acceptance_tests/<resource>_acceptance_test.go`) - E2E tests against real Capella API
3. **Internal schema types** (`internal/schema/<resource>.go`) - Go structs mapping to Terraform attributes
4. **Resource implementation** (`internal/resources/<resource>.go`) - Implements `resource.Resource` with CRUD methods
5. **Resource schema** (`internal/resources/<resource>_schema.go`) - Uses `SchemaBuilder` from `internal/schema`
6. **Datasource implementation** (`internal/datasources/<plural>.go`) - Implements `datasource.DataSource` interface
7. **Datasource schema** (`internal/datasources/<plural>_schema.go`) - Schema definition for datasource
8. **Provider registration** (`internal/provider/provider.go`) - Updates `Resources()` and `DataSources()` methods

**TDD Workflow:**

1. Generate unit tests (they fail - RED)
2. Generate resources/datasources
3. Unit tests pass (GREEN) - fast feedback cycle
4. Generate acceptance tests (E2E validation)
5. Acceptance tests pass - final verification
6. Validation: make fmt, make lint-fix, make build

**TODO - Future enhancements**:

- Example directory generation (`examples/<resource>/`)
- Documentation template generation
- Validation rule generation

#### Prerequisites

The droid requires `openapi.generated.yaml` in the repository root to extract field descriptions and types.

**Getting the OpenAPI spec:**

The spec lives in the sibling `couchbase-cloud` repository. The droid will automatically check for and create a symlink:

```bash
ln -s ../couchbase-cloud/openapi.yaml openapi.generated.yaml
```

If `../couchbase-cloud/openapi.yaml` doesn't exist, the droid will exit with an error message instructing you to clone the couchbase-cloud repository.

> **Note:** The OpenAPI spec is gitignored and never committed to this repository.

#### How to Use

1. **The droid will verify/create the openapi.generated.yaml symlink** (exits with error if source repo not found)

2. **Create a config file** following the example in `example-config.yaml`

3. **Call the droid** via Factory CLI:

   ```bash
   droid capella-terraform-provider-generator /path/to/config.yaml
   ```

4. **Review generated code** - The droid will show you all files it created

5. **Run validation** - The droid automatically runs:
   - `make fmt` - Format code
   - `make lint-fix` - Fix linting issues
   - `make build` - Verify compilation

#### Config File Structure

```yaml
resources:
  <resource_name>:
    create:
      path: /v4/organizations/{organizationId}/projects/{projectId}/...
      method: POST
    read:
      path: /v4/organizations/{organizationId}/projects/{projectId}/...
      method: GET
    update:
      path: /v4/organizations/{organizationId}/projects/{projectId}/...
      method: PUT
    delete:
      path: /v4/organizations/{organizationId}/projects/{projectId}/...
      method: DELETE

data_sources:
  <plural_resource_name>:
    datasource_type: single # or "list"
    read:
      path: /v4/organizations/{organizationId}/projects/{projectId}/...
      method: GET
      paginate: true # Optional: for list datasources with pagination
```

#### Key Features

- **Config-Driven**: You specify exactly what to generate
- **One Resource Per Run**: Creates one resource + one datasource per invocation
- **Parent Detection**: Extracts parent resources from path parameters automatically
- **OpenAPI Integration**: Uses `docs.GetOpenAPIDescription()` for field documentation
- **Template-Based**: References existing resources for consistent patterns
- **TDD-First**: Unit tests generated first for fast development cycle, then API resources
- **Dual Testing**: Unit tests (mocked, fast) + Acceptance tests (real API, E2E) for comprehensive validation
- **Simplicity First**: Follows "Simplicity First" and "Surgical Changes" principles
- **Validated Code**: Generated code compiles and passes linting automatically

#### What the Droid Generates

The droid generates **complete, working code** for resources and datasources:

1. **Internal schema types** - Go structs with proper field types and tfsdk tags
2. **Resource implementation** - Full CRUD methods with proper API client usage
3. **Resource schema** - Schema definition using SchemaBuilder
4. **Datasource implementation** - Read method with pagination support for list datasources
5. **Datasource schema** - Schema definition matching resource structure
6. **Provider registration** - Updates both Resources() and DataSources() methods

#### Example Workflow

```bash
# 1. Create config for a new API
cat > my-api-config.yaml << 'EOF'
resources:
  service_accounts:
    create:
      path: /v4/organizations/{organizationId}/serviceaccounts
      method: POST
    read:
      path: /v4/organizations/{organizationId}/serviceaccounts/{serviceAccountId}
      method: GET
    update:
      path: /v4/organizations/{organizationId}/serviceaccounts/{serviceAccountId}
      method: PUT
    delete:
      path: /v4/organizations/{organizationId}/serviceaccounts/{serviceAccountId}
      method: DELETE
data_sources:
  service_accounts:
    datasource_type: list
    interfaces:
      - configure
    read:
      path: /v4/organizations/{organizationId}/serviceaccounts
      method: GET
      paginate: true
EOF

# 2. Run the droid
droid capella-terraform-provider-generator my-api-config.yaml

# 3. Review the generated files
# - internal/schema/service_accounts.go
# - internal/resources/service_account.go
# - internal/resources/service_account_schema.go
# - internal/datasources/service_accounts.go
# - internal/datasources/service_accounts_schema.go
# - internal/provider/provider.go (updated)

# 4. Validation is automatic (done by droid):
# - make fmt
# - make lint-fix
# - make build
```

#### Notes

- IDs like `{alertId}` become the resource's identifier attribute
- API client: Uses v2 auto-generated client from `internal/generated/api/` (generated via oapi-codegen from `openapi.generated.yaml`)
- The droid maintains alphabetical ordering in provider registration
- Generated code follows existing codebase patterns
- After generation, verify code compiles and passes linting

#### TODO:

- **Create example configs** - Add example Terraform configurations in `examples/<resource>/`
- **Add documentation** - Update documentation templates

## Adding New Droids

To create a new custom droid for this project:

1. Create a new `.md` file in this directory
2. Use the GenerateDroid tool or manually create the droid configuration
3. Update this README with the new droid's documentation

See Factory documentation for more details on custom droids: https://docs.factory.ai/cli/configuration/custom-droids
