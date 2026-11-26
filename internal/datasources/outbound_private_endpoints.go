package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*OutboundPrivateEndpoints)(nil)
	_ datasource.DataSourceWithConfigure = (*OutboundPrivateEndpoints)(nil)
)

// OutboundPrivateEndpoints is the data source for every outbound private endpoint on a cluster.
type OutboundPrivateEndpoints struct {
	*providerschema.Data
}

// NewOutboundPrivateEndpoints is a helper function to simplify the provider implementation.
func NewOutboundPrivateEndpoints() datasource.DataSource {
	return &OutboundPrivateEndpoints{}
}

// Metadata returns the data source type name.
func (d *OutboundPrivateEndpoints) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_private_endpoints"
}

// Schema defines the schema for the data source.
func (d *OutboundPrivateEndpoints) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OutboundPrivateEndpointsSchema()
}

// Configure adds the provider configured client to the data source.
func (d *OutboundPrivateEndpoints) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read lists the cluster's outbound private endpoints. The route returns them all in one page.
func (d *OutboundPrivateEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.OutboundPrivateEndpoints
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = config.OrganizationId.ValueString()
		projectId      = config.ProjectId.ValueString()
		clusterId      = config.ClusterId.ValueString()
	)

	cfg := api.EndpointCfg{
		Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints",
			d.HostURL, organizationId, projectId, clusterId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Outbound Private Endpoints",
			fmt.Sprintf("Could not read the outbound private endpoints of cluster %s: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	var list outboundpe.ListOutboundPrivateEndpointsResponse
	if err := json.Unmarshal(response.Body, &list); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Outbound Private Endpoints",
			fmt.Sprintf("Could not unmarshal the outbound private endpoints of cluster %s: %s", clusterId, err),
		)
		return
	}

	// An empty list rather than null, so a cluster without endpoints reads as data = [].
	config.Data = make([]providerschema.OutboundPrivateEndpoint, 0, len(list.Data))
	for i := range list.Data {
		endpoint, err := providerschema.NewOutboundPrivateEndpoint(ctx, &list.Data[i], organizationId, projectId, clusterId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Outbound Private Endpoints",
				fmt.Sprintf("Could not convert outbound private endpoint %s of cluster %s: %s", list.Data[i].Id, clusterId, err),
			)
			return
		}
		config.Data = append(config.Data, *endpoint)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
