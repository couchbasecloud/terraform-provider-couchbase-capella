package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &BucketBackupExport{}
	_ datasource.DataSourceWithConfigure = &BucketBackupExport{}
)

// BucketBackupExport is the data source implementation.
type BucketBackupExport struct {
	*providerschema.Data
}

// NewBucketBackupExport is a helper function to simplify the provider implementation.
func NewBucketBackupExport() datasource.DataSource {
	return &BucketBackupExport{}
}

// Metadata returns the bucket backup export data source type name.
func (b *BucketBackupExport) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_backup_export"
}

// Schema defines the schema for the bucket backup export data source.
func (b *BucketBackupExport) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = BucketBackupExportSchema()
}

// Configure adds the provider configured client to the bucket backup export data source.
func (b *BucketBackupExport) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	b.Data = data
}

// Read refreshes the Terraform state with the latest data of the backup export.
func (b *BucketBackupExport) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.BucketBackupExportData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := b.validate(state); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket Backup Export",
			"Could not read backup export "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	exportResp, err := b.getBucketBackupExport(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Bucket Backup Export",
			"Could not read backup export "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	state = mapBucketBackupExport(exportResp, state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (b *BucketBackupExport) getBucketBackupExport(ctx context.Context, state providerschema.BucketBackupExportData) (*backupapi.GetBucketBackupExportResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups/%s/exports/%s",
		b.HostURL,
		state.OrganizationId.ValueString(),
		state.ProjectId.ValueString(),
		state.ClusterId.ValueString(),
		state.BucketId.ValueString(),
		state.BackupId.ValueString(),
		state.Id.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := b.ClientV1.ExecuteWithRetry(ctx, cfg, nil, b.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	exportResponse := backupapi.GetBucketBackupExportResponse{}
	if err = json.Unmarshal(response.Body, &exportResponse); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &exportResponse, nil
}

// validate is used to verify that all the fields in the datasource have been populated.
func (b *BucketBackupExport) validate(state providerschema.BucketBackupExportData) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if state.BucketId.IsNull() {
		return errors.ErrBucketIdMissing
	}
	if state.BackupId.IsNull() {
		return errors.ErrBackupIdMissing
	}
	if state.Id.IsNull() {
		return errors.ErrIdMissing
	}
	return nil
}

// mapBucketBackupExport copies the response onto the configured state. The archive details are only
// set once the export completes.
func mapBucketBackupExport(
	export *backupapi.GetBucketBackupExportResponse,
	state providerschema.BucketBackupExportData,
) providerschema.BucketBackupExportData {
	state.CycleId = types.StringValue(export.CycleId)
	state.BucketName = types.StringValue(export.BucketName)
	state.Status = types.StringValue(export.Status)
	state.CreatedAt = types.StringValue(export.CreatedAt.Format(time.RFC3339))

	if export.SizeInBytes != nil {
		state.SizeInBytes = types.Int64Value(*export.SizeInBytes)
	}
	if export.Sha256Checksum != nil {
		state.Sha256Checksum = types.StringValue(*export.Sha256Checksum)
	}
	if export.Expiration != nil {
		state.Expiration = types.StringValue(export.Expiration.Format(time.RFC3339))
	}
	if export.BackupDownloadURL != nil {
		state.BackupDownloadURL = types.StringValue(*export.BackupDownloadURL)
	}

	return state
}
