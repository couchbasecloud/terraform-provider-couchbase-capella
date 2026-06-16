package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
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

	networkPeerRequest := network_peer_api.CreateNetworkPeeringRequest{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}

	switch {
	case plan.ProviderConfig.AWSConfig != nil:
		providerConfigJSON, err := createAWSProviderConfig(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating network peer for AWS",
				errors.ErrConvertingProviderConfig.Error(),
			)
			return
		}
		networkPeerRequest.ProviderConfig = providerConfigJSON
	case plan.ProviderConfig.GCPConfig != nil:
		providerConfigJSON, err := createGCPProviderConfig(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating network peer for GCP",
				errors.ErrConvertingProviderConfig.Error(),
			)
			return
		}
		networkPeerRequest.ProviderConfig = providerConfigJSON
	case plan.ProviderConfig.AzureConfig != nil:
		providerConfigJSON, err := createAzureProviderConfig(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating network peer for Azure",
				errors.ErrConvertingProviderConfig.Error(),
			)
			return
		}
		networkPeerRequest.ProviderConfig = providerConfigJSON
	default:
		resp.Diagnostics.AddError(
			"Provider Config cannot be empty",
			errors.ErrProviderConfigCannotBeEmpty.Error(),
		)
		return
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

	response, err := n.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		networkPeerRequest,
		n.Token,
		nil,
	)
	if err != nil {
		// The API may persist a failed peering record even when returning a non-success status.
		// Attempt to find the resource so we can save it to state for lifecycle management.
		peerID := n.findNetworkPeerByPlan(ctx, organizationId, projectId, clusterId, plan)
		if peerID == "" {
			resp.Diagnostics.AddError(
				"Error creating network peer",
				errorMessageWhileNetworkPeerCreation+api.ParseError(err),
			)
			return
		}

		tflog.Info(ctx, "network peer was persisted despite creation error, saving to state", map[string]interface{}{
			"peer_id": peerID,
		})

		resp.Diagnostics.AddWarning(
			"Network peer created with errors",
			fmt.Sprintf(
				"The API returned an error during creation (%s) but the network peer %q was persisted. "+
					"The resource has been saved to state for lifecycle management. "+
					"Review the peer status in Capella and, if necessary, run terraform destroy to clean up.",
				api.ParseError(err), peerID,
			),
		)

		diags = resp.State.Set(ctx, initializeNetworkPeerPlanId(plan, peerID))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		n.setRefreshedStateOrWarn(ctx, resp, organizationId, projectId, clusterId, peerID, plan.ProviderType.ValueString(), err)
		return
	}

	networkPeerResponse := network_peer_api.GetNetworkPeeringRecordResponse{}
	err = json.Unmarshal(response.Body, &networkPeerResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network peer",
			errorMessageWhileNetworkPeerCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeNetworkPeerPlanId(plan, networkPeerResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	n.setRefreshedStateOrWarn(ctx, resp, organizationId, projectId, clusterId, networkPeerResponse.Id.String(), plan.ProviderType.ValueString(), nil)
}

// setRefreshedStateOrWarn attempts to read the full peer state and set it on the response.
// If the read fails (e.g. peer is in failed state with incomplete providerConfig), a warning
// is emitted instead of an error so that the partial state (containing the ID) is preserved.
func (n *NetworkPeer) setRefreshedStateOrWarn(ctx context.Context, resp *resource.CreateResponse, organizationId, projectId, clusterId, peerID, providerType string, originalErr error) {
	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, peerID, providerType)
	if err != nil {
		detail := "The network peer was created but its current status could not be fully read: " + err.Error()
		if originalErr != nil {
			detail += ". Original creation error: " + api.ParseError(originalErr)
		}
		resp.Diagnostics.AddWarning(
			"Network peer created with incomplete state",
			detail,
		)
		return
	}

	diags := resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// createAWSProviderConfig is the function to handle AWS configuration.
func createAWSProviderConfig(plan providerschema.NetworkPeer) ([]byte, error) {
	awsConfigForJSON := network_peer_api.AWSConfigData{
		AccountId: plan.ProviderConfig.AWSConfig.AccountId.ValueString(),
		Cidr:      plan.ProviderConfig.AWSConfig.Cidr.ValueString(),
		Region:    plan.ProviderConfig.AWSConfig.Region.ValueString(),
		VpcId:     plan.ProviderConfig.AWSConfig.VpcId.ValueString(),
	}

	providerConfigJSON, err := json.Marshal(awsConfigForJSON)
	if err != nil {
		return nil, err
	}

	plan.ProviderConfig.GCPConfig = nil
	plan.ProviderConfig.AzureConfig = nil

	return providerConfigJSON, nil
}

// createGCPProviderConfig is the function to handle GCP configuration.
func createGCPProviderConfig(plan providerschema.NetworkPeer) ([]byte, error) {
	gcpConfigJSON := network_peer_api.GCPConfigData{
		NetworkName:    plan.ProviderConfig.GCPConfig.NetworkName.ValueString(),
		ProjectId:      plan.ProviderConfig.GCPConfig.ProjectId.ValueString(),
		Cidr:           plan.ProviderConfig.GCPConfig.Cidr.ValueString(),
		ServiceAccount: plan.ProviderConfig.GCPConfig.ServiceAccount.ValueString(),
	}

	providerConfigJSON, err := json.Marshal(gcpConfigJSON)
	if err != nil {
		return nil, err
	}

	plan.ProviderConfig.AWSConfig = nil
	plan.ProviderConfig.AzureConfig = nil

	return providerConfigJSON, nil
}

// createAzureProviderConfig is the function to handle Azure configuration.
func createAzureProviderConfig(plan providerschema.NetworkPeer) ([]byte, error) {
	azureConfigJSON := network_peer_api.AzureConfigData{
		AzureTenantId:  plan.ProviderConfig.AzureConfig.AzureTenantId.ValueString(),
		ResourceGroup:  plan.ProviderConfig.AzureConfig.ResourceGroup.ValueString(),
		SubscriptionId: plan.ProviderConfig.AzureConfig.SubscriptionId.ValueString(),
		VnetId:         plan.ProviderConfig.AzureConfig.VnetId.ValueString(),
		Cidr:           plan.ProviderConfig.AzureConfig.Cidr.ValueString(),
	}

	providerConfigJSON, err := json.Marshal(azureConfigJSON)
	if err != nil {
		return nil, err
	}

	plan.ProviderConfig.AWSConfig = nil
	plan.ProviderConfig.GCPConfig = nil

	return providerConfigJSON, nil
}

// Read reads the NetworkPeer information.
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
			"Error Reading Network peer",
			"Could not validate Network peer "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		peerId         = IDs[providerschema.Id]
	)

	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, peerId, state.ProviderType.ValueString())
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
	// Couchbase Capella's v4 does not support a PUT endpoint for Network Peers.
	// Network Peers can only be created, read and deleted.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/Network-Peers
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
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
	_, err = n.ClientV1.ExecuteWithRetry(
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

func (n *NetworkPeer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

	if plan.ProviderConfig != nil && plan.ProviderConfig.AzureConfig != nil {
		if err := plan.ProviderConfig.AzureConfig.Validate(); err != nil {
			return err
		}
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

	if plan.Commands.IsNull() || plan.Commands.IsUnknown() {
		plan.Commands = types.SetNull(types.StringType)
	}
	types.SetNull(types.SetType{})

	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.ObjectNull(providerschema.PeeringStatus{}.AttributeTypes())
	}

	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return plan
}

// retrieveNetworkPeer retrieves network peer information for a specified organization, project, cluster, and peer ID.
func (n *NetworkPeer) retrieveNetworkPeer(
	ctx context.Context, organizationId, projectId, clusterId, peerId, providerType string,
) (*providerschema.NetworkPeer, error) {
	networkPeerResp, err := n.getNetworkPeer(ctx, organizationId, projectId, clusterId, peerId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrNotFound, err)
	}

	// Use the known provider type when the response doesn't include it
	// (common for failed peers with empty providerConfig).
	if networkPeerResp.ProviderType == "" || validateProviderTypeIsSameInPlanAndState(providerType, networkPeerResp.ProviderType) {
		networkPeerResp.ProviderType = providerType
	}

	audit := providerschema.NewCouchbaseAuditData(networkPeerResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %v", errors.ErrUnableToConvertAuditData, diags.Errors())
	}

	refreshedState, err := providerschema.NewNetworkPeer(ctx, networkPeerResp, organizationId, projectId, clusterId, auditObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}

	if providerType != "" {
		refreshedState.ProviderType = types.StringValue(providerType)
	}

	return refreshedState, nil
}

