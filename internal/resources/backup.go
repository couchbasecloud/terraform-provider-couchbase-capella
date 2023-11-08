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

	latestBackup, err := b.getLatestBackup(organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting latest bucket backup in a cluster",
			"Could not get the latest bucket backup : unexpected error "+err.Error(),
		)
		return
	}
	var backupFound bool
	if latestBackup != nil {
		backupFound = true
	}

	_, err = b.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups", b.HostURL, organizationId, projectId, clusterId, bucketId),
		http.MethodPost,
		BackupRequest,
		b.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not execute request, unexpected error: "+err.Error(),
		)
		return
	}

	BackupResponse, err := b.checkLatestBackupStatus(ctx, organizationId, projectId, clusterId, bucketId, backupFound, latestBackup)

	backupStats := providerschema.NewBackupStats(*BackupResponse.BackupStats)
	backupStatsObj, diags := types.ObjectValueFrom(ctx, backupStats.AttributeTypes(), backupStats)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Backup Stats",
			fmt.Sprintf("Could not read backup stats data in a backup record, unexpected error: %s", fmt.Errorf("error while backup stats conversion")),
		)
		return
	}

	scheduleInfo := providerschema.NewScheduleInfo(*BackupResponse.ScheduleInfo)
	scheduleInfoObj, diags := types.ObjectValueFrom(ctx, scheduleInfo.AttributeTypes(), scheduleInfo)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Error Reading Backup Schedule Info",
			fmt.Sprintf("Could not read backup schedule info in a backup record, unexpected error: %s", fmt.Errorf("error while backup schedule info conversion")),
		)
		return
	}

	refreshedState := providerschema.NewBackup(BackupResponse, organizationId, projectId, backupStatsObj, scheduleInfoObj)

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
	resourceNotFound, err := handleBackupError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup",
			"Could not read backup id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the Backup record.
func (b *Backup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me https://couchbasecloud.atlassian.net/browse/AV-66713
}

// Delete deletes the backup.
func (b *Backup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me https://couchbasecloud.atlassian.net/browse/AV-66712
}

// ImportState imports a remote backup that is not created by Terraform.
func (b *Backup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me https://couchbasecloud.atlassian.net/browse/AV-66714
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
			backupResp, err = b.getLatestBackup(organizationId, projectId, clusterId, bucketId)
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
	response, err := b.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s", b.HostURL, organizationId, projectId, clusterId, backupId),
		http.MethodGet,
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
func (b *Backup) getLatestBackup(organizationId, projectId, clusterId, bucketId string) (*backupapi.GetBackupResponse, error) {
	response, err := b.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups", b.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
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

// handleBackupError extracts error message if error is api.Error and
// also checks whether error is resource not found
func handleBackupError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}
