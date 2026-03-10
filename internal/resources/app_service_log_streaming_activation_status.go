package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AppServiceLogStreamingActivationStatus{}
	_ resource.ResourceWithConfigure   = &AppServiceLogStreamingActivationStatus{}
	_ resource.ResourceWithImportState = &AppServiceLogStreamingActivationStatus{}
)

// AppServiceLogStreamingActivationStatus is the resource implementation for managing
// the activation status (paused/enabled) of log streaming on an App Service.
type AppServiceLogStreamingActivationStatus struct {
	*providerschema.Data
}

// NewAppServiceLogStreamingActivationStatus is a helper function to simplify the provider implementation.
func NewAppServiceLogStreamingActivationStatus() resource.Resource {
	return &AppServiceLogStreamingActivationStatus{}
}

// Metadata returns the resource type name.
func (r *AppServiceLogStreamingActivationStatus) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_log_streaming_activation_status"
}

// Schema defines the schema for the resource.
func (r *AppServiceLogStreamingActivationStatus) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceLogStreamingActivationStatusSchema()
}

// Configure adds the provider configured client to the resource.
func (r *AppServiceLogStreamingActivationStatus) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.Data = data
}

// Create reads the current log streaming state and, if necessary, calls the
// pause or enable API to reach the desired state. If the current state
// already matches the desired state, no API call is made.
func (r *AppServiceLogStreamingActivationStatus) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceLogStreamingActivationStatus
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	desiredState := apigen.GetLogStreamingResponseConfigState(plan.State.ValueString())

	tflog.Debug(ctx, "creating log streaming activation state resource", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"app_service_id":  appServiceId,
		"state":           desiredState,
	})

	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	// Read the current state to determine if an API call is needed.
	currentConfigState, err := r.getCurrentConfigState(ctx, orgUUID, projUUID, clusterUUID, appServiceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading current Log Streaming activation status",
			"Could not read current Log Streaming config state: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "current log streaming config state during create", map[string]interface{}{
		"current_config_state": currentConfigState,
		"desired_state":        desiredState,
	})

	// If the current state already matches the desired state, skip the API call.
	if currentConfigState == desiredState {
		tflog.Info(ctx, "log streaming activation status already matches desired state")
	} else {
		err = r.applyActivationStatus(ctx, orgUUID, projUUID, clusterUUID, appServiceUUID, desiredState)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting Log Streaming activation status",
				"Could not set Log Streaming activation status: "+err.Error(),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data from the API.
func (r *AppServiceLogStreamingActivationStatus) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceLogStreamingActivationStatus
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Log Streaming activation status",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	organizationId := resourceIDs[providerschema.OrganizationId]
	projectId := resourceIDs[providerschema.ProjectId]
	clusterId := resourceIDs[providerschema.ClusterId]
	appServiceId := resourceIDs[providerschema.AppServiceId]

	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "reading log streaming activation state", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"app_service_id":  appServiceId,
	})

	currentConfigState, err := r.getCurrentConfigState(ctx, orgUUID, projUUID, clusterUUID, appServiceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading current Log Streaming activation status",
			"Could not read current Log Streaming config state: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "read log streaming config state", map[string]interface{}{
		"config_state": currentConfigState,
	})

	// If log streaming is disabled or disabling, the activation state concept is
	// meaningless - remove the resource from state.
	if currentConfigState == apigen.GetLogStreamingResponseConfigStateDisabled ||
		currentConfigState == apigen.GetLogStreamingResponseConfigStateDisabling {
		tflog.Info(ctx, "log streaming is disabled/disabling, removing activation status resource from state")
		resp.State.RemoveResource(ctx)
		return
	}

	refreshedState := providerschema.NewAppServiceLogStreamingActivationStatus(
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		currentConfigState,
	)
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update is triggered when state changes. It calls the pause or enable
// API and waits for the state transition to complete.
func (r *AppServiceLogStreamingActivationStatus) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppServiceLogStreamingActivationStatus
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := plan.Validate(); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Log Streaming activation status",
			"Could not validate plan: "+err.Error(),
		)
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()
	desiredState := apigen.GetLogStreamingResponseConfigState(plan.State.ValueString())

	tflog.Debug(ctx, "updating log streaming activation status", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"app_service_id":  appServiceId,
		"state":           desiredState,
	})

	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	err = r.applyActivationStatus(ctx, orgUUID, projUUID, clusterUUID, appServiceUUID, desiredState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Log Streaming activation status",
			"Could not update Log Streaming activation status: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete is a no-op. The activation state is not a real resource to destroy;
// removing it from Terraform state is sufficient.
func (r *AppServiceLogStreamingActivationStatus) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting log streaming activation status resource - removing from state")
}

