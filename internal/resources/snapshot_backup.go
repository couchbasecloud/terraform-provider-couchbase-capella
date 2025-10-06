package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	errorMessageConfigure = "Expected *ProviderSourceData, got:"
)

var (
	_ resource.Resource                = &SnapshotBackup{}
	_ resource.ResourceWithConfigure   = &SnapshotBackup{}
	_ resource.ResourceWithImportState = &SnapshotBackup{}
)

// SnapshotBackup is the Snapshot Backup resource implementation.
type SnapshotBackup struct {
	*providerschema.Data
}

// NewSnapshotBackup is a helper function to simplify the provider implementation.
func NewSnapshotBackup() resource.Resource {
	return &SnapshotBackup{}
}

// Metadata returns the Snapshot Backup resource type name.
func (s *SnapshotBackup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_backup"
}

// Schema defines the schema for the Snapshot Backup resource.
func (s *SnapshotBackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotBackupSchema()
}

// ImportState imports a remote backup that is not created by Terraform.
func (s *SnapshotBackup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates a new Snapshot Backup.
func (s *SnapshotBackup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.SnapshotBackup
	var refreshedState *providerschema.SnapshotBackup

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
		Retention:     plan.Retention.ValueInt64(),
		RegionsToCopy: providerschema.ConvertStringValueList(plan.RegionsToCopy),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	createResp, err := s.ClientV1.ExecuteWithRetry(
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
		return
	}

	backupResp, err := s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, createSnapshotBackupResponse.ID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error while checking latest snapshot backup status",
			errorMessageWhileSnapshotBackupCreation+api.ParseError(err),
		)
		refreshedState = &providerschema.SnapshotBackup{}
		refreshedState = setNullValues(refreshedState, clusterId, projectId, organizationId, createSnapshotBackupResponse.ID)
	} else {
		refreshedState, err = morphToTerraformCloudSnapshotBackup(ctx, backupResp, clusterId, projectId, organizationId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating snapshot backup",
				"Could not create snapshot backup: "+err.Error(),
			)
			return
		}
	}

	refreshedState.RegionsToCopy = plan.RegionsToCopy

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads snapshot backup information.
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
		if err == errors.ErrNotFound {
			resp.State.RemoveResource(ctx)
		}
		tflog.Debug(ctx, "error reading snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"err":            err,
		})
		return
	}

	refreshedState, err := morphToTerraformCloudSnapshotBackup(ctx, snapshotBackup, clusterId, projectId, organizationId)
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
	refreshedState.RestoreTimes = state.RestoreTimes
	refreshedState.CrossRegionRestorePreference = state.CrossRegionRestorePreference

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)

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

	if state.Retention.ValueInt64() != plan.Retention.ValueInt64() {
		err = s.updateRetention(ctx, organizationId, projectId, clusterId, Id, plan.Retention.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating snapshot backup retention",
				"Could not update snapshot backup id "+state.ID.String()+": "+err.Error(),
			)
			return
		}
	}

	if plan.RestoreTimes.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error restoring backup",
			"Could not restore backup id "+state.ID.String()+": plan restore times value is not set",
		)
		return
	}

	if !plan.RestoreTimes.IsNull() {
		if !state.RestoreTimes.IsNull() && !state.RestoreTimes.IsUnknown() {
			planRestoreTimes := *plan.RestoreTimes.ValueBigFloat()
			stateRestoreTimes := *state.RestoreTimes.ValueBigFloat()
			if planRestoreTimes.Cmp(&stateRestoreTimes) == -1 {
				resp.Diagnostics.AddError(
					"Error restoring backup",
					"Could not restore backup id "+state.ID.String()+": plan restore times value is not greater than state restore times value",
				)
				return
			} else if planRestoreTimes.Cmp(&stateRestoreTimes) == 1 {
				err = s.restoreSnapshotBackup(ctx, organizationId, projectId, clusterId, Id, providerschema.ConvertStringValueList(plan.CrossRegionRestorePreference))
				if err != nil {
					resp.Diagnostics.AddError(
						"Error restoring snapshot backup",
						"Could not restore snapshot backup id "+state.ID.String()+": "+err.Error(),
					)
					return
				}
			}
		}
	}

	refreshedState, err := s.getSnapshotBackup(ctx, organizationId, projectId, clusterId, Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot backup",
			"Could not read snapshot backup id "+Id+": "+err.Error(),
		)
		if err == errors.ErrNotFound {
			resp.State.RemoveResource(ctx)
		}
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

	state.Retention = types.Int64Value(refreshedState.Retention)
	state.Expiration = types.StringValue(refreshedState.Expiration)
	state.RegionsToCopy = plan.RegionsToCopy
	state.RestoreTimes = plan.RestoreTimes
	state.CrossRegionRestorePreference = plan.CrossRegionRestorePreference

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the snapshot backup.
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
	_, err = s.ClientV1.ExecuteWithRetry(
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

// getSnapshotBackups retrieves a list of snapshot backups for a cluster.
func (s *SnapshotBackup) getSnapshotBackups(ctx context.Context, organizationId, projectId, clusterId string) (*snapshot_backup.ListSnapshotBackupsResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	resp, err := s.ClientV1.ExecuteWithRetry(
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

// setNullValues sets the snapshot backup's computed values to null.
func setNullValues(refreshedState *providerschema.SnapshotBackup, clusterId, projectId, organizationId, Id string) *providerschema.SnapshotBackup {
	refreshedState.ID = types.StringValue(Id)
	refreshedState.ClusterID = types.StringValue(clusterId)
	refreshedState.ProjectID = types.StringValue(projectId)
	refreshedState.OrganizationId = types.StringValue(organizationId)
	refreshedState.CreatedAt = types.StringNull()
	refreshedState.Expiration = types.StringNull()
	refreshedState.Progress = types.ObjectNull(providerschema.Progress{}.AttributeTypes())
	refreshedState.Server = types.ObjectNull(providerschema.Server{}.AttributeTypes())
	refreshedState.CMEK = types.SetNull(types.ObjectType{AttrTypes: providerschema.CMEK{}.AttributeTypes()})
	refreshedState.CrossRegionCopies = types.SetNull(types.ObjectType{AttrTypes: providerschema.CrossRegionCopy{}.AttributeTypes()})
	refreshedState.Size = types.Int64Null()
	refreshedState.Type = types.StringNull()

	if refreshedState.Retention.IsNull() || refreshedState.Retention.IsUnknown() {
		refreshedState.Retention = types.Int64Null()
	}
	return refreshedState
}

// morphToTerraformCloudSnapshotBackup creates a snapshot backup from a snapshot backup response.
func morphToTerraformCloudSnapshotBackup(ctx context.Context, backupResp *snapshot_backup.SnapshotBackup, clusterId, projectId, organizationId string) (*providerschema.SnapshotBackup, error) {
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

	snapshotBackup := providerschema.NewSnapshotBackup(*backupResp, backupResp.ID, clusterId, projectId, organizationId, progressObj, serverObj, cmekSet, crossRegionCopySet)
	return &snapshotBackup, nil
}

func (s *SnapshotBackup) updateRetention(ctx context.Context, organizationId, projectId, clusterId, Id string, retention int64) error {
	updateSnapshotBackupRequest := snapshot_backup.EditBackupRetentionRequest{
		Retention: retention,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s", s.HostURL, organizationId, projectId, clusterId, Id)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := s.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		updateSnapshotBackupRequest,
		s.Token,
		nil,
	)
	if err != nil {
		tflog.Debug(ctx, "error executing update snapshot backup retention", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"retention":      retention,
			"err":            err,
		})
		return err
	}
	return nil
}

func (s *SnapshotBackup) restoreSnapshotBackup(ctx context.Context, organizationId, projectId, clusterId, Id string, crossRegionRestorePreference []string) error {
	var (
		resp *api.Response
		err  error
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/%s/restore", s.HostURL, organizationId, projectId, clusterId, Id)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}

	if len(crossRegionRestorePreference) > 0 {
		restoreSnapshotBackupRequest := snapshot_backup.RestoreSnapshotBackupRequest{
			CrossRegionRestorePreference: crossRegionRestorePreference,
		}
		resp, err = s.ClientV1.ExecuteWithRetry(
			ctx,
			cfg,
			restoreSnapshotBackupRequest,
			s.Token,
			nil,
		)
	} else {
		resp, err = s.ClientV1.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			s.Token,
			nil,
		)
	}
	if err != nil {
		tflog.Debug(ctx, "error executing restore snapshot backup", map[string]interface{}{
			"OrganizationId": organizationId,
			"ProjectID":      projectId,
			"ClusterID":      clusterId,
			"ID":             Id,
			"err":            err,
		})
		return err
	}

	var restoreResp snapshot_backup.RestoreSnapshotBackupResponse
	err = json.Unmarshal(resp.Body, &restoreResp)
	if err != nil {
		tflog.Debug(ctx, "error unmarshalling the restore snapshot backup response", map[string]interface{}{
			"resp": resp,
			"err":  err,
		})
		return err
	}

	return nil
}
