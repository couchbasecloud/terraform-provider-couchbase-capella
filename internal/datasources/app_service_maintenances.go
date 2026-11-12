package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = (*AppServiceMaintenances)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceMaintenances)(nil)
)

// AppServiceMaintenances is the App Service maintenances data source implementation.
type AppServiceMaintenances struct {
	*providerschema.Data
}

// NewAppServiceMaintenances is a helper function to simplify the provider implementation.
func NewAppServiceMaintenances() datasource.DataSource {
	return &AppServiceMaintenances{}
}

// Metadata returns the App Service maintenances data source type name.
func (d *AppServiceMaintenances) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_maintenances"
}

// Schema defines the schema for the App Service maintenances data source.
func (d *AppServiceMaintenances) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceMaintenancesSchema()
}

// Read refreshes the Terraform state with every maintenance record of the App Service.
func (d *AppServiceMaintenances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceMaintenances
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenances",
		d.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	// The list order is fixed by the API, which accepts no sortBy.
	response, err := api.GetPaginated[[]appservice_maintenance.GetMaintenanceResponse](ctx, d.ClientV1, d.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Maintenances",
			fmt.Sprintf("Could not read maintenances of app service %s, unexpected error: %s", appServiceId, api.ParseError(err)),
		)
		return
	}

	state.Data = make([]providerschema.AppServiceMaintenanceData, 0, len(response))
	for i := range response {
		audit := providerschema.NewCouchbaseAuditData(response[i].Audit)
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella App Service Maintenances",
				fmt.Sprintf("Could not read maintenances of app service %s, unexpected error: %s", appServiceId, errors.ErrUnableToConvertAuditData),
			)
			return
		}
		state.Data = append(state.Data, providerschema.NewAppServiceMaintenanceData(&response[i], auditObj))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the App Service maintenances data source.
func (d *AppServiceMaintenances) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
