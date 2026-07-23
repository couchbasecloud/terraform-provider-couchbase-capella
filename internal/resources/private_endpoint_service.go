package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PrivateEndpointService{}
	_ resource.ResourceWithConfigure   = &PrivateEndpointService{}
	_ resource.ResourceWithImportState = &PrivateEndpointService{}
)

const (
	errorMessageWhileEnablingPrivateEndpointService = "There is an error while enabling private endpoint service. Please check in Capella to see if there are any hanging resources that have been created, unexpected error: "
)

// Private endpoint service lifecycle states returned by the GET status API.
// enableFailed/disableFailed are terminal (no automatic retry); enabling,
// disabling, and unknown are transient; idle means no operation has run.
const (
	statusIdle          = "idle"
	statusEnabling      = "enabling"
	statusEnabled       = "enabled"
	statusEnableFailed  = "enableFailed"
	statusDisabling     = "disabling"
	statusDisabled      = "disabled"
	statusDisableFailed = "disableFailed"
	statusUnknown       = "unknown"
)

// These are vars rather than consts so unit tests can shorten them; production
// behavior is unchanged.
var (
	// cleanupTimeout bounds how long we wait for the backend to tear down a
	// failed enable before giving up and surfacing an escalation error.
	cleanupTimeout = 15 * time.Minute

	// pollInterval is how often the service status is polled while waiting for a
	// transition.
	pollInterval = 30 * time.Second

	// statusChangeTimeout bounds how long we wait for the service to reach the
	// desired lifecycle state, end-to-end.
	statusChangeTimeout = 60 * time.Minute
)

// PrivateEndpointService is the scope resource implementation.
type PrivateEndpointService struct {
	*providerschema.Data
}

// NewPrivateEndpointService is a helper function to simplify the provider implementation.
func NewPrivateEndpointService() resource.Resource {
	return &PrivateEndpointService{}
}

// Metadata returns the private endpoint service resource type name.
func (p *PrivateEndpointService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoint_service"
}

// Schema defines the schema for a private endpoint service resource.
func (p *PrivateEndpointService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PrivateEndpointServiceSchema()
}

