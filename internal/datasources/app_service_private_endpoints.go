package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServicePrivateEndpoints)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServicePrivateEndpoints)(nil)
)

// AppServicePrivateEndpoints is the App Service private endpoints data source implementation.
type AppServicePrivateEndpoints struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpoints is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpoints() datasource.DataSource {
	return &AppServicePrivateEndpoints{}
}

// Metadata returns the data source type name.
func (d *AppServicePrivateEndpoints) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoints"
}

// Schema defines the schema for the data source.
func (d *AppServicePrivateEndpoints) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointsSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServicePrivateEndpoints) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Data = data
}

// Read refreshes the Terraform state with the private endpoints of the App Service.
func (d *AppServicePrivateEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServicePrivateEndpointsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var statuses []string
	if !config.Status.IsNull() && !config.Status.IsUnknown() {
		diags = config.Status.ElementsAs(ctx, &statuses, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoints",
			fmt.Sprintf("Could not read private endpoints for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var endpointsResp apigen.GetPrivateEndpointsResponse
	if err := json.Unmarshal(response.Body, &endpointsResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoints",
			"Could not unmarshal private endpoints response: "+err.Error(),
		)
		return
	}

	// The list endpoint has no status query parameter, so the provider filters the result.
	config.Data = make([]providerschema.PrivateEndpointData, 0, len(endpointsResp.Endpoints))
	for _, e := range endpointsResp.Endpoints {
		if len(statuses) > 0 && !slices.Contains(statuses, string(e.Status)) {
			continue
		}
		config.Data = append(config.Data, providerschema.PrivateEndpointData{
			Id:          types.StringValue(e.Id),
			ServiceName: types.StringPointerValue(e.ServiceName),
			Status:      types.StringValue(string(e.Status)),
		})
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
