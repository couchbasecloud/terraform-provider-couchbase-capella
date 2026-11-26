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
	_ datasource.DataSource              = (*OutboundPrivateEndpointServiceVerification)(nil)
	_ datasource.DataSourceWithConfigure = (*OutboundPrivateEndpointServiceVerification)(nil)
)

// OutboundPrivateEndpointServiceVerification is the data source that checks a private link
// service before an outbound private endpoint is created to it.
type OutboundPrivateEndpointServiceVerification struct {
	*providerschema.Data
}

// NewOutboundPrivateEndpointServiceVerification is a helper function to simplify the provider implementation.
func NewOutboundPrivateEndpointServiceVerification() datasource.DataSource {
	return &OutboundPrivateEndpointServiceVerification{}
}

// Metadata returns the data source type name.
func (d *OutboundPrivateEndpointServiceVerification) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_private_endpoint_service_verification"
}

// Schema defines the schema for the data source.
func (d *OutboundPrivateEndpointServiceVerification) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OutboundPrivateEndpointServiceVerificationSchema()
}

// Configure adds the provider configured client to the data source.
func (d *OutboundPrivateEndpointServiceVerification) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read runs the service checks. The route is a POST because a GCP service name contains `/`,
// and it creates nothing, so running it on every read is safe. A service that isn't ready is a
// 200 and lands in state rather than failing the read.
func (d *OutboundPrivateEndpointServiceVerification) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.OutboundPrivateEndpointServiceVerification
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = config.OrganizationId.ValueString()
		projectId      = config.ProjectId.ValueString()
		clusterId      = config.ClusterId.ValueString()
		serviceName    = config.ServiceName.ValueString()
	)

	cfg := api.EndpointCfg{
		Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/serviceVerification",
			d.HostURL, organizationId, projectId, clusterId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusOK,
	}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, outboundpe.VerifyServiceRequest{ServiceName: serviceName}, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Verifying Outbound Private Endpoint Service",
			fmt.Sprintf("Could not verify service %s for cluster %s: %s", serviceName, clusterId, api.ParseError(err)),
		)
		return
	}

	var verification outboundpe.ServiceVerification
	if err := json.Unmarshal(response.Body, &verification); err != nil {
		resp.Diagnostics.AddError(
			"Error Verifying Outbound Private Endpoint Service",
			fmt.Sprintf("Could not unmarshal the verification of service %s for cluster %s: %s", serviceName, clusterId, err),
		)
		return
	}

	config.SetResult(&verification)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
