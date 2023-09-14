package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// Configure adds the provider configured client to the project resource.
func (r *Project) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
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

	projectRequest := api.CreateProjectRequest{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects", r.HostURL, plan.OrganizationId.ValueString()),
		http.MethodPost,
		projectRequest,
		r.Token,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not execute request, unexpected error: "+err.Error(),
		)
		return
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

	refreshedState, err := r.retrieveProject(ctx, plan.OrganizationId.ValueString(), projectResponse.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+projectResponse.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read Capella project ID "+projectResponse.Id.String()+": "+err.Error(),
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
	// todo
}

// Update updates the project.
func (r *Project) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo
}

// Delete deletes the project.
func (r *Project) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo
}

// ImportState imports a remote project that is not created by Terraform.
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
