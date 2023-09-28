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

	var privileges []string
	for _, a := range plan.Access {
		for _, p := range a.Privileges {
			privileges = append(privileges, p.ValueString())
		}
	}
	dbCredRequest.Access = []api.Access{{Privileges: privileges}}

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

	// store the password that was either auto-generated or supplied during credential creation request.
	refreshedState.Password = types.StringValue(dbResponse.Password)

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

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the database credential.
func (r *DatabaseCredential) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo in AV-62853
}

// Delete deletes the database credential.
func (r *DatabaseCredential) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo in AV-62166
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
	for i, access := range dbResp.Access {
		refreshedState.Access[i] = providerschema.Access{}
		for _, permission := range access.Privileges {
			refreshedState.Access[i].Privileges = append(refreshedState.Access[i].Privileges, types.StringValue(permission))
		}
	}

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
