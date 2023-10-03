package resources

import (
	"context"
	"fmt"
	providerschema "terraform-provider-capella/internal/schema"

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
	resp.Schema = ProjectSchema()
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
	// todo
}

// Read reads bucket information.
func (r *Bucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// todo
}

// Update updates the bucket.
func (r *Bucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo
}

// Delete deletes the bucket.
func (r *Bucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo
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
