---
name: tf-resource-gen
description: generate terraform resrouce from an openapi yaml spec.
---

# Terraform Resource Generator

## Instructions

-   before implementation check the terraform codebase if the resource for the feature already exists.
    search in internal/resources and see if there is any git worktree that has the feature already implemented.

    if the feature exists then do not generate code.  do not verify if it's a full implementation.
    if any part of the code exists, like provider registration, api, schema or resource implementation
    do not find or fix any gaps.

-   create a new git worktree as follows:

    git worktree add ~/resource_<feature> -b resource_<feature>

    then cd into the worktree directory ~/resource_<feature>/

    do not remove the worktree

-   use internal/resources/snapshot_backup.go as a reference.
    this includes api, terraform schema, and resource implementation itself.

-   Resource code should be in `internal/resources/`

-   Schema for the resource should be in its own file with format `<feature>_schema.go` in `internal/resources/`.

    Add validation for organization_id, project_id and cluster_id if present. For example with organization_id:

    capellaschema.AddAttr(attrs, "organization_id", builder, stringAttribute([]string{required, requiresReplace},validator.String(stringvalidator.LengthAtLeast(1))))

-   Create a struct with the feature name that embeds the Data struct. For example if the feature is SnapshotBackup:

    type SnapshotBackup struct {
        *providerschema.Data
    }

-  Need a New function. For example:

    func NewSnapshotBackup() resource.Resource {
     return &SnapshotBackup{}
    }

-.  Type should implement interfaces `resource.Resource`, `resource.ResourceWithConfigure`, and `resource.ResourceWithImportState`.
    Must use type conversion of nil to assert that the type implements the interfaces:

    var (
        _ resource.Resource                = (*SnapshotBackup)(nil)
        _ resource.ResourceWithConfigure   = (*SnapshotBackup)(nil)
        _ resource.ResourceWithImportState = (*SnapshotBackup)(nil)
    )

-  Need Metadata function that sets `resp.TypeName = req.ProviderTypeName + "_<resource_name>"`.

-  Need Schema function that calls the schema function from the schema file:

    func (s *SnapshotBackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
     resp.Schema = SnapshotBackupSchema()
    }

-  Need Configure function following this pattern:

     func (s *SnapshotBackup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
         if req.ProviderData == nil {
             return
         }
         data, ok := req.ProviderData.(*providerschema.Data)
         if !ok {
             resp.Diagnostics.AddError(
                 "Unexpected Resource Configure Type",
                 fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
             )
             return
         }
         s.Data = data
     }

-  Need ImportState function using `resource.ImportStatePassthroughID`.

- Implement CRUD methods (Create, Read, Update, Delete):
    - Create: use create endpoint, unmarshal response, call get to populate refreshed state, set state.
    - Read: use get endpoint.  validate IDs from state, call get endpoint, morph response to terraform state, handle ErrNotFound by removing resource from state.
    - Update: use update endpoint.  if there is no update endpoint then update handler has empty function body.
    - Delete: use delete endpoint, handle resource-not-found gracefully (just return without error).  if there is no delete endpoint then delete handler has empty function body.

-   Generate necessary API request/response structs in `internal/api/<feature>/`. Use `ClientV1` to make API calls with retry logic:

    response, err := s.ClientV1.ExecuteWithRetry(ctx, cfg, requestBody, s.Token, nil)

    Use `api.EndpointCfg` with appropriate URL, Method, and SuccessStatus.

-   Create a morph function to convert API response structs to terraform schema structs. Place terraform schema structs in `internal/schema/`.

-   Register the resource in `internal/provider/provider.go` in `func (p *capellaProvider) Resources`.

-   Create acceptance tests in `acceptance_tests/` with format `<feature>_acceptance_test.go`.
    - Tests should run in parallel using `resource.ParallelTest()`.
    - Include test steps for: Create+Read, ImportState, Update.
    - Use `globalProtoV6ProviderFactory` for provider factories.
    - Use helper functions like `randomStringWithPrefix` for resource names.
    - Verify key attributes with `resource.TestCheckResourceAttr` and `resource.TestCheckResourceAttrSet`.
    - do not run acceptance tests
