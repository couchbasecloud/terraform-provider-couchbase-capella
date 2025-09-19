package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Projects{}
	_ datasource.DataSourceWithConfigure = &Projects{}
)

// Projects is the project data source implementation.
type Projects struct {
	*providerschema.Data
}

// NewProjects is a helper function to simplify the provider implementation.
func NewProjects() datasource.DataSource {
	return &Projects{}
}

// Metadata returns the project data source type name.
func (d *Projects) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the project data source.
func (d *Projects) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to retrieve project details in an organization.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the project.",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the organization.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the project.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of a particular project.",
						},
						"audit": computedAuditAttribute,
						"if_match": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A precondition header that specifies the entity tag of a resource.",
						},
						"etag": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ETag header value returned by the server, used for optimistic concurrency control.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of projects.
func (d *Projects) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Projects
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = state.OrganizationId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects", d.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]api.GetProjectResponse](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read projects in organization "+state.OrganizationId.String()+": "+api.ParseError(err),
		)
		return
	}

	for _, project := range response {
		projectState := providerschema.OneProject{
			Id:             types.StringValue(project.Id.String()),
			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
			Name:           types.StringValue(project.Name),
			Description:    types.StringValue(project.Description),
			Audit: providerschema.CouchbaseAuditData{
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
func (d *Projects) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Data = data
}
