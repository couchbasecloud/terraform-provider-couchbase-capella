package resources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &AppServiceLogStreaming{}
	_ resource.ResourceWithConfigure      = &AppServiceLogStreaming{}
	_ resource.ResourceWithImportState    = &AppServiceLogStreaming{}
	_ resource.ResourceWithValidateConfig = &AppServiceLogStreaming{}
)

const (
	// errorMessageAfterLogStreamingCreation is the error message when log streaming creation
	// is initiated but the status check fails.
	errorMessageAfterLogStreamingCreation = "Log Streaming configuration is initiated, but encountered an error while checking the current" +
		" state. Please run `terraform plan` after a few minutes to know the current status." +
		" Additionally, run `terraform apply --refresh-only` to update the state from remote, unexpected error: "

	// errorMessageWhileLogStreamingCreation is the error message when log streaming creation fails.
	errorMessageWhileLogStreamingCreation = "There is an error during log streaming configuration. Please check in Capella to see Log Streaming " +
		"has been enabled, unexpected error: "
)

// AppServiceLogStreaming is the App Service Log Streaming resource implementation.
type AppServiceLogStreaming struct {
	*providerschema.Data
}

// NewAppServiceLogStreaming is a helper function to simplify the provider implementation.
func NewAppServiceLogStreaming() resource.Resource {
	return &AppServiceLogStreaming{}
}

// Metadata returns the resource type name.
func (r *AppServiceLogStreaming) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_log_streaming"
}

// Schema defines the schema for the resource.
func (r *AppServiceLogStreaming) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceLogStreamingSchema()
}

// Configure adds the provider configured client to the resource.
func (r *AppServiceLogStreaming) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new App Service Log Streaming configuration.
func (r *AppServiceLogStreaming) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceLogStreaming
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()
	appServiceId := plan.AppServiceId.ValueString()

	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	// Build the API request
	postReq, err := r.buildPostLogStreamingRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error building Log Streaming request",
			"Could not build Log Streaming configuration request: "+err.Error(),
		)
		return
	}

	// Call the API using ClientV2
	response, err := r.ClientV2.PostAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
		postReq,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Log Streaming configuration",
			errorMessageWhileLogStreamingCreation+err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.AddError(
			"Error creating Log Streaming configuration",
			fmt.Sprintf("Unexpected response while creating Log Streaming config: %s", string(response.Body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for the log streaming to be enabled
	err = r.waitForLogStreamingState(ctx, organizationId, projectId, clusterId, appServiceId, apigen.GetLogStreamingResponseConfigStateEnabled)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Log Streaming to be enabled",
			errorMessageAfterLogStreamingCreation+err.Error(),
		)
		return
	}

	// Refresh the state from the API
	refreshedState, err := r.refreshLogStreaming(ctx, organizationId, projectId, clusterId, appServiceId, plan.Credentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing Log Streaming state",
			errorMessageAfterLogStreamingCreation+err.Error(),
		)
		return
	}

	// Set the final state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *AppServiceLogStreaming) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceLogStreaming
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Log Streaming configuration",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	organizationId := resourceIDs[providerschema.OrganizationId]
	projectId := resourceIDs[providerschema.ProjectId]
	clusterId := resourceIDs[providerschema.ClusterId]
	appServiceId := resourceIDs[providerschema.AppServiceId]

	// Refresh the state from the API
	refreshedState, err := r.refreshLogStreaming(ctx, organizationId, projectId, clusterId, appServiceId, state.Credentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Log Streaming configuration",
			"Could not refresh Log Streaming configuration: "+err.Error(),
		)
		return
	}

	if refreshedState.ConfigState.ValueString() == string(apigen.GetLogStreamingResponseConfigStateDisabled) {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the App Service Log Streaming configuration.
func (r *AppServiceLogStreaming) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.AppServiceLogStreaming
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating log streaming configuration",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	organizationId := resourceIDs[providerschema.OrganizationId]
	projectId := resourceIDs[providerschema.ProjectId]
	clusterId := resourceIDs[providerschema.ClusterId]
	appServiceId := resourceIDs[providerschema.AppServiceId]

	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	// Note: output_type changes are handled by RequiresReplace, so we only get here
	// if credentials have changed but output_type stays the same.
	// The API uses the same POST endpoint for both create and update.

	// Build the API request
	postReq, err := r.buildPostLogStreamingRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error building Log Streaming update request",
			"Could not build Log Streaming configuration request: "+err.Error(),
		)
		return
	}

	// Call the API using ClientV2
	response, err := r.ClientV2.PostAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
		postReq,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Log Streaming configuration",
			"Could not update Log Streaming configuration: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.AddError(
			"Error updating Log Streaming configuration",
			fmt.Sprintf("Unexpected response while updating Log Streaming config: %s", string(response.Body)),
		)
		return
	}

	// Wait for the log streaming to be enabled
	err = r.waitForLogStreamingState(ctx, organizationId, projectId, clusterId, appServiceId, apigen.GetLogStreamingResponseConfigStateEnabled)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Log Streaming update to complete",
			"Could not update Log Streaming configuration: "+err.Error(),
		)
		return
	}

	// Refresh the state from the API
	refreshedState, err := r.refreshLogStreaming(ctx, organizationId, projectId, clusterId, appServiceId, plan.Credentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing Log Streaming state after update",
			"Could not refresh Log Streaming configuration: "+err.Error(),
		)
		return
	}

	// Set the final state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the App Service Log Streaming configuration.
