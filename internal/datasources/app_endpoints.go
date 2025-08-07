package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// NewAppEndpoint is used in (p *capellaProvider) DataSources for building the provider.
func NewAppEndpoint() datasource.DataSource {
	return &AppEndpoint{}
}

// Metadata returns the App Endpoints data source type name.
func (a *AppEndpoint) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoints"
}

// Configure defines the schema for the App Endpoints data source.
func (a *AppEndpoint) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}

// Read refreshes the Terraform state with the latest App Endpoints configs.
func (a *AppEndpoint) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppEndpoints
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading App Endpoints",
			"Could not read cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading App Endpoints",
			"Could not read cluster, unexpected error: project ID cannot be empty.",
		)
		return
	}

	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading App Endpoints",
			"Could not read cluster, unexpected error: cluster ID cannot be empty.",
		)
		return
	}

	if state.AppServiceId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading App Endpoints",
			"Could not read cluster, unexpected error: App Service ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]app_endpoints.GetAppEndpointResponse](ctx, a.Client, a.Token, cfg, api.SortByName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoints",
			fmt.Sprintf(
				"Could not read App Endpoints in organization %s and project %s, unexpected error: %s",
				organizationId, projectId, api.ParseError(err),
			),
		)
		return
	}

	for i := range response {
		newAppEndpoint, err := providerschema.NewAppEndpoint(ctx, &response[i])
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Capella Clusters",
				fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.Error()),
			)
		}
		state.Data = append(state.Data, *newAppEndpoint)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
