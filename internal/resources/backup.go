package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	internal_errors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

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

const errorMessageWhileBackupCreation = "There is an error during backup creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups", b.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	postResp, err := b.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		BackupRequest,
		b.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing create backup request",
			errorMessageWhileBackupCreation+api.ParseError(err),
		)
		return
	}

	// Capella's POST returns 202 with the new backup's id in the body. Polling
	// that exact id (rather than scanning the list endpoint for "is there a new
	// latest") avoids racing the API's undefined default sort order and the
	// per-bucket serialization window.
	var created backupapi.CreateBackupResponse
	if err := json.Unmarshal(postResp.Body, &created); err != nil || created.Id == "" {
		resp.Diagnostics.AddError(
			"Error parsing create backup response",
			fmt.Sprintf("POST returned %d but no backup id in body: %s; please check in Capella for hanging resources", postResp.Response.StatusCode, string(postResp.Body)),
		)
		return
	}

	backupResponse, err := b.waitForBackup(ctx, organizationId, projectId, clusterId, created.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while waiting for backup to reach a final state",
			"Could not confirm backup id "+created.Id+" reached a final state."+
				"Please check in Capella to see if any hanging resources have "+
				"been created, unexpected error: "+api.ParseError(err),
		)
		return
	}

	backupStats := providerschema.NewBackupStats(*backupResponse.BackupStats)
	backupStatsObj, diags := types.ObjectValueFrom(ctx, backupStats.AttributeTypes(), backupStats)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Backup Stats",
			fmt.Sprintf("Could not read backup stats data in a backup record, "+
				"please check in Capella to see if any hanging resources have been created, "+
				"unexpected error: %s", fmt.Errorf("error while backup stats conversion")),
		)
		return
	}

	scheduleInfo := providerschema.NewScheduleInfo(*backupResponse.ScheduleInfo)
	scheduleInfoObj, diags := types.ObjectValueFrom(ctx, scheduleInfo.AttributeTypes(), scheduleInfo)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Error Reading Backup Schedule Info",
			fmt.Sprintf("Could not read backup schedule info in a backup record, "+
				"please check in Capella to see if any hanging resources have been created, "+
				"unexpected error: %s", fmt.Errorf("error while backup schedule info conversion")),
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
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		backupId       = IDs[providerschema.Id]
	)

	refreshedState, err := b.retrieveBackup(ctx, organizationId, projectId, clusterId, backupId)
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
		return
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
	_, err = b.ClientV1.ExecuteWithRetry(
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
	_, err = b.ClientV1.ExecuteWithRetry(
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
		return internal_errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return internal_errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return internal_errors.ErrClusterIdCannotBeEmpty
	}
	if plan.BucketId.IsNull() {
		return internal_errors.ErrBucketIdCannotBeEmpty
	}
	if !plan.RestoreTimes.IsNull() && !plan.RestoreTimes.IsUnknown() {
		return internal_errors.ErrRestoreTimesMustNotBeSetWhileCreateBackup
	}
	return nil
}

// waitForBackup polls the specific backup record (looked up by its server-assigned
// id from the POST 202 response) until it reaches a final state (ready/failed) or
// the timeout expires. Polling by id avoids the previous design's race against the
// list endpoint's undefined default ordering and the per-bucket serialization
// window: we know exactly which record we're waiting on.
//
// 30 min is a deliberately tighter cap than the prior 90-min ceiling — under the
// 120m Go test budget, the suite can absorb several stuck backups instead of a
// single one consuming the entire run. Users with very large buckets that need
// longer should override this via a per-resource timeouts block (follow-up).
func (b *Backup) waitForBackup(ctx context.Context, organizationId, projectId, clusterId, backupId string) (*backupapi.GetBackupResponse, error) {
	const (
		timeout = 30 * time.Minute
		sleep   = 30 * time.Second
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s",
		b.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	timer := time.NewTimer(sleep)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, internal_errors.ErrBucketCreationStatusTimeout

		case <-timer.C:
			response, err := b.ClientV1.ExecuteWithRetry(ctx, cfg, nil, b.Token, nil)
			if err != nil {
				// Right after POST the new record may not yet be visible to GET;
				// treat 404 as transient and keep waiting (bounded by the outer timeout).
				if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
					tflog.Info(ctx, "backup not yet visible in API, retrying", map[string]any{"backupID": backupId})
					timer.Reset(sleep)
					continue
				}
				return nil, err
			}
			var backupResp backupapi.GetBackupResponse
			if err := json.Unmarshal(response.Body, &backupResp); err != nil {
				return nil, err
			}
			if backupapi.IsFinalState(backupResp.Status) {
				return &backupResp, nil
			}
			tflog.Info(ctx, "waiting for backup to reach final state", map[string]any{
				"backupID": backupId,
				"status":   backupResp.Status,
			})
			timer.Reset(sleep)
		}
	}
}

// retrieveBackup retrieves backup information from the specified organization and project
// using the provided backup ID by open-api call.
func (b *Backup) retrieveBackup(ctx context.Context, organizationId, projectId, clusterId, backupId string) (*providerschema.Backup, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s", b.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.ClientV1.ExecuteWithRetry(
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
		return nil, internal_errors.ErrUnableToConvertAuditData
	}

	sInfo := providerschema.NewScheduleInfo(*backupResp.ScheduleInfo)
	sInfoObj, diags := types.ObjectValueFrom(ctx, sInfo.AttributeTypes(), sInfo)
	if diags.HasError() {
		return nil, internal_errors.ErrUnableToConvertAuditData
	}

	refreshedState := providerschema.NewBackup(&backupResp, organizationId, projectId, bStatsObj, sInfoObj)
	return refreshedState, nil
}
