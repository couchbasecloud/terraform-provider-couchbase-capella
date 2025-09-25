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
	_ resource.Resource                = &SnapshotBackup{}
	_ resource.ResourceWithConfigure   = &SnapshotBackup{}
	_ resource.ResourceWithImportState = &SnapshotBackup{}
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

// ImportState imports a remote backup that is not created by Terraform.
func (s *SnapshotBackup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (s *SnapshotBackup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.SnapshotBackup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectID.ValueString()
		clusterId      = plan.ClusterID.ValueString()
	)

	createSnapshotBackupRequest := snapshot_backup.CreateSnapshotBackupRequest{
		Retention:     int(plan.Retention.ValueInt64()),
		RegionsToCopy: providerschema.ConvertStringValueList(plan.RegionsToCopy),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, organizationId, projectId, clusterId)
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
			"Error executing create snapshot backup",
			errorMessageWhileSnapshotBackupCreation+api.ParseError(err),
		)
		tflog.Debug(ctx, "error executing create snapshot backup", map[string]interface{}{
			"organizationId":              organizationId,
			"projectId":                   projectId,
			"clusterId":                   clusterId,
			"createSnapshotBackupRequest": createSnapshotBackupRequest,
			"createResp":                  createResp,
			"err":                         api.ParseError(err),
		})
		return
	}

	var createSnapshotBackupResponse snapshot_backup.CreateSnapshotBackupResponse
	err = json.Unmarshal(createResp.Body, &createSnapshotBackupResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling create snapshot backup response",
			"Could not unmarshal create snapshot backup response: "+err.Error(),
		)
		tflog.Debug(ctx, "error unmarshalling create snapshot backup response", map[string]interface{}{
			"organizationId":              organizationId,
			"projectId":                   projectId,
			"clusterId":                   clusterId,
			"createSnapshotBackupRequest": createSnapshotBackupRequest,
			"createResp":                  createResp,
			"err":                         api.ParseError(err),
		})
		return
	}

	// Checks the snapshot backup creation is complete.
	backupResp, err := s.checkSnapshotBackupStatus(ctx, organizationId, projectId, clusterId, createSnapshotBackupResponse.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while checking latest snapshot backup status",
			errorMessageWhileSnapshotBackupCreation+api.ParseError(err),
		)
		tflog.Debug(ctx, "error checking latest snapshot backup status", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"Id":             createSnapshotBackupResponse.ID,
			"err":            err,
		})
		return
	}

	refreshedState, err := createSnapshotBackup(ctx, backupResp, clusterId, projectId, organizationId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot backup",
			"Could not create snapshot backup: "+err.Error(),
		)
		tflog.Debug(ctx, "error creating snapshot backup", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"refreshedState": refreshedState,
			"err":            err,
		})
		return
	}

	refreshedState.RegionsToCopy = plan.RegionsToCopy

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
			"Could not read Capella Backup with ID "+state.ID.String()+": "+err.Error(),
		)
		tflog.Debug(ctx, "error validating snapshot backup IDs", map[string]interface{}{
			"state": state,
			"err":   err,
		})
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		Id             = IDs[providerschema.Id]
	)

	snapshotBackup, err := s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot backup",
			"Could not read snapshot backup id "+Id+": "+err.Error(),
		)
		tflog.Debug(ctx, "error reading snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"err":            err,
		})
		return
	}

	refreshedState, err := createSnapshotBackup(ctx, snapshotBackup, clusterId, projectId, organizationId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating snapshot backup",
			"Could not update snapshot backup id "+Id+": "+err.Error(),
		)
		tflog.Debug(ctx, "error updating snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"refreshedState": refreshedState,
			"err":            err,
		})
		return
	}

	refreshedState.RegionsToCopy = state.RegionsToCopy

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
			"Error Validating Snapshot Backup IDs",
			"Could not validate ids for snapshot backup id "+state.ID.String()+" unexpected error: "+err.Error(),
		)
		tflog.Debug(ctx, "error validating snapshot backup IDs", map[string]interface{}{
			"plan": plan,
			"err":  err,
		})
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		Id             = IDs[providerschema.Id]
	)

	updateSnapshotBackupRequest := snapshot_backup.EditBackupRetentionRequest{
		Retention: int(plan.Retention.ValueInt64()),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s", s.HostURL, organizationId, projectId, clusterId, Id)
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
			"Error executing update snapshot backup retention",
			"Could not update snapshot backup id "+Id+": "+api.ParseError(err),
		)
		tflog.Debug(ctx, "error executing update snapshot backup retention", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"retention":      plan.Retention.ValueInt64(),
			"err":            err,
		})
		return
	}

	refreshedState, err := s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot backup",
			"Could not read snapshot backup id "+Id+": "+err.Error(),
		)
		tflog.Debug(ctx, "error reading snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"refreshedState": refreshedState,
			"err":            err,
		})
		return
	}

	state.Retention = types.Int64Value(int64(refreshedState.Retention))
	state.Expiration = types.StringValue(refreshedState.Expiration)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "error setting snapshot backup state", map[string]interface{}{
			"state": state,
			"err":   err,
		})
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
		tflog.Debug(ctx, "error validating snapshot backup IDs", map[string]interface{}{
			"state.OrganizationId": state.OrganizationId.String(),
			"state.ProjectID":      state.ProjectID.String(),
			"state.ClusterID":      state.ClusterID.String(),
			"state.ID":             state.ID.String(),
			"err":                  err,
		})
		resp.Diagnostics.AddError(
			"Error validating snapshot backup IDs",
			"Could not validate ids for snapshot backup id "+state.ID.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		Id             = IDs[providerschema.Id]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s", s.HostURL, organizationId, projectId, clusterId, Id)
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
			tflog.Debug(ctx, "snapshot backup doesn't exist in remote server removing resource from state file", map[string]interface{}{
				"OrganizationId": organizationId,
				"ProjectID":      projectId,
				"ClusterID":      clusterId,
				"ID":             Id,
				"err":            err,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete backup id "+state.ID.String()+": "+errString,
		)
		tflog.Debug(ctx, "error deleting snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"err":            err,
		})
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

// checkLatestSnapshotBackupStatus monitors the status of a snapshot backup creation operation for a specified organization, project, and cluster ID.
// It periodically fetches the snapshot backup job status using the `getLatestSnapshotBackup` function and waits until the snapshot backup reaches a final state or until a specified timeout is reached.
// function and waits until the snapshot backup reaches a final state or until a specified timeout is reached.
func (s *SnapshotBackup) checkSnapshotBackupStatus(ctx context.Context, organizationId, projectId, clusterId, Id string) (*snapshot_backup.SnapshotBackup, error) {
	var (
		backupResp *snapshot_backup.SnapshotBackup
		err        error
	)

	// Assuming 60 minutes is the max time snapshot backup creation takes.
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 1

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			tflog.Debug(ctx, "snapshot backup creation status timeout", map[string]interface{}{
				"organizationId": organizationId,
				"projectId":      projectId,
				"clusterId":      clusterId,
			})
			return nil, errors.ErrSnapshotBackupCreationStatusTimeout

		case <-timer.C:
			backupResp, err = s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, Id)
			switch err {
			case nil:
				if snapshot_backup.IsFinalState(backupResp.Progress.Status) && s.checkCrossRegionCopyStatus(backupResp) == nil {
					return backupResp, nil
				}

				tflog.Debug(ctx, "waiting for snapshot backup to complete the execution")
			default:
				tflog.Debug(ctx, "error getting snapshot backup", map[string]interface{}{
					"organizationId": organizationId,
					"projectId":      projectId,
					"clusterId":      clusterId,
					"Id":             Id,
					"err":            err,
				})
				return nil, err
			}
			timer.Reset(sleep)
		}
	}
}

