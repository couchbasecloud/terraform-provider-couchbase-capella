package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
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

	if err := r.validateCreateProject(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	createReq := apigen.CreateProjectRequest{
		Description: func() *string {
			s := plan.Description.ValueString()
			if s == "" {
				return nil
			}
			return &s
		}(),
		Name: plan.Name.ValueString(),
	}

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", "invalid organization_id: "+err.Error())
		return
	}

	res, err := r.Apigen.PostProjectWithResponse(ctx, apigen.OrganizationId(orgUUID), createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			errorMessageWhileProjectCreation+api.ParseError(err),
		)
		return
	}
	if res.JSON201 == nil {
		resp.Diagnostics.AddError("Error creating project", "unexpected response status: "+res.Status())
		return
	}

	diags = resp.State.Set(ctx, initializeProjectWithPlanAndId(plan, res.JSON201.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.retrieveProject(ctx, organizationId, res.JSON201.Id.String())
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

	if err := r.validateProjectAttributesTrimmed(state); err != nil {
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

	putReq := apigen.PutProjectJSONRequestBody{
		Description: func() *string {
			s := state.Description.ValueString()
			if s == "" {
				return nil
			}
			return &s
		}(),
		Name: state.Name.ValueString(),
	}

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Capella Project", "invalid organization_id: "+err.Error())
		return
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Capella Project", "invalid project_id: "+err.Error())
		return
	}

	var params *apigen.PutProjectParams
	if !state.IfMatch.IsUnknown() && !state.IfMatch.IsNull() {
		ifMatch := apigen.IfMatch(state.IfMatch.ValueString())
		params = &apigen.PutProjectParams{IfMatch: &ifMatch}
	}

	_, err = r.Apigen.PutProjectWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), params, putReq)
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
			"Error Deleting Capella Project",
			"Could not delete Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.Id]
	)

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Capella Project", "invalid organization_id: "+err.Error())
		return
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Capella Project", "invalid project_id: "+err.Error())
		return
	}

	_, err = r.Apigen.DeleteProjectByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID))
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

func (r *Project) retrieveProject(ctx context.Context, organizationId, projectId string) (*providerschema.OneProject, error) {
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	res, err := r.Apigen.GetProjectByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("%s: unexpected status %s", errors.ErrExecutingRequest, res.Status())
	}

	projectResp := *res.JSON200

	etag := ""
	if res.HTTPResponse != nil {
		etag = res.HTTPResponse.Header.Get("ETag")
	}

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
		Etag: types.StringValue(etag),
	}

	return &refreshedState, nil
}

func (r *Project) validateCreateProject(plan providerschema.Project) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if !plan.IfMatch.IsNull() && !plan.IfMatch.IsUnknown() {
		return errors.ErrIfMatchCannotBeSetWhileCreate
	}

	return r.validateProjectAttributesTrimmed(plan)
}

func (r *Project) validateProjectAttributesTrimmed(plan providerschema.Project) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	if (!plan.Description.IsNull() && !plan.Description.IsUnknown()) && !providerschema.IsTrimmed(plan.Description.ValueString()) {
		return fmt.Errorf("description %s", errors.ErrNotTrimmed)
	}
	return nil
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
