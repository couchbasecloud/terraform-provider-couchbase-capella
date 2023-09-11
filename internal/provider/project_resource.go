package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-capella/client"
	"terraform-provider-capella/internal/capellaschema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// projectResource is the project resource implementation.
type projectResource struct {
	client *client.Client
}

// Metadata returns the project resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the project resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"if_match": schema.StringAttribute{
				Optional: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"audit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}

}

// Create creates a new project.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan capellaschema.ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectRequest := client.CreateProjectRequest{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	// Create new project
	response, err := r.client.PostProject(ctx, projectRequest, plan.OrganizationId.ValueString())
	switch err := err.(type) {
	case nil:
	case client.Error:
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.getProjectFromClient(ctx, plan.OrganizationId.ValueString(), response.Id.String())
	switch err := err.(type) {
	case nil:
	case client.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+response.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+response.Id.String()+": "+err.Error(),
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

// Configure It adds the provider configured client to the project resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *capella.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Read reads project information.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state capellaschema.ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var organizationId string
	if !state.OrganizationId.IsNull() {
		organizationId = state.OrganizationId.ValueString()
	}

	// Get refreshed project value from Capella
	refreshedState, err := r.getProjectFromClient(ctx, organizationId, state.Id.String())
	switch err := err.(type) {
	case nil:
	case client.Error:
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

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the project.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state capellaschema.ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectId := state.Id.ValueString()

	projectRequest := client.PutProjectRequest{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	var params client.PutProjectParams
	if !plan.IfMatch.IsUnknown() && !plan.IfMatch.IsNull() {
		ifMatch := plan.IfMatch.ValueString()
		params = client.PutProjectParams{
			IfMatch: &ifMatch,
		}
	}

	// Update existing project
	err := r.client.UpdateProject(ctx, projectRequest, plan.OrganizationId.ValueString(), projectId, &params)
	switch err := err.(type) {
	case nil:
	case client.Error:
		resp.Diagnostics.AddError(
			"Error Updating Capella Projects",
			"Could not update Capella project ID "+state.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Updating Capella Projects",
			"Could not update Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	currentState, err := r.getProjectFromClient(ctx, plan.OrganizationId.ValueString(), projectId)
	switch err := err.(type) {
	case nil:
	case client.Error:
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

	if !plan.IfMatch.IsUnknown() && !plan.IfMatch.IsNull() {
		currentState.IfMatch = plan.IfMatch
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the project.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state capellaschema.ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing project
	err := r.client.DeleteProject(ctx, state.OrganizationId.ValueString(), state.Id.ValueString())
	switch err := err.(type) {
	case nil:
	case client.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Projects",
				"Could not read Capella project ID "+state.Id.String()+": "+err.CompleteError(),
			)
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting project",
			"Could not delete project, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports a remote project that is not created by Terraform.
func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// It reads the project from Capella and converts it into a ProjectResponse object.
func (r *projectResource) getProjectFromClient(ctx context.Context, organizationId, projectId string) (*capellaschema.ProjectResponse, error) {
	projectResp, err := r.client.GetProject(ctx, organizationId, projectId)
	if err != nil {
		return nil, err
	}

	refreshedState := capellaschema.ProjectResponse{
		Id:             types.StringValue(projectResp.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(projectResp.Name),
		Description:    types.StringValue(projectResp.Description),
		Audit: capellaschema.CouchbaseAuditData{
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
