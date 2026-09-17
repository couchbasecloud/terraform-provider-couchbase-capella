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
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketBackupExport{}
	_ resource.ResourceWithConfigure   = &BucketBackupExport{}
	_ resource.ResourceWithImportState = &BucketBackupExport{}
)

const errorMessageAfterBucketBackupExportCreation = "Backup export job creation is successful, but encountered an error while checking the current" +
	" state of the backup export job. Please run `terraform plan` after 1-2 minutes to know the" +
	" current backup export job state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileBucketBackupExportCreation = "There is an error during backup export creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// BucketBackupExport is the resource implementation.
type BucketBackupExport struct {
	*providerschema.Data
}

func NewBucketBackupExport() resource.Resource {
	return &BucketBackupExport{}
}

// Metadata returns the bucket backup export resource type name.
func (b *BucketBackupExport) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_backup_export"
}

// Schema defines the schema for the bucket backup export resource.
func (b *BucketBackupExport) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BucketBackupExportSchema()
}

// Configure set provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (b *BucketBackupExport) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	b.Data = data
}

// Create queues a new backup export job.
//
// The export runs asynchronously on the backup infrastructure and a multi-terabyte backup can take
// a long time, so this does not wait for it to finish. It refreshes once to populate the computed
// attributes and leaves the export pending; the user polls by refreshing the resource.
func (b *BucketBackupExport) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.BucketBackupExport
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := b.validate(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup export job",
			"Could not create backup export job, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		bucketId       = plan.BucketId.ValueString()
		backupId       = plan.BackupId.ValueString()
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups/%s/export",
		b.HostURL,
		organizationId,
		projectId,
		clusterId,
		bucketId,
		backupId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	response, err := b.ClientV1.ExecuteWithRetry(ctx, cfg, nil, b.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup export job",
			errorMessageWhileBucketBackupExportCreation+api.ParseError(err),
		)
		return
	}

	createResponse := backupapi.CreateBucketBackupExportResponse{}
	if err = json.Unmarshal(response.Body, &createResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup export job",
			errorMessageWhileBucketBackupExportCreation+"An error occurred during unmarshalling: "+err.Error(),
		)
		return
	}

	// Record the export ID before refreshing. A refresh failure is only a warning, and without this
	// the ID would be lost and the export orphaned.
	plan.Id = types.StringValue(createResponse.ExportId)
	plan.CreatedAt = types.StringValue(createResponse.CreatedAt.Format(time.RFC3339))
	diags = resp.State.Set(ctx, initialBucketBackupExportState(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := b.refreshBucketBackupExport(ctx, organizationId, projectId, clusterId, bucketId, backupId, createResponse.ExportId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading backup export",
			errorMessageAfterBucketBackupExportCreation+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read gets backup export information.
func (b *BucketBackupExport) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.BucketBackupExport
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket Backup Export",
			"Could not read Capella bucket backup export: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.BucketId]
		backupId       = IDs[providerschema.BackupId]
		exportId       = IDs[providerschema.Id]
	)

	refreshedState, err := b.refreshBucketBackupExport(ctx, organizationId, projectId, clusterId, bucketId, backupId, exportId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket Backup Export",
			"Could not read Capella export id "+exportId+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update is not supported as the backup export API does not have an update endpoint.
func (b *BucketBackupExport) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Bucket Backup Export does not support update",
		"Bucket Backup Export does not support update",
	)
}

// Delete is intentionally a no-op. There is no endpoint to cancel or delete an export; the archive
// expires from cloud storage on its own and Capella drops the export record after about 7 days.
//
// The framework will automatically update the state file. See:
// https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#recommendations
func (b *BucketBackupExport) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (b *BucketBackupExport) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (b *BucketBackupExport) validate(plan providerschema.BucketBackupExport) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.BucketId.IsNull() {
		return errors.ErrBucketIdMissing
	}
	if plan.BackupId.IsNull() {
		return errors.ErrBackupIdMissing
	}
	return nil
}

func (b *BucketBackupExport) getBucketBackupExport(ctx context.Context, organizationId, projectId, clusterId, bucketId, backupId, exportId string) (*backupapi.GetBucketBackupExportResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups/%s/exports/%s",
		b.HostURL,
		organizationId,
		projectId,
		clusterId,
		bucketId,
		backupId,
		exportId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.ClientV1.ExecuteWithRetry(ctx, cfg, nil, b.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	exportResponse := backupapi.GetBucketBackupExportResponse{}
	if err = json.Unmarshal(response.Body, &exportResponse); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &exportResponse, nil
}

func (b *BucketBackupExport) refreshBucketBackupExport(ctx context.Context, organizationId, projectId, clusterId, bucketId, backupId, exportId string) (*providerschema.BucketBackupExport, error) {
	exportResp, err := b.getBucketBackupExport(ctx, organizationId, projectId, clusterId, bucketId, backupId, exportId)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.BucketBackupExport{
		Id:             types.StringValue(exportResp.Id),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		BucketId:       types.StringValue(bucketId),
		BackupId:       types.StringValue(backupId),
		CycleId:        types.StringValue(exportResp.CycleId),
		BucketName:     types.StringValue(exportResp.BucketName),
		Status:         types.StringValue(exportResp.Status),
		CreatedAt:      types.StringValue(exportResp.CreatedAt.Format(time.RFC3339)),
	}

	if exportResp.SizeInBytes != nil {
		refreshedState.SizeInBytes = types.Int64Value(*exportResp.SizeInBytes)
	}
	if exportResp.Sha256Checksum != nil {
		refreshedState.Sha256Checksum = types.StringValue(*exportResp.Sha256Checksum)
	}
	if exportResp.Expiration != nil {
		refreshedState.Expiration = types.StringValue(exportResp.Expiration.Format(time.RFC3339))
	}
	if exportResp.BackupDownloadURL != nil {
		refreshedState.BackupDownloadURL = types.StringValue(*exportResp.BackupDownloadURL)
	}

	return &refreshedState, nil
}

// initialBucketBackupExportState nulls the attributes only the get endpoint can populate, so the
// state is valid if the refresh that follows the create fails.
func initialBucketBackupExportState(plan providerschema.BucketBackupExport) providerschema.BucketBackupExport {
	plan.CycleId = types.StringNull()
	plan.BucketName = types.StringNull()
	plan.Status = types.StringNull()
	plan.SizeInBytes = types.Int64Null()
	plan.Sha256Checksum = types.StringNull()
	plan.Expiration = types.StringNull()
	plan.BackupDownloadURL = types.StringNull()

	return plan
}
