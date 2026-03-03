package datasources

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LoggingConfig{}
	_ datasource.DataSourceWithConfigure = &LoggingConfig{}
)

// LoggingConfig is the App Endpoint Logging Config data source implementation.
type LoggingConfig struct {
	*providerschema.Data
}

// NewLoggingConfig is a helper function to simplify the provider implementation.
func NewAppEndpointLoggingConfig() datasource.DataSource {
	return &LoggingConfig{}
}

// Metadata returns the LoggingConfig data source type name.
func (l *LoggingConfig) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_log_streaming_config"
}

// Schema defines the schema for the LoggingConfig data source.
func (l *LoggingConfig) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = LoggingConfigSchema()
}

// Read reads the Logging Config information for an App Endpoint.
func (l *LoggingConfig) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.LoggingConfig
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId  = state.OrganizationId.ValueString()
		projectId       = state.ProjectId.ValueString()
		clusterId       = state.ClusterId.ValueString()
		appServiceId    = state.AppServiceId.ValueString()
		appEndpointName = state.AppEndpointName.ValueString()
	)

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
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Logging Config",
			fmt.Sprintf("Could not read logging config for app endpoint %s, unexpected error: %s", appEndpointName, err.Error()),
		)
		tflog.Debug(ctx, "error getting app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
			"err":             err.Error(),
		})
		return
	}

	if getLoggingConfigResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected Status Reading App Endpoint Logging Config",
			"Could not read logging config for app endpoint with name "+state.AppEndpointName.String()+", unexpected status: "+string(getLoggingConfigResp.Body),
		)
		tflog.Debug(ctx, "unexpected status getting app endpoint logging config", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
		})
		return
	}

	state = *providerschema.NewLoggingConfig(*getLoggingConfigResp.JSON200, organizationId, projectId, clusterId, appServiceId, appEndpointName)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the Logging Config data source.
func (l *LoggingConfig) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	l.Data = data
}

func (l *LoggingConfig) mapIDsToUUIDs(organizationId, projectId, clusterId, appServiceId string) (organizationUUID, projectUUID, clusterUUID, appServiceUUID uuid.UUID) {
	organizationUUID, _ = uuid.Parse(organizationId)
	projectUUID, _ = uuid.Parse(projectId)
	clusterUUID, _ = uuid.Parse(clusterId)
	appServiceUUID, _ = uuid.Parse(appServiceId)

	return organizationUUID, projectUUID, clusterUUID, appServiceUUID
}
