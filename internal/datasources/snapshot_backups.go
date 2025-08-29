package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SnapshotBackups{}
	_ datasource.DataSourceWithConfigure = &SnapshotBackups{}
)

// SnapshotBackups is the SnapshotBackups data source implementation.
type SnapshotBackups struct {
	*providerschema.Data
}

// NewSnapshotBackups is a helper function to simplify the provider implementation.
func NewSnapshotBackups() datasource.DataSource {
	return &SnapshotBackups{}
}

// Metadata returns the snapshot backup data source type name.
func (d *SnapshotBackups) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_backups"
}

// Schema defines the schema for the SnapshotBackups data source.
func (d *SnapshotBackups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotBackupSchema()
}

func (d *SnapshotBackups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

// Configure adds the provider configured client to the snapshot backup data source.
func (d *SnapshotBackups) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.Data = data
}
