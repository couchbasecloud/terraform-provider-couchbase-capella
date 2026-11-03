package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceAuditLogStreamingState)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceAuditLogStreamingState)(nil)
)

// AppServiceAuditLogStreamingState is the App Service audit log streaming state data source implementation.
type AppServiceAuditLogStreamingState struct {
	*providerschema.Data
}

// NewAppServiceAuditLogStreamingState is a helper function to simplify the provider implementation.
func NewAppServiceAuditLogStreamingState() datasource.DataSource {
	return &AppServiceAuditLogStreamingState{}
}

// Metadata returns the data source type name.
func (d *AppServiceAuditLogStreamingState) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log_streaming_state"
}

// Schema defines the schema for the data source.
func (d *AppServiceAuditLogStreamingState) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogStreamingStateSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServiceAuditLogStreamingState) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the audit log streaming state of the App Service.
func (d *AppServiceAuditLogStreamingState) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServiceAuditLogStreamingStateData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogStreaming",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Streaming",
			fmt.Sprintf("Could not read audit log streaming state for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var streamingResp apigen.GetAuditLogStreamingResponse
	if err := json.Unmarshal(response.Body, &streamingResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Streaming",
			"Could not unmarshal audit log streaming response: "+err.Error(),
		)
		return
	}

	config.StreamingEnabled = types.BoolPointerValue(streamingResp.StreamingEnabled)
	config.LogStreamingState = types.StringPointerValue(streamingResp.LogStreamingState)
	config.OutputType = types.StringPointerValue(streamingResp.OutputType)

	disabled := []string{}
	if streamingResp.DisabledAppEndpoints != nil {
		disabled = *streamingResp.DisabledAppEndpoints
	}
	config.DisabledAppEndpoints, diags = types.SetValueFrom(ctx, types.StringType, disabled)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
