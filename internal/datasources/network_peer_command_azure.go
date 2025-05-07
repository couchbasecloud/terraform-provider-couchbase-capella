package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	network_peer_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/network_peer"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &AzureNetworkPeerCommand{}
	_ datasource.DataSourceWithConfigure = &AzureNetworkPeerCommand{}
)

// AzureNetworkPeerCommand is the data source implementation.
type AzureNetworkPeerCommand struct {
	*providerschema.Data
}

// NewAzureNetworkPeerCommand is a helper function to simplify the provider implementation.
func NewAzureNetworkPeerCommand() datasource.DataSource {
	return &AzureNetworkPeerCommand{}
}

// Metadata returns the data source type name.
func (a *AzureNetworkPeerCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_network_peer_command"
}

// Schema defines the schema for the azure network peer command data source.
func (a *AzureNetworkPeerCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to generate Azure CLI command for setting up VNet peering between your Azure VNet and a Capella cluster.Retrieves the role assignment command or script to be executed in the Azure CLI to assign a new network contributor role. It scopes only to the specified subscription and the virtual network within that subscription." +
			"\n\n Please check the github [repository](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/network_peer_command_azure) examples for the steps and more details.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization that owns the Capella cluster.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project containing the Capella cluster.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the Capella cluster to set up VNet peering with.",
			},
			"tenant_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Azure tenant ID where your VNet resides.",
			},
			"subscription_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Subscription ID is a GUID that uniquely identifies your subscription to use Azure services.",
			},
			"resource_group": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The resource group name holding the resource you’re connecting with Capella.",
			},
			"vnet_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The VNet ID is the name of the virtual network in Azure.",
			},
			"vnet_peering_service_principal": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The enterprise application object ID for the Capella service principal. You can find the enterprise application object ID in Azure by selecting Azure Active Directory -> Enterprise applications. Next, select the application name, the object ID is in the Object ID box.",
			},
			"command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The command to be run by the customer in is their external azure account in order to grant the service principal a network contributor role that is required for VNET peering.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of azure network peer command .
func (a *AzureNetworkPeerCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AzureVNetPeeringCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateAzurePeeringCommand(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Azure network peer command request",
			"Could not validate Azure network peer command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	AzurePeeringCommandRequest := network_peer_api.GetAzureVNetPeeringCommandRequest{
		ResourceGroup:               state.ResourceGroup.ValueString(),
		SubscriptionId:              state.SubscriptionId.ValueString(),
		TenantId:                    state.TenantId.ValueString(),
		VnetId:                      state.VnetId.ValueString(),
		VnetPeeringServicePrincipal: state.VnetPeeringServicePrincipal.ValueString(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers/networkPeerCommand", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		AzurePeeringCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Azure network peer command",
			"Could not read Azure network peer command: "+api.ParseError(err),
		)
		return
	}

	var AzurePeeringCommandResponse network_peer_api.GetAzureVNetPeeringCommandResponse
	err = json.Unmarshal(response.Body, &AzurePeeringCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling Azure network peering command response",
			"Could not unmarshall Azure network peering command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(AzurePeeringCommandResponse.Command)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the network peer command data source.
func (a *AzureNetworkPeerCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}

// validateAzurePeeringCommand ensures organization id, project id, cluster id, and all other required fields are valued.
func validateAzurePeeringCommand(config providerschema.AzureVNetPeeringCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.TenantId.IsNull() {
		return errors.ErrAzureTenantIdMissing
	}
	if config.SubscriptionId.IsNull() {
		return errors.ErrSubscriptionIdMissing
	}
	if config.VnetId.IsNull() {
		return errors.ErrVNetIdMissing
	}
	if config.ResourceGroup.IsNull() {
		return errors.ErrResourceGroup
	}
	if config.VnetPeeringServicePrincipal.IsNull() {
		return errors.ErrVnetPeeringServicePrincipal
	}

	return nil
}
