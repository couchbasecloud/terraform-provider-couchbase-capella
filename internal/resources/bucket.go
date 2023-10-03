package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func NewBucket() resource.Resource {
	return &Bucket{}
}

// Metadata returns the bucket resource type name.
func (r *Bucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Schema defines the schema for the bucket resource.
func (r *Bucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BucketSchema()
}

// Configure adds the provider configured client to the bucket resource.
func (r *Bucket) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Data = data
}

// Create creates a new bucket.
func (r *Bucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Bucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	if plan.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = plan.ClusterId.ValueString()

	bucketRequest := api.CreateBucketRequest{
		Name:                     plan.Name.ValueString(),
		Type:                     plan.Type.ValueString(),
		StorageBackend:           plan.StorageBackend.ValueString(),
		MemoryAllocationInMb:     plan.MemoryAllocationInMB.ValueInt64(),
		BucketConflictResolution: plan.BucketConflictResolution.ValueString(),
		DurabilityLevel:          plan.DurabilityLevel.ValueString(),
		Replicas:                 plan.Replicas.ValueInt64(),
		Flush:                    plan.Flush.ValueBool(),
		TimeToLiveInSeconds:      plan.TimeToLiveInSeconds.ValueInt64(),
	}

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", r.HostURL, organizationId, projectId, clusterId),
		http.MethodPost,
		bucketRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	bucketResponse := api.CreateBucketResponse{}
	err = json.Unmarshal(response.Body, &bucketResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveBucket(ctx, organizationId, projectId, clusterId, bucketResponse.Id)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Buckets",
			"Could not read Capella bucket with ID "+bucketResponse.Id+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Buckets",
			"Could not read Capella bucket with ID "+bucketResponse.Id+": "+err.Error(),
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

// Read reads bucket information.
func (r *Bucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Bucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bktId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Buckets in Capella",
			"Could not read Capella bucket with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	// Get refreshed Bucket value from Capella
	refreshedState, err := r.retrieveBucket(ctx, organizationId, projectId, clusterId, bktId)
	resourceNotFound, err := handleBucketError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading buckets",
			"Could not read bucket with id "+state.Id.String()+": "+err.Error(),
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

// Update updates the bucket.
func (r *Bucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state providerschema.Bucket
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bktId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Buckets in Capella",
			"Could not read Capella Bucket with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	bktCredRequest := api.PutBucketRequest{
		MemoryAllocationInMb: state.MemoryAllocationInMB.ValueInt64(),
		DurabilityLevel:      state.DurabilityLevel.ValueString(),
		Replicas:             state.Replicas.ValueInt64(),
		Flush:                state.Flush.ValueBool(),
		TimeToLiveInSeconds:  state.TimeToLiveInSeconds.ValueInt64(),
	}

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", r.HostURL, organizationId, projectId, clusterId, bktId),
		http.MethodPut,
		bktCredRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update an existing bucket, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update bucket, unexpected error: "+err.Error(),
		)
		return
	}

	currentState, err := r.retrieveBucket(ctx, organizationId, projectId, clusterId, bktId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket",
			"Could not read Capella Bucket with ID "+bktId+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket",
			"Could not read Capella Bucket with ID "+bktId+": "+err.Error(),
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

// Delete deletes the bucket.
func (r *Bucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.Bucket
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bktId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Buckets in Capella",
			"Could not read Capella Bucket with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", r.HostURL, organizationId, projectId, clusterId, bktId),
		http.MethodDelete,
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

// ImportState imports a remote bucket that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=proj123,organization_id=org123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *Bucket) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// retrieveBucket fetches the bucket by making a GET API call to the Capella V4 Public API.
// This usually helps retrieve the state of a newly created bucket that was created from Terraform.
func (b *Bucket) retrieveBucket(ctx context.Context, organizationId, projectId, clusterId, bucketId string) (*providerschema.OneBucket, error) {
	response, err := b.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", b.HostURL, organizationId, projectId, clusterId, bucketId),
		http.MethodGet,
		nil,
		b.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	bktResp := api.GetBucketResponse{}
	err = json.Unmarshal(response.Body, &bktResp)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.OneBucket{
		Id:                       types.StringValue(bktResp.Id),
		Name:                     types.StringValue(bktResp.Name),
		Type:                     types.StringValue(bktResp.Type),
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		StorageBackend:           types.StringValue(bktResp.StorageBackend),
		MemoryAllocationInMB:     types.Int64Value(bktResp.MemoryAllocationInMb),
		BucketConflictResolution: types.StringValue(bktResp.BucketConflictResolution),
		DurabilityLevel:          types.StringValue(bktResp.DurabilityLevel),
		Replicas:                 types.Int64Value(bktResp.Replicas),
		Flush:                    types.BoolValue(bktResp.Flush),
		TimeToLiveInSeconds:      types.Int64Value(bktResp.TimeToLiveInSeconds),
		EvictionPolicy:           types.StringValue(bktResp.EvictionPolicy),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(bktResp.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(bktResp.Audit.CreatedBy),
			ModifiedAt: types.StringValue(bktResp.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(bktResp.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(bktResp.Audit.Version)),
		},
		Stats: &providerschema.Stats{
			ItemCount:       types.Int64Value(bktResp.Stats.ItemCount),
			OpsPerSecond:    types.Int64Value(bktResp.Stats.OpsPerSecond),
			DiskUsedInMiB:   types.Int64Value(bktResp.Stats.DiskUsedInMiB),
			MemoryUsedInMiB: types.Int64Value(bktResp.Stats.MemoryUsedInMiB),
		},
	}
	return &refreshedState, nil
}

// this func extract error message if error is api.Error and also checks whether error is
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
