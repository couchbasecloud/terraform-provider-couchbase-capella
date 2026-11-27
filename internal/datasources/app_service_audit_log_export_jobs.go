package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceAuditLogExportJobs)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceAuditLogExportJobs)(nil)
)

// AppServiceAuditLogExportJobs is the App Service audit log export jobs data source implementation.
type AppServiceAuditLogExportJobs struct {
	*providerschema.Data
}

// NewAppServiceAuditLogExportJobs is a helper function to simplify the provider implementation.
func NewAppServiceAuditLogExportJobs() datasource.DataSource {
	return &AppServiceAuditLogExportJobs{}
}

// Metadata returns the data source type name.
func (d *AppServiceAuditLogExportJobs) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log_export_jobs"
}

// Schema defines the schema for the data source.
func (d *AppServiceAuditLogExportJobs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogExportJobsSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServiceAuditLogExportJobs) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the audit log export jobs of the App Service.
func (d *AppServiceAuditLogExportJobs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServiceAuditLogExportJobsData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogExports",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	jobs, err := api.GetPaginated[[]apigen.GetAuditExportDocResponse](ctx, d.ClientV1, d.Token, cfg, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Export Jobs",
			fmt.Sprintf("Could not read audit log export jobs for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	// The list endpoint has no export ID query parameter, so the provider filters the result.
	config.Data = make([]providerschema.AppServiceAuditLogExportJob, 0, len(jobs))
	for _, job := range jobs {
		if !config.ExportId.IsNull() && (job.Id == nil || *job.Id != config.ExportId.ValueString()) {
			continue
		}

		data, diags := providerschema.NewAppServiceAuditLogExportJob(ctx, job)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Data = append(config.Data, data)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
