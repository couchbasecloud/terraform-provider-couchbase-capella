package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
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
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_backups"
}

// Schema defines the schema for the SnapshotBackups data source.
func (d *SnapshotBackups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotBackupsSchema()
}

func (d *SnapshotBackups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.SnapshotBackups
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", d.HostURL, organizationId, projectId, clusterId)
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
			"Error Reading Capella Snapshot Backups",
			fmt.Sprintf("Could not read snapshot backups in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	var snapshotBackups snapshot_backup.ListSnapshotBackupsResponse
	err = json.Unmarshal(backupResps.Body, &snapshotBackups)
	if err != nil {
		diags.AddError(
			"Error Unmarshalling Capella Snapshot Backups",
			fmt.Sprintf("Could not unmarshal snapshot backups in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		tflog.Debug(ctx, "error unmarshalling snapshot backups", map[string]interface{}{
			"backupResps.Body": backupResps.Body,
			"err":              err,
		})
		return
	}

	for i := range snapshotBackups.Data {
		snapshotBackup := snapshotBackups.Data[i]

		progress := providerschema.NewProgress(snapshotBackup.Progress)
		progressObj, diags := types.ObjectValueFrom(ctx, progress.AttributeTypes(), progress)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during progress conversion",
				fmt.Sprintf("Could not convert progress to object: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during progress conversion", map[string]interface{}{
				"snapshotBackup": snapshotBackup,
				"progress":       progress,
			})
			return
		}

		server := providerschema.NewServer(snapshotBackup.Server)
		serverObj, diags := types.ObjectValueFrom(ctx, server.AttributeTypes(), server)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during server conversion",
				fmt.Sprintf("Could not convert server to object: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during server conversion", map[string]interface{}{
				"snapshotBackup": snapshotBackup,
				"server":         server,
			})
			return
		}

		cmekObjs := make([]basetypes.ObjectValue, 0)
		for _, cmek := range snapshotBackup.CMEK {
			cmek := providerschema.NewCMEK(cmek)
			cmekObj, diags := types.ObjectValueFrom(ctx, cmek.AttributeTypes(), cmek)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error during CMEK conversion",
					fmt.Sprintf("Could not convert CMEK to object: %s", diags.Errors()),
				)
				tflog.Debug(ctx, "error during CMEK conversion", map[string]interface{}{
					"snapshotBackup": snapshotBackup,
					"cmek":           cmek,
				})
				return
			}
			cmekObjs = append(cmekObjs, cmekObj)
		}

		cmekSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: providerschema.CMEK{}.AttributeTypes()}, cmekObjs)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during CMEK conversion",
				fmt.Sprintf("Could not convert CMEK to set: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during CMEK conversion", map[string]interface{}{
				"snapshotBackup": snapshotBackup,
				"cmekObjs":       cmekObjs,
			})
			return
		}

		crossRegionCopyObjs := make([]basetypes.ObjectValue, 0)
		for _, region := range snapshotBackup.CrossRegionCopies {
			crossRegionCopy := providerschema.NewCrossRegionCopy(region)
			crossRegionCopyObj, diags := types.ObjectValueFrom(ctx, crossRegionCopy.AttributeTypes(), crossRegionCopy)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error during cross region copy conversion",
					fmt.Sprintf("Could not convert cross region copy to object: %s", diags.Errors()),
				)
				tflog.Debug(ctx, "error during cross region copy conversion", map[string]interface{}{
					"snapshotBackup":  snapshotBackup,
					"crossRegionCopy": region,
				})
				return
			}
			crossRegionCopyObjs = append(crossRegionCopyObjs, crossRegionCopyObj)
		}

		crossRegionCopySet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: providerschema.CrossRegionCopy{}.AttributeTypes()}, crossRegionCopyObjs)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during cross region copy conversion",
				fmt.Sprintf("Could not convert cross region copy to set: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during cross region copy conversion", map[string]interface{}{
				"snapshotBackup":      snapshotBackup,
				"crossRegionCopyObjs": crossRegionCopyObjs,
			})
			return
		}

		newSnapshotBackupsData := providerschema.NewSnapshotBackupsData(snapshotBackup, snapshotBackup.ID, clusterId, projectId, organizationId, progressObj, serverObj, cmekSet, crossRegionCopySet)
		state.Data = append(state.Data, newSnapshotBackupsData)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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
