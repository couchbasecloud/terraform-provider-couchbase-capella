package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppServiceLogStreaming{}
	_ datasource.DataSourceWithConfigure = &AppServiceLogStreaming{}
)

// AppServiceLogStreaming is the app service log streaming data source implementation.
type AppServiceLogStreaming struct {
	*providerschema.Data
}

// NewAppServiceLogStreaming is a helper function to simplify the provider implementation.
func NewAppServiceLogStreaming() datasource.DataSource {
	return &AppServiceLogStreaming{}
}

// Metadata returns the data source type name.
func (d *AppServiceLogStreaming) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_log_streaming"
}

// Schema defines the schema for the data source.
func (d *AppServiceLogStreaming) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceLogStreamingSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServiceLogStreaming) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Data = data
}

// Read refreshes the Terraform state with the latest data from the API.
func (d *AppServiceLogStreaming) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServiceLogStreamingData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := config.OrganizationId.ValueString()
	projectId := config.ProjectId.ValueString()
	clusterId := config.ClusterId.ValueString()
	appServiceId := config.AppServiceId.ValueString()

	// Parse string IDs to UUIDs for the API client
	uuids, err := utils.ParseUUIDs(
		utils.IDField{Name: "organization_id", Value: organizationId},
		utils.IDField{Name: "project_id", Value: projectId},
		utils.IDField{Name: "cluster_id", Value: clusterId},
		utils.IDField{Name: "app_service_id", Value: appServiceId},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	orgUUID, projUUID, clusterUUID, appServiceUUID := uuids[0], uuids[1], uuids[2], uuids[3]

	response, err := d.ClientV2.GetAppServiceLogStreamingWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Log Streaming Configuration",
			fmt.Sprintf("Could not read log streaming configuration: %s: %s", errors.ErrExecutingRequest, err.Error()),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Reading App Service Log Streaming Configuration",
			fmt.Sprintf("Unexpected response while reading App Service Log Streaming Config: %s", string(response.Body)),
		)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Log Streaming Configuration",
			"API returned an empty response body.",
		)
		return
	}

	tflog.Info(ctx, "read app service log streaming configuration", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"app_service_id":  appServiceId,
	})

	// Map the API response to the datasource state.
	state := providerschema.NewAppServiceLogStreamingData(
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		response.JSON200,
	)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
