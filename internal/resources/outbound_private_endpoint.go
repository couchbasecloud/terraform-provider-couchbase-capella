package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*OutboundPrivateEndpoint)(nil)
	_ resource.ResourceWithConfigure   = (*OutboundPrivateEndpoint)(nil)
	_ resource.ResourceWithImportState = (*OutboundPrivateEndpoint)(nil)
)

// Variables rather than constants so tests can shorten them.
var (
	// outboundPrivateEndpointPollInterval is how often the endpoint is read while waiting on it.
	outboundPrivateEndpointPollInterval = 30 * time.Second

	// outboundPrivateEndpointTimeout bounds the wait for a create to leave pending and for a
	// delete to finish. Both normally take a few minutes.
	outboundPrivateEndpointTimeout = 60 * time.Minute

	// outboundPrivateEndpointConflictWait is the pause before resending a request the API
	// refused because the endpoint or the cluster was busy.
	outboundPrivateEndpointConflictWait = 30 * time.Second

	// outboundPrivateEndpointConflictAttempts bounds how often such a request is sent. The API
	// already waits up to 60 seconds for the cluster before it answers busy.
	outboundPrivateEndpointConflictAttempts = 5
)

// OutboundPrivateEndpoint is the resource for an endpoint in a Capella cluster's network that
// connects to a private link service the customer owns.
type OutboundPrivateEndpoint struct {
	*providerschema.Data
}

func NewOutboundPrivateEndpoint() resource.Resource {
	return &OutboundPrivateEndpoint{}
}

func (o *OutboundPrivateEndpoint) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_private_endpoint"
}

func (o *OutboundPrivateEndpoint) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OutboundPrivateEndpointSchema()
}

func (o *OutboundPrivateEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	o.Data = data
}

