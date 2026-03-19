---
applyTo: "acceptance_tests/**"
---

# Acceptance Tests Review Instructions

Guidelines for Copilot code reviews on acceptance tests in this repository.

# Review Focus

Acceptance tests should verify the real behaviour of Terraform resources and datasources
agaisnt the Capella API. Reviews should focus on test coverage completeness and correctness
of assertions over stylistic concerns.

# What NOT to Review

- Formatting or whitespace changes (handled by gofmt/golangci-lint)
- Import ordering or comment style preferences

---

# Coverage

For resource acceptance tests, flag if any of the following lifecycle operations are missing:

- **Create with required fields only**
- **Create with optional fields**
- **Read** (implicit in all tests but computed attributes should be asserted with `TestCheckResourceAttrSet`)
- **Update**
- **Delete** (handled automatically but flag if `CheckDestroy` is needed and missing)
- **ImportState** (if the resource supports import)

---

For datasource acceptance tests, flag if any of the following operations or assertions are missing:

Read (ensure all required and computed attributes are asserted with TestCheckResourceAttr or TestCheckResourceAttrSet)
Read with optional arguments (if the datasource supports optional arguments, test with and without them)
Read with invalid arguments (if applicable, verify proper error handling)
Multiple datasources (if relevant, test multiple instances in a single config)
ImportState (if the datasource supports import, verify import functionality)

---

## Test Case Quality

- Flag tests that only assert on input values (e.g. `name`) without also checking computed/server-set attributes (`id`, `etag`, `audit` fields).
- Flag test steps that use hardcoded IDs or tokens instead of `globalOrgId`, `globalProjectId`, `globalClusterId`, etc.
- Flag config generator functions that embed secrets or credentials as string literals.

---

## Test Structure

- Resources that don't depend on cluster state should use `resource.ParallelTest` instead of `resource.Test`. Flag serial tests that have no ordering dependency.
- Each test function should use `randomStringWithPrefix` for resource names to avoid collisions in parallel runs.
- Flag test files that define their own provider factory instead of using `globalProtoV6ProviderFactory`.

---

## Assertions

- Prefer `TestCheckResourceAttrSet` for server-generated values (IDs, timestamps) rather than asserting exact values that may change between runs.
- Use `TestCheckResourceAttr` for values that are deterministic from the config (name, description, cidr).
