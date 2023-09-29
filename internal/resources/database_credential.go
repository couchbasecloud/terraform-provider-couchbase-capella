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
	_ resource.Resource                = &DatabaseCredential{}
	_ resource.ResourceWithConfigure   = &DatabaseCredential{}
	_ resource.ResourceWithImportState = &DatabaseCredential{}
)

// DatabaseCredential is the database credential resource implementation.
type DatabaseCredential struct {
	*providerschema.Data
}

func NewDatabaseCredential() resource.Resource {
	return &DatabaseCredential{}
}

// Metadata returns the name that the database credential will follow in the terraform files.
// the name as per this function is capella_database_credential.
func (r *DatabaseCredential) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_credential"
}

// Schema defines the schema for the database credential resource.
func (r *DatabaseCredential) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseCredentialSchema()
}

// Configure adds the provider configured client to the database credential resource.
func (r *DatabaseCredential) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new database credential.
func (r *DatabaseCredential) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// todo in AV-61729
}

// Read reads database credential information.
func (r *DatabaseCredential) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// todo in AV-61729
}

// Update updates the database credential.
func (r *DatabaseCredential) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo in AV-62853
}

// Delete deletes the database credential.
func (r *DatabaseCredential) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo in AV-62166
}

// ImportState imports a remote database credential that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=user123,organization_id=org123,project_id=proj123,cluster_id=cluster123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *DatabaseCredential) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