func (o *OutboundPrivateEndpoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create sends the create and waits for the endpoint to leave pending, which is when the connection
// request has reached the customer's service. It doesn't wait for the customer to accept: the
// endpoint reads pendingAcceptance until they do, and a refresh picks up linked.
func (o *OutboundPrivateEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.OutboundPrivateEndpoint
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	body := outboundpe.CreateOutboundPrivateEndpointRequest{
		Name:         plan.Name.ValueString(),
		EndpointType: plan.EndpointType.ValueString(),
		ServiceName:  plan.ServiceName.ValueString(),
	}
	if !plan.Hostname.IsNull() && !plan.Hostname.IsUnknown() {
		body.Hostname = plan.Hostname.ValueString()
	}

	cfg := api.EndpointCfg{
		Url:           o.endpointsURL(organizationId, projectId, clusterId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusAccepted,
	}
	response, err := o.executeRetryingConflicts(ctx, cfg, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating outbound private endpoint",
			"Could not create the outbound private endpoint: "+api.ParseError(err),
		)
		return
	}

	var created outboundpe.CreateOutboundPrivateEndpointResponse
	if err := json.Unmarshal(response.Body, &created); err != nil {
		resp.Diagnostics.AddError(
			"Error creating outbound private endpoint",
			"The create was accepted but its response could not be read, so the endpoint may exist in Capella "+
				"without being tracked here: "+err.Error(),
		)
		return
	}

	// Track the endpoint before the wait, so a failed or interrupted wait leaves it in state for
	// the next apply or destroy rather than orphaned in Capella.
	resp.Diagnostics.Append(resp.State.Set(ctx, initialOutboundPrivateEndpointState(plan, created.Id))...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := o.waitForStatus(ctx, organizationId, projectId, clusterId, created.Id,
		func(endpoint *outboundpe.GetOutboundPrivateEndpointResponse) bool {
			return endpoint.Status != outboundpe.StatusPending
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for outbound private endpoint",
			fmt.Sprintf("The outbound private endpoint %s was created but did not leave pending: %s", created.Id, api.ParseError(err)),
		)
		return
	}

	refreshedState, err := o.newState(ctx, endpoint, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound private endpoint", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)

	switch endpoint.Status {
	case outboundpe.StatusFailed, outboundpe.StatusRejected, outboundpe.StatusDeleting:
		resp.Diagnostics.AddError(
			"Outbound private endpoint did not connect",
			fmt.Sprintf("The outbound private endpoint %s is %s%s. It is kept in state: apply again to replace it, or destroy it.",
				endpoint.Id, endpoint.Status, failureReasonSuffix(endpoint.FailureReason)),
		)
	case outboundpe.StatusUnrecognized:
		resp.Diagnostics.AddWarning(
			"Outbound private endpoint in an unrecognized state",
			fmt.Sprintf("Capella has no mapping for the cloud's status of the outbound private endpoint %s%s.",
				endpoint.Id, providerStatusSuffix(endpoint.ProviderStatus)),
		)
	}
}

func (o *OutboundPrivateEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.OutboundPrivateEndpoint
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading outbound private endpoint",
			"Could not validate the outbound private endpoint "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		endpointId     = IDs[providerschema.Id]
	)

	endpoint, err := o.getOutboundPrivateEndpoint(ctx, organizationId, projectId, clusterId, endpointId)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			tflog.Info(ctx, "outbound private endpoint doesn't exist in Capella, removing it from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading outbound private endpoint",
			"Could not read the outbound private endpoint "+endpointId+": "+api.ParseError(err),
		)
		return
	}

	refreshedState, err := o.newState(ctx, endpoint, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound private endpoint", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update renames the endpoint. The name is the only attribute that updates in place: every other
// configurable attribute replaces the endpoint.
func (o *OutboundPrivateEndpoint) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.OutboundPrivateEndpoint
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating outbound private endpoint",
			"Could not validate the outbound private endpoint "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		endpointId     = IDs[providerschema.Id]
	)

	cfg := api.EndpointCfg{
		Url:           o.endpointURL(organizationId, projectId, clusterId, endpointId),
		Method:        http.MethodPut,
		SuccessStatus: http.StatusNoContent,
	}
	body := outboundpe.UpdateOutboundPrivateEndpointRequest{Name: plan.Name.ValueString()}
	if _, err := o.executeRetryingConflicts(ctx, cfg, body); err != nil {
		resp.Diagnostics.AddError(
			"Error updating outbound private endpoint",
			"Could not rename the outbound private endpoint "+endpointId+": "+api.ParseError(err),
		)
		return
	}

	endpoint, err := o.getOutboundPrivateEndpoint(ctx, organizationId, projectId, clusterId, endpointId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading outbound private endpoint",
			"The outbound private endpoint "+endpointId+" was renamed but could not be read back: "+api.ParseError(err),
		)
		return
	}

	refreshedState, err := o.newState(ctx, endpoint, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound private endpoint", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Delete queues the delete and waits until the endpoint is gone, so that a cluster destroyed in
// the same apply isn't torn down under it.
func (o *OutboundPrivateEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.OutboundPrivateEndpoint
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting outbound private endpoint",
			"Could not validate the outbound private endpoint "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		endpointId     = IDs[providerschema.Id]
	)

	cfg := api.EndpointCfg{
		Url:           o.endpointURL(organizationId, projectId, clusterId, endpointId),
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusAccepted,
	}
	if _, err := o.executeRetryingConflicts(ctx, cfg, nil); err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting outbound private endpoint",
			"Could not delete the outbound private endpoint "+endpointId+": "+api.ParseError(err),
		)
		return
	}

	if err := o.waitUntilDeleted(ctx, organizationId, projectId, clusterId, endpointId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting outbound private endpoint",
			"The delete of the outbound private endpoint "+endpointId+" was queued but did not finish: "+api.ParseError(err),
		)
	}
}

// waitForStatus reads the endpoint until done accepts it. A read that fails for any reason other
// than the endpoint being gone is retried on the next tick.
func (o *OutboundPrivateEndpoint) waitForStatus(
	ctx context.Context, organizationId, projectId, clusterId, endpointId string,
	done func(*outboundpe.GetOutboundPrivateEndpointResponse) bool,
) (*outboundpe.GetOutboundPrivateEndpointResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, outboundPrivateEndpointTimeout)
	defer cancel()

	ticker := time.NewTicker(outboundPrivateEndpointPollInterval)
	defer ticker.Stop()

	for {
		endpoint, err := o.getOutboundPrivateEndpoint(ctx, organizationId, projectId, clusterId, endpointId)
		switch {
		case err == nil && done(endpoint):
			return endpoint, nil
		case err == nil:
			tflog.Info(ctx, "waiting for the outbound private endpoint", map[string]any{"status": endpoint.Status})
		default:
			if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
				return nil, err
			}
			tflog.Info(ctx, "retrying after an error reading the outbound private endpoint", map[string]any{"error": err.Error()})
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for the outbound private endpoint: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

// waitUntilDeleted reads the endpoint until it is gone. A teardown that fails leaves the endpoint
// failed with deleteFailed, which ends the wait with an error. Any other failed reason is the
// endpoint's earlier state, read before the delete took effect.
func (o *OutboundPrivateEndpoint) waitUntilDeleted(ctx context.Context, organizationId, projectId, clusterId, endpointId string) error {
	endpoint, err := o.waitForStatus(ctx, organizationId, projectId, clusterId, endpointId,
		func(endpoint *outboundpe.GetOutboundPrivateEndpointResponse) bool {
			return endpoint.Status == outboundpe.StatusFailed &&
				endpoint.FailureReason != nil && *endpoint.FailureReason == outboundpe.FailureReasonDeleteFailed
		},
	)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return nil
		}
		return err
	}
	return fmt.Errorf("the teardown failed%s, delete it again to retry", failureReasonSuffix(endpoint.FailureReason))
}

