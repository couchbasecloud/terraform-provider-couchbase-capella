package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaClient "terraform-provider-capella/client"
	"terraform-provider-capella/internal/capellaschema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

// NewProjectsDataSource is a helper function to simplify the provider implementation.
func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

// projectsDataSource is the project data source implementation.
type projectsDataSource struct {
	client *capellaClient.Client
}

// Metadata returns the project data source type name.
func (d *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the project data source.
func (d *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of projects.
func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state capellaschema.ProjectsResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	projects, err := d.client.GetProjects(ctx, state.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Capella Projects",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, project := range projects.Data {
		projectState := capellaschema.ProjectResponse{
			Id:             types.StringValue(project.Id.String()),
			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
			Name:           types.StringValue(project.Name),
			Description:    types.StringValue(project.Description),
			Audit: capellaschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(project.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(project.Audit.CreatedBy),
				ModifiedAt: types.StringValue(project.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(project.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(project.Audit.Version)),
			},
		}
		state.Data = append(state.Data, projectState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the project data source.
func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*capellaClient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *capella.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
