package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &FlushBucket{}
	_ resource.ResourceWithConfigure = &FlushBucket{}
)

const errorMessageFlushingBucket = "There is an error during execution of bucket flush. Please check in Capella to see if the documents for" +
	" have been deleted, unexpected error: "

// FlushBucket is the bucket resource implementation.
type FlushBucket struct {
	*providerschema.Data
}

// NewFlushBucket is a helper function to simplify the provider implementation.
func NewFlushBucket() resource.Resource {
	return &FlushBucket{}
}

// Metadata returns the flush Bucket resource type name.
func (c *FlushBucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flush"
}

// Schema defines the schema for the flush Bucket resource.
func (c *FlushBucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FlushBucketSchema()
}

// Create creates a new flush Bucket.
func (c *FlushBucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.FlushBucket
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	// Execute flush bucket. Nothing gets returned for it.
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/flush", c.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusOK}
	_, err := c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error flushing the bucket",
			errorMessageFlushingBucket+api.ParseError(err),
		)
		return
	}

	resp.State.Set(ctx, plan)
}

// Configure adds the provider configured api to the flush bucket resource.
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

func (c *FlushBucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Flush endpoint does not update the resource at any time other than the create.
}
