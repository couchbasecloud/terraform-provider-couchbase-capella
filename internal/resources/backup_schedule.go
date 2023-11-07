package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "terraform-provider-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BackupSchedule{}
	_ resource.ResourceWithConfigure   = &BackupSchedule{}
	_ resource.ResourceWithImportState = &BackupSchedule{}
)

// BackupSchedule is the BackupSchedule resource implementation.
type BackupSchedule struct {
	*providerschema.Data
}

// NewBackupSchedule is a helper function to simplify the provider implementation.
func NewBackupSchedule() resource.Resource {
	return &BackupSchedule{}
}

func (b BackupSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_schedule"

}

func (b BackupSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupScheduleSchema()
}

func (b BackupSchedule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
}

func (b BackupSchedule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
}

func (b BackupSchedule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
}

func (b BackupSchedule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
}

func (b BackupSchedule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
}

func (b BackupSchedule) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	//TODO implement me
}
