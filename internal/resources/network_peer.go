package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	networkPeerRequest := network_peer_api.CreateNetworkPeeringRequest{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}

	//Check if AWSConfig or GCPConfig is not nil and marshal it
	if plan.ProviderConfig.AWSConfig != nil {
		awsConfigForJSON := network_peer_api.AWSConfigData{
			AccountId: plan.ProviderConfig.AWSConfig.AccountId.ValueString(),
			Cidr:      plan.ProviderConfig.AWSConfig.Cidr.ValueString(),
			Region:    plan.ProviderConfig.AWSConfig.Region.ValueString(),
			VpcId:     plan.ProviderConfig.AWSConfig.VpcId.ValueString(),
		}

		providerConfigJSON, err := json.Marshal(awsConfigForJSON)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating network peer for AWS",
				errors.ErrConvertingProviderConfig.Error(),
			)
			return
		}

		networkPeerRequest.ProviderConfig = providerConfigJSON
		plan.ProviderConfig.GCPConfig = nil

	} else if plan.ProviderConfig.GCPConfig != nil {
		gcpConfigJSON := network_peer_api.GCPConfigData{
			NetworkName:    plan.ProviderConfig.GCPConfig.NetworkName.ValueString(),
			ProjectId:      plan.ProviderConfig.GCPConfig.ProjectId.ValueString(),
			Cidr:           plan.ProviderConfig.GCPConfig.Cidr.ValueString(),
			ServiceAccount: plan.ProviderConfig.GCPConfig.ServiceAccount.ValueString(),
		}
		providerConfigJSON, err := json.Marshal(gcpConfigJSON)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating network peer for GCP",
				errors.ErrConvertingProviderConfig.Error(),
			)
			return
		}
		networkPeerRequest.ProviderConfig = providerConfigJSON
		plan.ProviderConfig.AWSConfig = nil
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
		networkPeerRequest,
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

	// Set the ProviderType based on the request
	if networkPeerRequest.ProviderType == "aws" {
		networkPeerResponse.ProviderType = "aws"
	} else if networkPeerRequest.ProviderType == "gcp" {
		networkPeerResponse.ProviderType = "gcp"
	}

	log.Print("***********PAULOMEE PROVIDER TYPE************", networkPeerResponse.ProviderType)

	diags = resp.State.Set(ctx, initializeNetworkPeerPlanId(plan, networkPeerResponse.Id.String(), networkPeerResponse.ProviderType))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, networkPeerResponse.Id.String(), networkPeerResponse.ProviderType)
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

	refreshedState, err := n.retrieveNetworkPeer(ctx, organizationId, projectId, clusterId, peerId, state.ProviderType.String())
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
func initializeNetworkPeerPlanId(plan providerschema.NetworkPeer, id, providerType string) providerschema.NetworkPeer {
	plan.Id = types.StringValue(id)
	log.Print("*****************  initializeNetworkPeerPlanId *************", providerType)
	plan.ProviderType = types.StringValue(providerType)

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
func (n *NetworkPeer) retrieveNetworkPeer(ctx context.Context, organizationId, projectId, clusterId, peerId, providerType string) (*providerschema.NetworkPeer, error) {
	networkPeerResp, err := n.getNetworkPeer(ctx, organizationId, projectId, clusterId, peerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(networkPeerResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
	}

	refreshedState, err := providerschema.NewNetworkPeer(ctx, networkPeerResp, organizationId, projectId, clusterId, auditObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}
	refreshedState.ProviderType = types.StringValue(providerType)
	log.Print("*****************  retrieveNetworkPeer *************", providerType)
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
