package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*Project)(nil)
	_ datasource.DataSourceWithConfigure = (*Project)(nil)
)

type Project struct {
	*providerschema.Data
}

func NewProject() datasource.DataSource {
	return &Project{}
}

func (d *Project) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *Project) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectSchema()
}

func (d *Project) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Project
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := d.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Project",
			"Could not read project: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.Id.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", d.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Project",
			"Could not read project "+projectId+": "+api.ParseError(err),
		)
		return
	}

	var projectResponse api.GetProjectResponse
	err = json.Unmarshal(response.Body, &projectResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			"Could not read project, unexpected error: "+err.Error(),
		)
		return
	}

	audit := providerschema.NewCouchbaseAuditData(projectResponse.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error while audit conversion",
			"Could not perform audit conversion",
		)
		return
	}

	projectState := providerschema.Project{
		Id:             types.StringValue(projectResponse.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(projectResponse.Name),
		Description:    types.StringValue(projectResponse.Description),
		Etag:           types.StringValue(projectResponse.Etag),
		Audit:          auditObj,
	}
	state = projectState

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *Project) validate(state providerschema.Project) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.Id.IsNull() {
		return errors.ErrProjectIdMissing
	}
	return nil
}

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