func (s *SnapshotBackup) checkCrossRegionCopyStatus(backupResp *snapshot_backup.SnapshotBackup) error {
	for _, crossRegionCopy := range backupResp.CrossRegionCopies {
		if !snapshot_backup.IsFinalState(crossRegionCopy.Status) {
			return fmt.Errorf("cross region copy status is not final")
		}
	}
	return nil
}

func (s *SnapshotBackup) getSnapshotBackups(ctx context.Context, organizationId, projectId, clusterId string) (*snapshot_backup.ListSnapshotBackupsResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	resp, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		tflog.Debug(ctx, "error reading snapshot backups", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		return nil, err
	}

	var snapshotBackups snapshot_backup.ListSnapshotBackupsResponse
	err = json.Unmarshal(resp.Body, &snapshotBackups)
	if err != nil {
		tflog.Debug(ctx, "error unmarshalling the list of snapshot backups", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"resp":           resp,
			"err":            err,
		})
		return nil, err
	}

	return &snapshotBackups, nil
}

// getSnapshotBackup retrieves a snapshot backup by its ID.
func (s *SnapshotBackup) getSnapshotBackup(ctx context.Context, organizationId, projectId, clusterId, Id string) (*snapshot_backup.SnapshotBackup, error) {
	resp, err := s.getSnapshotBackups(ctx, organizationId, projectId, clusterId)
	if err != nil {
		tflog.Debug(ctx, "error reading snapshot backups", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		return nil, err
	}

	if len(resp.Data) > 0 {
		for _, backup := range resp.Data {
			if backup.ID == Id {
				return &backup, nil
			}
		}
	}
	tflog.Debug(ctx, "snapshot backup not found", map[string]interface{}{
		"organizationId": organizationId,
		"projectId":      projectId,
		"clusterId":      clusterId,
		"Id":             Id,
	})
	return nil, errors.ErrNotFound
}

// createSnapshotBackup creates a snapshot backup from a snapshot backup response.
func createSnapshotBackup(ctx context.Context, backupResp *snapshot_backup.SnapshotBackup, clusterId, projectId, organizationId string) (*providerschema.SnapshotBackup, error) {
	progress := providerschema.NewProgress(backupResp.Progress)
	progressObj, diags := types.ObjectValueFrom(ctx, progress.AttributeTypes(), progress)
	if diags.HasError() {
		tflog.Debug(ctx, "error during progress conversion", map[string]interface{}{
			"backupResp": backupResp,
			"progress":   progress,
		})
		return nil, fmt.Errorf("error during progress conversion")
	}

	server := providerschema.NewServer(backupResp.Server)
	serverObj, diags := types.ObjectValueFrom(ctx, server.AttributeTypes(), server)
	if diags.HasError() {
		tflog.Debug(ctx, "error during server conversion", map[string]interface{}{
			"backupResp": backupResp,
			"server":     server,
		})
		return nil, fmt.Errorf("error during server conversion")
	}

	cmekObjs := make([]basetypes.ObjectValue, 0)
	for _, cmek := range backupResp.CMEK {
		cmek := providerschema.NewCMEK(cmek)
		cmekObj, diags := types.ObjectValueFrom(ctx, cmek.AttributeTypes(), cmek)
		if diags.HasError() {
			tflog.Debug(ctx, "error during CMEK conversion", map[string]interface{}{
				"backupResp": backupResp,
				"cmek":       cmek,
			})
			return nil, fmt.Errorf("error during CMEK conversion")
		}
		cmekObjs = append(cmekObjs, cmekObj)
	}

	cmekSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: providerschema.CMEK{}.AttributeTypes()}, cmekObjs)
	if diags.HasError() {
		tflog.Debug(ctx, "error during CMEK conversion", map[string]interface{}{
			"backupResp": backupResp,
			"cmekObjs":   cmekObjs,
		})
		return nil, fmt.Errorf("error during CMEK conversion")
	}

	crossRegionCopyObjs := make([]basetypes.ObjectValue, 0)
	for _, region := range backupResp.CrossRegionCopies {
		crossRegionCopy := providerschema.NewCrossRegionCopy(region)
		crossRegionCopyObj, diags := types.ObjectValueFrom(ctx, crossRegionCopy.AttributeTypes(), crossRegionCopy)
		if diags.HasError() {
			tflog.Debug(ctx, "error during cross region copy conversion", map[string]interface{}{
				"backupResp": backupResp,
				"region":     region,
			})
			return nil, fmt.Errorf("error during cross region copy conversion")
		}
		crossRegionCopyObjs = append(crossRegionCopyObjs, crossRegionCopyObj)
	}

	crossRegionCopySet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: providerschema.CrossRegionCopy{}.AttributeTypes()}, crossRegionCopyObjs)
	if diags.HasError() {
		tflog.Debug(ctx, "error during cross region copy conversion", map[string]interface{}{
			"backupResp":          backupResp,
			"crossRegionCopyObjs": crossRegionCopyObjs,
		})
		return nil, fmt.Errorf("error during cross region copy conversion")
	}

	snapshotBackup := providerschema.NewSnapshotBackup(ctx, *backupResp, backupResp.ID, clusterId, projectId, organizationId, progressObj, serverObj, cmekSet, crossRegionCopySet)
	return &snapshotBackup, nil
}
