package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Project{}
	_ datasource.DataSourceWithConfigure = &Project{}
)

// Project is the project data source implementation.
type Project struct {
	*providerschema.Data
}

// NewProject is a helper function to simplify the provider implementation.
func NewProject() datasource.DataSource {
	return &Project{}
}

// Metadata returns the project data source type name.
func (d *Project) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the project data source.
func (d *Project) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              computedStringAttribute,
						"organization_id": computedStringAttribute,
						"name":            computedStringAttribute,
						"description":     computedStringAttribute,
						"audit":           computedAuditAttribute,
						"if_match":        computedStringAttribute,
						"etag":            computedStringAttribute,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of projects.
func (d *Project) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects", d.HostURL, organizationId),
		http.MethodGet,
		nil,
		d.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Projects",
				"Could not read projects in organization "+state.OrganizationId.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read projects in organization "+state.OrganizationId.String()+": "+err.Error(),
		)
		return
	}

	projectResp := api.GetProjectsResponse{}
	err = json.Unmarshal(response.Body, &projectResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	for _, project := range projectResp.Data {
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
func (d *Project) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
