package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
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

	freeTierBucketRequest := api.CreateFreeTierBucketRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.MemoryAllocationInMB.IsNull() && !plan.MemoryAllocationInMB.IsUnknown() {
		freeTierBucketRequest.MemoryAllocationInMb = plan.MemoryAllocationInMB.ValueInt64Pointer()
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/freeTier", f.HostURL, organizationId, projectId, clusterId)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	resp, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		freeTierBucketRequest,
		f.Token,
		nil,
	)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating free-tier bucket",
			errors.ErrorMessageWhileFreeTierBucketCreation.Error()+api.ParseError(err),
		)
		return
	}

	freeTierBucketResponse := bucketapi.CreateBucketResponse{}
	err = json.Unmarshal(resp.Body, &freeTierBucketResponse)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating free-tier bucket",
			errors.ErrorMessageWhileFreeTierBucketCreation.Error()+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	refreshedState, err := f.retrieveFreeTierBucket(ctx, organizationId, projectId, clusterId, freeTierBucketResponse.Id)
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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, bucketId)

	updateBucketRequest := api.UpdateFreeTierBucketRequest{
		MemoryAllocationInMb: plannedState.MemoryAllocationInMB.ValueInt64(),
	}

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		updateBucketRequest,
		f.Token,
		nil,
	)
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
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, freeTierBucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
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
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/freeTier/%s", f.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	bucketResp := bucketapi.GetBucketResponse{}
	err = json.Unmarshal(response.Body, &bucketResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	refreshedState := providerschema.OneBucket{
		Id:                       types.StringValue(bucketResp.Id),
		Name:                     types.StringValue(bucketResp.Name),
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		Type:                     types.StringValue(bucketResp.Type),
		StorageBackend:           types.StringValue(bucketResp.StorageBackend),
		MemoryAllocationInMB:     types.Int64Value(bucketResp.MemoryAllocationInMb),
		BucketConflictResolution: types.StringValue(bucketResp.BucketConflictResolution),
		DurabilityLevel:          types.StringValue(bucketResp.DurabilityLevel),
		Replicas:                 types.Int64Value(bucketResp.Replicas),
		Flush:                    types.BoolValue(bucketResp.Flush),
		TimeToLiveInSeconds:      types.Int64Value(bucketResp.TimeToLiveInSeconds),
		EvictionPolicy:           types.StringValue(bucketResp.EvictionPolicy),
		Stats: &providerschema.Stats{
			ItemCount:       types.Int64Value(bucketResp.Stats.ItemCount),
			OpsPerSecond:    types.Int64Value(bucketResp.Stats.OpsPerSecond),
			DiskUsedInMiB:   types.Int64Value(bucketResp.Stats.DiskUsedInMib),
			MemoryUsedInMiB: types.Int64Value(bucketResp.Stats.MemoryUsedInMib),
		},
	}

	return &refreshedState, nil

}
