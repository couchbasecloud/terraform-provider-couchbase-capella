package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	errorMessageWhileSnapshotBackupCreation = "There is an error during snapshot backup creation. Please check in Capella to see if any hanging resources" +
		" have been created, unexpected error: "
	errorMessageWhileSnapshotBackupUpdate = "There is an error during snapshot backup update. Please check in Capella to see if any hanging resources" +
		" have been created, unexpected error: "
	errorMessageConfigure = "Expected *ProviderSourceData, got:"
)

var (
	_ resource.Resource              = &SnapshotBackup{}
	_ resource.ResourceWithConfigure = &SnapshotBackup{}
)

type ID struct {
	ID string `json:"id"`
}

type IDList struct {
	Data []ID `json:"data"`
}

type SnapshotBackup struct {
	*providerschema.Data
}

func NewSnapshotBackup() resource.Resource {
	return &SnapshotBackup{}
}

func (s *SnapshotBackup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_backup"
}

func (s *SnapshotBackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotBackupSchema()
}

func (s *SnapshotBackup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.SnapshotBackup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := s.validateCreateSnapshotBackupRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create snapshot backup request",
			"Could not create snapshot backup "+err.Error(),
		)
		return
	}

	var (
		tenantId  = plan.TenantID.ValueString()
		projectId = plan.ProjectID.ValueString()
		clusterId = plan.ClusterID.ValueString()
	)

	createSnapshotBackupRequest := snapshot_backup.CreateSnapshotBackupRequest{
		Retention: int(plan.Retention.ValueInt64()),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, tenantId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	createResp, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		createSnapshotBackupRequest,
		s.Token,
		nil,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing create snapshot backup request",
			errorMessageWhileSnapshotBackupCreation+api.ParseError(err),
		)
		return
	}

	var createSnapshotBackupResponse snapshot_backup.CreateSnapshotBackupResponse
	err = json.Unmarshal(createResp.Body, &createSnapshotBackupResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling create snapshot backup response",
			"Could not unmarshal create snapshot backup response: "+err.Error(),
		)
	}

	// Checks the snapshot backup creation is complete.
	backupResp, err := s.checkSnapshotBackupStatus(ctx, tenantId, projectId, clusterId, createSnapshotBackupResponse.BackupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while checking latest snapshot backup status",
			errorMessageWhileSnapshotBackupCreation+api.ParseError(err),
		)
		return
	}

	refreshedState, err := createSnapshotBackup(ctx, backupResp, clusterId, projectId, tenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot backup",
			"Could not create snapshot backup: "+err.Error(),
		)
		return
	}

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *SnapshotBackup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.SnapshotBackup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup in Capella",
			"Could not read Capella Backup with ID "+state.BackupID.String()+": "+err.Error(),
		)
		return
	}
	var (
		tenantId  = IDs[providerschema.OrganizationId]
		projectId = IDs[providerschema.ProjectId]
		clusterId = IDs[providerschema.ClusterId]
		backupId  = state.BackupID.ValueString()
	)

	snapshotBackup, err := s.getSnapshotBackup(ctx, tenantId, projectId, clusterId, backupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot backup",
			"Could not read snapshot backup id "+backupId+": "+err.Error(),
		)
		return
	}

	refreshedState, err := createSnapshotBackup(ctx, snapshotBackup, clusterId, projectId, tenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating snapshot backup",
			"Could not update snapshot backup id "+backupId+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the retention of a snapshot backup.
func (s *SnapshotBackup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan providerschema.SnapshotBackup

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
			"Error Reading Snapshot Backup in Capella",
			"Could not read Capella Snapshot Backup with ID "+state.BackupID.String()+": "+err.Error(),
		)
		return
	}
	var (
		tenantId  = IDs[providerschema.OrganizationId]
		projectId = IDs[providerschema.ProjectId]
		clusterId = IDs[providerschema.ClusterId]
		backupId  = state.BackupID.ValueString()
	)

	updateSnapshotBackupRequest := snapshot_backup.EditBackupRetentionRequest{
		Retention: int(plan.Retention.ValueInt64()),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s", s.HostURL, tenantId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		updateSnapshotBackupRequest,
		s.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing update snapshot backup request retention",
			"Could not update snapshot backup id "+backupId+": "+api.ParseError(err),
		)
		return
	}

	backupResp, err := s.getSnapshotBackup(ctx, tenantId, projectId, clusterId, backupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot backup",
			"Could not get snapshot backup id "+backupId+": "+err.Error(),
		)
		return
	}

	state.Retention = types.Int64Value(int64(backupResp.Retention))
	state.Expiration = types.StringValue(backupResp.Expiration)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *SnapshotBackup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.SnapshotBackup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete backup id "+state.BackupID.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		tenantId  = IDs[providerschema.OrganizationId]
		projectId = IDs[providerschema.ProjectId]
		clusterId = IDs[providerschema.ClusterId]
		backupId  = state.BackupID.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s", s.HostURL, tenantId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err = s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
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
			"Could not delete backup id "+state.BackupID.String()+": "+errString,
		)
		return
	}
}

