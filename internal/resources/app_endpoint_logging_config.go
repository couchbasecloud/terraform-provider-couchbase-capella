package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	errorMessageWhileLoggingConfigUpdate = "There is an error during app endpoint logging config update. Please check the logging config in capella"
)

var (
	_ resource.Resource                = &LoggingConfig{}
	_ resource.ResourceWithConfigure   = &LoggingConfig{}
	_ resource.ResourceWithImportState = &LoggingConfig{}
)

// Logging Config is the App Endpoint Logging Config resource implementation.
type LoggingConfig struct {
	*providerschema.Data
}

// NewLoggingConfig is a helper function to simplify the provider implementation.
func NewLoggingConfig() resource.Resource {
	return &LoggingConfig{}
}

// Metadata returns the Logging Config resource type name.
func (l *LoggingConfig) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_log_streaming_config"
}

// Schema defines the schema for the Logging Config resource.
func (l *LoggingConfig) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LoggingConfigSchema()
}

// ImportState imports a remote logging config that is not created by Terraform.
func (l *LoggingConfig) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to app_endpoint_name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

// Create sets up the Logging Config for an App Endpoint.
func (l *LoggingConfig) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.LoggingConfig
	var refreshedState *providerschema.LoggingConfig

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
	)

	err := l.upsertLoggingConfig(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing upsert app endpoint logging config",
			errorMessageWhileLoggingConfigUpdate+err.Error(),
		)
		return
	}

	loggingConfig, err := l.getLoggingConfig(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting App Endpoint Logging Config in Capella",
			"Could not get Capella App Endpoint Logging Config for app endpoint with name "+appEndpointName+": "+err.Error(),
		)
		return
	}

	refreshedState = providerschema.NewLoggingConfig(*loggingConfig, organizationId, projectId, clusterId, appServiceId, appEndpointName)

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads the Logging Config information for an App Endpoint.
func (l *LoggingConfig) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.LoggingConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	IDs, err := state.Validate()
	if err != nil {
		tflog.Debug(ctx, "Error validating app endpoint logging config", map[string]interface{}{
			"state": state,
			"err":   err,
		})
		resp.Diagnostics.AddError(
			"Error Validating App Endpoint Logging Config in Capella",
			"Could not validate Capella App Endpoint Logging Config for app endpoint with name "+state.AppEndpointName.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId  = IDs[providerschema.OrganizationId]
		projectId       = IDs[providerschema.ProjectId]
		clusterId       = IDs[providerschema.ClusterId]
		appServiceId    = IDs[providerschema.AppServiceId]
		appEndpointName = IDs[providerschema.AppEndpointName]
	)

	loggingConfig, err := l.getLoggingConfig(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting App Endpoint Logging Config in Capella",
			"Could not get Capella App Endpoint Logging Config for app endpoint with name "+state.AppEndpointName.String()+": "+err.Error(),
		)
		return
	}

	refreshedState := providerschema.NewLoggingConfig(*loggingConfig, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the Logging Config for an App Endpoint.
func (l *LoggingConfig) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan providerschema.LoggingConfig

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
	)

	err := l.upsertLoggingConfig(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing upsert app endpoint logging config",
			errorMessageWhileLoggingConfigUpdate+err.Error(),
		)
		return
	}

	loggingConfig, err := l.getLoggingConfig(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting App Endpoint Logging Config in Capella",
			"Could not get Capella App Endpoint Logging Config for app endpoint with name "+appEndpointName+": "+err.Error(),
		)
		return
	}

	refreshedState := providerschema.NewLoggingConfig(*loggingConfig, organizationId, projectId, clusterId, appServiceId, appEndpointName)

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete removes the App Endpoint Logging Config from the state file.
// The App Endpoint Logging Config in Capella will remain in its last configuration.
func (l *LoggingConfig) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.LoggingConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

}

// Configure adds the provider configured api to the app endpoint logging config resource.
func (l *LoggingConfig) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			errorMessageConfigure+fmt.Sprintf("%T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	l.Data = data
}

// upsertLoggingConfig creates or updates a Logging Config for an App Endpoint.
func (l *LoggingConfig) upsertLoggingConfig(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string, plan providerschema.LoggingConfig) error {

	organizationUUID, projectUUID, clusterUUID, appServiceUUID := l.mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId)

	putLoggingConfigRequest := api.PutAppEndpointLogStreamingConfigJSONRequestBody{
		LogLevel: plan.LogLevel.ValueStringPointer(),
		LogKeys:  providerschema.BaseStringsToStringsPointer(plan.LogKeys),
	}

	_, err := l.ClientV2.PutAppEndpointLogStreamingConfigWithResponse(
		ctx,
		organizationUUID,
		projectUUID,
		clusterUUID,
		appServiceUUID,
		appEndpointName,
		putLoggingConfigRequest,
	)
	if err != nil {
		tflog.Debug(ctx, "error executing update app endpoint logging config", map[string]interface{}{
			"organizationId":             organizationId,
			"projectId":                  projectId,
			"clusterId":                  clusterId,
			"appServiceId":               appServiceId,
			"appEndpointName":            appEndpointName,
			"updateLoggingConfigRequest": putLoggingConfigRequest,
			"err":                        err.Error(),
		})
		return err
	}

	return nil
}

// getLoggingConfig retrieves the Logging Config for an App Endpoint.
func (l *LoggingConfig) getLoggingConfig(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string) (*api.ConsoleLoggingConfig, error) {

	organizationUUID, projectUUID, clusterUUID, appServiceUUID := l.mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId)

	getLoggingConfigResp, err := l.ClientV2.GetAppEndpointLogStreamingConfigWithResponse(
		ctx,
		organizationUUID,
		projectUUID,
		clusterUUID,
		appServiceUUID,
		appEndpointName,
	)
	if err != nil {
		tflog.Debug(ctx, "error getting update app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
			"err":             err.Error(),
		})
		return nil, err
	}

	if getLoggingConfigResp.JSON200 == nil {
		tflog.Debug(ctx, "unexpected status getting app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
		})
		return nil, errors.ErrUnexpectedStatusGettingAppEndpointLoggingConfig
	}

	var loggingConfig api.ConsoleLoggingConfig
	err = json.Unmarshal(getLoggingConfigResp.Body, &loggingConfig)
	if err != nil {
		tflog.Debug(ctx, "error unmarshalling app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
		})
	}

	return &loggingConfig, nil
}

func (l *LoggingConfig) mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId string) (organizationUUID, projectUUID, clusterUUID, appServiceUUID uuid.UUID) {
	organizationUUID, _ = uuid.Parse(organizationId)
	projectUUID, _ = uuid.Parse(projectId)
	clusterUUID, _ = uuid.Parse(clusterId)
	appServiceUUID, _ = uuid.Parse(appServiceId)

	return organizationUUID, projectUUID, clusterUUID, appServiceUUID
}
