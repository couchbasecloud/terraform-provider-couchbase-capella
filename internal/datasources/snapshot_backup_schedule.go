package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup_schedule"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SnapshotBackupSchedule{}
	_ datasource.DataSourceWithConfigure = &SnapshotBackupSchedule{}
)

// SnapshotBackupSchedule is the SnapshotBackupSchedule data source implementation.
type SnapshotBackupSchedule struct {
	*providerschema.Data
}

// NewSnapshotBackupSchedule is a helper function to simplify the provider implementation.
func NewSnapshotBackupSchedule() datasource.DataSource {
	return &SnapshotBackupSchedule{}
}

// Metadata returns the snapshot backup schedule data source type name.
func (d *SnapshotBackupSchedule) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_backup_schedule"
}

// Schema defines the schema for the SnapshotBackupSchedule data source.
func (d *SnapshotBackupSchedule) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotBackupScheduleSchema()
}

func (d *SnapshotBackupSchedule) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.SnapshotBackupSchedule
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationID.ValueString()
	projectId := state.ProjectID.ValueString()
	clusterId := state.ClusterID.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackupschedule", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	backupResps, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Snapshot Backup Schedule",
			fmt.Sprintf("Could not read snapshot backup schedule in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	var snapshotBackupSchedule snapshot_backup_schedule.SnapshotBackupSchedule
	err = json.Unmarshal(backupResps.Body, &snapshotBackupSchedule)
	if err != nil {
		diags.AddError(
			"Error Unmarshalling Capella Snapshot Backup Schedule",
			fmt.Sprintf("Could not unmarshal snapshot backup schedule in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		tflog.Debug(ctx, "error unmarshalling snapshot backups", map[string]interface{}{
			"backupResps.Body": backupResps.Body,
			"err":              err,
		})
		return
	}

	state = providerschema.NewSnapshotBackupSchedule(snapshotBackupSchedule, organizationId, projectId, clusterId)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the snapshot backup schedule data source.
func (d *SnapshotBackupSchedule) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
