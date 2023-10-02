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
	_ resource.Resource                = &User{}
	_ resource.ResourceWithConfigure   = &User{}
	_ resource.ResourceWithImportState = &User{}
)

// User is the User resource implementation
type User struct {
	*providerschema.Data
}

func NewUsers() resource.Resource {
	return &User{}
}

// Metadata returns the users resource type name
func (r *User) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the allowlist resource.
func (r *User) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UserSchema()
}

// Configure sets provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (r *User) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create creates a new user
func (r *User) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// todo (AV-69624):
}

// Read reads user information
func (r *User) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// todo (AV-69625):
}

// Update updates the user
func (r *User) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo (AV-69626):
}

// Delete deletes the user
func (r *User) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo (AV-69627):
}

// ImportState imports a remote user that was not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=cluster123,project_id=proj123,organization_id=org123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *User) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
