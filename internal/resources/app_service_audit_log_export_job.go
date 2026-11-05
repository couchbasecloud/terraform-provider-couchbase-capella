package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                   = (*AppServiceAuditLogExportJob)(nil)
	_ resource.ResourceWithConfigure      = (*AppServiceAuditLogExportJob)(nil)
	_ resource.ResourceWithValidateConfig = (*AppServiceAuditLogExportJob)(nil)
)

// AppServiceAuditLogExportJob is the App Service audit log export job resource implementation.
// It does not support import: the API does not return start and end, so an imported job
// would have them null in state, and the next plan would replace it with a new job.
type AppServiceAuditLogExportJob struct {
	*providerschema.Data
}

// NewAppServiceAuditLogExportJob is a helper function to simplify the provider implementation.
func NewAppServiceAuditLogExportJob() resource.Resource {
	return &AppServiceAuditLogExportJob{}
}

// Metadata returns the resource type name.
func (a *AppServiceAuditLogExportJob) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log_export_job"
}

// Schema defines the schema for the resource.
func (a *AppServiceAuditLogExportJob) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogExportJobSchema()
}

// ValidateConfig rejects an export window where end is before start.
func (a *AppServiceAuditLogExportJob) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config providerschema.AppServiceAuditLogExportJobResource
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	start, err := time.Parse(time.RFC3339, config.Start.ValueString())
	if err != nil {
		return
	}
	end, err := time.Parse(time.RFC3339, config.End.ValueString())
	if err != nil {
		return
	}

	if end.Before(start) {
		resp.Diagnostics.AddAttributeError(
			path.Root("end"),
			"Invalid Audit Log Export Window",
			"end must not be earlier than start",
		)
	}
}

// Configure adds the provider configured client to the resource.
func (a *AppServiceAuditLogExportJob) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	a.Data = data
}

// Create initiates an App Service audit log export job.
func (a *AppServiceAuditLogExportJob) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceAuditLogExportJobResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validators skip values that are unknown at plan time, so check again here.
	start, err := time.Parse(time.RFC3339, plan.Start.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service audit log export job",
			"Could not parse start time: "+err.Error(),
		)
		return
	}
	end, err := time.Parse(time.RFC3339, plan.End.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service audit log export job",
			"Could not parse end time: "+err.Error(),
		)
		return
	}
	if end.Before(start) {
		resp.Diagnostics.AddError(
			"Error creating App Service audit log export job",
			"end must not be earlier than start",
		)
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogExports",
		a.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, apigen.CreateClusterAuditLogExportRequest{Start: start, End: end}, a.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service audit log export job",
			"Could not create App Service audit log export job: "+api.ParseError(err),
		)
		return
	}

	var createResp apigen.CreateClusterAuditLogExportResponse
	if err := json.Unmarshal(response.Body, &createResp); err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service audit log export job",
			"Could not unmarshal create response: "+err.Error(),
		)
		return
	}

	// Save the ID first, so the job is tracked even if the read below fails.
	plan.Id = types.StringValue(createResp.ExportId)
	plan.Status = types.StringNull()
	plan.DownloadId = types.StringNull()
	plan.DownloadExpires = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.refresh(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading App Service audit log export job",
			"The export job "+createResp.ExportId+" was created, but could not be read. Run `terraform apply --refresh-only` to update the state: "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read refreshes the Terraform state with the latest export job data.
func (a *AppServiceAuditLogExportJob) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceAuditLogExportJobResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service audit log export job",
			"Could not validate App Service audit log export job: "+err.Error(),
		)
		return
	}
	state.OrganizationId = types.StringValue(IDs[providerschema.OrganizationId])
	state.ProjectId = types.StringValue(IDs[providerschema.ProjectId])
	state.ClusterId = types.StringValue(IDs[providerschema.ClusterId])
	state.AppServiceId = types.StringValue(IDs[providerschema.AppServiceId])
	state.Id = types.StringValue(IDs[providerschema.Id])

	refreshedState, err := a.refresh(ctx, state)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service audit log export job",
			"Could not read App Service audit log export job "+state.Id.ValueString()+": "+errString,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update is not supported. Every input attribute requires replacement, so the framework never calls it.
func (a *AppServiceAuditLogExportJob) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"App Service audit log export job does not support update",
		"App Service audit log export job does not support update",
	)
}

// Delete only removes the export job from the state, because the API has no delete endpoint.
func (a *AppServiceAuditLogExportJob) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// refresh gets the export job and copies the computed fields into state.
// The API does not return start and end, so they keep their state values.
func (a *AppServiceAuditLogExportJob) refresh(ctx context.Context, state providerschema.AppServiceAuditLogExportJobResource) (*providerschema.AppServiceAuditLogExportJobResource, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogExports/%s",
		a.HostURL,
		state.OrganizationId.ValueString(),
		state.ProjectId.ValueString(),
		state.ClusterId.ValueString(),
		state.AppServiceId.ValueString(),
		state.Id.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, err
	}

	var exportResp apigen.GetAuditExportDocResponse
	if err := json.Unmarshal(response.Body, &exportResp); err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	state.Status = types.StringPointerValue(exportResp.Status)
	state.DownloadId = types.StringPointerValue(exportResp.DownloadId)
	state.DownloadExpires = types.StringNull()
	if exportResp.DownloadExpires != nil {
		state.DownloadExpires = types.StringValue(exportResp.DownloadExpires.Format(time.RFC3339))
	}

	state.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	if exportResp.Audit != nil {
		audit := providerschema.NewCouchbaseAuditData(api.CouchbaseAuditData(*exportResp.Audit))
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			return nil, fmt.Errorf("%s: %v", errors.ErrUnableToConvertAuditData, diags.Errors())
		}
		state.Audit = auditObj
	}

	return &state, nil
}
