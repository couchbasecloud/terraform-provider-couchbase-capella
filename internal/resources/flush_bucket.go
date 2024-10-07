package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &FlushBucket{}
	_ resource.ResourceWithConfigure   = &FlushBucket{}
	_ resource.ResourceWithImportState = &FlushBucket{}
)

// Bucket is the bucket resource implementation.
type FlushBucket struct {
	*providerschema.Data
}

// NewBucket is a helper function to simplify the provider implementation.
func NewFlushBucket() resource.Resource {
	return &FlushBucket{}
}

// Metadata returns the Bucket resource type name.
func (c *FlushBucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flush"
}

// Schema defines the schema for the Bucket resource.
func (c *FlushBucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FlushBucketSchema()
}

// Create creates a new Bucket.
func (c *FlushBucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.FlushBucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := c.validateFlushBucketRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error executing bucket flush",
			"Could not flush the bucket, unexpected error: "+err.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	// Execute flush bucke. Nothing gets returned for it.
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/flush", c.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusOK}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error flushing the bucket",
			errorMessageWhileBucketCreation+api.ParseError(err),
		)
		return
	}

	resp.State.Set(ctx, providerschema.FlushBucket{
		BucketId:       types.StringValue(bucketId),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
	})
}

// Configure adds the provider configured api to the project resource.
func (c *FlushBucket) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (c *FlushBucket) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Couchbase Capella's v4 does not support a READ operation for flush endpoint.
}

func (r *FlushBucket) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Flush endpoint does not support delete as there is nothing on Capella to delete for this resource.
}

func (c *FlushBucket) ImportState(_ context.Context, _ resource.ImportStateRequest, _ *resource.ImportStateResponse) {
	// Flush endpoint is not a managed resource on capella. It is purely managed by terraform.
}

// Flushes the bucket. Only updates, if we change the flush bucket resource. Can't re-execute the flush on the same bucket.
func (c *FlushBucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.FlushBucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := c.validateFlushBucketRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error executing bucket flush",
			"Could not flush the bucket, unexpected error: "+err.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	// Execute flush bucket. Nothing gets returned for it.
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/flush", c.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusOK}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error flushing the bucket",
			errorMessageWhileBucketCreation+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(diags...)
	resp.State.Set(ctx, providerschema.FlushBucket{
		BucketId:       types.StringValue(bucketId),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
	})
}

func (r *FlushBucket) validateFlushBucketRequest(plan providerschema.FlushBucket) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.BucketId.IsNull() {
		return errors.ErrBucketIdMissing
	}
	return nil
}
