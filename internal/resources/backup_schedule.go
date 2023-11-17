package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	scheduleapi "terraform-provider-capella/internal/api/backup_schedule"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BackupSchedule{}
	_ resource.ResourceWithConfigure   = &BackupSchedule{}
	_ resource.ResourceWithImportState = &BackupSchedule{}
)

// BackupSchedule is the BackupSchedule resource implementation.
type BackupSchedule struct {
	*providerschema.Data
}

// NewBackupSchedule is a helper function to simplify the provider implementation.
func NewBackupSchedule() resource.Resource {
	return &BackupSchedule{}
}

// Metadata returns the BackupSchedule resource type name.
func (b *BackupSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_schedule"

}

// Schema defines the schema for the BackupSchedule resource.
func (b *BackupSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupScheduleSchema()
}

// Create creates a new BackupSchedule.
func (b *BackupSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.BackupSchedule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := b.validateCreateBackupScheduleRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create backup schedule request",
			"Could not create backup schedule "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	var weeklySchedule *providerschema.WeeklySchedule
	diags.Append(req.Config.GetAttribute(ctx, path.Root("weekly_schedule"), &weeklySchedule)...)

	BackupScheduleRequest := scheduleapi.CreateBackupScheduleRequest{
		Type: plan.Type.ValueString(),
		WeeklySchedule: scheduleapi.WeeklySchedule{
			DayOfWeek:              weeklySchedule.DayOfWeek.ValueString(),
			StartAt:                weeklySchedule.StartAt.ValueInt64(),
			IncrementalEvery:       weeklySchedule.IncrementalEvery.ValueInt64(),
			RetentionTime:          weeklySchedule.RetentionTime.ValueString(),
			CostOptimizedRetention: weeklySchedule.CostOptimizedRetention.ValueBool(),
		},
	}
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/schedules", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	_, err = b.Client.Execute(
		cfg,
		BackupScheduleRequest,
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

	refreshedState, err := b.retrieveBackupSchedule(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Backup Schedule",
			"Could not read Capella Backup Schedule for the bucket: %s "+bucketId+": "+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read reads BackupSchedule information.
func (b *BackupSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.BackupSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup schedule",
			"Could not read backup schedule for bucket"+state.BucketId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
	)

	// Get refreshed backup schedule from Capella
	refreshedState, err := b.retrieveBackupSchedule(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading backup schedule",
			"Could not read backup schedule for bucket"+state.BucketId.String()+": "+errString,
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

// Update updates the BackupSchedule.
func (b *BackupSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.BackupSchedule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup schedule",
			"Could not update backup schedule for bucket"+plan.BucketId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
	)

	var weeklySchedule *providerschema.WeeklySchedule
	diags.Append(req.Config.GetAttribute(ctx, path.Root("weekly_schedule"), &weeklySchedule)...)

	BackupScheduleRequest := scheduleapi.UpdateBackupScheduleRequest{
		Type: plan.Type.ValueString(),
		WeeklySchedule: scheduleapi.WeeklySchedule{
			DayOfWeek:              weeklySchedule.DayOfWeek.ValueString(),
			StartAt:                weeklySchedule.StartAt.ValueInt64(),
			IncrementalEvery:       weeklySchedule.IncrementalEvery.ValueInt64(),
			RetentionTime:          weeklySchedule.RetentionTime.ValueString(),
			CostOptimizedRetention: weeklySchedule.CostOptimizedRetention.ValueBool(),
		},
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/schedules", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = b.Client.Execute(
		cfg,
		BackupScheduleRequest,
		b.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup schedule",
			"Could not update backup schedule for bucket"+plan.BucketId.String()+": "+api.ParseError(err),
		)
		return
	}

	currentState, err := b.retrieveBackupSchedule(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading backup schedule",
			"Could not read backup schedule for bucket"+plan.BucketId.String()+": "+errString,
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the BackupSchedule.
func (b *BackupSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.BackupSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup schedule",
			"Could not delete backup schedule with bucket id "+state.BucketId.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/schedules", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	// Delete existing backup schedule
	_, err = b.Client.Execute(
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
			"Error deleting backup schedule",
			"Could not delete backup schedule with bucket id "+state.BucketId.String()+" unexpected error: "+errString,
		)
		return
	}
}

func (b *BackupSchedule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
}

func (b *BackupSchedule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (a *BackupSchedule) validateCreateBackupScheduleRequest(plan providerschema.BackupSchedule) error {
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

// retrieveBackupSchedule retrieves backup schedule information from the specified organization and project
// using the provided bucket ID by open-api call
func (b *BackupSchedule) retrieveBackupSchedule(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.BackupSchedule, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/schedules", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.Client.Execute(
		cfg,
		nil,
		b.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	backupScheduleResp := scheduleapi.GetBackupScheduleResponse{}
	err = json.Unmarshal(response.Body, &backupScheduleResp)
	if err != nil {
		return nil, err
	}

	scheduleInfo := providerschema.NewWeeklySchedule(*backupScheduleResp.WeeklySchedule)
	scheduleObj, diags := types.ObjectValueFrom(ctx, scheduleInfo.AttributeTypes(), scheduleInfo)
	if diags.HasError() {
		return nil, errors.ErrUnableToConvertAuditData
	}

	refreshedState := providerschema.NewBackupSchedule(&backupScheduleResp, organizationId, projectId, scheduleObj)
	return refreshedState, nil
}