// Create enables private endpoint service.
func (p *PrivateEndpointService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.PrivateEndpointService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateCreateEndpointService(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating private endpoint service request",
			"Could not validate private endpoint service request, unexpected error: "+err.Error(),
		)
		return
	}
	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err = p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error enabling private endpoint service",
			errorMessageWhileEnablingPrivateEndpointService+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializePrivateEndpointServicePlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = p.waitUntilStatusChanges(ctx, true, organizationId, projectId, clusterId)
	if err != nil {
		// Terminal enableFailed: clean up the orphaned infra and remove the
		// resource so the next apply performs a clean re-create.
		if stderrors.Is(err, errors.ErrPrivateEndpointServiceEnableFailed) {
			p.handleFailedEnable(ctx, &resp.State, &resp.Diagnostics, organizationId, projectId, clusterId, err)
			return
		}
		resp.Diagnostics.AddError(
			"Error could not enable private endpoint service",
			"Error could not enable private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private endpoint service status",
			"Error reading private endpoint service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the private endpoint service status.
func (p *PrivateEndpointService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.PrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Private Endpoint Service in Capella",
			"Could not read Capella private endpoint service on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading private endpoint service status",
			"Error reading private endpoint service status, unexpected error: "+errString,
		)

		return
	}

	if refreshedState.Status.ValueString() == statusEnableFailed {
		tflog.Info(ctx, "private endpoint service is in enableFailed state; removing from state to force re-create")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update will enable/disable the private endpoint service.
func (p *PrivateEndpointService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config providerschema.PrivateEndpointService
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	status := "enabling"

	if !config.Enabled.ValueBool() {
		cfg.Method = http.MethodDelete
		status = "disabling"
	}

	_, err := p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error "+status+" private endpoint service",
			"Error "+status+" private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx,
		config.Enabled.ValueBool(),
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString())
	if err != nil {
		// An enable-flavored update that fails terminally is recovered the same
		// way as Create: clean up and remove from state for a clean retry.
		if config.Enabled.ValueBool() && stderrors.Is(err, errors.ErrPrivateEndpointServiceEnableFailed) {
			p.handleFailedEnable(ctx, &resp.State, &resp.Diagnostics,
				config.OrganizationId.ValueString(),
				config.ProjectId.ValueString(),
				config.ClusterId.ValueString(),
				err)
			return
		}
		resp.Diagnostics.AddError(
			"Error "+status+" private endpoint service",
			"Error "+status+"private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := p.getServiceState(ctx,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private endpoint service status",
			"Error reading private endpoint service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete disables private endpoint service on the cluster.
func (p *PrivateEndpointService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.PrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Private Endpoint Service in Capella",
			"Could not validate Capella private endpoint service on cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	// If private endpoint service is already disabled, just remove the resource from the state file.
	if !state.Enabled.ValueBool() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error disabling private endpoint service",
			"Could not disable private endpoint service for cluster "+clusterId+" unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx, false, organizationId, projectId, clusterId)
	if err != nil {
		// On a terminal disableFailed we fail fast and keep the resource in state
		// (Terraform retains a resource whose Delete errors). The correct retry
		// for a failed disable is another destroy, which Terraform performs
		// naturally on the next run.
		if stderrors.Is(err, errors.ErrPrivateEndpointServiceDisableFailed) {
			resp.Diagnostics.AddError(
				"Private endpoint service disable failed",
				fmt.Sprintf(
					"Disable failed for cluster %s: %s. The resource has been kept in state; "+
						"re-run terraform destroy to retry, and contact Couchbase Capella Support if it persists.",
					clusterId, err.Error(),
				),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error could not disable private endpoint service",
			"Error could not disable private endpoint service, unexpected error: "+err.Error(),
		)
	}
}

// Configure It adds the provider configured api to the private endpoint service resource.
func (p *PrivateEndpointService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	p.Data = data
}

// ImportState imports a private endpoint service status.
func (p *PrivateEndpointService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// validateCreateEndpointService ensures organization id, project id and cluster id are valued.
func validateCreateEndpointService(plan providerschema.PrivateEndpointService) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}

	return nil
}

// initializePlan initializes an instance of providerschema.PrivateEndpointService
// with the specified plan. It marks all computed fields as null.
func initializePrivateEndpointServicePlan(plan providerschema.PrivateEndpointService) providerschema.PrivateEndpointService {
	if plan.Enabled.IsNull() || plan.Enabled.IsUnknown() {
		plan.Enabled = types.BoolNull()
	}
	// status is computed; never persist an unknown value to state.
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.StringNull()
	}
	// service_name is computed; never persist an unknown value to state.
	if plan.ServiceName.IsNull() || plan.ServiceName.IsUnknown() {
		plan.ServiceName = types.StringNull()
	}
	return plan
}

// waitUntilStatusChanges waits until the service reaches the desired state on the
// cluster. When the API reports an explicit lifecycle status it keeps polling on
// transient ones (enabling/disabling/unknown) and fails fast on terminal ones
// (enableFailed/disableFailed) — but only once it has seen evidence the current
// operation is in flight, so a residual terminal state from a prior attempt is
// not mistaken for failure of the operation we just issued. When the status is
// absent — which happens on GCP, when the private endpoint status feature flag
// is disabled, or on older control planes — it falls back to the Enabled
// boolean, preserving the previous behavior. The overall timeout remains as a
// backstop.
func (p *PrivateEndpointService) waitUntilStatusChanges(ctx context.Context, finalState bool, organizationId, projectId, clusterId string) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, statusChangeTimeout)
	defer cancel()

	timer := time.NewTimer(0)
	defer timer.Stop()

	// sawInFlight flips true once we observe any status other than the
	// terminal-failure ones, which is our evidence that the backend has
	// progressed past whatever residual state was present before our POST.
	// Until then, a reported enableFailed/disableFailed could be stale and is
	// treated as transient.
	var sawInFlight bool

	// lastTerminalFailure records a terminal-failure status we observed but
	// deferred on because sawInFlight was still false (i.e. it might have been
	// residual from a prior attempt). If the backend never transitions and we
	// hit the overall timeout, this lets us surface the typed failure error so
	// the caller still routes to cleanup, instead of a generic timeout that
	// would leave orphaned infra behind and skip state removal.
	var lastTerminalFailure error

	// transientFailure gates whether a terminal-failure status observed while
	// sawInFlight is true gets returned right away or deferred once more. A
	// single enableFailed/disableFailed sighting can be a transient blip
	// (e.g. a timeout on the backend's side) that clears up on the backend's
	// own automatic retry, so we don't want to fail the resource on it
	// immediately. It starts true, so the first sighting is given the benefit
	// of the doubt and just flips it to false; only if the same terminal
	// status is still being reported on the next poll (after another
	// pollInterval backoff) do we treat it as a real, non-transient failure
	// and return the error.
	var transientFailure = true

	for {
		select {
		case <-ctx.Done():
			// If the operation appeared stuck at a terminal failure the whole
			// time (sawInFlight never flipped), trust that failure now rather
			// than returning a generic timeout — the generic timeout does not
			// route to handleFailedEnable, so cleanup/state-removal would be
			// skipped.
			if lastTerminalFailure != nil {
				return lastTerminalFailure
			}
			return errors.ErrPrivateEndpointServiceTimeout

		case <-timer.C:
			response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId)
			if err != nil {
				// If our deadline expired mid-request, surface the typed
				// timeout rather than the raw context error — the select may
				// race with ctx.Done() when both are ready.
				if ctx.Err() != nil {
					if lastTerminalFailure != nil {
						return lastTerminalFailure
					}
					return errors.ErrPrivateEndpointServiceTimeout
				}
				return err
			}

			// Status is absent on GCP, when the status feature flag is disabled,
			// or on older control planes: fall back to the Enabled boolean,
			// preserving the previous behavior.
			if response.Status == nil {
				if response.Enabled == finalState {
					return nil
				}
				timer.Reset(pollInterval)
				continue
			}

			switch *response.Status {
			case statusEnableFailed:
				if sawInFlight {
					// Already given the benefit of the doubt once; seeing the
					// same failure again means it wasn't a transient blip.
					if !transientFailure {
						return errors.ErrPrivateEndpointServiceEnableFailed
					}
					transientFailure = false
				}
				// Possibly-stale: defer, but remember it so a never-transitioning
				// enableFailed is still routed to cleanup on timeout.
				lastTerminalFailure = errors.ErrPrivateEndpointServiceEnableFailed
			case statusDisableFailed:
				if sawInFlight {
					// Already given the benefit of the doubt once; seeing the
					// same failure again means it wasn't a transient blip.
					if !transientFailure {
						return errors.ErrPrivateEndpointServiceDisableFailed
					}
					transientFailure = false
				}
				lastTerminalFailure = errors.ErrPrivateEndpointServiceDisableFailed
			case statusEnabling, statusDisabling, statusUnknown:
				sawInFlight = true
				// Reset so a future terminal-failure sighting is treated as a
				// fresh, isolated one (needs to repeat again before we trust
				// it) rather than immediately confirming a stale prior one.
				transientFailure = true
				// The backend progressed, so any earlier terminal status was
				// genuinely stale; a later stall is a real transition timeout.
				lastTerminalFailure = nil
			case statusEnabled, statusDisabled, statusIdle:
				sawInFlight = true
				lastTerminalFailure = nil
				if response.Enabled == finalState {
					return nil
				}
			}
			timer.Reset(pollInterval)
		}
	}
}

