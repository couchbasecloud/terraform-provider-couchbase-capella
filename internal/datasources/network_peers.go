package datasources

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	network_peer_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/network_peer"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &NetworkPeers{}
	_ datasource.DataSourceWithConfigure = &NetworkPeers{}
)

// NetworkPeers is the NetworkPeers data source implementation.
type NetworkPeers struct {
	*providerschema.Data
}

// NewNetworkPeers is a helper function to simplify the provider implementation.
func NewNetworkPeers() datasource.DataSource {
	return &NetworkPeers{}
}

// Metadata returns the network peers data source type name.
func (n *NetworkPeers) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_peers"
}

// Schema defines the schema for the NetworkPeers data source.
func (n *NetworkPeers) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NetworkPeerSchema()
}

// Read refreshes the Terraform state with the latest data of network peers.
func (n *NetworkPeers) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.NetworkPeers
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading network peer",
			"Could not read network peer, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading network peer",
			"Could not read network peer, unexpected error: project ID cannot be empty.",
		)
		return
	}

	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading network peer",
			"Could not read network peer, unexpected error: cluster ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers", n.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]network_peer_api.GetNetworkPeeringRecordResponse](ctx, n.Client, n.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Network Peering",
			fmt.Sprintf(
				"Could not read network peer in organization %s, project %s, and cluster %s, unexpected error: %s",
				organizationId, projectId, clusterId, api.ParseError(err),
			),
		)
		return
	}

	log.Print("***********************RESPONSE******************", response)
	for i := range response {
		networkPeer := response[i]
		audit := providerschema.NewCouchbaseAuditData(networkPeer.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Network Peer",
				fmt.Sprintf("Could not read network peer list in organization %s, project %s, and cluster %s, unexpected error: %s", organizationId, projectId, clusterId, errors.ErrUnableToConvertAuditData),
			)
		}

		newNetworkPeerData, err := providerschema.NewNetworkPeerData(&networkPeer, organizationId, projectId, clusterId, auditObj)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Network peers",
				fmt.Sprintf("Could not read network peers in organization %s, project %s,and cluster %s, unexpected error: %s", organizationId, projectId, clusterId, err.Error()),
			)
		}
		state.Data = append(state.Data, *newNetworkPeerData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the cluster data source.
func (n *NetworkPeers) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	n.Data = data
}
