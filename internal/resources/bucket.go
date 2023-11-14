package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	bucketapi "terraform-provider-capella/internal/api/bucket"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

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

	BucketRequest := bucketapi.CreateBucketRequest{
		Name: plan.Name.ValueString(),
	}

	// Check for optional fields
	if !plan.StorageBackend.IsNull() && !plan.StorageBackend.IsUnknown() {
		BucketRequest.StorageBackend = plan.StorageBackend.ValueStringPointer()
	}

	if !plan.MemoryAllocationInMB.IsNull() && !plan.MemoryAllocationInMB.IsUnknown() {
		BucketRequest.MemoryAllocationInMb = plan.MemoryAllocationInMB.ValueInt64Pointer()
	}

	if !plan.BucketConflictResolution.IsNull() && !plan.BucketConflictResolution.IsUnknown() {
		BucketRequest.BucketConflictResolution = plan.BucketConflictResolution.ValueStringPointer()
	}

	if !plan.DurabilityLevel.IsNull() && !plan.DurabilityLevel.IsUnknown() {
		BucketRequest.DurabilityLevel = plan.DurabilityLevel.ValueStringPointer()
	}

	if !plan.Replicas.IsNull() && !plan.Replicas.IsUnknown() {
		BucketRequest.Replicas = plan.Replicas.ValueInt64Pointer()
	}

	if !plan.Flush.IsNull() && !plan.Flush.IsUnknown() {
		BucketRequest.Flush = plan.Flush.ValueBoolPointer()
	}

	if !plan.TimeToLiveInSeconds.IsNull() && !plan.TimeToLiveInSeconds.IsUnknown() {
		BucketRequest.TimeToLiveInSeconds = plan.TimeToLiveInSeconds.ValueInt64Pointer()
	}

	if !plan.EvictionPolicy.IsNull() && !plan.EvictionPolicy.IsUnknown() {
		BucketRequest.EvictionPolicy = plan.EvictionPolicy.ValueStringPointer()
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		BucketRequest.Type = plan.Type.ValueStringPointer()
	}

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrOrganizationIdCannotBeEmpty.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	if plan.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	if plan.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = plan.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := c.Client.Execute(
		cfg,
		BucketRequest,
		c.Token,
		nil,
	)
	_, err = handleBucketError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	BucketResponse := bucketapi.GetBucketResponse{}
	err = json.Unmarshal(response.Body, &BucketResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, error during unmarshalling:"+err.Error(),
		)
		return
	}

	refreshedState, err := c.retrieveBucket(ctx, organizationId, projectId, clusterId, BucketResponse.Id)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket",
			"Could not read Capella bucket with ID "+BucketResponse.Id+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket",
			"Could not read Capella bucket with ID "+BucketResponse.Id+": "+err.Error(),
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
	resourceNotFound, err := handleBucketError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading bucket",
			"Could not read bucket with id "+state.Id.String()+": "+err.Error(),
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
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = state.ProjectId.ValueString()

	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = state.ClusterId.ValueString()

	if state.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var bucketId = state.Id.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", r.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := r.Client.Execute(
		cfg,
		nil,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Deleting the Bucket",
				"Could not delete Bucket associated with cluster "+clusterId+": "+err.CompleteError(),
			)
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Deleting Bucket",
			"Could not delete Bucket associated with cluster "+clusterId+": "+err.Error(),
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
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", c.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.Client.Execute(
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	bucketResp := bucketapi.GetBucketResponse{}
	err = json.Unmarshal(response.Body, &bucketResp)
	if err != nil {
		return nil, err
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
			ItemCount:       types.Int64Value(int64(bucketResp.Stats.ItemCount)),
			OpsPerSecond:    types.Int64Value(int64(bucketResp.Stats.OpsPerSecond)),
			DiskUsedInMiB:   types.Int64Value(int64(bucketResp.Stats.DiskUsedInMib)),
			MemoryUsedInMiB: types.Int64Value(int64(bucketResp.Stats.MemoryUsedInMib)),
		},
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

	bucketUpdateRequest := bucketapi.PutBucketRequest{
		MemoryAllocationInMb: state.MemoryAllocationInMB.ValueInt64(),
		DurabilityLevel:      state.DurabilityLevel.ValueString(),
		Replicas:             state.Replicas.ValueInt64(),
		Flush:                state.Flush.ValueBool(),
		TimeToLiveInSeconds:  state.TimeToLiveInSeconds.ValueInt64(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", c.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	response, err := c.Client.Execute(
		cfg,
		bucketUpdateRequest,
		c.Token,
		nil,
	)

	_, err = handleBucketError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update bucket, unexpected error: "+string(response.Body),
		)
		return
	}

	currentState, err := c.retrieveBucket(ctx, organizationId, projectId, clusterId, bucketId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update Capella bucket with ID "+bucketId+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update Capella bucket with ID "+bucketId+": "+err.Error(),
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

// handleBucketError extracts error message if error is api.Error and also checks whether error is
// resource not found
func handleBucketError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}
