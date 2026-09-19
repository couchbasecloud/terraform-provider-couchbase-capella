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
	_ resource.Resource                = &AppServicePrivateEndpointService{}
	_ resource.ResourceWithConfigure   = &AppServicePrivateEndpointService{}
	_ resource.ResourceWithImportState = &AppServicePrivateEndpointService{}
)

const (
	errorMessageWhileEnablingAppServicePrivateEndpointService = "There is an error while enabling the App Service private endpoint service. Please check in Capella to see if there are any hanging resources that have been created, unexpected error: "

	// statusFailed is the single terminal-failure state reported by the App Service
	// private endpoint service status API. Unlike the operational cluster's API,
	// there is no separate enableFailed/disableFailed status; which operation
	// failed is instead disambiguated using TargetState.
	statusFailed = "failed"
)

// AppServicePrivateEndpointService is the App Service scoped private endpoint service
// resource implementation.
type AppServicePrivateEndpointService struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpointService is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpointService() resource.Resource {
	return &AppServicePrivateEndpointService{}
}

// Metadata returns the App Service private endpoint service resource type name.
func (p *AppServicePrivateEndpointService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoint_service"
}

// Schema defines the schema for an App Service private endpoint service resource.
func (p *AppServicePrivateEndpointService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointServiceSchema()
}

