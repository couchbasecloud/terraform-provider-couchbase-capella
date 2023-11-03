package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Project{}
	_ resource.ResourceWithConfigure   = &Project{}
	_ resource.ResourceWithImportState = &Project{}
)

// Project is the project resource implementation.
type Project struct {
	*providerschema.Data
}

func NewProject() resource.Resource {
	return &Project{}
}

// Metadata returns the project resource type name.
func (r *Project) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the project resource.
func (r *Project) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProjectSchema()
}

// Configure adds the provider configured client to the project resource.
func (r *Project) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new project.
func (r *Project) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	projectRequest := api.CreateProjectRequest{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects", r.HostURL, organizationId),
		http.MethodPost,
		projectRequest,
		r.Token,
		nil,
	)
	if err != nil {
		err := CheckApiError(err)
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err,
		)
	}

	projectResponse := api.GetProjectResponse{}
	err = json.Unmarshal(response.Body, &projectResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveProject(ctx, organizationId, projectResponse.Id.String())
	if err != nil {
		err := CheckApiError(err)
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err,
		)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads project information.
func (r *Project) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state providerschema.Project
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.Id]
	)

	// Get refreshed project value from Capella
	refreshedState, err := r.retrieveProject(ctx, organizationId, projectId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Projects",
				"Could not read Capella project ID "+projectId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+projectId+": "+err.Error(),
		)
		return
	}

	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		refreshedState.IfMatch = state.IfMatch
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the project.
func (r *Project) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state providerschema.Project
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Capella Project",
			"Could not update Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.Id]
	)

	projectRequest := api.PutProjectRequest{
		Description: state.Description.ValueString(),
		Name:        state.Name.ValueString(),
	}

	var headers = make(map[string]string)
	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		headers["If-Match"] = state.IfMatch.ValueString()
	}

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId),
		http.MethodPut,
		projectRequest,
		r.Token,
		headers,
	)
	if err != nil {
		err := CheckApiError(err)
		resp.Diagnostics.AddError(
			"Error Updating Capella Projects",
			"Could not update Capella project ID "+state.Id.String()+": "+err,
		)
	}

	currentState, err := r.retrieveProject(ctx, organizationId, projectId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Projects",
				"Could not read Capella project ID "+state.Id.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		currentState.IfMatch = state.IfMatch
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the project.
func (r *Project) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.Project
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Capella Project",
			"Could not update Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.Id]
	)

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId),
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
				"Error Deleting Capella Projects",
				"Could not delete Capella project ID "+projectId+": "+err.CompleteError(),
			)
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Deleting Capella Projects",
			"Could not delete Capella project ID "+projectId+": "+err.Error(),
		)
		return
	}
}

// ImportState imports a remote project that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=proj123,organization_id=org123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *Project) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *Project) retrieveProject(ctx context.Context, organizationId, projectId string) (*providerschema.OneProject, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId),
		http.MethodGet,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	projectResp := api.GetProjectResponse{}
	err = json.Unmarshal(response.Body, &projectResp)
	if err != nil {
		return nil, err
	}

	projectResp.Etag = response.Response.Header.Get("ETag")

	refreshedState := providerschema.OneProject{
		Id:             types.StringValue(projectResp.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(projectResp.Name),
		Description:    types.StringValue(projectResp.Description),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(projectResp.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(projectResp.Audit.CreatedBy),
			ModifiedAt: types.StringValue(projectResp.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(projectResp.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(projectResp.Audit.Version)),
		},
		Etag: types.StringValue(projectResp.Etag),
	}

	return &refreshedState, nil
}
