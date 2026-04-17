package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppEndpointResync{}
	_ datasource.DataSourceWithConfigure = &AppEndpointResync{}
)

// AppEndpointResync is the AppEndpointResync data source implementation.
type AppEndpointResync struct {
	*providerschema.Data
}

// NewAppEndpointResync is a helper function to simplify the provider implementation.
func NewAppEndpointResync() datasource.DataSource {
	return &AppEndpointResync{}
}

// Metadata returns the AppEndpointResync data source type name.
func (a *AppEndpointResync) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_resync"
}

// Schema defines the schema for the AppEndpointResync data source.
func (a *AppEndpointResync) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointResyncSchema()
}

// Read reads the Resync information for an App Endpoint.
func (a *AppEndpointResync) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpointResyncData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId  = config.OrganizationId.ValueString()
		projectId       = config.ProjectId.ValueString()
		clusterId       = config.ClusterId.ValueString()
		appServiceId    = config.AppServiceId.ValueString()
		appEndpointName = config.AppEndpoint.ValueString()
	)

	uuids, err := utils.ParseUUIDs(
		utils.IDField{Name: "organization_id", Value: organizationId},
		utils.IDField{Name: "project_id", Value: projectId},
		utils.IDField{Name: "cluster_id", Value: clusterId},
		utils.IDField{Name: "app_service_id", Value: appServiceId},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing IDs",
			fmt.Sprintf("Could not parse resource IDs: %s", err.Error()),
		)
		return
	}
	organizationUUID, projectUUID, clusterUUID, appServiceUUID := uuids[0], uuids[1], uuids[2], uuids[3]

	getResyncStatusResp, err := a.ClientV2.GetAppEndpointResyncWithResponse(ctx, organizationUUID, projectUUID, clusterUUID, appServiceUUID, appEndpointName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Endpoint Resync",
			fmt.Sprintf("Could not read the resync status of app endpoint %s, unexpected error: %s", appEndpointName, err.Error()),
		)
		tflog.Debug(ctx, "error getting App Endpoint Resync status", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
			"err":             err.Error(),
		})
		return
	}

	if getResyncStatusResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected Status Reading App Endpoint Resync",
			"Could not read the resync status for app endpoint with name "+config.AppEndpoint.String()+", unexpected status: "+string(getResyncStatusResp.Body),
		)
		tflog.Debug(ctx, "unexpected status getting app endpoint resync status", map[string]interface{}{
			"organizationId":  organizationId,
			"projectId":       projectId,
			"clusterId":       clusterId,
			"appServiceId":    appServiceId,
			"appEndpointName": appEndpointName,
		})
		return
	}

	state := &providerschema.AppEndpointResyncData{
		OrganizationId: config.OrganizationId,
		ProjectId:      config.ProjectId,
		ClusterId:      config.ClusterId,
		AppServiceId:   config.AppServiceId,
		AppEndpoint:    config.AppEndpoint,
		DocsChanged:    types.Int64Value(int64(getResyncStatusResp.JSON200.DocsChanged)),
		DocsProcessed:  types.Int64Value(int64(getResyncStatusResp.JSON200.DocsProcessed)),
		LastError:      types.StringValue(getResyncStatusResp.JSON200.LastError),
		StartTime:      types.StringValue(getResyncStatusResp.JSON200.StartTime.Format("2006-01-02T15:04:05Z")),
		State:          types.StringValue(string(getResyncStatusResp.JSON200.State)),
	}

	state.CollectionsProcessing = types.MapNull(types.SetType{
		ElemType: types.StringType,
	})

	if getResyncStatusResp.JSON200.CollectionsProcessing != nil && len(*getResyncStatusResp.JSON200.CollectionsProcessing) > 0 {
		mapValue, diags := types.MapValueFrom(
			ctx,
			types.SetType{ElemType: types.StringType},
			getResyncStatusResp.JSON200.CollectionsProcessing,
		)
		if diags.HasError() {
			return
		}

		state.CollectionsProcessing = mapValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (a *AppEndpointResync) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
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

	a.Data = data
}
