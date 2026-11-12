package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceAuditLogExportJob)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceAuditLogExportJob)(nil)
)

// AppServiceAuditLogExportJob is the App Service audit log export job data source implementation.
type AppServiceAuditLogExportJob struct {
	*providerschema.Data
}

// NewAppServiceAuditLogExportJob is a helper function to simplify the provider implementation.
func NewAppServiceAuditLogExportJob() datasource.DataSource {
	return &AppServiceAuditLogExportJob{}
}

// Metadata returns the data source type name.
func (d *AppServiceAuditLogExportJob) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log_export_job"
}

// Schema defines the schema for the data source.
func (d *AppServiceAuditLogExportJob) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogExportJobSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServiceAuditLogExportJob) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with an audit log export job of the App Service.
func (d *AppServiceAuditLogExportJob) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServiceAuditLogExportJobData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogExports/%s",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
		config.Id.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Export Job",
			fmt.Sprintf("Could not read audit log export job %s for App Service %s: %s", config.Id.ValueString(), config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var job apigen.GetAuditExportDocResponse
	if err := json.Unmarshal(response.Body, &job); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Export Job",
			"Could not unmarshal audit log export job response: "+err.Error(),
		)
		return
	}

	config.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	if job.Audit != nil {
		audit := providerschema.NewCouchbaseAuditData(api.CouchbaseAuditData(*job.Audit))
		var diags diag.Diagnostics
		config.Audit, diags = types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	config.DownloadExpires = types.StringNull()
	if job.DownloadExpires != nil {
		config.DownloadExpires = types.StringValue(job.DownloadExpires.Format(time.RFC3339))
	}

	config.TenantId = types.StringPointerValue(job.TenantId)
	config.Status = types.StringPointerValue(job.Status)
	config.DownloadId = types.StringPointerValue(job.DownloadId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
