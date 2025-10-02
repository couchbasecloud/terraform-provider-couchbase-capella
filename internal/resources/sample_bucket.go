package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	samplebucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/sample_bucket"
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

const errorMessageAfterSampleBucketCreation = "Sample bucket loading is successful, but encountered an error while checking the current" +
	" state of the sample bucket. Please run `terraform plan` after 1-2 minutes to know the" +
	" current sample bucket state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileSampleBucketCreation = "There is an error during sample bucket loading. Please check in Capella to see if any hanging resources" +
	" have been loaded, unexpected error: "

// SampleBucket is the sample bucket resource implementation.
type SampleBucket struct {
	*providerschema.Data
}

// NewSampleBucket is a helper function to simplify the provider implementation.
func NewSampleBucket() resource.Resource {
	return &SampleBucket{}
}

// Metadata returns the SampleBucket resource type name.
func (s *SampleBucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sample_bucket"
}

// Configure adds the configured client to the SampleBucket resource.
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

// ImportState imports a remote sample bucket that is not created by Terraform.
func (s *SampleBucket) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the SampleBucket resource.
func (s *SampleBucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SampleBucketSchema()
}

// Create loads a new sample bucket.
func (s *SampleBucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.SampleBucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	sampleBucketRequest := samplebucketapi.CreateSampleBucketRequest{
		Name: plan.Name.ValueString(),
	}
	if err := s.validateCreateSampleBucket(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Loading Sample Bucket",
			"Could not load sample bucket, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets", s.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := s.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		sampleBucketRequest,
		s.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Loading Sample Bucket",
			errorMessageWhileSampleBucketCreation+api.ParseError(err),
		)
		return
	}

	time.Sleep(10 * time.Second)

	sampleBucketResponse := samplebucketapi.CreateSampleBucketResponse{}
	err = json.Unmarshal(response.Body, &sampleBucketResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Loading Sample Bucket",
			errorMessageWhileSampleBucketCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeSampleBucketWithPlanAndId(plan, sampleBucketResponse.Id))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := s.retrieveSampleBucket(ctx, organizationId, projectId, clusterId, sampleBucketResponse.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Loading Sample Bucket "+sampleBucketResponse.Id,
			errorMessageAfterSampleBucketCreation+api.ParseError(err),
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

// Read reads SampleBucket information.
func (s *SampleBucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.SampleBucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Sample Bucket in Capella",
			"Could not read Capella sample bucket with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.Id]
	)

	refreshedState, err := s.retrieveSampleBucket(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Sample Bucket in Capella",
			"Could not read sample bucket with id "+state.Id.String()+": "+errString,
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
	// SampleBuckets can only be created, read and deleted.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/sampleBucket
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

// Delete deletes the SampleBucket.
func (s *SampleBucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.SampleBucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error Deleting Sample Bucket",
			"Could not delete sample bucket, unexpected error: "+errors.ErrOrganizationIdCannotBeEmpty.Error(),
		)
		return
	}
	var organizationId = state.OrganizationId.ValueString()

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error Deleting Sample Bucket",
			"Could not delete sample bucket, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = state.ProjectId.ValueString()

	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error Deleting Sample Bucket",
			"Could not delete sample bucket, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = state.ClusterId.ValueString()

	if state.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Error Deleting Sample Bucket",
			"Could not delete sample bucket, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var bucketId = state.Id.ValueString()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets/%s", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := s.ClientV1.ExecuteWithRetry(
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
			"Error Deleting Sample Bucket",
			"Could not delete sample bucket associated with cluster "+clusterId+": "+errString,
		)
		return
	}
}

func (r *SampleBucket) validateCreateSampleBucket(plan providerschema.SampleBucket) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return r.validateSampleBucketName(plan)
}

func (r *SampleBucket) validateSampleBucketName(plan providerschema.SampleBucket) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}

	if !isValidSampleName(plan.Name.ValueString()) {
		return errors.ErrInvalidSampleBucketName
	}

	return nil
}

// retrieveSampleBucket retrieves sample bucket information for a specified organization, project, cluster and sample bucket ID.
func (s *SampleBucket) retrieveSampleBucket(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.SampleBucket, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets/%s", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := s.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	sampleBucketResp := samplebucketapi.GetSampleBucketResponse{}
	err = json.Unmarshal(response.Body, &sampleBucketResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	var sampleStats providerschema.Stats
	if sampleBucketResp.Stats != nil {
		sampleStats = providerschema.NewStats(*sampleBucketResp.Stats)
	}

	sampleBucketStatsObj, diags := types.ObjectValueFrom(ctx, sampleStats.AttributeTypes(), sampleStats)
	if diags.HasError() {
		return nil, errors.ErrUnableToConvertAuditData
	}

	refreshedState := providerschema.SampleBucket{
		Id:                       types.StringValue(sampleBucketResp.Id),
		Name:                     types.StringValue(sampleBucketResp.Name),
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		Type:                     types.StringValue(sampleBucketResp.Type),
		StorageBackend:           types.StringValue(sampleBucketResp.StorageBackend),
		MemoryAllocationInMB:     types.Int64Value(sampleBucketResp.MemoryAllocationInMb),
		BucketConflictResolution: types.StringValue(sampleBucketResp.BucketConflictResolution),
		DurabilityLevel:          types.StringValue(sampleBucketResp.DurabilityLevel),
		Replicas:                 types.Int64Value(sampleBucketResp.Replicas),
		Flush:                    types.BoolValue(sampleBucketResp.Flush),
		TimeToLiveInSeconds:      types.Int64Value(sampleBucketResp.TimeToLiveInSeconds),
		EvictionPolicy:           types.StringValue(sampleBucketResp.EvictionPolicy),
		Stats:                    sampleBucketStatsObj,
	}

	return &refreshedState, nil
}

func isValidSampleName(category string) bool {
	switch category {
	case
		"travel-sample",
		"beer-sample",
		"gamesim-sample":
		return true
	}
	return false
}

// initializeBucketWithPlanAndId initializes an instance of providerschema.Bucket
// with the specified plan and ID. It marks all computed fields as null.
func initializeSampleBucketWithPlanAndId(plan providerschema.SampleBucket, id string) providerschema.SampleBucket {
	plan.Id = types.StringValue(id)

	plan.Type = types.StringNull()

	plan.StorageBackend = types.StringNull()

	plan.MemoryAllocationInMB = types.Int64Null()

	plan.BucketConflictResolution = types.StringNull()

	plan.DurabilityLevel = types.StringNull()

	plan.Replicas = types.Int64Null()

	plan.Flush = types.BoolNull()

	plan.TimeToLiveInSeconds = types.Int64Null()

	plan.EvictionPolicy = types.StringNull()

	plan.Stats = types.ObjectNull(providerschema.Stats{}.AttributeTypes())
	return plan
}
