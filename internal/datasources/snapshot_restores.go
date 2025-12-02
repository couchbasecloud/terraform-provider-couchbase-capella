package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource                   = &SnapshotRestores{}
	_ datasource.DataSourceWithConfigure      = &SnapshotRestores{}
	_ datasource.DataSourceWithValidateConfig = &SnapshotRestores{}
)

// SnapshotRestores is the SnapshotRestores data source implementation.
type SnapshotRestores struct {
	*providerschema.Data
}

// NewSnapshotRestores is a helper function to simplify the provider implementation.
func NewSnapshotRestores() datasource.DataSource {
	return &SnapshotRestores{}
}

// Metadata returns the snapshot restores data source type name.
func (d *SnapshotRestores) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot_restores"
}

// Schema defines the schema for the SnapshotRestores data source.
func (d *SnapshotRestores) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotRestoresSchema()
}

func (d *SnapshotRestores) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.SnapshotRestores
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationID.ValueString()
	projectId := state.ProjectID.ValueString()
	clusterId := state.ClusterID.ValueString()

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
	var names []string

	// Since the list API doesn't implement query parameters useful for filtering,
	// filtering is done by provider.
	if state.Filters != nil {
		diags := state.Filters.Values.ElementsAs(ctx, &names, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	for i := range snapshotRestores.Data {
		snapshotRestore := snapshotRestores.Data[i]

		if slices.Contains(names, string(snapshotRestore.Status)) || len(names) == 0 {
			newSnapshotRestoreData := providerschema.NewSnapshotRestoreData(snapshotRestore, clusterId, projectId, organizationId)
			state.Data = append(state.Data, newSnapshotRestoreData)
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the snapshot restores data source.
func (d *SnapshotRestores) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ValidateConfig checks that if 'name' or 'values' is set in filter block', then both are set.
func (d *SnapshotRestores) ValidateConfig(
	ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse,
) {
	var config providerschema.SnapshotRestores
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Filters != nil {
		if (config.Filters.Name.IsNull() && !config.Filters.Values.IsNull()) ||
			(!config.Filters.Name.IsNull() && config.Filters.Values.IsNull()) {
			resp.Diagnostics.AddError(
				"Invalid Filters Configuration",
				"Both 'name' and 'values' in filter block must be configured.",
			)
		}
	}
}