// Configure adds the provider configured api to the snapshot backup resource.
func (s *SnapshotBackup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			errorMessageConfigure+fmt.Sprintf("%T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	s.Data = data
}

func (s *SnapshotBackup) validateCreateSnapshotBackupRequest(plan providerschema.SnapshotBackup) error {
	if plan.TenantID.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectID.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterID.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}
	return nil
}

// checkLatestSnapshotBackupStatus monitors the status of a snapshot backup creation operation for a specified organization, project, and cluster ID.
// It periodically fetches the snapshot backup job status using the `getLatestSnapshotBackup` function and waits until the snapshot backup reaches a final state or until a specified timeout is reached.
// function and waits until the snapshot backup reaches a final state or until a specified timeout is reached.
func (s *SnapshotBackup) checkSnapshotBackupStatus(ctx context.Context, organizationId, projectId, clusterId, backupId string) (*snapshot_backup.SnapshotBackup, error) {
	var (
		backupResp *snapshot_backup.SnapshotBackup
		err        error
	)

	// Assuming 60 minutes is the max time backup completion takes.
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 1

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, errors.ErrSnapshotBackupCreationStatusTimeout

		case <-timer.C:
			backupResp, err = s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, backupId)
			switch err {
			case nil:
				if snapshot_backup.IsFinalState(backupResp.Progress.Status) {
					return backupResp, nil
				}
				const msg = "waiting for snapshot backup to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return nil, err
			}
			timer.Reset(sleep)
		}
	}
}

func (s *SnapshotBackup) getSnapshotBackups(ctx context.Context, tenantId, projectId, clusterId string) (*snapshot_backup.ListSnapshotBackupsResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, tenantId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := UnmarshalJSON(&snapshot_backup.ListSnapshotBackupsResponse{}, response.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// getSnapshotBackup retrieves a snapshot backup by its BackupID.
func (s *SnapshotBackup) getSnapshotBackup(ctx context.Context, tenantId, projectId, clusterId, backupId string) (*snapshot_backup.SnapshotBackup, error) {
	resp, err := s.getSnapshotBackups(ctx, tenantId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) > 0 {
		for _, backup := range resp.Data {
			if backup.BackupID == backupId {
				return &backup, nil
			}
		}
	}
	return nil, errors.ErrNotFound
}

// UnmarshalJSON unmarshals the JSON response into a ListSnapshotBackupsResponse so 'id' is assigned to the BackupID field.
func UnmarshalJSON(s *snapshot_backup.ListSnapshotBackupsResponse, data []byte) (*snapshot_backup.ListSnapshotBackupsResponse, error) {
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	idList := IDList{}
	err = json.Unmarshal(data, &idList)
	if err != nil {
		return nil, err
	}
	for i, id := range idList.Data {
		s.Data[i].BackupID = id.ID
	}
	return s, nil
}

// createSnapshotBackup creates a snapshot backup from a snapshot backup response.
func createSnapshotBackup(ctx context.Context, backupResp *snapshot_backup.SnapshotBackup, clusterId, projectId, tenantId string) (*providerschema.SnapshotBackup, error) {
	progress := providerschema.NewProgress(backupResp.Progress)
	progressObj, diags := types.ObjectValueFrom(ctx, progress.AttributeTypes(), progress)
	if diags.HasError() {
		return nil, fmt.Errorf("error during progress conversion")
	}

	server := providerschema.NewServer(backupResp.Server)
	serverObj, diags := types.ObjectValueFrom(ctx, server.AttributeTypes(), server)
	if diags.HasError() {
		return nil, fmt.Errorf("error during server conversion")
	}

	cmekObjs := make([]basetypes.ObjectValue, 0)
	for _, cmek := range backupResp.CMEK {
		cmek := providerschema.NewCMEK(cmek)
		cmekObj, diags := types.ObjectValueFrom(ctx, cmek.AttributeTypes(), cmek)
		if diags.HasError() {
			return nil, fmt.Errorf("error during CMEK conversion")
		}
		cmekObjs = append(cmekObjs, cmekObj)
	}

	cmekSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: providerschema.CMEK{}.AttributeTypes()}, cmekObjs)
	if diags.HasError() {
		return nil, fmt.Errorf("error during CMEK conversion")
	}

	snapshotBackup := providerschema.NewSnapshotBackup(ctx, *backupResp, backupResp.BackupID, clusterId, projectId, tenantId, progressObj, serverObj, cmekSet)
	return &snapshotBackup, nil
}
