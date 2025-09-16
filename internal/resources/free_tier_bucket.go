package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
)

var (
	_ resource.Resource                = &FreeTierBucket{}
	_ resource.ResourceWithConfigure   = &FreeTierBucket{}
	_ resource.ResourceWithImportState = &FreeTierBucket{}
)

type FreeTierBucket struct {
	*providerschema.Data
}

func NewFreeTierBucket() resource.Resource {
	return &FreeTierBucket{}
}
func (f *FreeTierBucket) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (f *FreeTierBucket) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	data, ok := request.ProviderData.(*providerschema.Data)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	f.Data = data
}

func (f *FreeTierBucket) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_free_tier_bucket"
}

func (f *FreeTierBucket) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = FreeTierBucketSchema()
}

func (f *FreeTierBucket) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.Bucket
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	freeTierBucketRequest := apigen.CreateFreeTierBucketRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.MemoryAllocationInMB.IsNull() && !plan.MemoryAllocationInMB.IsUnknown() {
		v := int(plan.MemoryAllocationInMB.ValueInt64())
		freeTierBucketRequest.MemoryAllocationInMb = &v
	}

	orgUUID, _ := uuid.Parse(plan.OrganizationId.ValueString())
	projUUID, _ := uuid.Parse(plan.ProjectId.ValueString())
	cluUUID, _ := uuid.Parse(plan.ClusterId.ValueString())

	res, err := f.ClientV2.CreateFreeTierBucketWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), freeTierBucketRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating free-tier bucket",
			errors.ErrorMessageWhileFreeTierBucketCreation.Error()+api.ParseError(err),
		)
		return
	}
	if res.JSON201 == nil {
		response.Diagnostics.AddError("Error creating free-tier bucket", "unexpected status: "+res.Status())
		return
	}

	refreshedState, err := f.retrieveFreeTierBucket(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), res.JSON201.Id)
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error fetching free-tier bucket",
			errors.ErrorMessageWhileFreeTierBucketCreation.Error()+api.ParseError(err),
		)
	}
	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f *FreeTierBucket) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var currState providerschema.Bucket
	diags := request.State.Get(ctx, &currState)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	IDs, err := currState.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading free-tier bucket in Capella",
			"Could not read free-tier Bucket with ID "+currState.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.Id]
	)

	refreshedState, err := f.retrieveFreeTierBucket(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading free-tier bucket",
			"could not read the free-tier bucket with ID "+currState.Id.String()+": "+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierBucket) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plannedState providerschema.Bucket
	diags := request.Plan.Get(ctx, &plannedState)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	var currState providerschema.Bucket
	diags = request.State.Get(ctx, &currState)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	organizationId := currState.OrganizationId.ValueString()
	projectId := currState.ProjectId.ValueString()
	clusterId := currState.ClusterId.ValueString()
	bucketId := currState.Id.ValueString()

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)

	updateBucketRequest := apigen.UpdateFreeTierBucketRequest{
		MemoryAllocationInMb: int(plannedState.MemoryAllocationInMB.ValueInt64()),
	}

	_, err := f.ClientV2.UpdateFreeTierBucketWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(bucketId), updateBucketRequest)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error updating free-tier bucket",
			"Could not update free-tier bucket, unexpected error: "+bucketId+": "+errString,
		)
		return
	}

	updatedState, err := f.retrieveFreeTierBucket(ctx, organizationId, projectId, clusterId, bucketId)

	if err != nil {
		response.Diagnostics.AddWarning(
			"Bucket updated but could not retrieve the bucket info",
			"Could not fetch updated free-tier bucket info, bucket ID - "+bucketId+": "+api.ParseError(err),
		)
	}

	diags = response.State.Set(ctx, updatedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierBucket) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var currentState providerschema.Bucket
	diags := request.State.Get(ctx, &currentState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	organizationId := currentState.OrganizationId.ValueString()
	projectId := currentState.ProjectId.ValueString()
	clusterId := currentState.ClusterId.ValueString()
	freeTierBucketId := currentState.Id.ValueString()

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)

	_, err := f.ClientV2.DeleteFreeTierBucketByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(freeTierBucketId))
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error deleting the bucket",
			"Could not delete free-tier bucket associated with cluster "+clusterId+": "+errString,
		)
		return
	}
}

func (f *FreeTierBucket) retrieveFreeTierBucket(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.OneBucket, error) {
	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)

	res, err := f.ClientV2.GetFreeTierBucketByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(bucketId))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("%s: unexpected status %s", errors.ErrExecutingRequest, res.Status())
	}

	bucketResp := *res.JSON200

	refreshedState := providerschema.OneBucket{
		Id:                       types.StringValue(bucketResp.Id),
		Name:                     types.StringValue(bucketResp.Name),
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		Type:                     types.StringValue(bucketResp.Type),
		StorageBackend:           types.StringValue(bucketResp.StorageBackend),
		MemoryAllocationInMB:     types.Int64Value(int64(bucketResp.MemoryAllocationInMb)),
		BucketConflictResolution: types.StringValue(bucketResp.BucketConflictResolution),
		DurabilityLevel:          types.StringValue(bucketResp.DurabilityLevel),
		Replicas:                 types.Int64Value(int64(bucketResp.Replicas)),
		TimeToLiveInSeconds:      types.Int64Value(int64(bucketResp.TimeToLiveInSeconds)),
		EvictionPolicy:           types.StringValue(string(bucketResp.EvictionPolicy)),
	}
	refreshedState.Stats = &providerschema.Stats{
		ItemCount:       types.Int64Value(int64(bucketResp.Stats.ItemCount)),
		OpsPerSecond:    types.Int64Value(int64(bucketResp.Stats.OpsPerSecond)),
		DiskUsedInMiB:   types.Int64Value(int64(bucketResp.Stats.DiskUsedInMib)),
		MemoryUsedInMiB: types.Int64Value(int64(bucketResp.Stats.MemoryUsedInMib)),
	}

	return &refreshedState, nil

}
