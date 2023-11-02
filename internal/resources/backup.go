package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "terraform-provider-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Backup{}
	_ resource.ResourceWithConfigure   = &Backup{}
	_ resource.ResourceWithImportState = &Backup{}
)

// Backup is the Backup resource implementation.
type Backup struct {
	*providerschema.Data
}

// NewBackup is a helper function to simplify the provider implementation.
func NewBackup() resource.Resource {
	return &Backup{}
}

func (b Backup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (b Backup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupSchema()
}

func (b Backup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b Backup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (b Backup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b Backup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (b Backup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b Backup) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}
