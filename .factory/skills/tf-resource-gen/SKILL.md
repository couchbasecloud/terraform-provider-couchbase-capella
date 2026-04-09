---
name: tf-resource-gen
description: generate terraform resources based on openapi spec.
---

# Terraform Resource Generator

## Instructions

0.  First inspect the repo for existing resource code, schema files, api structs, provider registrations, and acceptance tests for the feature.
    For each step 1 through 12 if the resource already satisfies the user request, skip the step completely, make no edits and proceed to the next step.
    Do not make any minor edits or fixes to existing code if the resource already satisfies the user request.  
    For example:
    - If api structs exist, skip creating them and use/update the existing structs
    - If a resource already exists, skip recreating it and add the missing schema, registration, or test coverage
    - If an acceptance test already exists, skip creating a duplicate and extend coverage only if needed

1.  Resource code should be in `internal/resources/`.

2.  Schema for the resource should be in its own file with format `<feature>_schema.go` in `internal/resources/`.

 Add validation for organization_id, project_id and cluster_id if present. For example with organization_id:

```
 capellaschema.AddAttr(attrs, "organization_id", builder, stringAttribute([]string{required, requiresReplace},
validator.String(stringvalidator.LengthAtLeast(1))))
```

3.  Create a struct with the feature name that embeds the Data struct. For example if the feature is SnapshotBackup:

```
 type SnapshotBackup struct {
     *providerschema.Data
 }
```

4.  Need a New function. For example:

```
 func NewSnapshotBackup() resource.Resource {
     return &SnapshotBackup{}
 }
```

5.  Type should implement interfaces `resource.Resource`, `resource.ResourceWithConfigure`, and `resource.ResourceWithImportState`.
 Must use type conversion of nil to assert that the type implements the interfaces:

```
 var (
     _ resource.Resource                = (*SnapshotBackup)(nil)
     _ resource.ResourceWithConfigure   = (*SnapshotBackup)(nil)
     _ resource.ResourceWithImportState = (*SnapshotBackup)(nil)
 )
 ```

6.  Need Metadata function that sets `resp.TypeName = req.ProviderTypeName + "_<resource_name>"`.

7.  Need Schema function that calls the schema function from the schema file:

```
 func (s *SnapshotBackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
     resp.Schema = SnapshotBackupSchema()
 }
```

8.  Need Configure function following this pattern:

```
 func (s *SnapshotBackup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
     if req.ProviderData == nil {
         return
     }
     data, ok := req.ProviderData.(*providerschema.Data)
     if !ok {
         resp.Diagnostics.AddError(
             "Unexpected Resource Configure Type",
             fmt.Sprintf("Expected *providerschema, got: %T. Please report this issue to the provider developers.", req.ProviderData),
         )
         return
     }
     s.Data = data
 }
```

9.  Need ImportState function using `resource.ImportStatePassthroughID`.

10. Implement CRUD methods (Create, Read, Update, Delete):
 - Create: use create endpoint, unmarshal response, call get to populate refreshed state, set state.
 - Read: use get endpoint.  validate IDs from state, call get endpoint, morph response to terraform state, handle ErrNotFound by removing resource from state.
 - Update: use update endpoint.  if there is no update endpoint, update handler should have empty function body.
   leave a comment explaining that the resource does not support in-place updates, so an empty update handler will cause
   the framework to recreate the resource.  mark all user-configurable attributes that can change with the `RequiresReplace` plan modifier so that any change forces recreation.
 - Delete: use delete endpoint, handle resource-not-found gracefully (just return without error).  if there is no delete endpoint then delete handler has empty function body.

11. Generate necessary API request/response structs in `internal/api/<feature>/`. Use `ClientV1` to make API calls with retry logic:

```
 response, err := s.ClientV1.ExecuteWithRetry(ctx, cfg, requestBody, s.Token, nil)
```

 Use `api.EndpointCfg` with appropriate URL, Method, and SuccessStatus.

12. Create a morph function to convert API response structs to terraform schema structs. Place terraform schema structs in `internal/schema/`.

13. Register the resource in `internal/provider/provider.go` in `func (p *capellaProvider) Resources`.

14. Create acceptance tests in `acceptance_tests/` with format `<feature>_acceptance_test.go`.
 - Tests should run in parallel using `resource.ParallelTest()`.
 - Include test steps for: Create+Read, ImportState, Update, and Delete (implicit).
 - Use `globalProtoV6ProviderFactory` for provider factories.
 - Use helper functions like `randomStringWithPrefix` for resource names.
 - Verify key attributes with `resource.TestCheckResourceAttr` and `resource.TestCheckResourceAttrSet`.
