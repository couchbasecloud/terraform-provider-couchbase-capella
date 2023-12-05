package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"terraform-provider-capella/internal/api"
	backupapi "terraform-provider-capella/internal/api/backup"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Backup{}
	_ resource.ResourceWithConfigure   = &Backup{}
	_ resource.ResourceWithImportState = &Backup{}
)

// Backup is the Backup resource implementation.
type Backup struct {
	*providerschema.Data
}

// NewBackup is a helper function to simplify the provider implementation.
func NewBackup() resource.Resource {
	return &Backup{}
}

// Metadata returns the Backup resource type name.
func (b *Backup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

// Schema defines the schema for the Backup resource.
func (b *Backup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupSchema()
}

// Create creates a new Backup.
func (b *Backup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Backup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := b.validateCreateBackupRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create backup request",
			"Could not create backup "+err.Error(),
		)
		return
	}

	BackupRequest := backupapi.CreateBackupRequest{}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	latestBackup, err := b.getLatestBackup(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting latest bucket backup in a cluster",
			"Could not get the latest bucket backup : unexpected error "+api.ParseError(err),
		)
		return
	}

	var backupFound bool
	if latestBackup != nil {
		backupFound = true
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err = b.Client.ExecuteWithRetry(
		ctx,
		cfg,
		BackupRequest,
		b.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing create backup request",
			"Could not execute create backup request : unexpected error "+api.ParseError(err),
		)
		return
	}

	backupResponse, err := b.checkLatestBackupStatus(ctx, organizationId, projectId, clusterId, bucketId, backupFound, latestBackup)
	if err != nil {
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error whiling checking latest backup status",
				fmt.Sprintf("Could not read check latest backup status, unexpected error: "+api.ParseError(err)),
			)
			return
		}
	}

	backupStats := providerschema.NewBackupStats(*backupResponse.BackupStats)
	backupStatsObj, diags := types.ObjectValueFrom(ctx, backupStats.AttributeTypes(), backupStats)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Backup Stats",
			fmt.Sprintf("Could not read backup stats data in a backup record, unexpected error: %s", fmt.Errorf("error while backup stats conversion")),
		)
		return
	}

	scheduleInfo := providerschema.NewScheduleInfo(*backupResponse.ScheduleInfo)
	scheduleInfoObj, diags := types.ObjectValueFrom(ctx, scheduleInfo.AttributeTypes(), scheduleInfo)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Error Reading Backup Schedule Info",
			fmt.Sprintf("Could not read backup schedule info in a backup record, unexpected error: %s", fmt.Errorf("error while backup schedule info conversion")),
		)
		return
	}

	refreshedState := providerschema.NewBackup(backupResponse, organizationId, projectId, backupStatsObj, scheduleInfoObj)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read reads backup information.
func (b *Backup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Backup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup in Capella",
			"Could not read Capella Backup with ID "+state.Id.String()+": "+err.Error(),
		)
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.BucketId]
		backupId       = IDs[providerschema.Id]
	)

	refreshedState, err := b.retrieveBackup(ctx, organizationId, projectId, clusterId, bucketId, backupId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading backup",
			"Could not read backup id "+state.Id.String()+": "+errString,
		)
		return
	}

	refreshedState.Restore = state.Restore
	refreshedState.RestoreTimes = state.RestoreTimes

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the Backup record.
func (b *Backup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan providerschema.Backup

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup in Capella",
			"Could not read Capella Backup with ID "+state.Id.String()+": "+err.Error(),
		)
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		backupId       = IDs[providerschema.Id]
	)

	var restore *providerschema.Restore
	diags.Append(req.Config.GetAttribute(ctx, path.Root("restore"), &restore)...)

	if plan.RestoreTimes.IsNull() || plan.RestoreTimes.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error restoring backup",
			"Could not restore backup id "+state.Id.String()+": plan restore times value is not set",
		)
		return
	}

	if !state.RestoreTimes.IsNull() && !state.RestoreTimes.IsUnknown() {
		planRestoreTimes := *plan.RestoreTimes.ValueBigFloat()
		stateRestoreTimes := *state.RestoreTimes.ValueBigFloat()
		if planRestoreTimes.Cmp(&stateRestoreTimes) != 1 {
			resp.Diagnostics.AddError(
				"Error restoring backup",
				"Could not restore backup id "+state.Id.String()+": plan restore times value is not greater than state restore times value",
			)
			return
		}
	}

	var newServices []backupapi.Service
	for _, service := range restore.Services {
		newService := service.ValueString()
		newServices = append(newServices, backupapi.Service(newService))
	}

	restoreRequest := backupapi.CreateRestoreRequest{
		TargetClusterId:       restore.TargetClusterId.ValueString(),
		SourceClusterId:       restore.SourceClusterId.ValueString(),
		BackupId:              backupId,
		Services:              &newServices,
		ForceUpdates:          restore.ForceUpdates.ValueBool(),
		AutoRemoveCollections: restore.AutoRemoveCollections.ValueBool(),
		FilterKeys:            restore.FilterKeys.ValueString(),
		FilterValues:          restore.FilterValues.ValueString(),
		IncludeData:           restore.IncludeData.ValueString(),
		ExcludeData:           restore.ExcludeData.ValueString(),
		MapData:               restore.MapData.ValueString(),
		ReplaceTTL:            restore.ReplaceTTL.ValueString(),
		ReplaceTTLWith:        restore.ReplaceTTLWith.ValueString(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s/restore", b.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err = b.Client.ExecuteWithRetry(
		ctx,
		cfg,
		restoreRequest,
		b.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error restoring backup",
			"Could not restore backup id "+state.Id.String()+": plan restore times value is not greater than state restore times value",
		)
		return
	}

	if !plan.Restore.IsUnknown() && !plan.Restore.IsNull() {
		restore.Status = types.StringValue("RESTORE INITIATED")
		restoreObj, diags := types.ObjectValueFrom(ctx, restore.AttributeTypes(), restore)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.Restore = restoreObj
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the backup.
func (b *Backup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.Backup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete backup id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		backupId       = resourceIDs[providerschema.Id]
	)

	// Delete existing Backup
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s", b.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = b.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		b.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete backup id "+state.Id.String()+": "+errString,
		)
		return
	}
}

