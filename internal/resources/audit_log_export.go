package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AuditLogExport{}
	_ resource.ResourceWithConfigure   = &AuditLogExport{}
	_ resource.ResourceWithImportState = &AuditLogExport{}
)

const errorMessageAfterAuditLogExportCreation = "Audit log export job creating is successful, but encountered an error while checking the current" +
	" state of the audit log export job. Please run `terraform plan` after 1-2 minutes to know the" +
	" current audit log export job state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileAuditLogExportCreation = "There is an error during audit log export creating. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// AuditLogExport is the resource implementation.
type AuditLogExport struct {
	*providerschema.Data
}

func NewAuditLogExport() resource.Resource {
	return &AuditLogExport{}
}

// Metadata returns the audit log export resource type name.
func (a *AuditLogExport) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_log_export"
}

// Schema defines the schema for the audit log export resource.
func (a *AuditLogExport) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AuditLogExportSchema()
}

// Configure set provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (a *AuditLogExport) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new audit log export job.
func (a *AuditLogExport) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AuditLogExport
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := a.validate(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log export job",
			"Could not create audit log export job, unexpected error: "+err.Error(),
		)
		return
	}

	start, err := time.Parse(time.RFC3339, plan.Start.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log export job",
			"Could not parse start time, unexpected error: "+err.Error(),
		)
		return
	}

	end, err := time.Parse(time.RFC3339, plan.End.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log export job",
			"Could not parse end time, unexpected error: "+err.Error(),
		)
		return
	}

	auditLogExportRequest := api.CreateClusterAuditLogExportRequest{
		Start: start,
		End:   end,
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/auditLogExports",
		a.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		auditLogExportRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log export job",
			errorMessageWhileAuditLogExportCreation+api.ParseError(err),
		)
		return
	}

	auditLogExportResponse := api.CreateClusterAuditLogExportResponse{}
	err = json.Unmarshal(response.Body, &auditLogExportResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log export job",
			errorMessageWhileAuditLogExportCreation+"An error occurred during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initialState(plan, auditLogExportResponse.ExportId))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.refreshAuditLogExport(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), auditLogExportResponse.ExportId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading audit log export",
			errorMessageAfterAuditLogExportCreation+api.ParseError(err),
		)
		return
	}

	// API server returns time using offset.  If user specifies UTC in zulu format,
	// this will cause state mismatch.  We overwrite API response with what's in the plan.
	refreshedState.Start = types.StringValue(strings.Trim(plan.Start.String(), "\""))
	refreshedState.End = types.StringValue(strings.Trim(plan.End.String(), "\""))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read gets audit log export information.
func (a *AuditLogExport) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AuditLogExport
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Export",
			"Could not read Capella audit log export: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		exportId       = IDs[providerschema.Id]
	)

	// refresh the existing allow list
	refreshedState, err := a.refreshAuditLogExport(ctx, organizationId, projectId, clusterId, exportId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Export",
			"Could not read Capella export id "+exportId+": "+errString,
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update is not supported as audit log export API does not have update endpoint.
func (a *AuditLogExport) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Audit Log Export does not support update",
		"Audit Log Export does not support update",
	)
}

// Delete is not supported as audit log export API does not have delete endpoint.
func (a *AuditLogExport) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Audit Log Export does not support delete",
		"Audit Log Export does not support delete",
	)
}

func (a *AuditLogExport) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *AuditLogExport) validate(plan providerschema.AuditLogExport) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}

func (a *AuditLogExport) getAuditLogExport(ctx context.Context, organizationId, projectId, clusterId, exportId string) (*api.GetClusterAuditLogExportResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/auditLogExports/%s",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		exportId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	auditExportResponse := api.GetClusterAuditLogExportResponse{}
	err = json.Unmarshal(response.Body, &auditExportResponse)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &auditExportResponse, nil
}

func (a *AuditLogExport) refreshAuditLogExport(ctx context.Context, organizationId, projectId, clusterId, exportId string) (*providerschema.AuditLogExport, error) {
	auditLogExportResp, err := a.getAuditLogExport(ctx, organizationId, projectId, clusterId, exportId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	refreshedState := providerschema.AuditLogExport{
		Id:             types.StringValue(exportId),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Start:          types.StringValue(auditLogExportResp.Start.String()),
		End:            types.StringValue(auditLogExportResp.End.String()),
		CreatedAt:      types.StringValue(auditLogExportResp.CreatedAt.String()),
	}

	if auditLogExportResp.AuditLogDownloadURL != nil {
		refreshedState.AuditLogDownloadURL = types.StringValue(*auditLogExportResp.AuditLogDownloadURL)
	}
	if auditLogExportResp.Expiration != nil {
		refreshedState.Expiration = types.StringValue(auditLogExportResp.Expiration.String())
	}

	return &refreshedState, nil
}

// initialState initializes an instance of providerschema.AuditLogExport
// with the specified plan and ID. It marks all computed fields as null.
func initialState(plan providerschema.AuditLogExport, exportId string) providerschema.AuditLogExport {
	plan.Id = types.StringValue(exportId)

	if plan.AuditLogDownloadURL.IsNull() || plan.AuditLogDownloadURL.IsUnknown() {
		plan.AuditLogDownloadURL = types.StringNull()
	}
	if plan.Expiration.IsNull() || plan.Expiration.IsUnknown() {
		plan.Expiration = types.StringNull()
	}
	if plan.Start.IsNull() || plan.Start.IsUnknown() {
		plan.Start = types.StringNull()
	}
	if plan.End.IsNull() || plan.End.IsUnknown() {
		plan.End = types.StringNull()
	}
	if plan.CreatedAt.IsNull() || plan.CreatedAt.IsUnknown() {
		plan.CreatedAt = types.StringNull()
	}
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.StringNull()
	}

	return plan
}
