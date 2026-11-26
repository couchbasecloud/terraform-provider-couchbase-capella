package datasources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*OutboundPrivateEndpointCapacity)(nil)
	_ datasource.DataSourceWithConfigure = (*OutboundPrivateEndpointCapacity)(nil)
)

// OutboundPrivateEndpointCapacity is the data source for a cluster's outbound private endpoint capacity.
type OutboundPrivateEndpointCapacity struct {
	*providerschema.Data
}

// NewOutboundPrivateEndpointCapacity is a helper function to simplify the provider implementation.
func NewOutboundPrivateEndpointCapacity() datasource.DataSource {
	return &OutboundPrivateEndpointCapacity{}
}

// Metadata returns the data source type name.
func (d *OutboundPrivateEndpointCapacity) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_private_endpoint_capacity"
}

// Schema defines the schema for the data source.
func (d *OutboundPrivateEndpointCapacity) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OutboundPrivateEndpointCapacitySchema()
}

// Configure adds the provider configured client to the data source.
func (d *OutboundPrivateEndpointCapacity) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the cluster's outbound private endpoint capacity.
func (d *OutboundPrivateEndpointCapacity) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.OutboundPrivateEndpointCapacityData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := config.OrganizationId.ValueString()
	projectId := config.ProjectId.ValueString()
	clusterId := config.ClusterId.ValueString()

	cfg := api.EndpointCfg{
		Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/capacity",
			d.HostURL, organizationId, projectId, clusterId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Outbound Private Endpoint Capacity",
			capacityReadErrorDetail(clusterId, err),
		)
		return
	}

	var capacity outboundpe.GetOutboundPrivateEndpointCapacityResponse
	if err := json.Unmarshal(response.Body, &capacity); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Outbound Private Endpoint Capacity",
			fmt.Sprintf("Could not unmarshal the outbound private endpoint capacity of cluster %s: %s", clusterId, err),
		)
		return
	}

	config.ClusterEndpointLimit = &providerschema.OutboundPrivateEndpointCapacityUsage{
		Used:  types.Int64Value(capacity.ClusterEndpointLimit.Used),
		Limit: types.Int64Value(capacity.ClusterEndpointLimit.Limit),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// capacityReadErrorDetail explains a failed capacity read. The route answers 422 while the
// cluster is turned off, which reads as a bad request unless the diagnostic says so.
func capacityReadErrorDetail(clusterId string, err error) string {
	detail := fmt.Sprintf("Could not read the outbound private endpoint capacity of cluster %s: %s", clusterId, api.ParseError(err))

	var apiErr *api.Error
	if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusUnprocessableEntity {
		detail += "\n\nThe capacity can only be read while the cluster is on. If the cluster is turned off, turn it on and read again."
	}
	return detail
}
