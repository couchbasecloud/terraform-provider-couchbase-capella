package resources

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*ReplicationActivationStatus)(nil)
	_ resource.ResourceWithConfigure   = (*ReplicationActivationStatus)(nil)
	_ resource.ResourceWithImportState = (*ReplicationActivationStatus)(nil)
)

// replicationActivationStatusPollInterval and replicationActivationStatusTimeout are overridden in tests.
var (
	replicationActivationStatusPollInterval = 10 * time.Second
	replicationActivationStatusTimeout      = 15 * time.Minute
)

// ReplicationActivationStatus is the resource that pauses or resumes a replication.
type ReplicationActivationStatus struct {
	*providerschema.Data
}

// NewReplicationActivationStatus is a helper function to simplify the provider implementation.
func NewReplicationActivationStatus() resource.Resource {
	return &ReplicationActivationStatus{}
}

// Metadata returns the replication activation status resource type name.
func (r *ReplicationActivationStatus) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_activation_status"
}

// Schema defines the schema for the replication activation status resource.
func (r *ReplicationActivationStatus) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ReplicationActivationStatusSchema()
}

// Configure adds the provider configured client to the replication activation status resource.
func (r *ReplicationActivationStatus) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Data = data
}

// Create pauses or resumes the replication, unless it already has the requested status.
func (r *ReplicationActivationStatus) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ReplicationActivationStatus
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.setActivationStatus(ctx, plan); err != nil {
		resp.Diagnostics.AddError("Error setting replication activation status", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read writes the current status. A replication that is starting is shown as running, and one that is
// pausing as paused. A failed replication is shown as failed, so the next plan shows a change.
func (r *ReplicationActivationStatus) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ReplicationActivationStatus
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication activation status",
			"Could not read the activation status of replication "+state.ReplicationId.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	replicationId := IDs[providerschema.ReplicationId]

	replicationResp, err := getReplication(ctx, r.Data, organizationId, projectId, clusterId, replicationId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading replication activation status",
			"Could not read the activation status of replication "+replicationId+": "+errString,
		)
		return
	}

	refreshedState := providerschema.NewReplicationActivationStatus(
		organizationId, projectId, clusterId, replicationId, activationStatus(replicationResp.Status),
	)
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update pauses or resumes the replication, unless it already has the requested status.
func (r *ReplicationActivationStatus) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ReplicationActivationStatus
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.setActivationStatus(ctx, plan); err != nil {
		resp.Diagnostics.AddError("Error updating replication activation status", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete does not call the API. The v4 API has no delete for the activation status, only pause
// (DELETE) and resume (POST), so destroy only removes the resource from the state.
func (r *ReplicationActivationStatus) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Info(ctx, "removing the replication activation status from the state, the replication does not change")
}

// ImportState imports the activation status of a replication.
// example: replication_id=<replication_id>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>
func (r *ReplicationActivationStatus) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("replication_id"), req, resp)
}

// setActivationStatus pauses or resumes the replication, unless it already has the requested status,
// and waits until it has that status.
func (r *ReplicationActivationStatus) setActivationStatus(ctx context.Context, plan providerschema.ReplicationActivationStatus) error {
	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		replicationId  = plan.ReplicationId.ValueString()
		desired        = plan.State.ValueString()
	)

	ctx, cancel := context.WithTimeout(ctx, replicationActivationStatusTimeout)
	defer cancel()

	status, err := r.getStatus(ctx, organizationId, projectId, clusterId, replicationId)
	if err != nil {
		return replicationActivationStatusError(replicationId, err)
	}
	if activationStatus(status) == desired {
		return nil
	}

	err = retryReplicationRequest(ctx, replicationId, isClusterNotHealthyError, func() error {
		return r.sendActivationStatus(ctx, organizationId, projectId, clusterId, replicationId, desired)
	})
	// The API returns 409 when the replication is already paused or running. Every replication error
	// has the same code, so a GET confirms the status.
	if isConflictError(err) && r.hasActivationStatus(ctx, organizationId, projectId, clusterId, replicationId, desired) {
		return nil
	}
	if err != nil {
		return replicationActivationStatusError(replicationId, err)
	}

	return r.waitForActivationStatus(ctx, organizationId, projectId, clusterId, replicationId, desired)
}

