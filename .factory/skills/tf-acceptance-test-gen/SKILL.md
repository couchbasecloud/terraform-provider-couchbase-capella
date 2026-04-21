---
name: tf-acceptance-test-gen
description: Generate acceptance tests for a named Terraform resource or data source in the Couchbase Capella provider.
---

# Terraform Acceptance Test Generator

Acceptance tests live in `acceptance_tests/` and are named `<feature>_acceptance_test.go`.

## Test structure

Every resource test follows this pattern:

```go
func TestAcc<Feature>Resource(t *testing.T) {
    resourceName := randomStringWithPrefix("tf_acc_<feature>_")
    resourceReference := "couchbase-capella_<type_name>." + resourceName

    resource.ParallelTest(t, resource.TestCase{
        ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
        Steps: []resource.TestStep{
            // Create and Read
            {
                Config: testAcc<Feature>ResourceConfig(resourceName, <values...>),
                Check: resource.ComposeAggregateTestCheckFunc(
                    testAccExists<Feature>Resource(t, resourceReference),
                    resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
                    resource.TestCheckResourceAttr(resourceReference, "<required_field>", "<expected_value>"),
                    resource.TestCheckResourceAttrSet(resourceReference, "id"),
                    resource.TestCheckResourceAttrSet(resourceReference, "<computed_field>"),
                ),
            },
            // ImportState
            {
                ResourceName:      resourceReference,
                ImportStateIdFunc: generate<Feature>ImportIdForResource(resourceReference),
                ImportState:       true,
            },
            // Update (omit if resource has no updatable fields)
            {
                Config: testAcc<Feature>ResourceConfig(resourceName, <updated_values...>),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceReference, "<updated_field>", "<new_value>"),
                ),
            },
        },
    })
}
```

- Always use `resource.ParallelTest()`.
- Always include an ImportState step.
- Use `TestCheckResourceAttr` for fields with known values; `TestCheckResourceAttrSet` for computed fields.
- Never hardcode IDs — use `globalOrgId`, `globalProjectId`, `globalClusterId`, `globalBucketId`, `globalAppServiceId`.

## Writing meaningful assertions

The goal of each assertion is to guard a specific behaviour, not to confirm the test ran. Apply these rules:

**Assert values, not just presence.** Every field you set in the config must be asserted with `TestCheckResourceAttr` using the exact value you set. `TestCheckResourceAttrSet` only proves a field is non-empty — it does not prove the provider stored the right value. Reserve `TestCheckResourceAttrSet` for fields whose values are genuinely unknowable at test time (server-generated IDs, timestamps, status fields set by the API).

**After an update, assert the full state.** The update step must check every field — not just the one that changed. This catches regressions where updating field A silently resets field B to its default.

**Error tests must be specific.** The `ExpectError` regex must match a substring of the real error message for that specific invalid input, not a generic pattern like `"error"` or `"invalid"`. Read the resource's `Create` implementation and `internal/errors/` to find the actual message. A test that accepts any error is not guarding the right behaviour.

**The API existence check must verify values, not just existence.** The `retrieve*FromServer` function should unmarshal the response and assert that the key field values on the API match what is in Terraform state. A resource that creates successfully but stores wrong values will pass an existence-only check.

**Optional fields need their own test.** If a resource has optional fields that change API behaviour (not just metadata), write a dedicated test that sets those fields and asserts the resulting state. Don't rely on the happy-path test covering them incidentally.

**Default values must be verified.** If a field has a default, write a test that omits it and asserts the default value appears in state. This confirms the provider applies and reads back the default correctly.

## Config builder

```go
func testAcc<Feature>ResourceConfig(resourceName string, field1 <type>) string {
    return fmt.Sprintf(`
    %[1]s

    resource "couchbase-capella_<type_name>" "%[2]s" {
        organization_id = "%[3]s"
        project_id      = "%[4]s"
        cluster_id      = "%[5]s"
        <field>         = <format_verb>
    }
    `, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, field1)
}
```

Always start with `%[1]s` for `globalProviderBlock`. Number format verbs sequentially.

## Import ID function

Read the `ImportState()` function in `internal/resources/<feature>.go` to find the exact composite key format, then implement:

```go
func generate<Feature>ImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
    return func(state *terraform.State) (string, error) {
        var rawState map[string]string
        for _, m := range state.Modules {
            if len(m.Resources) > 0 {
                if v, ok := m.Resources[resourceReference]; ok {
                    rawState = v.Primary.Attributes
                }
            }
        }
        return fmt.Sprintf(
            "id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
            rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
        ), nil
    }
}
```

## API existence check

Verifies the resource exists on the API, not just in Terraform state:

```go
func testAccExists<Feature>Resource(t *testing.T, resourceReference string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        var rawState map[string]string
        for _, m := range s.Modules {
            if len(m.Resources) > 0 {
                if v, ok := m.Resources[resourceReference]; ok {
                    rawState = v.Primary.Attributes
                }
            }
        }
        data := newTestClient(t)
        return retrieve<Feature>FromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
    }
}
```

Use `data.ClientV1.ExecuteWithRetry` for the API call. See `snapshot_backup_acceptance_test.go` for a full example.

## Error-case test

```go
func TestAcc<Feature>ResourceInvalid<Field>(t *testing.T) {
    resourceName := randomStringWithPrefix("tf_acc_<feature>_")
    resource.ParallelTest(t, resource.TestCase{
        ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
        Steps: []resource.TestStep{
            {
                Config:      testAcc<Feature>ResourceConfig(resourceName, <invalid_value>),
                ExpectError: regexp.MustCompile("<substring of expected error>"),
            },
        },
    })
}
```

## Data source tests

Use a `data` block instead of `resource`, reference as `data.couchbase-capella_<type_name>.<name>`, and omit the ImportState and Update steps.

## Verifying and running the tests

### Without credentials — compile check only

When credentials are not available the tests cannot be executed. The only verification possible is that the file compiles:

```bash
go build ./acceptance_tests/
```

This catches syntax errors and missing imports. It does not confirm the tests are correct or that the resource behaves as expected. Do not use `go test -c` — it produces the same result but writes an unnecessary binary to disk.

### With credentials — running the tests

Set the three required environment variables. The Makefile validates these and exits with a clear error if any are missing:

```bash
export TF_VAR_host="https://cloudapi.cloud.couchbase.com"
export TF_VAR_auth_token="<your Capella API key>"
export TF_VAR_organization_id="<your organization ID>"
```

Run the full suite:

```bash
make testacc
```

Run only the tests for the resource you just wrote (do this during development to avoid running the full suite):

```bash
TF_ACC=1 go test -timeout=120m -v ./acceptance_tests/ -run TestAcc<Feature>
```

> ⚠️ Acceptance tests create real resources in Couchbase Capella and cost money to run.

### Skipping expensive setup

By default `TestMain` creates a project, cluster, bucket, and app service before any test runs — cluster creation alone takes ~15 minutes. Point at existing resources to skip that:

```bash
export TF_VAR_project_id="<existing project ID>"
export TF_VAR_cluster_id="<existing cluster ID>"
export TF_VAR_bucket_id="<existing bucket ID>"
export TF_VAR_app_service_id="<existing app service ID>"
```

Any variable that is unset causes that resource to be created by setup and destroyed when the suite finishes.

## Reference examples

| Scenario | File |
|---|---|
| Simple resource with update | `acceptance_tests/snapshot_backup_acceptance_test.go` |
| Resource with many optional fields | `acceptance_tests/apikey_acceptance_test.go` |
| Complex nested resource | `acceptance_tests/cluster_acceptance_test.go` |
