package datasources

import (
	"context"
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	datasource_projects "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/tf/datasource_projects"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
func (d *Projects) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_projects.ProjectsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data of projects.
func (d *Projects) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state datasource_projects.ProjectsModel
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

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Capella Projects", "invalid organization_id: "+err.Error())
		return
	}

	listResp, err := d.ClientV2.ListProjectsWithResponse(ctx, orgUUID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Projects",
			"Could not read projects in organization "+state.OrganizationId.String()+": "+api.ParseError(err),
		)
		return
	}
	if listResp.JSON200 == nil {
		resp.Diagnostics.AddError("Error Reading Capella Projects", "unexpected response status: "+listResp.Status())
		return
	}

	// Build list of project objects using generated DataValue
	dataItems := make([]attr.Value, 0, len(listResp.JSON200.Data))
	for _, project := range listResp.JSON200.Data {
		audit := datasource_projects.NewAuditValueMust(
			datasource_projects.AuditValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"created_at":  types.StringValue(project.Audit.CreatedAt.String()),
				"created_by":  types.StringValue(project.Audit.CreatedBy),
				"modified_at": types.StringValue(project.Audit.ModifiedAt.String()),
				"modified_by": types.StringValue(project.Audit.ModifiedBy),
				"version":     types.Int64Value(int64(project.Audit.Version)),
			},
		)
		// Convert AuditValue to ObjectValue
		auditObj, _ := audit.ToObjectValue(ctx)
		item := datasource_projects.DataValue{
			Audit:       auditObj,
			Description: types.StringValue(project.Description),
			Id:          types.StringValue(project.Id.String()),
			Name:        types.StringValue(project.Name),
		}
		dataItems = append(dataItems, item)
	}
	list, di := types.ListValue(datasource_projects.DataValue{}.Type(ctx), dataItems)
	resp.Diagnostics.Append(di...)
	state.Data = list

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
