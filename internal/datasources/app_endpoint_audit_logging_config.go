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
	_ datasource.DataSource              = (*AppEndpointAuditLoggingConfig)(nil)
	_ datasource.DataSourceWithConfigure = (*AppEndpointAuditLoggingConfig)(nil)
)

// AppEndpointAuditLoggingConfig is the App Endpoint audit logging config data source implementation.
type AppEndpointAuditLoggingConfig struct {
	*providerschema.Data
}

// NewAppEndpointAuditLoggingConfig is a helper function to simplify the provider implementation.
func NewAppEndpointAuditLoggingConfig() datasource.DataSource {
	return &AppEndpointAuditLoggingConfig{}
}

// Metadata returns the data source type name.
func (d *AppEndpointAuditLoggingConfig) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_audit_logging_config"
}

// Schema defines the schema for the data source.
func (d *AppEndpointAuditLoggingConfig) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointAuditLoggingConfigSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppEndpointAuditLoggingConfig) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the audit logging config of the App Endpoint.
func (d *AppEndpointAuditLoggingConfig) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpointAuditLoggingConfigData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/auditLog",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
		config.AppEndpointName.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Audit Log Config",
			fmt.Sprintf("Could not read audit log config for App Endpoint %s: %s", config.AppEndpointName.ValueString(), api.ParseError(err)),
		)
		return
	}

	var auditLogResp apigen.GetAppEndpointAuditLogResponse
	if err := json.Unmarshal(response.Body, &auditLogResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Audit Log Config",
			"Could not unmarshal audit log config response: "+err.Error(),
		)
		return
	}

	config.AuditEnabled = types.BoolPointerValue(auditLogResp.AuditEnabled)

	config.EnabledEventIds = []types.Int64{}
	if auditLogResp.EnabledEventIds != nil {
		for _, e := range *auditLogResp.EnabledEventIds {
			if e.Id != nil {
				config.EnabledEventIds = append(config.EnabledEventIds, types.Int64Value(int64(*e.Id)))
			}
		}
	}

	config.DisabledUsers = toAppEndpointAuditDisabledRoles(auditLogResp.DisabledUsers)
	config.DisabledRoles = toAppEndpointAuditDisabledRoles(auditLogResp.DisabledRoles)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func toAppEndpointAuditDisabledRoles(in *apigen.DisabledUserRoles) []providerschema.AppEndpointAuditDisabledRole {
	out := []providerschema.AppEndpointAuditDisabledRole{}
	if in == nil {
		return out
	}
	for _, r := range *in {
		out = append(out, providerschema.AppEndpointAuditDisabledRole{
			Domain: types.StringValue(r.Domain),
			Name:   types.StringValue(r.Name),
		})
	}
	return out
}