// Create enables the private endpoint service on an App Service.
func (p *AppServicePrivateEndpointService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServicePrivateEndpointService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateCreateAppServiceEndpointService(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service private endpoint service request",
			"Could not validate App Service private endpoint service request, unexpected error: "+err.Error(),
		)
		return
	}
	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
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
			"Error enabling App Service private endpoint service",
			errorMessageWhileEnablingAppServicePrivateEndpointService+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeAppServicePrivateEndpointServicePlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = p.waitUntilStatusChanges(ctx, true, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		// Terminal failed-with-targetState-enabled: clean up the orphaned infra and
		// remove the resource so the next apply performs a clean re-create.
		if stderrors.Is(err, errors.ErrAppServicePrivateEndpointServiceEnableFailed) {
			p.handleFailedEnable(ctx, &resp.State, &resp.Diagnostics, organizationId, projectId, clusterId, appServiceId, err)
			return
		}
		resp.Diagnostics.AddError(
			"Error could not enable App Service private endpoint service",
			"Error could not enable App Service private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint service status",
			"Error reading App Service private endpoint service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the App Service private endpoint service status.
func (p *AppServicePrivateEndpointService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoint Service in Capella",
			"Could not read Capella App Service private endpoint service on app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
	)

	refreshedState, err := p.getServiceState(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint service status",
			"Error reading App Service private endpoint service status, unexpected error: "+errString,
		)

		return
	}

	if refreshedState.State.ValueString() == statusFailed && refreshedState.TargetState.ValueString() == statusEnabled {
		tflog.Info(ctx, "App Service private endpoint service is in a failed enablement state; removing from state to force re-create")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update will enable/disable the private endpoint service on an App Service.
func (p *AppServicePrivateEndpointService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config providerschema.AppServicePrivateEndpointService
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		p.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
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
			"Error "+status+" App Service private endpoint service",
			"Error "+status+" App Service private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx,
		config.Enabled.ValueBool(),
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString())
	if err != nil {
		// An enable-flavored update that fails terminally is recovered the same
		// way as Create: clean up and remove from state for a clean retry.
		if config.Enabled.ValueBool() && stderrors.Is(err, errors.ErrAppServicePrivateEndpointServiceEnableFailed) {
			p.handleFailedEnable(ctx, &resp.State, &resp.Diagnostics,
				config.OrganizationId.ValueString(),
				config.ProjectId.ValueString(),
				config.ClusterId.ValueString(),
				config.AppServiceId.ValueString(),
				err)
			return
		}
		resp.Diagnostics.AddError(
			"Error "+status+" App Service private endpoint service",
			"Error "+status+" App Service private endpoint service, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := p.getServiceState(ctx,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service private endpoint service status",
			"Error reading App Service private endpoint service status, unexpected error: "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete disables the private endpoint service on an App Service.
func (p *AppServicePrivateEndpointService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServicePrivateEndpointService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service Private Endpoint Service in Capella",
			"Could not validate Capella App Service private endpoint service on app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
	)

	// If the private endpoint service is already disabled, just remove the resource from the state file.
	if !state.Enabled.ValueBool() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
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
			"Error disabling App Service private endpoint service",
			"Could not disable App Service private endpoint service for app service "+appServiceId+" unexpected error: "+err.Error(),
		)
		return
	}

	err = p.waitUntilStatusChanges(ctx, false, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		// On a terminal failed-with-targetState-disabled we fail fast and keep the
		// resource in state (Terraform retains a resource whose Delete errors). The
		// correct retry for a failed disable is another destroy, which Terraform
		// performs naturally on the next run.
		if stderrors.Is(err, errors.ErrAppServicePrivateEndpointServiceDisableFailed) {
			resp.Diagnostics.AddError(
				"App Service private endpoint service disable failed",
				fmt.Sprintf(
					"Disable failed for app service %s: %s. The resource has been kept in state; "+
						"re-run terraform destroy to retry, and contact Couchbase Capella Support if it persists.",
					appServiceId, err.Error(),
				),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error could not disable App Service private endpoint service",
			"Error could not disable App Service private endpoint service, unexpected error: "+err.Error(),
		)
	}
}

// Configure adds the provider configured api to the App Service private endpoint service resource.
func (p *AppServicePrivateEndpointService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports an App Service private endpoint service status.
func (p *AppServicePrivateEndpointService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// validateCreateAppServiceEndpointService ensures organization id, project id, cluster id, and
// app service id are valued.
func validateCreateAppServiceEndpointService(plan providerschema.AppServicePrivateEndpointService) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}

	return nil
}

// initializeAppServicePrivateEndpointServicePlan initializes an instance of
// providerschema.AppServicePrivateEndpointService with the specified plan. It marks all
// computed fields as null.
func initializeAppServicePrivateEndpointServicePlan(plan providerschema.AppServicePrivateEndpointService) providerschema.AppServicePrivateEndpointService {
	if plan.Enabled.IsNull() || plan.Enabled.IsUnknown() {
		plan.Enabled = types.BoolNull()
	}
	// state is computed; never persist an unknown value to state.
	if plan.State.IsNull() || plan.State.IsUnknown() {
		plan.State = types.StringNull()
	}
	// target_state is computed; never persist an unknown value to state.
	if plan.TargetState.IsNull() || plan.TargetState.IsUnknown() {
		plan.TargetState = types.StringNull()
	}
	return plan
}

// waitUntilStatusChanges waits until the private endpoint service reaches the desired
// state on the App Service. It keeps polling on the transient states (enabling/
// disabling) and fails fast on the terminal failed state — but only once it has seen
// evidence the current operation is in flight, so a residual failed state from a
// prior attempt is not mistaken for failure of the operation we just issued. Because
// the App Service status API reports a single "failed" state rather than separate
// enableFailed/disableFailed states, TargetState is used to disambiguate which
// operation failed. The overall timeout remains as a backstop.
func (p *AppServicePrivateEndpointService) waitUntilStatusChanges(ctx context.Context, finalState bool, organizationId, projectId, clusterId, appServiceId string) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, statusChangeTimeout)
	defer cancel()

	timer := time.NewTimer(0)
	defer timer.Stop()

	// sawInFlight flips true once we observe a transient state (enabling/disabling),
	// which is our evidence that the backend has progressed past whatever residual
	// state was present before our POST/DELETE. Until then, a reported failed state
	// could be stale and is treated as transient.
	var sawInFlight bool

	// lastTerminalFailure records a terminal failure we observed but deferred on
	// because sawInFlight was still false (i.e. it might have been residual from a
	// prior attempt). If the backend never transitions and we hit the overall
	// timeout, this lets us surface the typed failure error so the caller still
	// routes to cleanup, instead of a generic timeout that would leave orphaned
	// infra behind and skip state removal.
	var lastTerminalFailure error

	for {
		select {
		case <-ctx.Done():
			if lastTerminalFailure != nil {
				return lastTerminalFailure
			}
			return errors.ErrAppServicePrivateEndpointServiceTimeout

		case <-timer.C:
			response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId, appServiceId)
			if err != nil {
				if ctx.Err() != nil {
					if lastTerminalFailure != nil {
						return lastTerminalFailure
					}
					return errors.ErrAppServicePrivateEndpointServiceTimeout
				}
				return err
			}

			// State is best-effort and may be absent; keep polling until it appears
			// or the overall timeout is hit.
			if response.State == nil {
				timer.Reset(pollInterval)
				continue
			}

			switch *response.State {
			case statusFailed:
				failure := errors.ErrAppServicePrivateEndpointServiceDisableFailed
				if response.TargetState != nil && *response.TargetState == statusEnabled {
					failure = errors.ErrAppServicePrivateEndpointServiceEnableFailed
				}
				if sawInFlight {
					return failure
				}
				// Possibly-stale: defer, but remember it so a never-transitioning
				// failed state is still routed to cleanup on timeout.
				lastTerminalFailure = failure
			case statusEnabling, statusDisabling:
				sawInFlight = true
				// The backend progressed, so any earlier terminal status was
				// genuinely stale; a later stall is a real transition timeout.
				lastTerminalFailure = nil
			case statusEnabled, statusDisabled:
				sawInFlight = true
				lastTerminalFailure = nil
				if (*response.State == statusEnabled) == finalState {
					return nil
				}
			}
			timer.Reset(pollInterval)
		}
	}
}