// ImportState imports a remote log streaming activation state into Terraform.
func (r *AppServiceLogStreamingActivationStatus) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// applyActivationStatus calls the pause or enable API based on the desired state and
// waits for the state transition to complete.
func (r *AppServiceLogStreamingActivationStatus) applyActivationStatus(
	ctx context.Context,
	orgUUID, projUUID, clusterUUID, appServiceUUID uuid.UUID,
	desiredState apigen.GetLogStreamingResponseConfigState,
) error {
	switch desiredState {
	case apigen.GetLogStreamingResponseConfigStatePaused:
		tflog.Debug(ctx, "pausing log streaming")
		response, err := r.ClientV2.PauseAppServiceLogStreamingWithResponse(
			ctx,
			orgUUID,
			projUUID,
			clusterUUID,
			appServiceUUID,
		)
		if err != nil {
			return fmt.Errorf("error calling pause API: %w", err)
		}
		if response.StatusCode() != http.StatusAccepted {
			return fmt.Errorf("failed to pause log streaming: %s", string(response.Body))
		}

		tflog.Debug(ctx, "log streaming now pausing, waiting for paused state")
		return waitForLogStreamingState(ctx, r.ClientV2, orgUUID, projUUID, clusterUUID, appServiceUUID, apigen.GetLogStreamingResponseConfigStatePaused)

	case apigen.GetLogStreamingResponseConfigStateEnabled:
		tflog.Debug(ctx, "enabling log streaming")
		response, err := r.ClientV2.ResumeAppServiceLogStreamingWithResponse(
			ctx,
			orgUUID,
			projUUID,
			clusterUUID,
			appServiceUUID,
		)
		if err != nil {
			return fmt.Errorf("error calling resume API: %w", err)
		}
		if response.StatusCode() != http.StatusAccepted {
			return fmt.Errorf("failed to enable log streaming: %s", string(response.Body))
		}

		tflog.Debug(ctx, "log streaming now enabling, waiting for enabled state")
		return waitForLogStreamingState(ctx, r.ClientV2, orgUUID, projUUID, clusterUUID, appServiceUUID, apigen.GetLogStreamingResponseConfigStateEnabled)

	default:
		return fmt.Errorf("unsupported state: %s", desiredState)
	}
}

// getCurrentConfigState reads the current log streaming config state from the API.
func (r *AppServiceLogStreamingActivationStatus) getCurrentConfigState(
	ctx context.Context,
	orgUUID, projUUID, clusterUUID, appServiceUUID uuid.UUID,
) (apigen.GetLogStreamingResponseConfigState, error) {
	response, err := r.ClientV2.GetAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	if response.JSON200 == nil || response.JSON200.ConfigState == nil {
		return "", fmt.Errorf("failed to get current log streaming config state: %s", string(response.Body))
	}

	return *response.JSON200.ConfigState, nil
}

// parseUUIDs parses the string IDs into UUID types for the generated API client.
func (r *AppServiceLogStreamingActivationStatus) parseUUIDs(organizationId, projectId, clusterId, appServiceId string) (uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, error) {
	orgUUID, err := utils.ParseUUID("organization_id", organizationId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}

	projUUID, err := utils.ParseUUID("project_id", projectId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}

	clusterUUID, err := utils.ParseUUID("cluster_id", clusterId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}

	appServiceUUID, err := utils.ParseUUID("app_service_id", appServiceId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}

	return orgUUID, projUUID, clusterUUID, appServiceUUID, nil
}
