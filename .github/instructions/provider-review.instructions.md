---
applyTo: "internal/**"
---

# Provider Code Review Instructions

Guidelines for Copilot code reviews on Go provider code in this repository.

## Review Focus

When reviewing code, prefer identifying bugs and regressions over stylistic concerns.

## What NOT to Review

- Formatting or whitespace changes (handled by gofmt/golangci-lint)
- Generated files (`*.gen.go`)
- Import ordering or comment style preferences

---

## API

- Ensure changes use the V1 API client rather than the V2 client.
- Flag API URLs which do not use the /v4/ path prefix.

---

## Resources

- Must contain compile time-interface checks e.g. `var _ resource.Resource = &Type{}`
- Schema should be define in a separate `*_schema.go` file
- Blank methods included to satisfy interfaces should have a comment referencing the v4 documentation for the resource.

---

## Datasources

- Must contain compile time-interface checks e.g. `var _ datasource.DataSource = &Type{}`
- List-style datasources should use api.GetPaginated instead of
  implemeting manual pagination logic.

---

## Schema

- Schema attributes should not set MarkdownDescription manually. Instead they should use the AddAttr + SchemaBuilder pattern to pull description from the OpenAPI spec.

---

## Code Quality Issues to Flag

- Incorrect or missing error handling. Discourage use of `panic`.
- Functions over 100 lines or deeply nested (>3 levels)
- Commented-out code or unjustified `interface{}`
- Hardcoded IDs or tokens
- Manual retry logic instead of using `ExecuteWithRetry`
- Raw error strings instead of `fmt.Errorf` or wrapped errors with `%w`
