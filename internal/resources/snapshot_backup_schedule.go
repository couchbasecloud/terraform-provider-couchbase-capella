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
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup_schedule"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	errorMessageWhileSnapshotBackupScheduleCreation = "There is an error during snapshot backup schedule creation. Please check in Capella to see if any hanging resources" +
		" have been created, unexpected error: "
)

var (
	_ resource.Resource                = &SnapshotBackupSchedule{}
	_ resource.ResourceWithConfigure   = &SnapshotBackupSchedule{}
	_ resource.ResourceWithImportState = &SnapshotBackupSchedule{}
)

// SnapshotBackupSchedule is the Snapshot Backup Schedule resource implementation.
type SnapshotBackupSchedule struct {
	*providerschema.Data
}

// NewSnapshotBackupSchedule is a helper function to simplify the provider implementation.
func NewSnapshotBackupSchedule() resource.Resource {
	return &SnapshotBackupSchedule{}
}

// Metadata returns the Snapshot Backup Schedule resource type name.
func (s *SnapshotBackupSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_backup_schedule"
}

// Schema defines the schema for the Snapshot Backup Schedule resource.
func (s *SnapshotBackupSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotBackupScheduleSchema()
}

// ImportState imports a remote snapshot backup schedule that is not created by Terraform.
func (s *SnapshotBackupSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to cluster_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// Create creates a new Snapshot Backup Schedule.
func (s *SnapshotBackupSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.SnapshotBackupSchedule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationID.ValueString()
		projectId      = plan.ProjectID.ValueString()
		clusterId      = plan.ClusterID.ValueString()
	)

	err := s.upsertSnapshotBackupSchedule(ctx, organizationId, projectId, clusterId, plan)
	if err != nil {
		tflog.Debug(ctx, "Error upserting snapshot backup schedule", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		resp.Diagnostics.AddError(
			"Error Upserting Snapshot Backup Schedule in Capella",
			errorMessageWhileSnapshotBackupScheduleCreation+err.Error(),
		)
		return
	}

	var refreshedState *providerschema.SnapshotBackupSchedule
	snapshotBackupSchedule, err := s.getSnapshotBackupSchedule(ctx, organizationId, projectId, clusterId, plan.StartTime.ValueString())
	if err != nil {
		tflog.Debug(ctx, "Error getting snapshot backup schedule after upsert", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		resp.Diagnostics.AddWarning(
			"Error Getting Snapshot Backup Schedule in Capella",
			"Could not get Capella Snapshot Backup Schedule for cluster with ID "+plan.ClusterID.String()+": "+err.Error(),
		)
		refreshedState = &providerschema.SnapshotBackupSchedule{}
		refreshedState.CopyToRegions = []types.String{}
	} else {
		newSnapshotBackupSchedule := providerschema.NewSnapshotBackupSchedule(*snapshotBackupSchedule, organizationId, projectId, clusterId)
		refreshedState = &newSnapshotBackupSchedule
	}

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read reads snapshot backup schedule information.
func (s *SnapshotBackupSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.SnapshotBackupSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	IDs, err := state.Validate()
	if err != nil {
		tflog.Debug(ctx, "Error validating snapshot backup schedule", map[string]interface{}{
			"state": state,
			"err":   err,
		})
		resp.Diagnostics.AddError(
			"Error Validating Backup Schedule in Capella",
			"Could not validate Capella Backup Schedule for cluster with ID "+state.ClusterID.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	snapshotBackupSchedule, err := s.getSnapshotBackupSchedule(ctx, organizationId, projectId, clusterId, state.StartTime.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Snapshot Backup Schedule in Capella",
			"Could not get Capella Snapshot Backup Schedule for cluster with ID "+state.ClusterID.String()+": "+err.Error(),
		)
		return
	}

	refreshedState := providerschema.NewSnapshotBackupSchedule(*snapshotBackupSchedule, organizationId, projectId, clusterId)
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the snapshot backup schedule.
func (s *SnapshotBackupSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan providerschema.SnapshotBackupSchedule

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationID.ValueString()
		projectId      = plan.ProjectID.ValueString()
		clusterId      = plan.ClusterID.ValueString()
	)

	err := s.upsertSnapshotBackupSchedule(ctx, organizationId, projectId, clusterId, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Upserting Snapshot Backup Schedule in Capella",
			"Could not upsert Capella Snapshot Backup Schedule for cluster with ID "+clusterId+": "+err.Error(),
		)
		return
	}
	snapshotBackupSchedule, err := s.getSnapshotBackupSchedule(ctx, organizationId, projectId, clusterId, plan.StartTime.ValueString())
	if err != nil {
		tflog.Debug(ctx, "Error getting snapshot backup schedule after upsert", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		resp.Diagnostics.AddError(
			"Error Getting Snapshot Backup Schedule in Capella",
			"Could not get Capella Snapshot Backup Schedule for cluster with ID "+plan.ClusterID.String()+": "+err.Error(),
		)
	}

	refreshedState := providerschema.NewSnapshotBackupSchedule(*snapshotBackupSchedule, organizationId, projectId, clusterId)

	// Sets state to fully populated data.
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the snapshot backup schedule.
func (s *SnapshotBackupSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.SnapshotBackupSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		tflog.Debug(ctx, "Error validating snapshot backup schedule", map[string]interface{}{
			"state": state,
			"err":   err,
		})

		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete backup id "+state.ClusterID.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackupschedule", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
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
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting backup",
			"Could not delete snapshot backup schedule for cluster with ID "+state.ClusterID.String()+": "+errString,
		)
		return
	}
}

// upsertSnapshotBackupSchedule creates or updates the snapshot backup schedule.
func (s *SnapshotBackupSchedule) upsertSnapshotBackupSchedule(ctx context.Context, organizationId, projectId, clusterId string, plan providerschema.SnapshotBackupSchedule) error {

	createSnapshotBackupScheduleRequest := snapshot_backup_schedule.SnapshotBackupSchedule{
		Interval:      plan.Interval.ValueInt64(),
		Retention:     plan.Retention.ValueInt64(),
		StartTime:     plan.StartTime.ValueString(),
		CopyToRegions: providerschema.ConvertStringValueList(plan.CopyToRegions),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackupschedule", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := s.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		createSnapshotBackupScheduleRequest,
		s.Token,
		nil,
	)

	if err != nil {
		tflog.Debug(ctx, "Error upserting snapshot backup schedule", map[string]interface{}{
			"organizationId":                      organizationId,
			"projectId":                           projectId,
			"clusterId":                           clusterId,
			"createSnapshotBackupScheduleRequest": createSnapshotBackupScheduleRequest,
			"err":                                 err,
		})
		return err
	}
	return nil
}

// getSnapshotBackupSchedule retrieves the snapshot backup schedule for a cluster.
func (s *SnapshotBackupSchedule) getSnapshotBackupSchedule(ctx context.Context, organizationId, projectId, clusterId string, stateTimeString string) (*snapshot_backup_schedule.SnapshotBackupSchedule, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackupschedule", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	backupScheduleResp, err := s.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)

	if err != nil {
		tflog.Debug(ctx, "Error getting snapshot backup schedule", map[string]interface{}{
			"organizationId": organizationId,
			"projectId":      projectId,
			"clusterId":      clusterId,
			"err":            err,
		})
		return nil, err
	}

	snapshotBackupSchedule := snapshot_backup_schedule.SnapshotBackupSchedule{}
	err = json.Unmarshal(backupScheduleResp.Body, &snapshotBackupSchedule)
	if err != nil {
		tflog.Debug(ctx, "Error unmarshalling snapshot backup schedule", map[string]interface{}{
			"backupScheduleResp.Body": backupScheduleResp.Body,
			"err":                     err,
		})
		return nil, err
	}

	snapshotBackupSchedule.StartTime, err = s.getStartTime(ctx, stateTimeString, &snapshotBackupSchedule)
	if err != nil {
		return nil, err
	}

	return &snapshotBackupSchedule, nil
}

// getStartTime compares the start time of the current state with the start time of the actual resource, and returns the resource's start time only if it is different.
// This ensures that the resource storing an equivalent start time does not cause the state to be unnecessarily updated.
func (s *SnapshotBackupSchedule) getStartTime(ctx context.Context, currentStartTimeString string, snapshotBackupSchedule *snapshot_backup_schedule.SnapshotBackupSchedule) (string, error) {
	newStartTime, err := time.Parse(time.RFC3339, snapshotBackupSchedule.StartTime)
	if err != nil {
		tflog.Debug(ctx, "Error parsing updated start time", map[string]interface{}{
			"snapshotBackupSchedule.StartTime": snapshotBackupSchedule.StartTime,
			"err":                              err,
		})
		return "", err
	}
	if currentStartTimeString == "" {
		return newStartTime.Format(time.RFC3339), nil
	}
	currentStartTime, err := time.Parse(time.RFC3339, currentStartTimeString)
	if err != nil {
		tflog.Debug(ctx, "Error parsing current start time", map[string]interface{}{
			"currentStartTimeString": currentStartTimeString,
			"err":                    err,
		})
		return "", err
	}
	if currentStartTime.Equal(newStartTime) {
		return currentStartTimeString, nil
	}
	return newStartTime.Format(time.RFC3339), nil
}

// Configure adds the provider configured api to the snapshot backup schedule resource.
func (s *SnapshotBackupSchedule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
