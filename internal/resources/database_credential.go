package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DatabaseCredential{}
	_ resource.ResourceWithConfigure   = &DatabaseCredential{}
	_ resource.ResourceWithImportState = &DatabaseCredential{}
)

// DatabaseCredential is the database credential resource implementation.
type DatabaseCredential struct {
	*providerschema.Data
}

func NewDatabaseCredential() resource.Resource {
	return &DatabaseCredential{}
}

// Metadata returns the name that the database credential will follow in the terraform files.
// the name as per this function is capella_database_credential.
func (r *DatabaseCredential) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_credential"
}

// Schema defines the schema for the database credential resource.
func (r *DatabaseCredential) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseCredentialSchema()
}

// Configure adds the provider configured client to the database credential resource.
func (r *DatabaseCredential) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Data = data
}

// Create creates a new database credential. This function will validate the mandatory fields in the resource.CreateRequest
// before invoking the Capella V4 API.
func (r *DatabaseCredential) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.DatabaseCredential
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrOrganizationIdCannotBeEmpty.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	if plan.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	if plan.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = plan.ClusterId.ValueString()

	dbCredRequest := api.CreateDatabaseCredentialRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Password.IsNull() {
		dbCredRequest.Password = plan.Password.ValueString()
	}

	dbCredRequest.Access = createAccess(plan)

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users", r.HostURL, organizationId, projectId, clusterId),
		http.MethodPost,
		dbCredRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+err.Error(),
		)
		return
	}

	dbResponse := api.CreateDatabaseCredentialResponse{}
	err = json.Unmarshal(response.Body, &dbResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveDatabaseCredential(ctx, organizationId, projectId, clusterId, dbResponse.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Credentials",
			"Could not read Capella database credential with ID "+dbResponse.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Credentials",
			"Could not read Capella database credential with ID "+dbResponse.Id.String()+": "+err.Error(),
		)
		return
	}

	refreshedState.Password = types.StringValue(dbResponse.Password)
	// store the password that was either auto-generated or supplied during credential creation request.
	// todo: there is a bug in the V4 public APIs where the API returns the password in the response only if it is auto-generated.
	// This will be fixed in AV-62867.
	// For now, we are working around this issue.
	if dbResponse.Password == "" {
		// this means the customer had provided a password in the terraform file during creation, store that.
		refreshedState.Password = plan.Password
	}

	// todo: there is a bug in cp-open-api where the access field is empty in the GET API response,
	// we are going to work around this for private preview.
	// The fix will be done in SURF-7366
	// For now, we are appending same permissions that the customer passed in the terraform files and not relying on the GET API response.
	refreshedState.Access = mapAccess(plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads database credential information.
func (r *DatabaseCredential) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.DatabaseCredential
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Credentials in Capella",
			"Could not read Capella database credential with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	// Get refreshed Cluster value from Capella
	refreshedState, err := r.retrieveDatabaseCredential(ctx, organizationId, projectId, clusterId, dbId)
	resourceNotFound, err := handleDatabaseCredentialError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading database credential",
			"Could not read database credential with id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	// if the user had provided the password in the input, we store that in the terraform state file.
	refreshedState.Password = state.Password

	// todo: there is a bug in cp-open-api where the access field is empty in the GET API response,
	// we are going to work around this for private preview.
	// The fix will be done in SURF-7366
	// For now, we are appending same permissions that the customer passed in the terraform files and not relying on the GET API response.
	refreshedState.Access = mapAccess(state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the database credential.
func (r *DatabaseCredential) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state providerschema.DatabaseCredential
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Credentials in Capella",
			"Could not read Capella database credential with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	dbCredRequest := api.PutDatabaseCredentialRequest{
		// it is expected that the password in the state file will never be empty.
		Password: state.Password.ValueString(),
	}

	dbCredRequest.Access = createAccess(state)

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s", r.HostURL, organizationId, projectId, clusterId, dbId),
		http.MethodPut,
		dbCredRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error updating database credential",
			"Could not update an existing database credential, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error updating database credential",
			"Could not update database credential, unexpected error: "+err.Error(),
		)
		return
	}

	currentState, err := r.retrieveDatabaseCredential(ctx, organizationId, projectId, clusterId, dbId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Credentials",
			"Could not read Capella database credential with ID "+dbId+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Credentials",
			"Could not read Capella database credential with ID "+dbId+": "+err.Error(),
		)
		return
	}

	// this will ensure that the state file stores the new updated password, if password is not to be updated, it will retain the older one.
	currentState.Password = state.Password

	// todo: there is a bug in cp-open-api where the access field is empty in the GET API response,
	// we are going to work around this for private preview.
	// The fix will be done in SURF-7366
	// For now, we are appending same permissions that the customer passed in the terraform files and not relying on the GET API response.
	currentState.Access = mapAccess(state)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the database credential.
func (r *DatabaseCredential) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.DatabaseCredential
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Credentials in Capella",
			"Could not read Capella database credential with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s", r.HostURL, organizationId, projectId, clusterId, dbId),
		http.MethodDelete,
		nil,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Deleting the Database Credential",
				"Could not delete Database Credential associated with cluster "+clusterId+": "+err.CompleteError(),
			)
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Deleting Database Credential",
			"Could not delete Database Credential associated with cluster "+clusterId+": "+err.Error(),
		)
		return
	}
}