// waitUntilCleanedUp waits for the backend to finish tearing down a failed enable
// after a disable (DELETE) has been issued. It succeeds once the service reaches
// disabled and fails fast if the teardown itself reports a failed state. It is
// bounded by cleanupTimeout so a stuck cleanup does not block apply indefinitely.
func (p *AppServicePrivateEndpointService) waitUntilCleanedUp(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, cleanupTimeout)
	defer cancel()

	// Fire immediately so a terminal failed state (or post-teardown 404) is
	// observed without paying a pollInterval delay.
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.ErrAppServicePrivateEndpointServiceTimeout

		case <-timer.C:
			response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId, appServiceId)
			if err != nil {
				// A post-teardown 404 means cleanup finished — the service no
				// longer exists on the backend, which is the success condition.
				if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
					return nil
				}
				return err
			}

			if response.State != nil {
				switch *response.State {
				case statusFailed:
					// Unlike the cluster-level API, there is only one failed
					// state, so TargetState is needed to tell a genuine cleanup
					// (disable) failure apart from a residual failed state left
					// over from the original enable failure that triggered this
					// cleanup — the backend may not have caught up to our DELETE
					// yet. Only the former is terminal here.
					if response.TargetState != nil && *response.TargetState == statusDisabled {
						return errors.ErrAppServicePrivateEndpointServiceDisableFailed
					}
				case statusDisabled:
					return nil
				}
				// enabling/disabling/enabled, or a stale failed state: cleanup
				// still in progress, keep polling.
			}
			timer.Reset(pollInterval)
		}
	}
}

// cleanupFailedEnable issues a disable (DELETE) to trigger backend teardown of a
// failed enable and waits for the teardown to complete. The backend allows disable
// from the failed state specifically so this orphaned infra can be cleaned up.
func (p *AppServicePrivateEndpointService) cleanupFailedEnable(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		p.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	if _, err := p.ClientV1.ExecuteWithRetry(ctx, cfg, nil, p.Token, nil); err != nil {
		return fmt.Errorf("could not trigger cleanup of failed enable: %w", err)
	}

	return p.waitUntilCleanedUp(ctx, organizationId, projectId, clusterId, appServiceId)
}

// handleFailedEnable performs recovery for a terminal enablement failure: it
// triggers backend cleanup of the orphaned resources, removes the resource from
// state so the next apply performs a clean re-create, and surfaces an actionable
// error. State is removed even when cleanup itself fails, because leaving a
// permanently-failed resource in state recreates the stuck-pipeline problem; the
// loud error directs escalation for the rare orphaned-infra case.
func (p *AppServicePrivateEndpointService) handleFailedEnable(ctx context.Context, state *tfsdk.State, diags *diag.Diagnostics, organizationId, projectId, clusterId, appServiceId string, cause error) {
	tflog.Error(ctx, "App Service private endpoint service enablement failed; triggering cleanup and removing from state")

	cleanupErr := p.cleanupFailedEnable(ctx, organizationId, projectId, clusterId, appServiceId)
	state.RemoveResource(ctx)

	if cleanupErr != nil {
		diags.AddError(
			"App Service private endpoint service enablement failed and automatic cleanup did not complete",
			fmt.Sprintf(
				"Enablement failed for app service %s: %s. Automatic cleanup of the failed resources did not complete: %s. "+
					"There may be orphaned resources in your cloud account; please contact Couchbase Capella Support. "+
					"The resource has been removed from state; re-run terraform apply to retry enablement.",
				appServiceId, cause.Error(), cleanupErr.Error(),
			),
		)
		return
	}

	diags.AddError(
		"App Service private endpoint service enablement failed",
		fmt.Sprintf(
			"Enablement failed for app service %s: %s. The failed resources were cleaned up automatically and the "+
				"resource has been removed from state; re-run terraform apply to retry enablement.",
			appServiceId, cause.Error(),
		),
	)
}

// getServiceStatus retrieves the current App Service private endpoint service status.
func (p *AppServicePrivateEndpointService) getServiceStatus(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*api.GetAppServicePrivateEndpointServiceStatusResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService", p.HostURL, organizationId, projectId, clusterId, appServiceId)
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

	status := api.GetAppServicePrivateEndpointServiceStatusResponse{}
	err = json.Unmarshal(response.Body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// getServiceState morphs the App Service private endpoint service status into a
// terraform schema.
func (p *AppServicePrivateEndpointService) getServiceState(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) (*providerschema.AppServicePrivateEndpointService, error) {
	response, err := p.getServiceStatus(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, err
	}

	state := providerschema.AppServicePrivateEndpointService{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
		Enabled:        types.BoolValue(response.State != nil && *response.State == statusEnabled),
		State:          types.StringNull(),
		TargetState:    types.StringNull(),
	}
	if response.State != nil {
		state.State = types.StringValue(*response.State)
	}
	if response.TargetState != nil {
		state.TargetState = types.StringValue(*response.TargetState)
	}

	return &state, nil
}
