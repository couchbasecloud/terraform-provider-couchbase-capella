package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	backupapi "terraform-provider-capella/internal/api/backup"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Backups{}
	_ datasource.DataSourceWithConfigure = &Backups{}
)

// Backups is the Backups data source implementation.
type Backups struct {
	*providerschema.Data
}

// NewBackups is a helper function to simplify the provider implementation.
func NewBackups() datasource.DataSource {
	return &Backups{}
}

// Metadata returns the backup data source type name.
func (d *Backups) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups"
}

// Schema defines the schema for the Backups data source.
func (d *Backups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = BackupSchema()
}

// Read refreshes the Terraform state with the latest data of backups.
func (d *Backups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Backups
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backups in Capella",
			"Could not read Capella Backups in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	// Get all the cycles
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/cycles", d.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.Client.Execute(
		cfg,
		nil,
		d.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Backup Cycles",
			fmt.Sprintf("Could not read backup cycles in a cluster %s, unexpected error: %s", clusterId, err.CompleteError()),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Backup Cycles",
			fmt.Sprintf("Could not read backup cycles in a cluster %s, unexpected error: %s", clusterId, err.Error()),
		)
		return
	}

	cyclesResp := backupapi.GetCyclesResponse{}
	err = json.Unmarshal(response.Body, &cyclesResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup cycles",
			"Could not read backup cycles, unexpected error: "+err.Error(),
		)
		return
	}

	// Loop through the cycles to fetch all backups within a cycle
	for _, cycle := range cyclesResp.Data {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/cycles/%s", d.HostURL, organizationId, projectId, clusterId, bucketId, cycle.CycleId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := d.Client.Execute(
			cfg,
			nil,
			d.Token,
			nil,
		)
		switch err := err.(type) {
		case nil:
		case api.Error:
			resp.Diagnostics.AddError(
				"Error Reading Capella Backups in a Cycle",
				fmt.Sprintf("Could not read backups for a cycle %s in a bucket %s, unexpected error: %s", cycle.CycleId, bucketId, err.CompleteError()),
			)
			return
		default:
			resp.Diagnostics.AddError(
				"Error Reading Capella Backups in a Cycle",
				fmt.Sprintf("Could not read backups for a cycle %s in a bucket %s, unexpected error: %s", cycle.CycleId, bucketId, err.Error()),
			)
			return
		}

		backupsResp := backupapi.GetBackupsResponse{}
		err = json.Unmarshal(response.Body, &backupsResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading backups",
				"Could not read backups, unexpected error: "+err.Error(),
			)
			return
		}

		for _, backup := range backupsResp.Data {
			backupStats := providerschema.NewBackupStats(*backup.BackupStats)
			backupStatsObj, diags := types.ObjectValueFrom(ctx, backupStats.AttributeTypes(), backupStats)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error Reading Backup Stats",
					fmt.Sprintf("Could not read backup stats data in a backup record, unexpected error: %s", fmt.Errorf("error while backup stats conversion")),
				)
				return
			}

			scheduleInfo := providerschema.NewScheduleInfo(*backup.ScheduleInfo)
			scheduleInfoObj, diags := types.ObjectValueFrom(ctx, scheduleInfo.AttributeTypes(), scheduleInfo)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error Error Reading Backup Schedule Info",
					fmt.Sprintf("Could not read backup schedule info in a backup record, unexpected error: %s", fmt.Errorf("error while backup schedule info conversion")),
				)
				return
			}

			newBackupData := providerschema.NewBackupData(&backup, organizationId, projectId, backupStatsObj, scheduleInfoObj)
			state.Data = append(state.Data, *newBackupData)
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the cluster data source.
func (d *Backups) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
