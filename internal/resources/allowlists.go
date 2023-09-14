package resources

import (
	"context"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// AllowList is the AllowList resource implementation.
type AllowList struct {
	*providerschema.Data
}

func NewAllowList() resource.Resource {
	return &AllowList{}
}

// Metadata returns the project resource type name.
func (r *AllowList) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allowlist"
}

// Schema defines the schema for the project resource.
func (r *AllowList) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

}

func (r *AllowList) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

}

func (r *AllowList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

// Read reads project information.
func (r *AllowList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// todo
}

// Update updates the project.
func (r *AllowList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// todo
}

// Delete deletes the project.
func (r *AllowList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// todo
}

// ImportState imports a remote allowlist that is not created by Terraform.
func (r *AllowList) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AllowList) retrieveAllowList(ctx context.Context, organizationId, projectId string) (*providerschema.OneAllowList, error) {
	return nil, nil
}