func (r *AppServiceLogStreaming) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceLogStreaming
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Log Streaming configuration",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	organizationId := resourceIDs[providerschema.OrganizationId]
	projectId := resourceIDs[providerschema.ProjectId]
	clusterId := resourceIDs[providerschema.ClusterId]
	appServiceId := resourceIDs[providerschema.AppServiceId]

	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	// Call the delete API
	response, err := r.ClientV2.DeleteAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Log Streaming configuration",
			"Could not delete Log Streaming configuration: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.AddError(
			"Error deleting Log Streaming configuration",
			fmt.Sprintf("Unexpected response while disabling Log Streaming: %s", string(response.Body)),
		)
		return
	}

	// Wait for the log streaming to be disabled
	err = r.waitForLogStreamingState(ctx, organizationId, projectId, clusterId, appServiceId, apigen.GetLogStreamingResponseConfigStateDisabled)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Log Streaming deletion to complete",
			"Could not delete Log Streaming configuration: "+err.Error(),
		)
		return
	}
}

// ValidateConfig validates the resource configuration.
// It checks the output type is valid and matches the correct credentials object and that no other credentials objects are set.
func (r *AppServiceLogStreaming) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config providerschema.AppServiceLogStreaming
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If credentials is unknown (e.g. passed via a variable), we can't validate.
	if config.Credentials.IsUnknown() || config.Credentials.IsNull() {
		return
	}

	// If output_type is unknown, we can't validate the match.
	if config.OutputType.IsUnknown() {
		return
	}

	outputType := config.OutputType.ValueString()

	// Extract the LogStreamingCredentials from the types.Object
	creds, diags := config.AsLogStreamingCredentials(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if creds == nil {
		return
	}

	// Map each output_type to a check for whether its matching credential block is provided.
	type credentialCheck struct {
		outputType string
		isPresent  bool
	}

	checks := []credentialCheck{
		{string(apigen.GetLogStreamingResponseOutputTypeDatadog), !creds.Datadog.IsNull() && !creds.Datadog.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeDynatrace), !creds.Dynatrace.IsNull() && !creds.Dynatrace.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeElastic), !creds.Elastic.IsNull() && !creds.Elastic.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeGenericHttp), !creds.GenericHttp.IsNull() && !creds.GenericHttp.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeLoki), !creds.Loki.IsNull() && !creds.Loki.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeSplunk), !creds.Splunk.IsNull() && !creds.Splunk.IsUnknown()},
		{string(apigen.GetLogStreamingResponseOutputTypeSumologic), !creds.Sumologic.IsNull() && !creds.Sumologic.IsUnknown()},
	}

	// Validate output_type is a supported value and find the matching credential check.
	var matchFound bool
	for _, check := range checks {
		if check.outputType == outputType {
			matchFound = true
			if !check.isPresent {
				resp.Diagnostics.AddAttributeError(
					path.Root("credentials"),
					"Missing Credential Configuration",
					fmt.Sprintf("credentials.%s must be configured when output_type is %q", check.outputType, outputType),
				)
				return
			}
		} else if check.isPresent {
			resp.Diagnostics.AddAttributeError(
				path.Root("credentials"),
				"Invalid Credential Configuration",
				fmt.Sprintf("credentials.%s must not be configured when output_type is %q", check.outputType, outputType),
			)
			return
		}
	}

	if !matchFound {
		resp.Diagnostics.AddAttributeError(
			path.Root("output_type"),
			"Invalid Attribute Configuration",
			fmt.Sprintf("Unsupported output_type %q. Please read the documentation for supported values.", outputType),
		)
		return
	}
}

