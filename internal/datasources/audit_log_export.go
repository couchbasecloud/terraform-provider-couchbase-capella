package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AuditLogExport{}
	_ datasource.DataSourceWithConfigure = &AuditLogExport{}
)

// AuditLogExport is the data source implementation.
type AuditLogExport struct {
	*providerschema.Data
}

// NewAuditLogExport is a helper function to simplify the provider implementation.
func NewAuditLogExport() datasource.DataSource {
	return &AuditLogExport{}
}

// Metadata returns the audit log export data source type name.
func (a *AuditLogExport) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_log_export"
}

// Schema defines the schema for the auditlogexport data source.
func (a *AuditLogExport) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                    computedStringAttribute,
						"organization_id":       computedStringAttribute,
						"project_id":            computedStringAttribute,
						"cluster_id":            computedStringAttribute,
						"auditlog_download_url": computedStringAttribute,
						"expiration":            computedStringAttribute,
						"start":                 computedStringAttribute,
						"end":                   computedStringAttribute,
						"created_at":            computedStringAttribute,
						"status":                computedStringAttribute,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the auditlogexport data source.
func (a *AuditLogExport) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	a.Data = data
}

// Read refreshes the Terraform state with the latest data of auditlogexport.
func (a *AuditLogExport) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AuditLogExports
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	err := a.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Exports",
			"Could not read audit log exports in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	auditLogExports, err := a.listAuditLogExports(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Exports",
			"Could not read audit log exports in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	state = a.mapResponseBody(auditLogExports, &state)

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// listAuditLogExports executes calls to the list audit log export endpoint. It handles pagination and
// returns a slice of individual audit log export responses retrieved from multiple pages.
func (a *AuditLogExport) listAuditLogExports(ctx context.Context, organizationId, projectId, clusterId string) ([]api.GetClusterAuditLogExportResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/auditLogExports",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	return api.GetPaginated[[]api.GetClusterAuditLogExportResponse](ctx, a.Client, a.Token, cfg, "")
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (a *AuditLogExport) validate(state providerschema.AuditLogExports) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}

func (a *AuditLogExport) mapResponseBody(
	auditLogExports []api.GetClusterAuditLogExportResponse,
	state *providerschema.AuditLogExports,
) providerschema.AuditLogExports {
	for _, export := range auditLogExports {
		exportState := providerschema.AuditLogExport{
			Id:             types.StringValue(export.AuditLogExportId),
			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
			ProjectId:      types.StringValue(state.ProjectId.ValueString()),
			ClusterId:      types.StringValue(state.ClusterId.ValueString()),
			Start:          types.StringValue(export.Start.String()),
			End:            types.StringValue(export.End.String()),
			CreatedAt:      types.StringValue(export.CreatedAt.String()),
			Status:         types.StringValue(export.Status),
		}

		if export.AuditLogDownloadURL != nil {
			exportState.AuditLogDownloadURL = types.StringValue(*export.AuditLogDownloadURL)
			exportState.Expiration = types.StringValue(export.Expiration.String())
		}

		state.Data = append(state.Data, exportState)
	}

	return *state
}
