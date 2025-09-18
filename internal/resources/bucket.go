package resources

import (
	"context"
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Bucket{}
	_ resource.ResourceWithConfigure   = &Bucket{}
	_ resource.ResourceWithImportState = &Bucket{}
)

const errorMessageAfterBucketCreation = "Bucket creation is successful, but encountered an error while checking the current" +
	" state of the bucket. Please run `terraform plan` after 1-2 minutes to know the" +
	" current bucket state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileBucketCreation = "There is an error during bucket creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// Bucket is the bucket resource implementation.
type Bucket struct {
	*providerschema.Data
}

// NewBucket is a helper function to simplify the provider implementation.
func NewBucket() resource.Resource {
	return &Bucket{}
}

// Metadata returns the Bucket resource type name.
func (c *Bucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Schema defines the schema for the Bucket resource.
func (c *Bucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BucketSchema()
}

// Create creates a new Bucket.
func (c *Bucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Bucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	BucketRequest := apigen.CreateBucketRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.StorageBackend.IsNull() && !plan.StorageBackend.IsUnknown() {
		v := apigen.StorageBackend(plan.StorageBackend.ValueString())
		BucketRequest.StorageBackend = &v
	}
	if !plan.MemoryAllocationInMB.IsNull() && !plan.MemoryAllocationInMB.IsUnknown() {
		v := int(plan.MemoryAllocationInMB.ValueInt64())
		BucketRequest.MemoryAllocationInMb = &v
	}
	if !plan.BucketConflictResolution.IsNull() && !plan.BucketConflictResolution.IsUnknown() {
		v := apigen.BucketConflictResolution(plan.BucketConflictResolution.ValueString())
		BucketRequest.BucketConflictResolution = &v
	}
	if !plan.DurabilityLevel.IsNull() && !plan.DurabilityLevel.IsUnknown() {
		v := apigen.DurabilityLevel(plan.DurabilityLevel.ValueString())
		BucketRequest.DurabilityLevel = &v
	}
	if !plan.Replicas.IsNull() && !plan.Replicas.IsUnknown() {
		v := apigen.Replicas(plan.Replicas.ValueInt64())
		BucketRequest.Replicas = &v
	}
	if !plan.Flush.IsNull() && !plan.Flush.IsUnknown() {
		v := plan.Flush.ValueBool()
		BucketRequest.Flush = &v
		BucketRequest.FlushEnabled = &v
	}
	if !plan.TimeToLiveInSeconds.IsNull() && !plan.TimeToLiveInSeconds.IsUnknown() {
		v := int(plan.TimeToLiveInSeconds.ValueInt64())
		BucketRequest.TimeToLiveInSeconds = &v
	}
	if !plan.EvictionPolicy.IsNull() && !plan.EvictionPolicy.IsUnknown() {
		v := apigen.EvictionPolicy(plan.EvictionPolicy.ValueString())
		BucketRequest.EvictionPolicy = &v
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		v := apigen.Type(plan.Type.ValueString())
		BucketRequest.Type = &v
	}

	if err := c.validateCreateBucket(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	orgUUID, _ := uuid.Parse(plan.OrganizationId.ValueString())
	projUUID, _ := uuid.Parse(plan.ProjectId.ValueString())
	cluUUID, _ := uuid.Parse(plan.ClusterId.ValueString())

	res, err := c.ClientV2.PostBucketWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), BucketRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			errorMessageWhileBucketCreation+api.ParseError(err),
		)
		return
	}
	if res.JSON201 == nil {
		resp.Diagnostics.AddError("Error creating bucket", "unexpected status: "+res.Status())
		return
	}

	diags = resp.State.Set(ctx, initializeBucketWithPlanAndId(plan, res.JSON201.Id))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveBucket(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), res.JSON201.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating bucket",
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

// Configure It adds the provider configured api to the project resource.
func (c *Bucket) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	c.Data = data
}

// Read reads the bucket information.
func (c *Bucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	refreshedState, err := c.retrieveBucket(ctx, organizationId, projectId, clusterId, bucketId)
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

// Delete deletes the bucket.
func (r *Bucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.Bucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting bucket",
			"Could not delete bucket, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.Id]
	)

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)

	_, err = r.ClientV2.DeleteBucketByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(bucketId))
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

// ImportState imports a remote cluster that is not created by Terraform.
func (c *Bucket) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// retrieveBucket retrieves bucket information for a specified organization, project, cluster and bucket ID.
func (c *Bucket) retrieveBucket(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.OneBucket, error) {
	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)
	bktUUID := bucketId // bucketId is string in spec

	res, err := c.ClientV2.GetBucketByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(bktUUID))
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
	if bucketResp.Flush != nil {
		refreshedState.Flush = types.BoolValue(*bucketResp.Flush)
	} else {
		refreshedState.Flush = types.BoolNull()
	}
	refreshedState.Stats = &providerschema.Stats{
		ItemCount:       types.Int64Value(int64(bucketResp.Stats.ItemCount)),
		OpsPerSecond:    types.Int64Value(int64(bucketResp.Stats.OpsPerSecond)),
		DiskUsedInMiB:   types.Int64Value(int64(bucketResp.Stats.DiskUsedInMib)),
		MemoryUsedInMiB: types.Int64Value(int64(bucketResp.Stats.MemoryUsedInMib)),
	}

	return &refreshedState, nil
}

// Update updates the bucket.
func (c *Bucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state providerschema.Bucket
	diags := req.Plan.Get(ctx, &state)
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

	bucketUpdateRequest := apigen.UpdateBucketRequest{
		MemoryAllocationInMb: int(state.MemoryAllocationInMB.ValueInt64()),
		DurabilityLevel:      apigen.UpdateBucketRequestDurabilityLevel(state.DurabilityLevel.ValueString()),
		Replicas:             apigen.UpdateBucketRequestReplicas(state.Replicas.ValueInt64()),
		TimeToLiveInSeconds:  int(state.TimeToLiveInSeconds.ValueInt64()),
	}
	if !state.Flush.IsNull() && !state.Flush.IsUnknown() {
		v := state.Flush.ValueBool()
		bucketUpdateRequest.Flush = &v
	}

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)

	_, err = c.ClientV2.PutBucketWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(bucketId), bucketUpdateRequest)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update bucket, unexpected error: "+bucketId+": "+errString,
		)
		return
	}

	currentState, err := c.retrieveBucket(ctx, organizationId, projectId, clusterId, bucketId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update Capella bucket with ID "+bucketId+": "+api.ParseError(err),
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

// initializeBucketWithPlanAndId initializes an instance of providerschema.Bucket
// with the specified plan and ID. It marks all computed fields as null.
func initializeBucketWithPlanAndId(plan providerschema.Bucket, id string) providerschema.Bucket {
	plan.Id = types.StringValue(id)
	if plan.StorageBackend.IsNull() || plan.StorageBackend.IsUnknown() {
		plan.StorageBackend = types.StringNull()
	}
	if plan.EvictionPolicy.IsNull() || plan.EvictionPolicy.IsUnknown() {
		plan.EvictionPolicy = types.StringNull()
	}
	plan.Stats = types.ObjectNull(providerschema.Stats{}.AttributeTypes())
	return plan
}

func (r *Bucket) validateCreateBucket(plan providerschema.Bucket) error {
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

func (r *Bucket) validateBucketAttributesTrimmed(plan providerschema.Bucket) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	return nil
}
