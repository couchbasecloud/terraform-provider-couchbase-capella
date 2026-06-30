package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*DatabasePrivileges)(nil)
	_ datasource.DataSourceWithConfigure = (*DatabasePrivileges)(nil)
)

// DatabasePrivileges is the list database privileges datasource implementation.
type DatabasePrivileges struct {
	*providerschema.Data
}

func NewDatabasePrivileges() datasource.DataSource {
	return &DatabasePrivileges{}
}

func (d *DatabasePrivileges) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_privileges"
}

func (d *DatabasePrivileges) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DatabasePrivilegesSchema()
}

func (d *DatabasePrivileges) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.DatabasePrivileges
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Privileges",
			"Could not read Capella database privileges in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	privileges, err := d.listPrivileges(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Privileges",
			"Could not read database privileges in cluster "+clusterId+": "+api.ParseError(err),
		)
		return
	}

	for _, priv := range privileges {
		item := providerschema.DatabasePrivilegeItem{
			Name:  types.StringValue(priv.Name),
			Group: types.StringValue(priv.Group),
		}
		if priv.Resources != nil && priv.Resources.Buckets != nil {
			item.Resources = &providerschema.Resources{
				Buckets: providerschema.MapBucketsFromAPI(priv.Resources.Buckets),
			}
		}
		state.Data = append(state.Data, item)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *DatabasePrivileges) listPrivileges(ctx context.Context, organizationId, projectId, clusterId string) ([]api.GetDatabasePrivilegeResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privileges",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	var privileges []api.GetDatabasePrivilegeResponse
	if err := json.Unmarshal(response.Body, &privileges); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return privileges, nil
}

func (d *DatabasePrivileges) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