// ImportState imports a remote resource that is not managed by Terraform.
func (r *AppServiceLogStreaming) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use passthrough ID which will be parsed by Validate()
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// buildPostLogStreamingRequest builds the API request from the Terraform plan.
func (r *AppServiceLogStreaming) buildPostLogStreamingRequest(ctx context.Context, plan providerschema.AppServiceLogStreaming) (apigen.PostLogStreamingRequest, error) {
	outputType := apigen.PostLogStreamingRequestOutputType(plan.OutputType.ValueString())

	var apiCredentials apigen.PostLogStreamingRequest_Credentials

	if plan.Credentials.IsNull() || plan.Credentials.IsUnknown() {
		return apigen.PostLogStreamingRequest{}, fmt.Errorf("credentials are required")
	}

	creds, diags := plan.AsLogStreamingCredentials(ctx)
	if diags.HasError() {
		return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract credentials: %s", diags.Errors())
	}
	if creds == nil {
		return apigen.PostLogStreamingRequest{}, fmt.Errorf("credentials are required")
	}

	// Set credentials based on which provider is configured
	switch outputType {
	case apigen.PostLogStreamingRequestOutputTypeDatadog:
		dd, diags := creds.AsDatadogCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract datadog credentials: %s", diags.Errors())
		}
		if dd == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("datadog credentials are required when output_type is 'datadog'")
		}
		err := apiCredentials.FromDatadog(apigen.Datadog{
			Url:    dd.Url.ValueString(),
			ApiKey: dd.ApiKey.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set datadog credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeDynatrace:
		dt, diags := creds.AsDynatraceCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract dynatrace credentials: %s", diags.Errors())
		}
		if dt == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("dynatrace credentials are required when output_type is 'dynatrace'")
		}
		err := apiCredentials.FromDynatrace(apigen.Dynatrace{
			Url:      dt.Url.ValueString(),
			ApiToken: dt.ApiToken.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set dynatrace credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeElastic:
		el, diags := creds.AsElasticCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract elastic credentials: %s", diags.Errors())
		}
		if el == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("elastic credentials are required when output_type is 'elastic'")
		}
		err := apiCredentials.FromElastic(apigen.Elastic{
			Url:      el.Url.ValueString(),
			User:     el.User.ValueString(),
			Password: el.Password.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set elastic credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeGenericHttp:
		gh, diags := creds.AsGenericHttpCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract generic_http credentials: %s", diags.Errors())
		}
		if gh == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("generic_http credentials are required when output_type is 'generic_http'")
		}
		genericHttp := apigen.GenericHttp{
			Url: gh.Url.ValueString(),
		}
		// User and password are optional for generic_http
		if !gh.User.IsNull() && !gh.User.IsUnknown() {
			user := gh.User.ValueString()
			genericHttp.User = &user
		}
		if !gh.Password.IsNull() && !gh.Password.IsUnknown() {
			password := gh.Password.ValueString()
			genericHttp.Password = &password
		}
		err := apiCredentials.FromGenericHttp(genericHttp)
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set generic_http credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeLoki:
		lk, diags := creds.AsLokiCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract loki credentials: %s", diags.Errors())
		}
		if lk == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("loki credentials are required when output_type is 'loki'")
		}
		err := apiCredentials.FromLoki(apigen.Loki{
			Url:      lk.Url.ValueString(),
			User:     lk.User.ValueString(),
			Password: lk.Password.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set loki credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeSplunk:
		sp, diags := creds.AsSplunkCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract splunk credentials: %s", diags.Errors())
		}
		if sp == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("splunk credentials are required when output_type is 'splunk'")
		}
		err := apiCredentials.FromSplunk(apigen.Splunk{
			Url:         sp.Url.ValueString(),
			SplunkToken: sp.SplunkToken.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set splunk credentials: %w", err)
		}

	case apigen.PostLogStreamingRequestOutputTypeSumologic:
		sl, diags := creds.AsSumologicCredentials(ctx)
		if diags.HasError() {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to extract sumologic credentials: %s", diags.Errors())
		}
		if sl == nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("sumologic credentials are required when output_type is 'sumologic'")
		}
		err := apiCredentials.FromSumologic(apigen.Sumologic{
			Url: sl.Url.ValueString(),
		})
		if err != nil {
			return apigen.PostLogStreamingRequest{}, fmt.Errorf("failed to set sumologic credentials: %w", err)
		}

	default:
		return apigen.PostLogStreamingRequest{}, fmt.Errorf("unsupported output_type: %s", outputType)
	}

	return apigen.PostLogStreamingRequest{
		OutputType:  outputType,
		Credentials: apiCredentials,
	}, nil
}

