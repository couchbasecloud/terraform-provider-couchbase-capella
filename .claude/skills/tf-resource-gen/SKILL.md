---
name: tf-resource-gen
description: generate terraform resources based on openapi spec.
---

# Terraform Resource Generator

## Instructions

0.  First inspect the repo for existing resource code, schema files, api structs, 
    provider registrations, and acceptance tests for the feature. For each step 1 through 
    14 if the resource already satisfies the user request skip the step 
    completely, make no edits and proceed to the next step. Do not make any minor edits or fixes 
    to existing code if the resource already satisfies the user request. 
    
    For example:
    - If api structs exist, skip creating them and use the existing structs
    - If a resource already exists, skip recreating it unless explicitly asked to edit it. 
    - If an acceptance test already exists, skip creating a duplicate and extend coverage only if asked. 

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

   Check if the create endpoint is async. In the spec look for 202 response code. If async follow steps in
   Polling Resources section.

 - Read: use get endpoint. validate IDs from state, call get endpoint, morph response to terraform state, handle ErrNotFound by removing resource from state.
 - Update: use update endpoint. if there is no update endpoint, update handler should have empty function body.
   leave a comment explaining that the resource does not support in-place updates, so an empty update handler will cause
   the framework to recreate the resource. mark all user-configurable attributes that can change with the `RequiresReplace` plan modifier so that any change forces recreation.

   If the create endpoint is async, follow the steps in Polling Resources section for the update endpoint as well.

 - Delete: use delete endpoint, handle resource-not-found gracefully (just return without error).

   If there is no delete endpoint then delete handler has empty function body.

   If the create endpoint is async, follow the steps in Polling Resources section for the delete endpoint as well.

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



## Polling Resources

After calling the create, update or delete endpoint poll the resource until it reaches a final state. The function should be
called checkFeatureStatus. For example if the feature is Cluster then the function is checkClusterStatus.
Return the refreshed terraform state.

The pattern is:

```
func (c *Cluster) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

// get plan

// call create endpoint

refreshedState, err = c.checkClusterStatus(ctx, organizationId, projectId, clusterResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			errorMessageAfterClusterCreationInitiation+api.ParseError(err),
		)
		return
	}

// set state

diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)

}
```

Do not call checkClusterStatus and getCluster/retrieveCluster back-to-back as this is redundant.
Use this pattern for Update and Delete handlers.


Here is an example on how to poll a resource. It must return a terraform object that can be set to state.
It should use a ticker to poll every 1 minute and a context timeout of 60 minutes. If the context times out return an error.

If the resource reaches a final state return the refreshed state. If there is an error calling getCluster just try again.
```
func (c *Cluster) checkClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) (*providerschema.Cluster, error) {
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cluster creation status transition timed out after initiation: %w", ctx.Err())
		case <-ticker.C:
			clusterResp, err := c.getCluster(ctx, organizationId, projectId, ClusterId)
			if err != nil {
				tflog.Info(ctx, "retrying after error polling cluster status", map[string]interface{}{"error": err.Error()})
				continue
			}

			if clusterapi.IsFinalState(clusterResp.CurrentState) {
				audit := providerschema.NewCouchbaseAuditData(clusterResp.Audit)
				auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
				if diags.HasError() {
					return nil, fmt.Errorf("%s: %s", errors.ErrUnableToConvertAuditData, diags.Errors())
				}
				refreshedState, err := providerschema.NewCluster(ctx, clusterResp, organizationId, projectId, auditObj)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
				}
				return refreshedState, nil
			}

			tflog.Info(ctx, "waiting for cluster to complete the execution")
		}
	}
}
```