// executeRetryingConflicts sends a request and resends it while the API answers that the endpoint
// or the cluster is busy, which happens when several endpoints on one cluster change in one apply.
func (o *OutboundPrivateEndpoint) executeRetryingConflicts(ctx context.Context, cfg api.EndpointCfg, body any) (*api.Response, error) {
	for attempt := 1; ; attempt++ {
		response, err := o.ClientV1.ExecuteWithRetry(ctx, cfg, body, o.Token, nil)
		if err == nil || !isOutboundPrivateEndpointBusy(err) || attempt == outboundPrivateEndpointConflictAttempts {
			return response, err
		}

		tflog.Info(ctx, "outbound private endpoint busy, retrying", map[string]any{"attempt": attempt, "error": err.Error()})
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %w", err, ctx.Err())
		case <-time.After(outboundPrivateEndpointConflictWait):
		}
	}
}

func isOutboundPrivateEndpointBusy(err error) bool {
	var apiErr *api.Error
	if !stderrors.As(err, &apiErr) {
		return false
	}
	return apiErr.Code == outboundpe.ErrCodeBusy || apiErr.Code == outboundpe.ErrCodeOperationInProgress
}

func (o *OutboundPrivateEndpoint) getOutboundPrivateEndpoint(
	ctx context.Context, organizationId, projectId, clusterId, endpointId string,
) (*outboundpe.GetOutboundPrivateEndpointResponse, error) {
	cfg := api.EndpointCfg{
		Url:           o.endpointURL(organizationId, projectId, clusterId, endpointId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := o.ClientV1.ExecuteWithRetry(ctx, cfg, nil, o.Token, nil)
	if err != nil {
		return nil, err
	}

	var endpoint outboundpe.GetOutboundPrivateEndpointResponse
	if err := json.Unmarshal(response.Body, &endpoint); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &endpoint, nil
}

func (o *OutboundPrivateEndpoint) newState(
	ctx context.Context, endpoint *outboundpe.GetOutboundPrivateEndpointResponse, organizationId, projectId, clusterId string,
) (*providerschema.OutboundPrivateEndpoint, error) {
	state, err := providerschema.NewOutboundPrivateEndpoint(ctx, endpoint, organizationId, projectId, clusterId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrRefreshingState, err)
	}
	return state, nil
}

func (o *OutboundPrivateEndpoint) endpointsURL(organizationId, projectId, clusterId string) string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints",
		o.HostURL, organizationId, projectId, clusterId)
}

func (o *OutboundPrivateEndpoint) endpointURL(organizationId, projectId, clusterId, endpointId string) string {
	return o.endpointsURL(organizationId, projectId, clusterId) + "/" + endpointId
}

// initialOutboundPrivateEndpointState is the plan with the new ID and every value the API hasn't
// returned yet set to null, since state can't hold unknowns.
func initialOutboundPrivateEndpointState(plan providerschema.OutboundPrivateEndpoint, id string) providerschema.OutboundPrivateEndpoint {
	plan.Id = types.StringValue(id)
	if plan.Hostname.IsUnknown() {
		plan.Hostname = types.StringNull()
	}
	plan.CloudProvider = types.ObjectNull(providerschema.OutboundPrivateEndpointCloudProvider{}.AttributeTypes())
	plan.ProviderEndpointId = types.StringNull()
	plan.Status = types.StringNull()
	plan.FailureReason = types.StringNull()
	plan.ProviderStatus = types.StringNull()
	plan.IsCustomHostname = types.BoolValue(!plan.Hostname.IsNull())
	plan.PrivateDnsEnabled = types.BoolNull()
	plan.IpAddresses = types.ListNull(types.StringType)
	plan.ProviderConfig = types.ObjectNull(providerschema.OutboundPrivateEndpointProviderConfig{}.AttributeTypes())
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return plan
}

func failureReasonSuffix(reason *string) string {
	if reason == nil || *reason == "" {
		return ""
	}
	return " (" + *reason + ")"
}

func providerStatusSuffix(status *string) string {
	if status == nil || *status == "" {
		return ""
	}
	return ": " + *status
}
