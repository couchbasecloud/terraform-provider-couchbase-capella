package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppEndpointResync{}
	_ datasource.DataSourceWithConfigure = &AppEndpointResync{}
)

// Collections is the collections data source implementation.
type AppEndpointResync struct {
	*providerschema.Data
}

// NewCollections is a helper function to simplify the provider implementation.
func NewAppEndpointResync() datasource.DataSource {
	return &AppEndpointResync{}
}

// Metadata returns the collection data source type name.
func (a *AppEndpointResync) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_resync"
}

// Schema defines the schema for the collection data source.
func (a *AppEndpointResync) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointResyncSchema()
}

// Read refreshes the Terraform state with the latest data of collections.
func (a *AppEndpointResync) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpointResyncData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/resync",
		a.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
		config.AppEndpoint.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Endpoint Resync",
			fmt.Sprintf("Could not read resync status of app endpoint %s, unexpected error: %s", config.AppEndpoint.ValueString(), api.ParseError(err)),
		)
		return
	}

	var resyncResponse api.GetResyncResponse
	if err = json.Unmarshal(response.Body, &resyncResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing app endpoint resync response",
			"error during unmarshalling: "+err.Error(),
		)
		return
	}

	state := &providerschema.AppEndpointResyncData{
		OrganizationId: config.OrganizationId,
		ProjectId:      config.ProjectId,
		ClusterId:      config.ClusterId,
		AppServiceId:   config.AppServiceId,
		AppEndpoint:    config.AppEndpoint,
		DocsChanged:    types.Int64Value(resyncResponse.DocsChanged),
		DocsProcessed:  types.Int64Value(resyncResponse.DocsProcessed),
		LastError:      types.StringValue(resyncResponse.LastError),
		StartTime:      types.StringValue(resyncResponse.StartTime.Format("2006-01-02T15:04:05Z")),
		State:          types.StringValue(string(resyncResponse.State)),
	}

	if len(resyncResponse.CollectionsProcessing) > 0 {
		mapValue, diags := types.MapValueFrom(
			ctx,
			types.SetType{ElemType: types.StringType},
			resyncResponse.CollectionsProcessing,
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
