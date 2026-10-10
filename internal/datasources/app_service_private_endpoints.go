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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppServicePrivateEndpoints{}
	_ datasource.DataSourceWithConfigure = &AppServicePrivateEndpoints{}
)

// AppServicePrivateEndpoints is the data source implementation.
type AppServicePrivateEndpoints struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpoints is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpoints() datasource.DataSource {
	return &AppServicePrivateEndpoints{}
}

// Metadata returns the data source type name.
func (p *AppServicePrivateEndpoints) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoints"
}

// Schema defines the schema for the App Service private endpoint data source.
func (p *AppServicePrivateEndpoints) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointsSchema()
}

// Read refreshes the Terraform state with the latest data of the App Service's private endpoints.
func (p *AppServicePrivateEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpoints
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := p.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Private Endpoints",
			"Could not read App Service private endpoints on app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints", p.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Private Endpoints",
			"Could not read App Service private endpoints on app service "+state.AppServiceId.String()+": "+api.ParseError(err),
		)
		return
	}

	privateEndpointsResp := api.GetAppServicePrivateEndpointsResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointsResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling App Service private endpoints",
			"Could not read App Service private endpoints, unexpected error: "+err.Error(),
		)
		return
	}

	for _, e := range privateEndpointsResp.Endpoints {
		endpointData := providerschema.AppServicePrivateEndpointData{}
		endpointData.Id = types.StringValue(e.Id)
		endpointData.Status = types.StringValue(e.Status)
		endpointData.ServiceName = types.StringValue(e.ServiceName)
		state.Data = append(state.Data, endpointData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint data source.
func (p *AppServicePrivateEndpoints) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	p.Data = data
}

// validate ensures organization id, project id, cluster id, and app service id are valued.
func (p *AppServicePrivateEndpoints) validate(state providerschema.AppServicePrivateEndpoints) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if state.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}
	return nil
}
