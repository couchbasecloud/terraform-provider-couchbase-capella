package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	network_peer_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/network_peer"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NetworkPeer{}
	_ resource.ResourceWithConfigure   = &NetworkPeer{}
	_ resource.ResourceWithImportState = &NetworkPeer{}
)

const errorMessageWhileNetworkPeerCreation = "There is an error during network peer creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

type NetworkPeer struct {
	*providerschema.Data
}

func NewNetworkPeer() resource.Resource {
	return &NetworkPeer{}
}

func (n *NetworkPeer) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_peer"
}

func (n *NetworkPeer) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkPeerSchema()
}

// Create a network peer for the CSP.
func (n *NetworkPeer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.NetworkPeer
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := n.validateCreateNetworkPeer(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating create network peering request",
			"Could not validate network peering create request: "+err.Error(),
		)
		return
	}

	var providerConfigAWS network_peer_api.AWSConfig
	var providerConfigGCP network_peer_api.GCPConfig

	networkPeerRequest := network_peer_api.CreateNetworkPeeringRequest{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}
	//check type conversion here
	if plan.ProviderType.ValueString() == "aws" {
		providerConfigAWS = network_peer_api.AWSConfig{
			AccountId: plan.ProviderConfig.AWSConfig.AccountId.ValueString(),
			Cidr:      plan.ProviderConfig.AWSConfig.Cidr.ValueString(),
			Region:    plan.ProviderConfig.AWSConfig.Region.ValueString(),
			VpcId:     plan.ProviderConfig.AWSConfig.VpcId.ValueString(),
		}
		//err := networkPeerRequest.FromAWS(providerConfigAWS)
		//if err != nil {
		//	fmt.Errorf("for AWS %s: %w", errors.ErrConvertingProviderConfig, err)
		//}
		networkPeerRequest.AWSConfig = providerConfigAWS
	} else if plan.ProviderType.ValueString() == "gcp" {
		providerConfigGCP = network_peer_api.GCPConfig{
			NetworkName:    plan.ProviderConfig.GCPConfig.NetworkName.ValueString(),
			ProjectId:      plan.ProviderConfig.GCPConfig.ProjectId.ValueString(),
			Cidr:           plan.ProviderConfig.GCPConfig.Cidr.ValueString(),
			ServiceAccount: plan.ProviderConfig.GCPConfig.ServiceAccount.ValueString(),
		}
		//err := networkPeerRequest.FromGCP(providerConfigGCP)
		//if err != nil {
		//	fmt.Errorf("for GCP %s: %w", errors.ErrConvertingProviderConfig, err)
		//}
		networkPeerRequest.GCPConfig = providerConfigGCP
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers",
		n.HostURL,
		organizationId,
		projectId,
		clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	response, err := n.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		n.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network peer",
			errorMessageWhileNetworkPeerCreation+api.ParseError(err),
		)
		return
	}

	networkPeerResponse := network_peer_api.GetNetworkPeeringRecordResponse{}
	err = json.Unmarshal(response.Body, &networkPeerResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network peer",
			errorMessageWhileNetworkPeerCreation+"error during unmarshalling:"+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeNetworkPeerPlanId(plan, networkPeerResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, networkPeerResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading network peering service status",
			"Error reading network peering service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *NetworkPeer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.NetworkPeer
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Private Endpoint",
			"Could not validate private endpoint "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		peerId         = IDs[providerschema.Id]
	)

	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, peerId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading private endpoint status",
			"Error reading private endpoint status, unexpected error: "+errString,
		)

		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update there is no update API so returns an error.
func (n *NetworkPeer) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"No update API for network peering",
		"There doesn't exist an update API for network peering",
	)
}

// Delete deletes a network peer of the CSP.
func (n *NetworkPeer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.NetworkPeer
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting network peer ",
			"Could not delete network peer id:"+state.Id.String()+"due to validation error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		peerId         = IDs[providerschema.Id]
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers/%s",
		n.HostURL,
		organizationId,
		projectId,
		clusterId,
		peerId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = n.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		n.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting network peer",
			"Could not delete network peer id "+state.Id.String()+": "+errString,
		)
		return
	}
}

func (n *NetworkPeer) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	n.Data = data
}

func (n *NetworkPeer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_id"), req, resp)
}

func (n *NetworkPeer) validateCreateNetworkPeer(plan providerschema.NetworkPeer) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.Id.IsNull() {
		return errors.ErrPeerIdMissing
	}

	return n.validateNetworkPeerAttributesTrimmed(plan)
}

func (n *NetworkPeer) validateNetworkPeerAttributesTrimmed(plan providerschema.NetworkPeer) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	if (!plan.ProviderType.IsNull() && !plan.ProviderType.IsUnknown()) && !providerschema.IsTrimmed(plan.ProviderType.ValueString()) {
		return fmt.Errorf("providerType %s", errors.ErrNotTrimmed)
	}
	return nil
}

// initializeNetworkPeerPlanId initializes an instance of providerschema.NetworkPeer
// with the specified plan. It marks all computed fields as null.
func initializeNetworkPeerPlanId(plan providerschema.NetworkPeer, id string) providerschema.NetworkPeer {
	plan.Id = types.StringValue(id)
	for _, command := range plan.Commands {
		if command.IsNull() || command.IsUnknown() {
			command = types.StringNull()
		}
	}
	if plan.ProviderConfig.ProviderId.IsNull() || plan.ProviderConfig.ProviderId.IsUnknown() {
		plan.ProviderConfig.ProviderId = types.StringNull()
	}
	if plan.Status.Reasoning.IsNull() || plan.Status.Reasoning.IsUnknown() {
		plan.Status.Reasoning = types.StringNull()
	}
	if plan.Status.State.IsNull() || plan.Status.State.IsUnknown() {
		plan.Status.State = types.StringNull()
	}
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return plan
}

// retrieveNetworkPeer retrieves network peer information for a specified organization, project, cluster, and peer ID.
func (n *NetworkPeer) retrieveNetworkPeer(ctx context.Context, organizationId, projectId, clusterId, peerId string) (*providerschema.NetworkPeer, error) {
	networkPeerResp, err := n.getNetworkPeer(ctx, organizationId, projectId, clusterId, peerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(networkPeerResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
	}

	//refreshedState, err := providerschema.NewNetworkPeer(networkPeerResp, organizationId, projectId, clusterId, providerschema.MorphCommands(networkPeerResp.Commands), auditObj)
	refreshedState, err := providerschema.NewNetworkPeer(networkPeerResp, organizationId, projectId, clusterId, auditObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}

// getNetworkPeer retrieves cluster information from the specified organization and project
// using the provided cluster ID by open-api call.
func (n *NetworkPeer) getNetworkPeer(ctx context.Context, organizationId, projectId, clusterId, peerId string) (*network_peer_api.GetNetworkPeeringRecordResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers/%s", n.HostURL, organizationId, projectId, clusterId, peerId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := n.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		n.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	networkResp := network_peer_api.GetNetworkPeeringRecordResponse{}
	err = json.Unmarshal(response.Body, &networkResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &networkResp, nil
}
