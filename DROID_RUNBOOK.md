# Droid Runbook: Generating Terraform Data Sources and Acceptance Tests

A practical guide for using AI droids (coding agents) to generate Terraform datasources, resources, and acceptance tests for the Couchbase Capella provider.

## Overview

This provider uses a **skill-based droid system** (located in `.factory/skills/`) to automate the creation of Terraform data sources, resources, and acceptance tests from the OpenAPI spec. Instead of hand-writing boilerplate, you describe what you need and the droid generates the relevant code, schema, tests, and provider registration.

Available skills:

| Skill | Purpose |
|---|---|
| `tf-datasource-gen` | Generate a Terraform data source from an OpenAPI spec endpoint |
| `tf-resource-gen` | Generate a Terraform resource from an OpenAPI spec endpoint |
| `tf-acceptance-test-gen` | Generate acceptance tests for a named resource or data source |

> **QE engineers:** If you are here to write acceptance tests, go straight to [Writing Acceptance Tests](#writing-acceptance-tests).

## Requirements

- [Git](https://git-scm.com/)
- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.2
- A coding agent / AI droid with access to the workspace
- The OpenAPI spec (`openapi.generated.yaml`) at the project root

## Quick Start

### 0. Verify the Droid is Installed

Before starting, confirm your coding agent / AI droid is available and can access the workspace:

```bash
# Check that the skill file exists
ls .factory/skills/tf-datasource-gen/SKILL.md


# Check that Go is installed
go version

# Check that Terraform is installed
terraform version

# Check that goimports is installed (needed for code review steps)
which goimports || go install golang.org/x/tools/cmd/goimports@latest
```

If any of these checks fail, install the missing tool before proceeding. If the skill file (`.factory/skills/tf-datasource-gen/SKILL.md`) is missing, the droid will not know how to generate data sources — restore it from the repository first.

> **Tip:** If you're using a specific droid platform (e.g., GitHub Copilot, Cursor, Cline), make sure the agent has workspace-level file access and can read/write files in `internal/`, `acceptance_tests/`, and `.factory/`.

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

To invoke the droid and generate code:

1. Run `droid` in your terminal.
2. Type `/settings` to configure the droid. We suggest the following settings:
   - **Spec mode model:** Claude Opus 4.6
   - **Default mode model:** Gemini 3 Flash
   - **Reasoning effort:** High
   - **Autonomy:** Medium
3. Switch to **spec mode** by pressing `Shift + Tab`.
4. Enter your prompt (see examples below).
5. Review, edit, and approve the generated plan.

---

## Run-the-Droid Examples

Copy any prompt below and paste it into your coding agent to kick off the full generation workflow. Each example targets a real Capella API feature.

<!-- Example 1 — Snapshot Backups (GET + LIST) -->
<details>
<summary>▶️ <strong>Example 1 — Snapshot Backups</strong> (GET single backup + LIST all backups)</summary>

**Droid prompt — paste this directly:**

```text
Generate Terraform data sources for Snapshot Backups using the tf-datasource-gen skill and the OpenAPI spec.

Feature  : Snapshot Backups
GET path : GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}/backups/{backupId}
LIST path: GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}/backups
Spec path: openapi.generated.yaml (see paths starting with /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}/backups)
```

</details>

<!-- Example 2 — Clusters (LIST only) -->
<details>
<summary>▶️ <strong>Example 2 — Clusters</strong> (LIST only — no single-resource GET)</summary>

**Droid prompt — paste this directly:**

```text
Generate Terraform data sources for Clusters using the tf-datasource-gen skill and the OpenAPI spec.

Feature  : Clusters
LIST path: GET /v4/organizations/{organizationId}/projects/{projectId}/clusters
Spec path: openapi.generated.yaml (see paths starting with /v4/organizations/{organizationId}/projects/{projectId}/clusters)

There is no single-resource GET endpoint — only generate the list data source.
```

</details>

<!-- Example 3 — App Endpoints (GET + LIST) -->
<details>
<summary>▶️ <strong>Example 3 — App Endpoints</strong> (GET single endpoint + LIST all endpoints)</summary>

**Droid prompt — paste this directly:**

```text
Generate Terraform data sources for App Endpoints using the tf-datasource-gen skill and the OpenAPI spec.

Feature  : App Endpoints
GET path : GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/appservices/{appServiceId}/endpoints/{endpointId}
LIST path: GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/appservices/{appServiceId}/endpoints
Spec path: openapi.generated.yaml (see paths starting with /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/appservices/{appServiceId}/endpoints)
```

</details>

<!-- Example 4 — Bucket Resource (CRUD) -->
<details>
<summary>▶️ <strong>Example 4 — Bucket Resource</strong> (full CRUD resource, not just a data source)</summary>

**Droid prompt — paste this directly:**

```text
Generate a Terraform resource for Buckets using the tf-resource-gen skill and the OpenAPI spec.

Feature    : Bucket
CREATE path: POST /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets
GET path   : GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}
UPDATE path: PUT /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}
DELETE path: DELETE /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets/{bucketId}
Spec path  : openapi.generated.yaml (see paths starting with /v4/organizations/{organizationId}/projects/{projectId}/clusters/{clusterId}/buckets)
```

</details>

<!-- Example 5 — Private Endpoint DNS - Updating An Existing Resource -->
<details>
<summary>▶️ <strong>Example 4 — Private Endpoint DNS</strong> (Single field addition to existing API)</summary>

**Droid prompt — paste this directly:**

```text
Update Terraform resource for private endpoint dns using the tf-resource-gen skill and the OpenAPI spec. This is a new field in the existing LIST private endpoints API. 

Feature  : Private Endpoints
GET path: GET /v4/organizations/{organizationId}/projects/{projectId}/clusters/privateEndpointService/endpoints
Spec path: openapi.generated.yaml (see paths starting with /v4/organizations/{organizationId}/projects/{projectId}/clusters/privateEndpointService)
```

</details>

<!-- Example 6 — Acceptance Tests -->
<details>
<summary>▶️ <strong>Example 6 — Acceptance Tests</strong> (generate acceptance tests for an existing resource or data source)</summary>

**Droid prompt — paste this directly:**

Use the tf-acceptance-test-gen skill to generate acceptance tests for the Bucket resource.

```text
Use the tf-acceptance-test-gen skill to generate acceptance tests for the Bucket resource.

Resource type name : couchbase-capella_bucket
Feature            : Bucket
Implementation file: internal/resources/bucket.go
```

Adjust `Resource type name`, `Feature`, and `Implementation file` for your target. For a data source, use the data source type name (e.g. `couchbase-capella_cloud_snapshot_backup`) and point at the file in `internal/datasources/`.

</details>

### 3. Implementation Steps (What the Droid Generates)

The droid follows the steps defined in `.factory/skills/tf-datasource-gen/SKILL.md`. For full details on what gets generated and how, read the skill file directly:

```bash
cat .factory/skills/tf-datasource-gen/SKILL.md
```

### 4. Review the Generated Code

After the droid finishes, verify the output matches the steps above and project conventions:

- [ ] **Struct** embeds `*providerschema.Data`
- [ ] **Interface assertions** exist for `datasource.DataSource` and `datasource.DataSourceWithConfigure`
- [ ] **`Metadata`** returns `req.ProviderTypeName + "_feature"`
- [ ] **`Configure`** extracts `*providerschema.Data` with proper error handling
- [ ] **Schema validators** are present for `organization_id`, `project_id`, `cluster_id`
- [ ] **API structs** in `internal/api/` match the OpenAPI spec
- [ ] **`ClientV1`** is used for all API calls (not an older client)
- [ ] **Provider registration** in `provider.go` includes both `New` functions
- [ ] **Acceptance tests** exist and use `resource.ParallelTest()`

### 5. Run Acceptance Tests

> ⚠️ **Acceptance tests create real resources and may cost money.**

Run acceptance tests only for the newly generated feature instead of the full suite:

```bash
TF_ACC=1 go test -timeout=60m -v ./acceptance_tests/ -run <regex>
```

For example, to run only snapshot backup tests:

```bash
TF_ACC=1 go test -timeout=60m -v ./acceptance_tests/ -run TestAccSnapshotBackup
```

Tests should use `resource.ParallelTest()` for parallel execution.

#### Prerequisites

Store credentials in a `.env` file at the repo root (this file is gitignored — never commit it):

```bash
# .env
export TF_VAR_host="https://cloudapi.cloud.couchbase.com"
export TF_VAR_auth_token="<your Capella API key>"
export TF_VAR_organization_id="<your organization ID>"
```

Source it before starting the droid:

```bash
source .env
```

The droid checks these variables at session start and tells you whether tests can be run.

#### Skipping expensive cluster setup

By default the suite creates a project, cluster, bucket, and app service (~15 min). Add existing resource IDs to your `.env` to skip that:

```bash
# .env (append)
export TF_VAR_project_id="<existing project ID>"
export TF_VAR_cluster_id="<existing cluster ID>"
export TF_VAR_bucket_id="<existing bucket ID>"
export TF_VAR_app_service_id="<existing app service ID>"
```

Anything unset will be created by setup and torn down when the suite finishes.

#### Writing Acceptance Tests

Use the `tf-acceptance-test-gen` skill. See [Example 6](#run-the-droid-examples) above for a ready-to-paste prompt. Minimal form:

```text
Use the tf-acceptance-test-gen skill to generate acceptance tests for the Bucket resource.

Resource type name : couchbase-capella_bucket
Feature            : Bucket
Implementation file: internal/resources/bucket.go
```
---

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
    ├── tf-datasource-gen/
    │   └── SKILL.md
    ├── tf-resource-gen/
    │   └── SKILL.md
    └── tf-acceptance-test-gen/
        └── SKILL.md
```

---

## Troubleshooting

| Problem | Solution |
|---|---|
| Skills not available | Type `/skills` in the droid. You should see 3 skills: `tf-datasource-gen`, `tf-resource-gen`, `tf-acceptance-test-gen` |
| Droid uses old API client | Tell it to use `ClientV1` explicitly |
| Missing schema validators | Add `requiredStringWithValidator()` for org/project/cluster IDs |
| Build fails after generation | Run `goimports` then `go vet`, fix, repeat up to 5 times |
| Data source not appearing in Terraform | Check it's registered in `provider.go` `DataSources()` |
| `openapi.gen.go` consuming context | Never read that file — it's excluded per `AGENTS.md` |
| Debugging the agent | Review the log file in `~/.factory/sessions/`. Look for a `.jsonl` file for the current session |

## Jira MCP Integration

To enable the agent to read Jira tickets, we recommend adding the Jira MCP server:

1. Launch Droid.
2. Type `/mcp`.
3. Select the option to manually add an MCP server.
4. Provide a name, e.g., `jira`.
5. Select option **1** for HTTP (streamable).
6. Enter the URL: `https://mcp.atlassian.com/v1/mcp`
7. Proceed with authentication.

Once configured, the droid can fetch ticket details directly from Jira when you reference them in prompts.

## Tips

- **One feature at a time.** Don't ask the droid to generate multiple unrelated data sources in one pass. Implementing multiple features with multiple agents in parallel is in progress.
- **Provide the endpoint paths.** The more specific you are about which OpenAPI endpoints to use, the better the output.
- **Check the diff.** Always review `git diff` before committing generated code.
- **Iterate.** If the first pass isn't perfect, point the droid at the specific issue and ask it to fix just that.

