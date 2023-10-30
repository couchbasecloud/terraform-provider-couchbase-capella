package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/errors"

	backupapi "terraform-provider-capella/internal/api/backup"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "terraform-provider-capella/internal/schema"
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

func (b *Backup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (b *Backup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupSchema()
}

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

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		BackupRequest.Type = plan.Type.ValueStringPointer()
		BackupRequest.WeeklySchedule = &backupapi.WeeklySchedule{
			DayOfWeek:              plan.WeeklySchedule.DayOfWeek.ValueString(),
			StartAt:                plan.WeeklySchedule.StartAt.ValueInt64(),
			IncrementalEvery:       plan.WeeklySchedule.IncrementalEvery.ValueInt64(),
			RetentionTime:          plan.WeeklySchedule.RetentionTime.ValueString(),
			CostOptimizedRetention: plan.WeeklySchedule.CostOptimizedRetention.ValueBool(),
		}
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	response, err := b.Client.Execute(
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

	BackupResponse := backupapi.GetBackupResponse{}
	err = json.Unmarshal(response.Body, &BackupResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup",
			"Could not create backup, error during unmarshalling:"+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	//diags = resp.State.Set(ctx, nil)

	//resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (b *Backup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (b *Backup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b *Backup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (b *Backup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

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
