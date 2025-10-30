package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	_ datasource.DataSource              = &ProjectSnapshotBackups{}
	_ datasource.DataSourceWithConfigure = &ProjectSnapshotBackups{}
)

// ProjectSnapshotBackups is the ProjectSnapshotBackups data source implementation.
type ProjectSnapshotBackups struct {
	*providerschema.Data
}

// NewProjectSnapshotBackups is a helper function to simplify the provider implementation.
func NewProjectSnapshotBackups() datasource.DataSource {
	return &ProjectSnapshotBackups{}
}

// Metadata returns the project snapshot backup data source type name.
func (d *ProjectSnapshotBackups) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_snapshot_backups"
}

// Schema defines the schema for the ProjectSnapshotBackups data source.
func (d *ProjectSnapshotBackups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectSnapshotBackupSchema()
}

func (d *ProjectSnapshotBackups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ProjectSnapshotBackups
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()

	queryParam := buildQueryParams(&state)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/cloudsnapshotbackups", d.HostURL, organizationId, projectId)
	if len(queryParam) > 0 {
		url = url + BuildQueryParams(queryParam)
	}

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
			"Error Reading Capella Project Snapshot Backups",
			fmt.Sprintf("Could not read snapshot backups in project %s, unexpected error: %s", projectId, api.ParseError(err)),
		)
		return
	}

	var ProjectSnapshotBackups snapshot_backup.ListProjectSnapshotBackupsResponse
	err = json.Unmarshal(backupResps.Body, &ProjectSnapshotBackups)
	if err != nil {
		diags.AddError(
			"Error Unmarshalling Capella Snapshot Backups",
			fmt.Sprintf("Could not unmarshal snapshot backups in project %s, unexpected error: %s", projectId, api.ParseError(err)),
		)
		tflog.Debug(ctx, "error unmarshalling snapshot backups", map[string]interface{}{
			"backupResps.Body": backupResps.Body,
			"err":              err,
		})
		return
	}

	for i := range ProjectSnapshotBackups.Data {
		data := ProjectSnapshotBackups.Data[i]
		mostRecentSnapshot := data.MostRecentSnapshot
		oldestSnapshot := ProjectSnapshotBackups.Data[i].OldestSnapshot

		newMostRecentSnapshot := morphProjectSnapshotBackup(ctx, resp, mostRecentSnapshot, projectId, organizationId)
		newOldestSnapshot := morphProjectSnapshotBackup(ctx, resp, oldestSnapshot, projectId, organizationId)
		if newOldestSnapshot == nil || newMostRecentSnapshot == nil {
			return
		}

		newProjectSnapshotBackups := providerschema.NewProjectSnapshotBackupData(
			data,
			*newMostRecentSnapshot,
			*newOldestSnapshot,
		)

		state.Data = append(state.Data, newProjectSnapshotBackups)
	}

	state.Cursor = &providerschema.Cursor{
		Hrefs: &providerschema.Hrefs{
			First:    types.StringValue(ProjectSnapshotBackups.Cursor.Hrefs.First),
			Last:     types.StringValue(ProjectSnapshotBackups.Cursor.Hrefs.Last),
			Next:     types.StringValue(ProjectSnapshotBackups.Cursor.Hrefs.Next),
			Previous: types.StringValue(ProjectSnapshotBackups.Cursor.Hrefs.Previous),
		},
		Pages: &providerschema.Pages{
			Last:       types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.Last)),
			Next:       types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.Next)),
			Page:       types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.Page)),
			PerPage:    types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.PerPage)),
			Previous:   types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.Previous)),
			TotalItems: types.Int64Value(int64(ProjectSnapshotBackups.Cursor.Pages.TotalItems)),
		},
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func morphProjectSnapshotBackup(ctx context.Context, resp *datasource.ReadResponse, projectSnapshotBackup snapshot_backup.ProjectSnapshot, projectId, tenantId string) *basetypes.ObjectValue {
	progress := providerschema.NewProgress(projectSnapshotBackup.Progress)
	progressObj, diags := types.ObjectValueFrom(ctx, progress.AttributeTypes(), progress)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error during progress conversion",
			fmt.Sprintf("Could not convert progress to object: %s", diags.Errors()),
		)
		tflog.Debug(ctx, "error during progress conversion", map[string]interface{}{
			"projectSnapshotBackup": projectSnapshotBackup,
			"progress":              progress,
		})
		return nil
	}

	server := providerschema.NewServer(projectSnapshotBackup.Server)
	serverObj, diags := types.ObjectValueFrom(ctx, server.AttributeTypes(), server)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error during server conversion",
			fmt.Sprintf("Could not convert server to object: %s", diags.Errors()),
		)
		tflog.Debug(ctx, "error during server conversion", map[string]interface{}{
			"projectSnapshotBackup": projectSnapshotBackup,
			"server":                server,
		})
		return nil
	}

	cmekObjs := make([]basetypes.ObjectValue, 0)
	for _, cmek := range projectSnapshotBackup.CMEK {
		cmek := providerschema.NewCMEK(cmek)
		cmekObj, diags := types.ObjectValueFrom(ctx, cmek.AttributeTypes(), cmek)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during CMEK conversion",
				fmt.Sprintf("Could not convert CMEK to object: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during CMEK conversion", map[string]interface{}{
				"projectSnapshotBackup": projectSnapshotBackup,
				"cmek":                  cmek,
			})
			return nil
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
			"projectSnapshotBackup": projectSnapshotBackup,
			"cmekObjs":              cmekObjs,
		})
		return nil
	}

	crossRegionCopyObjs := make([]basetypes.ObjectValue, 0)
	for _, region := range projectSnapshotBackup.CrossRegionCopies {
		crossRegionCopy := providerschema.NewCrossRegionCopy(region)
		crossRegionCopyObj, diags := types.ObjectValueFrom(ctx, crossRegionCopy.AttributeTypes(), crossRegionCopy)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error during cross region copy conversion",
				fmt.Sprintf("Could not convert cross region copy to object: %s", diags.Errors()),
			)
			tflog.Debug(ctx, "error during cross region copy conversion", map[string]interface{}{
				"projectSnapshotBackup": projectSnapshotBackup,
				"crossRegionCopy":       region,
			})
			return nil
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
			"projectSnapshotBackup": projectSnapshotBackup,
			"crossRegionCopyObjs":   crossRegionCopyObjs,
		})
		return nil
	}

	projectSnapshot := providerschema.NewProjectSnapshot(
		projectSnapshotBackup,
		tenantId,
		projectId,
		progressObj,
		serverObj,
		cmekSet,
		crossRegionCopySet,
	)

	projectSnapshotFormated, diags := types.ObjectValueFrom(ctx, providerschema.ProjectSnapshot{}.AttributeTypes(), projectSnapshot)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error during project snapshot backup conversion",
			fmt.Sprintf("Could not convert project snapshot backup to set: %s", diags.Errors()),
		)
		tflog.Debug(ctx, "error during  project snapshot backup conversion", map[string]interface{}{
			"projectSnapshotBackup": projectSnapshotFormated,
		})
		return nil
	}
	return &projectSnapshotFormated
}

func buildQueryParams(state *providerschema.ProjectSnapshotBackups) map[string][]string {
	queryParam := make(map[string][]string)
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		page := int(state.Page.ValueInt64())
		queryParam["page"] = []string{strconv.Itoa(page)}
	}
	if !state.PerPage.IsNull() && !state.PerPage.IsUnknown() {
		perPage := int(state.PerPage.ValueInt64())
		queryParam["perPage"] = []string{strconv.Itoa(perPage)}
	}
	if !state.SortBy.IsNull() && !state.SortBy.IsUnknown() {
		sortBy := state.SortBy.ValueString()
		queryParam["sortBy"] = []string{sortBy}
	}
	if !state.SortDirection.IsNull() && !state.SortDirection.IsUnknown() {
		sortDir := state.SortDirection.ValueString()
		queryParam["sortDirection"] = []string{sortDir}
	}
	return queryParam
}

// Configure adds the provider configured client to the snapshot backup data source.
func (d *ProjectSnapshotBackups) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