// getNetworkPeer retrieves network peer information from the specified organization and project
// using the provided cluster ID and peer ID by open-api call.
func (n *NetworkPeer) getNetworkPeer(
	ctx context.Context, organizationId, projectId, clusterId, peerId string,
) (*network_peer_api.GetNetworkPeeringRecordResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers/%s", n.HostURL, organizationId, projectId, clusterId, peerId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := n.ClientV1.ExecuteWithRetry(
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

	// Best-effort provider type detection from response config fields.
	// For failed peers, providerConfig may be empty so this detection is not fatal.
	_ = defineProviderForResponse(&networkResp)

	return &networkResp, nil
}

// defineProviderForResponse sets the provider type on the response based on populated
// provider config fields. Returns an error if provider type cannot be determined (e.g.
// when a failed peer has an empty providerConfig), but callers may choose to ignore it.
func defineProviderForResponse(networkResp *network_peer_api.GetNetworkPeeringRecordResponse) error {
	if len(networkResp.ProviderConfig) == 0 || bytes.Equal(networkResp.ProviderConfig, json.RawMessage("null")) {
		return fmt.Errorf("%w: providerConfig is empty", errors.ErrReadingProviderConfig)
	}
	azure, err := networkResp.AsAZURE()
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrReadingAzureConfig, err)
	}

	gcp, err := networkResp.AsGCP()
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrReadingGCPConfig, err)
	}

	aws, err := networkResp.AsAWS()
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrReadingAWSConfig, err)
	}

	switch {
	case azure.AzureConfigData.AzureTenantId != "":
		networkResp.ProviderType = "azure"
	case gcp.GCPConfigData.ProjectId != "":
		networkResp.ProviderType = "gcp"
	case aws.AWSConfigData.VpcId != "":
		networkResp.ProviderType = "aws"
	default:
		return fmt.Errorf("%w: unable to determine provider type from config fields", errors.ErrReadingProviderConfig)
	}

	return nil
}