// waitUntilCleanedUp waits for the backend to finish tearing down a failed enable
// after a disable (DELETE) has been issued. It succeeds once the service reaches
// disabled/idle (or reports disabled via the boolean when status is absent — on
// GCP, with the status feature flag disabled, or on older control planes) and
// fails fast if the teardown itself reports disableFailed. It is bounded by
// cleanupTimeout so a stuck cleanup does not block apply indefinitely.
func (p *PrivateEndpointService) waitUntilCleanedUp(ctx context.Context, organizationId, projectId, clusterId string) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, cleanupTimeout)
	defer cancel()

	// Fire immediately so a terminal disableFailed (or post-teardown 404) is
	// observed without paying a pollInterval delay.
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.ErrPrivateEndpointServiceTimeout

		case <-timer.C:
			response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId)
			if err != nil {
				// A post-teardown 404 means cleanup finished — the service no
				// longer exists on the backend, which is the success condition.
				if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
					return nil
				}
				return err
			}

			if response.Status != nil {
				switch *response.Status {
				case statusDisableFailed:
					return errors.ErrPrivateEndpointServiceDisableFailed
				case statusDisabled, statusIdle:
					return nil
				}
				// enableFailed/disabling/etc: cleanup still in progress, keep polling.
			} else if !response.Enabled {
				return nil
			}
			timer.Reset(pollInterval)
		}
	}
}