// ImportState imports a remote database credential that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=user123,organization_id=org123,project_id=proj123,cluster_id=cluster123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *DatabaseCredential) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// retrieveDatabaseCredential fetches the database credential by making a GET API call to the Capella V4 Public API.
// This usually helps retrieve the state of a newly created database credential that was created from Terraform.
func (r *DatabaseCredential) retrieveDatabaseCredential(ctx context.Context, organizationId, projectId, clusterId, dbId string) (*providerschema.OneDatabaseCredential, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s", r.HostURL, organizationId, projectId, clusterId, dbId),
		http.MethodGet,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	dbResp := api.GetDatabaseCredentialResponse{}
	err = json.Unmarshal(response.Body, &dbResp)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.OneDatabaseCredential{
		Id:             types.StringValue(dbResp.Id.String()),
		Name:           types.StringValue(dbResp.Name),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(dbResp.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(dbResp.Audit.CreatedBy),
			ModifiedAt: types.StringValue(dbResp.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(dbResp.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(dbResp.Audit.Version)),
		},
	}
	// todo: there is a bug in cp-open-api where the access field is empty in the GET API response,
	// we are going to work around this for private preview.
	// The fix will be done in SURF-7366
	// For now, we are appending same permissions that the customer passed in the terraform files and not relying on the GET API response.
	// the below code will be uncommented once the bug is fixed.
	/*	for i, access := range dbResp.Access {
			refreshedState.Access[i] = providerschema.Access{}
			for _, permission := range access.Privileges {
				refreshedState.Access[i].Privileges = append(refreshedState.Access[i].Privileges, types.StringValue(permission))
			}
		}
	*/
	return &refreshedState, nil
}

// this func extract error message if error is api.Error and also checks whether error is
// resource not found
func handleDatabaseCredentialError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}

// todo: add a unit test for this, tracking under: https://couchbasecloud.atlassian.net/browse/AV-63401
func createAccess(input providerschema.DatabaseCredential) []api.Access {
	var access = make([]api.Access, len(input.Access))

	for i, acc := range input.Access {
		access[i] = api.Access{Privileges: make([]string, len(acc.Privileges))}
		for j, permission := range acc.Privileges {
			access[i].Privileges[j] = permission.ValueString()
		}
		if acc.Resources != nil {
			if acc.Resources.Buckets != nil {
				access[i].Resources = &api.AccessibleResources{Buckets: make([]api.Bucket, len(acc.Resources.Buckets))}
				for k, bucket := range acc.Resources.Buckets {
					access[i].Resources.Buckets[k].Name = acc.Resources.Buckets[k].Name.ValueString()
					if bucket.Scopes != nil {
						access[i].Resources.Buckets[k].Scopes = make([]api.Scope, len(bucket.Scopes))
						for s, scope := range bucket.Scopes {
							access[i].Resources.Buckets[k].Scopes[s].Name = scope.Name.ValueString()
							if scope.Collections != nil {
								access[i].Resources.Buckets[k].Scopes[s].Collections = make([]string, len(scope.Collections))
								for c, coll := range scope.Collections {
									access[i].Resources.Buckets[k].Scopes[s].Collections[c] = coll.ValueString()
								}
							}
						}
					}
				}
			}
		} else {
			// todo: There is a bug in the PUT V4 API where we cannot pass empty buckets list as it leads to a nil pointer exception.
			// to workaround this bug, I have temporarily added a fix where we pass an empty list of buckets if the terraform input field doesn't contain any buckets.
			// fix for the V4 API bug will come as part of https://couchbasecloud.atlassian.net/browse/AV-63388

			access[i].Resources = &api.AccessibleResources{Buckets: make([]api.Bucket, 0)}
		}
	}

	return access
}

// mapAccess needs a 1:1 mapping when we store the output as the refreshed state.
// todo: add a unit test, tracking under: https://couchbasecloud.atlassian.net/browse/AV-63401
func mapAccess(plan providerschema.DatabaseCredential) []providerschema.Access {
	var access = make([]providerschema.Access, len(plan.Access))

	for i, acc := range plan.Access {
		access[i] = providerschema.Access{Privileges: make([]types.String, len(acc.Privileges))}
		for j, permission := range acc.Privileges {
			access[i].Privileges[j] = permission
		}
		if acc.Resources != nil {
			if acc.Resources.Buckets != nil {
				access[i].Resources = &providerschema.Resources{Buckets: make([]providerschema.BucketResource, len(acc.Resources.Buckets))}
				for k, bucket := range acc.Resources.Buckets {
					access[i].Resources.Buckets[k].Name = acc.Resources.Buckets[k].Name
					if bucket.Scopes != nil {
						access[i].Resources.Buckets[k].Scopes = make([]providerschema.Scope, len(bucket.Scopes))
						for s, scope := range bucket.Scopes {
							access[i].Resources.Buckets[k].Scopes[s].Name = scope.Name
							if scope.Collections != nil {
								access[i].Resources.Buckets[k].Scopes[s].Collections = make([]types.String, len(scope.Collections))
								for c, coll := range scope.Collections {
									access[i].Resources.Buckets[k].Scopes[s].Collections[c] = coll
								}
							}
						}
					}
				}
			}
		}
	}

	return access
}
