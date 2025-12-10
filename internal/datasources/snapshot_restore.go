package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SnapshotRestore{}
	_ datasource.DataSourceWithConfigure = &SnapshotRestore{}
)

// SnapshotRestores is the SnapshotRestores data source implementation.
type SnapshotRestore struct {
	*providerschema.Data
}

// NewSnapshotRestores is a helper function to simplify the provider implementation.
func NewSnapshotRestore() datasource.DataSource {
	return &SnapshotRestore{}
}

// Metadata returns the snapshot restore data source type name.
func (d *SnapshotRestore) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_restore"
}

// Schema defines the schema for the SnapshotRestores data source.
func (d *SnapshotRestore) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotRestoreSchema()
}

func (d *SnapshotRestore) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.SnapshotRestore
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationID.ValueString()
	projectId := state.ProjectID.ValueString()
	clusterId := state.ClusterID.ValueString()
	restoreId := state.ID.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups/restores", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	restoreResps, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Snapshot Restores",
			fmt.Sprintf("Could not read snapshot restores in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	var snapshotRestores snapshot_backup.ListSnapshotRestoresResponse
	err = json.Unmarshal(restoreResps.Body, &snapshotRestores)
	if err != nil {
		diags.AddError(
			"Error Unmarshalling Capella Snapshot Restores",
			fmt.Sprintf("Could not unmarshal snapshot restores in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		tflog.Debug(ctx, "error unmarshalling snapshot restores", map[string]interface{}{
			"snapshotRestoresResps.Body": restoreResps.Body,
			"err":                        err,
		})
		return
	}

	for i := range snapshotRestores.Data {
		if snapshotRestores.Data[i].ID == restoreId {
			newSnapshotRestore := providerschema.NewSnapshotRestore(snapshotRestores.Data[i], organizationId, projectId, clusterId)
			diags = resp.State.Set(ctx, &newSnapshotRestore)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.AddError(
		"Snapshot Restore Not Found",
		fmt.Sprintf("Could not find snapshot restore with ID %s in cluster %s", restoreId, clusterId),
	)
}

// Configure adds the provider configured client to the snapshot restore data source.
func (d *SnapshotRestore) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