// ImportState imports a remote backup that is not created by Terraform.
func (b *Backup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured api to the backup resource.
func (b *Backup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (a *Backup) validateCreateBackupRequest(plan providerschema.Backup) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	if plan.BucketId.IsNull() {
		return errors.ErrBucketIdCannotBeEmpty
	}
	if !plan.RestoreTimes.IsNull() && !plan.RestoreTimes.IsUnknown() {
		return errors.ErrRestoreTimesMustNotBeSetWhileCreateBackup
	}
	return nil
}

// checkLatestBackupStatus monitors the status of a backup creation operation for a specified
// organization, project, and cluster ID. It periodically fetches the backup job status using the `getLatestBackup`
// function and waits until the backup reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (b *Backup) checkLatestBackupStatus(ctx context.Context, organizationId, projectId, clusterId, bucketId string, backupFound bool, latestBackup *backupapi.GetBackupResponse) (*backupapi.GetBackupResponse, error) {
	var (
		backupResp *backupapi.GetBackupResponse
		err        error
	)

	// Assuming 60 minutes is the max time backup completion takes, can change after discussion
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 1

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			const msg = "bucket backup creation status transition timed out after initiation"
			return nil, fmt.Errorf(msg)

		case <-timer.C:
			backupResp, err = b.getLatestBackup(ctx, organizationId, projectId, clusterId, bucketId)
			switch err {
			case nil:
				// If there is no existing backup for a bucket, check for a new backup record to be created.
				// If a backup record exists already, wait for a backup record with a new ID to be created.
				if !backupFound && backupResp != nil && backupapi.IsFinalState(backupResp.Status) {
					return backupResp, nil
				} else if backupFound && backupResp != nil && latestBackup.Id != backupResp.Id && backupapi.IsFinalState(backupResp.Status) {
					return backupResp, nil
				}
				const msg = "waiting for backup to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return nil, err
			}
			timer.Reset(sleep)
		}
	}
}

// retrieveBackup retrieves backup information from the specified organization and project
// using the provided backup ID by open-api call
func (b *Backup) retrieveBackup(ctx context.Context, organizationId, projectId, clusterId, bucketId, backupId string) (*providerschema.Backup, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s", b.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		b.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	backupResp := backupapi.GetBackupResponse{}
	err = json.Unmarshal(response.Body, &backupResp)
	if err != nil {
		return nil, err
	}

	bStats := providerschema.NewBackupStats(*backupResp.BackupStats)
	bStatsObj, diags := types.ObjectValueFrom(ctx, bStats.AttributeTypes(), bStats)
	if diags.HasError() {
		return nil, errors.ErrUnableToConvertAuditData
	}

	sInfo := providerschema.NewScheduleInfo(*backupResp.ScheduleInfo)
	sInfoObj, diags := types.ObjectValueFrom(ctx, sInfo.AttributeTypes(), sInfo)
	if diags.HasError() {
		return nil, errors.ErrUnableToConvertAuditData
	}

	refreshedState := providerschema.NewBackup(&backupResp, organizationId, projectId, bStatsObj, sInfoObj)
	return refreshedState, nil
}

// getLatestBackup retrieves the latest backup information for a specified bucket in a cluster
// from the specified organization, project and cluster using the provided bucket ID by open-api call
func (b *Backup) getLatestBackup(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*backupapi.GetBackupResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups", b.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		b.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	clusterResp := backupapi.GetBackupsResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}

	// check for a backup record of the specified bucket
	for _, backup := range clusterResp.Data {
		if backup.BucketId == bucketId {
			return &backup, nil
		}
	}
	return nil, nil
}
