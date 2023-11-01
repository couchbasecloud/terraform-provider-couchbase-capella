package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	backupapi "terraform-provider-capella/internal/api/backup"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"
	"time"
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

	//if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
	//	BackupRequest.Type = plan.Type.ValueStringPointer()
	//	//BackupRequest.WeeklySchedule = &backupapi.WeeklySchedule{
	//	//	DayOfWeek:              plan.WeeklySchedule.DayOfWeek.ValueString(),
	//	//	StartAt:                plan.WeeklySchedule.StartAt.ValueInt64(),
	//	//	IncrementalEvery:       plan.WeeklySchedule.IncrementalEvery.ValueInt64(),
	//	//	RetentionTime:          plan.WeeklySchedule.RetentionTime.ValueString(),
	//	//	CostOptimizedRetention: plan.WeeklySchedule.CostOptimizedRetention.ValueBool(),
	//	//}
	//}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	fmt.Printf("//////////////bucketID: %s", bucketId)
	latestBackup, err := b.getLatestBackup(organizationId, projectId, clusterId, bucketId)
	if err != nil {
		fmt.Print("//////////////ERRR")
		fmt.Print(err.Error())
	}
	fmt.Printf("//////////////latestBackup: %q", latestBackup)
	var backupFound bool
	if latestBackup != nil {
		backupFound = true
	}
	fmt.Printf("//////////////backupFound: %t", backupFound)

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
	//err = json.Unmarshal(response.Body, BackupResponse)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error creating backup",
	//		"Could not create backup, error during unmarshalling:"+err.Error(),
	//	)
	//	return
	//}

	refreshedState := providerschema.NewBackup(BackupResponse, organizationId, projectId)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
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

// checkClusterStatus monitors the status of a cluster creation, update and deletion operation for a specified
// organization, project, and cluster ID. It periodically fetches the cluster status using the `getCluster`
// function and waits until the cluster reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (b *Backup) checkLatestBackupStatus(ctx context.Context, organizationId, projectId, clusterId, bucketId string, backupFound bool, latestBackup *backupapi.GetBackupResponse) (*backupapi.GetBackupResponse, error) {
	var (
		backupResp *backupapi.GetBackupResponse
		err        error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 1

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			const msg = "cluster creation status transition timed out after initiation"
			return nil, fmt.Errorf(msg)

		case <-timer.C:
			backupResp, err = b.getLatestBackup(organizationId, projectId, clusterId, bucketId)
			switch err {
			case nil:
				fmt.Println("%%%%%%%%%%%%%%")
				fmt.Println(backupResp)
				fmt.Println(")))))))))))))))))))))))))))")
				if !backupFound && backupResp != nil && backupapi.IsFinalState(backupResp.Status) {
					fmt.Println("^^^^^^^^^^^^^^^^^^^^")
					fmt.Println(backupResp.Id)
					return backupResp, nil
				} else if backupFound && backupResp != nil && latestBackup.Id == backupResp.Id && backupapi.IsFinalState(backupResp.Status) {
					fmt.Println("#######################")
					fmt.Println(backupResp.Id)
					return backupResp, nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return nil, err
			}
			timer.Reset(sleep)
		}
	}
}

// getCluster retrieves cluster information from the specified organization and project
// using the provided cluster ID by open-api call
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

	for _, backup := range clusterResp.Data {
		if backup.BucketId == bucketId {
			fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&")
			fmt.Print(backup.Id)
			return &backup, nil
		}
	}

	//fmt.Print(clusterResp.Data)
	return nil, nil
}