// getStatus gets the replication status. It retries server errors, network errors and the 422 that
// the API returns while a cluster is not healthy.
func (r *ReplicationActivationStatus) getStatus(ctx context.Context, organizationId, projectId, clusterId, replicationId string) (replicationapi.Status, error) {
	var status replicationapi.Status
	err := retryReplicationRequest(ctx, replicationId, isRetryableGetError, func() error {
		replicationResp, err := getReplication(ctx, r.Data, organizationId, projectId, clusterId, replicationId)
		if err != nil {
			return err
		}
		status = replicationResp.Status
		return nil
	})
	return status, err
}

// waitForActivationStatus polls the replication until it has the requested status. It stops at once
// when the replication is failed.
func (r *ReplicationActivationStatus) waitForActivationStatus(ctx context.Context, organizationId, projectId, clusterId, replicationId, desired string) error {
	for {
		status, err := r.getStatus(ctx, organizationId, projectId, clusterId, replicationId)
		if err != nil {
			return replicationActivationStatusError(replicationId, err)
		}
		if activationStatus(status) == desired {
			return nil
		}
		if status == replicationapi.StatusFailed {
			return fmt.Errorf("%w: replication %s", errors.ErrReplicationActivationStatusFailed, replicationId)
		}
		if err := waitForNextPoll(ctx, replicationId); err != nil {
			return err
		}
	}
}

// retryReplicationRequest calls request until it succeeds or fails with an error that retryable rejects.
// A new replication can leave both clusters in the Peering state for some minutes.
func retryReplicationRequest(ctx context.Context, replicationId string, retryable func(error) bool, request func() error) error {
	for {
		err := request()
		if err == nil || !retryable(err) {
			return err
		}
		tflog.Warn(ctx, "replication request failed, retrying", map[string]interface{}{
			"replication_id": replicationId,
			"error":          api.ParseError(err),
		})
		if err := waitForNextPoll(ctx, replicationId); err != nil {
			return err
		}
	}
}

// waitForNextPoll waits for the poll interval, or returns a timeout error when the context ends first.
func waitForNextPoll(ctx context.Context, replicationId string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: replication %s: %w", errors.ErrReplicationActivationStatusTimeout, replicationId, ctx.Err())
	case <-time.After(replicationActivationStatusPollInterval):
		return nil
	}
}

// isRetryableGetError reports whether a failed GET can be tried again: a server error, a network
// error, or the 422 that the API returns while a cluster is not healthy.
func isRetryableGetError(err error) bool {
	return !isClientError(err) || isClusterNotHealthyError(err)
}

// sendActivationStatus sends the pause (DELETE) or resume (POST) request.
func (r *ReplicationActivationStatus) sendActivationStatus(ctx context.Context, organizationId, projectId, clusterId, replicationId, desired string) error {
	method := http.MethodPost
	if desired == string(replicationapi.StatusPaused) {
		method = http.MethodDelete
	}
	cfg := api.EndpointCfg{
		Url:           replicationapi.ActivationStatusURL(r.HostURL, organizationId, projectId, clusterId, replicationId),
		Method:        method,
		SuccessStatus: http.StatusNoContent,
	}
	_, err := r.ClientV1.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	return err
}

// hasActivationStatus reports whether a GET shows the requested status. A failed GET reports false.
func (r *ReplicationActivationStatus) hasActivationStatus(ctx context.Context, organizationId, projectId, clusterId, replicationId, desired string) bool {
	replicationResp, err := getReplication(ctx, r.Data, organizationId, projectId, clusterId, replicationId)
	if err != nil {
		tflog.Warn(ctx, "could not read the replication to confirm its activation status", map[string]interface{}{
			"replication_id": replicationId,
			"error":          api.ParseError(err),
		})
		return false
	}
	return activationStatus(replicationResp.Status) == desired
}

// activationStatus maps a replication status to the state attribute: pending is shown as running,
// and pausing as paused. Other statuses are not changed.
func activationStatus(status replicationapi.Status) string {
	switch status {
	case replicationapi.StatusPending:
		return string(replicationapi.StatusRunning)
	case replicationapi.StatusPausing:
		return string(replicationapi.StatusPaused)
	default:
		return string(status)
	}
}

// replicationActivationStatusError adds the replication ID to an API error, and a hint for a 404.
// A timeout error already has the replication ID, so it is returned as it is.
func replicationActivationStatusError(replicationId string, err error) error {
	if stderrors.Is(err, errors.ErrReplicationActivationStatusTimeout) {
		return err
	}
	var hint string
	if isNotFoundError(err) {
		hint = "\n\nCheck that cluster_id is the source cluster of this replication. For the reverse direction " +
			"of a twoWay replication, use reverse_replication_id and the target cluster ID."
	}
	return fmt.Errorf("could not set the activation status of replication %s: %s%s", replicationId, api.ParseError(err), hint)
}
