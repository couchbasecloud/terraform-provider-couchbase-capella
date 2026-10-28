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
	_ datasource.DataSource              = (*AppServiceAuditLog)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceAuditLog)(nil)
)

// AppServiceAuditLog is the App Service audit log data source implementation.
type AppServiceAuditLog struct {
	*providerschema.Data
}

// NewAppServiceAuditLog is a helper function to simplify the provider implementation.
func NewAppServiceAuditLog() datasource.DataSource {
	return &AppServiceAuditLog{}
}

// Metadata returns the data source type name.
func (d *AppServiceAuditLog) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log"
}

// Schema defines the schema for the data source.
func (d *AppServiceAuditLog) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServiceAuditLog) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the audit logging state of the App Service.
func (d *AppServiceAuditLog) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServiceAuditLogData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLog",
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
			"Error Reading App Service Audit Log",
			fmt.Sprintf("Could not read audit log state for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var auditLogResp apigen.CreateAppServiceAuditLogRequest
	if err := json.Unmarshal(response.Body, &auditLogResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log",
			"Could not unmarshal audit log state response: "+err.Error(),
		)
		return
	}

	config.AuditEnabled = types.BoolValue(auditLogResp.AuditEnabled)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