// refreshLogStreaming retrieves the current state of the log streaming configuration from the API.
func (r *AppServiceLogStreaming) refreshLogStreaming(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	existingCredentials types.Object,
) (*providerschema.AppServiceLogStreaming, error) {
	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IDs: %w", err)
	}

	response, err := r.ClientV2.GetAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", response.StatusCode(), string(response.Body))
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("API returned empty response body")
	}

	// Extract values from response
	// Note: Credentials are not returned by the API, so we preserve them from the plan/state
	return providerschema.NewAppServiceLogStreaming(
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		response.JSON200,
		existingCredentials,
	), nil
}

// waitForLogStreamingState waits for the log streaming configuration to no longer be transitioning.
func (r *AppServiceLogStreaming) waitForLogStreamingState(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	targetState apigen.GetLogStreamingResponseConfigState,
) error {
	// Log Streaming state transition should usually only take up to a minute when all nodes are healthy, but allow for extra time in case a node is having transient issues
	const timeout = time.Minute * 3
	const sleepDuration = time.Second * 3

	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, appServiceUUID, err := r.parseUUIDs(organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return fmt.Errorf("failed to parse IDs: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	timer := time.NewTimer(sleepDuration)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for log streaming to reach state '%s'", string(targetState))

		case <-timer.C:
			response, err := r.ClientV2.GetAppServiceLogStreamingWithResponse(
				ctx,
				orgUUID,
				projUUID,
				clusterUUID,
				appServiceUUID,
			)
			if err != nil {
				return fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
			}

			if response.JSON200 == nil || response.JSON200.ConfigState == nil {
				tflog.Debug(ctx, fmt.Sprintf("No config state field detected. API returned status %d: %s", response.StatusCode(), string(response.Body)))
				return fmt.Errorf("API returned empty response body or missing config state")
			}

			currentState := string(*response.JSON200.ConfigState)
			tflog.Info(ctx, fmt.Sprintf("log streaming config state: %s, waiting for: %s", currentState, string(targetState)))

			// Check if we've reached the target state
			if currentState == string(targetState) {
				tflog.Debug(ctx, "target log streaming state reached: "+currentState)
				return nil
			}

			// Check if we're in a final state that's not our target
			if isFinalLogStreamingState(currentState) && currentState != string(targetState) {
				return fmt.Errorf("log streaming reached final state '%s' instead of expected '%s'", currentState, string(targetState))
			}

			timer.Reset(sleepDuration)
		}
	}
}

// isFinalLogStreamingState returns true if the state is a non-transitioning state.
func isFinalLogStreamingState(state string) bool {
	switch state {
	case string(apigen.GetLogStreamingResponseConfigStateEnabled),
		string(apigen.GetLogStreamingResponseConfigStateDisabled),
		string(apigen.GetLogStreamingResponseConfigStatePaused),
		string(apigen.GetLogStreamingResponseConfigStateErrored):
		return true
	default:
		return false
	}
}

// parseUUIDs parses the string IDs into UUID types for the generated API client.
func (r *AppServiceLogStreaming) parseUUIDs(organizationId, projectId, clusterId, appServiceId string) (uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, error) {
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid organization_id: %w", err)
	}

	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid project_id: %w", err)
	}

	clusterUUID, err := uuid.Parse(clusterId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid cluster_id: %w", err)
	}

	appServiceUUID, err := uuid.Parse(appServiceId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid app_service_id: %w", err)
	}

	return orgUUID, projUUID, clusterUUID, appServiceUUID, nil
}