// cleanupFailedEnable issues a disable (DELETE) to trigger backend teardown of a
// failed enable and waits for the teardown to complete. The backend allows
// disable from the enableFailed state specifically so this orphaned infra can be
// cleaned up.
func (p *PrivateEndpointService) cleanupFailedEnable(ctx context.Context, organizationId, projectId, clusterId string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	if _, err := p.ClientV1.ExecuteWithRetry(ctx, cfg, nil, p.Token, nil); err != nil {
		return fmt.Errorf("could not trigger cleanup of failed enable: %w", err)
	}

	return p.waitUntilCleanedUp(ctx, organizationId, projectId, clusterId)
}

// handleFailedEnable performs recovery for a terminal enableFailed:
// it triggers backend cleanup of the orphaned resources, removes the resource
// from state so the next apply performs a clean re-create, and surfaces an
// actionable error. State is removed even when cleanup itself fails, because
// leaving a permanently-failed resource in state recreates the stuck-pipeline
// problem; the loud error directs escalation for the rare orphaned-infra case.
func (p *PrivateEndpointService) handleFailedEnable(ctx context.Context, state *tfsdk.State, diags *diag.Diagnostics, organizationId, projectId, clusterId string, cause error) {
	tflog.Error(ctx, "private endpoint service enablement failed; triggering cleanup and removing from state")

	cleanupErr := p.cleanupFailedEnable(ctx, organizationId, projectId, clusterId)
	state.RemoveResource(ctx)

	if cleanupErr != nil {
		diags.AddError(
			"Private endpoint service enablement failed and automatic cleanup did not complete",
			fmt.Sprintf(
				"Enablement failed for cluster %s: %s. Automatic cleanup of the failed resources did not complete: %s. "+
					"There may be orphaned resources in your cloud account; please contact Couchbase Capella Support. "+
					"The resource has been removed from state; re-run terraform apply to retry enablement.",
				clusterId, cause.Error(), cleanupErr.Error(),
			),
		)
		return
	}

	diags.AddError(
		"Private endpoint service enablement failed",
		fmt.Sprintf(
			"Enablement failed for cluster %s: %s. The failed resources were cleaned up automatically and the resource "+
				"has been removed from state; re-run terraform apply to retry enablement.",
			clusterId, cause.Error(),
		),
	)
}

// getServiceStatus retrieves current private endpoint service status.
func (p *PrivateEndpointService) getServiceStatus(ctx context.Context, organizationId, projectId, clusterId string) (*api.GetPrivateEndpointServiceStatusResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService", p.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	status := api.GetPrivateEndpointServiceStatusResponse{}
	err = json.Unmarshal(response.Body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// getServiceState morphs service status into a terraform schema.
func (p *PrivateEndpointService) getServiceState(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.PrivateEndpointService, error) {
	response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	state := providerschema.PrivateEndpointService{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Enabled:        types.BoolValue(response.Enabled),
		Status:         types.StringNull(),
		ServiceName:    types.StringNull(),
	}
	if response.Status != nil {
		state.Status = types.StringValue(*response.Status)
	}
	if response.ServiceName != nil {
		state.ServiceName = types.StringValue(*response.ServiceName)
	}

	return &state, nil
}
