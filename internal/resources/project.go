package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

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

const errorMessageAfterProjectCreation = "Project creation is successful, but encountered an error while checking the current" +
	" state of the project. Please run `terraform plan` after 1-2 minutes to know the" +
	" current project state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileProjectCreation = "There is an error during project creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects", r.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		projectRequest,
		r.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			errorMessageWhileProjectCreation+api.ParseError(err),
		)
		return
	}

	projectResponse := api.GetProjectResponse{}
	err = json.Unmarshal(response.Body, &projectResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			errorMessageWhileProjectCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeProjectWithPlanAndId(plan, projectResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.retrieveProject(ctx, organizationId, projectResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating project",
			errorMessageAfterProjectCreation+api.ParseError(err),
		)
		return
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
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Project",
			"Could not read Capella project ID "+projectId+": "+errString,
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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		projectRequest,
		r.Token,
		headers,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Updating Capella Projects",
			"Could not update Capella project ID "+state.Id.String()+": "+errString,
		)
		return
	}

	currentState, err := r.retrieveProject(ctx, organizationId, projectId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Updating Capella Project",
			"Could not update Capella project ID "+projectId+": "+errString,
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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Capella Project",
			"Could not delete Capella project ID "+projectId+": "+errString,
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

func (r *Project) retrieveProject(_ context.Context, organizationId, projectId string) (*providerschema.OneProject, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", r.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	projectResp := api.GetProjectResponse{}
	err = json.Unmarshal(response.Body, &projectResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
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

// initializeProjectWithPlanAndId initializes an instance of providerschema.Project
// with the specified plan and ID. It marks all computed fields as null.
func initializeProjectWithPlanAndId(plan providerschema.Project, id string) providerschema.Project {
	plan.Id = types.StringValue(id)
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Etag = types.StringNull()
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}

	return plan
}