// findNetworkPeerByPlan lists all network peers for a cluster and returns the ID
// of the peer whose name, provider type, and provider-specific identifiers match
// the plan. Network peer names are not unique, so matching on name alone could
// return a different peer.
func (n *NetworkPeer) findNetworkPeerByPlan(ctx context.Context, organizationId, projectId, clusterId string, plan providerschema.NetworkPeer) string {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/networkPeers", n.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	peers, err := api.GetPaginated[[]network_peer_api.GetNetworkPeeringRecordResponse](ctx, n.ClientV1, n.Token, cfg, api.SortById)
	if err != nil {
		tflog.Warn(ctx, "failed to list network peers while searching for persisted record", map[string]interface{}{
			"error": err.Error(),
		})
		return ""
	}

	name := plan.Name.ValueString()
	providerType := strings.ToLower(plan.ProviderType.ValueString())

	for i := range peers {
		if peers[i].Name != name {
			continue
		}
		if !strings.EqualFold(peers[i].ProviderType, providerType) {
			continue
		}
		if !matchesProviderConfig(peers[i], plan) {
			continue
		}
		return peers[i].Id.String()
	}
	return ""
}

// matchesProviderConfig checks whether a peer response matches the provider-specific
// identifiers in the plan (VPC ID for AWS, network name for GCP, VNet ID for Azure).
func matchesProviderConfig(peer network_peer_api.GetNetworkPeeringRecordResponse, plan providerschema.NetworkPeer) bool {
	if len(peer.ProviderConfig) == 0 || bytes.Equal(peer.ProviderConfig, json.RawMessage("null")) {
		// Cannot verify config fields; accept name+providerType match only.
		return true
	}

	switch {
	case plan.ProviderConfig.AWSConfig != nil:
		aws, err := peer.AsAWS()
		if err != nil {
			return false
		}
		return aws.AWSConfigData.VpcId == plan.ProviderConfig.AWSConfig.VpcId.ValueString()
	case plan.ProviderConfig.GCPConfig != nil:
		gcp, err := peer.AsGCP()
		if err != nil {
			return false
		}
		return gcp.GCPConfigData.NetworkName == plan.ProviderConfig.GCPConfig.NetworkName.ValueString()
	case plan.ProviderConfig.AzureConfig != nil:
		azure, err := peer.AsAZURE()
		if err != nil {
			return false
		}
		return azure.AzureConfigData.VnetId == plan.ProviderConfig.AzureConfig.VnetId.ValueString()
	default:
		return false
	}
}

func validateProviderTypeIsSameInPlanAndState(planProviderType, stateProviderType string) bool {
	return strings.EqualFold(planProviderType, stateProviderType)
}
