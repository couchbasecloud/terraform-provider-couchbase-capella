package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	samplebucket "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/sample_bucket"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SampleBucket{}
	_ resource.ResourceWithConfigure   = &SampleBucket{}
	_ resource.ResourceWithImportState = &SampleBucket{}
)

// Samples is the samples resource implementation.
type SampleBucket struct {
	*providerschema.Data
}

// NewSamples is a helper function to simplify the provider implementation.
func NewSampleBucket() resource.Resource {
	return &SampleBucket{}
}

// Metadata returns the samples resource type name.
func (s *SampleBucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_samplebucket"
}

// Configure It adds the provider configured api to the project resource.
func (s *SampleBucket) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	s.Data = data
}

// ImportState imports a remote sample cluster that is not created by Terraform.
func (s *SampleBucket) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the samples resource.
func (s *SampleBucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SampleBucketSchema()
}

func (s *SampleBucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Bucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	BucketRequest := samplebucket.CreateSampleBucketRequest{
		Name: plan.Name.ValueString(),
	}
	if err := s.validateCreateBucket(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		BucketRequest,
		s.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			errorMessageWhileBucketCreation+api.ParseError(err),
		)
		return
	}

	BucketResponse := samplebucket.CreateSampleBucketResponse{}
	err = json.Unmarshal(response.Body, &BucketResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			errorMessageWhileBucketCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}
	// Add validation on name. If name doesn't equal the sample names then you should throw an error here or
	// you could just let the create throw an error for it.
	plan.Id = types.StringValue(BucketResponse.Id)
	//diags = resp.State.Set(ctx, initializeSampleBucketWithPlanAndId(plan, BucketResponse.Id))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// failed to retrieve the bucket
	refreshedState, err := s.retrieveBucket(ctx, organizationId, projectId, clusterId, BucketResponse.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating bucket "+BucketResponse.Id,
			errorMessageAfterBucketCreation+api.ParseError(err),
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

func (s *SampleBucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Bucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Bucket in Capella",
			"Could not read Capella Bucket with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.Id]
	)

	refreshedState, err := s.retrieveBucket(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading bucket",
			"Could not read bucket with id "+state.Id.String()+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *SampleBucket) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for sample buckets.
	// Allowlists can only be created, read and deleted.
	// http://cbc-cp-api.s3-website-us-east-1.amazonaws.com/#tag/sampleBucket
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

func (s *SampleBucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.Bucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrOrganizationIdCannotBeEmpty.Error(),
		)
		return
	}
	var organizationId = state.OrganizationId.ValueString()

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = state.ProjectId.ValueString()

	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = state.ClusterId.ValueString()

	if state.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var bucketId = state.Id.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets/%s", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := s.Client.ExecuteWithRetry(
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
			"Error Deleting the Bucket",
			"Could not delete Bucket associated with cluster "+clusterId+": "+errString,
		)
		return
	}
}

func (r *SampleBucket) validateCreateBucket(plan providerschema.Bucket) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return r.validateBucketAttributesTrimmed(plan)
}

// Add extra validaiton for
func (r *SampleBucket) validateBucketAttributesTrimmed(plan providerschema.Bucket) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	return nil
}

// retrieveBucket retrieves bucket information for a specified organization, project, cluster and bucket ID.
func (s *SampleBucket) retrieveBucket(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.OneBucket, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets/%s", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	bucketResp := samplebucket.GetSampleBucketResponse{}
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

// initializeBucketWithPlanAndId initializes an instance of providerschema.Bucket
// with the specified plan and ID. It marks all computed fields as null.
func initializeSampleBucketWithPlanAndId(plan providerschema.Bucket, id string) providerschema.Bucket {
	plan.Id = types.StringValue(id)
	// Do I need this?
	/*
		if plan.StorageBackend.IsNull() || plan.StorageBackend.IsUnknown() {
			plan.StorageBackend = types.StringNull()
		}
		if plan.EvictionPolicy.IsNull() || plan.EvictionPolicy.IsUnknown() {
			plan.EvictionPolicy = types.StringNull()
		}
		plan.Stats = types.ObjectNull(providerschema.Stats{}.AttributeTypes())*/
	return plan
}
