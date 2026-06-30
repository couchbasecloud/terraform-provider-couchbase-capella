package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*DatabaseRoles)(nil)
	_ datasource.DataSourceWithConfigure = (*DatabaseRoles)(nil)
)

// DatabaseRoles is the database roles list data source implementation.
type DatabaseRoles struct {
	*providerschema.Data
}

func NewDatabaseRoles() datasource.DataSource {
	return &DatabaseRoles{}
}

func (d *DatabaseRoles) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_roles"
}

func (d *DatabaseRoles) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DatabaseRolesSchema()
}

func (d *DatabaseRoles) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.DatabaseRoles
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Roles in Capella",
			"Could not read Capella database roles in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/roles", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]api.GetDatabaseRoleResponse](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Roles",
			"Could not read database roles in cluster "+clusterId+": "+api.ParseError(err),
		)
		return
	}

	for _, role := range response {
		item := providerschema.DatabaseRoleItem{
			Id:             types.StringValue(role.Id.String()),
			Name:           types.StringValue(role.Name),
			Description:    types.StringValue(role.Description),
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			Audit:          providerschema.NewCouchbaseAuditData(role.Audit),
			Access:         mapAccessFromAPI(role.Access),
		}
		state.Data = append(state.Data, item)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *DatabaseRoles) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
