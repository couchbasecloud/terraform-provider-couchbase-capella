package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// NewAppEndpoint is used in (p *capellaProvider) DataSources for building the provider.
func NewAppEndpoint() datasource.DataSource {
	return &AppEndpoint{}
}

// Metadata returns the App Service CIDRs data source type name.
func (a *AppEndpoint) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

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

func (a *AppEndpoint) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppEndpoints
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster, unexpected error: project ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", a.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]clusterapi.GetClusterResponse](ctx, a.Client, a.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Clusters",
			fmt.Sprintf(
				"Could not read clusters in organization %s and project %s, unexpected error: %s",
				organizationId, projectId, api.ParseError(err),
			),
		)
		return
	}

	for i := range response {
		cluster := response[i]

		newClusterData, err := providerschema.NewClusterData(&cluster, organizationId, projectId, auditObj)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Capella Clusters",
				fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.Error()),
			)
		}
		state.Data = append(state.Data, *newClusterData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
