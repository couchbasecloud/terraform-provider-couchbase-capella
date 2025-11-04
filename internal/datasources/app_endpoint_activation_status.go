package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	appservice "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppEndpointActivationStatusDS)(nil)
	_ datasource.DataSourceWithConfigure = (*AppEndpointActivationStatusDS)(nil)
)

// AppEndpointActivationStatusDS fetches activation status (state) of an App Endpoint.
type AppEndpointActivationStatusDS struct {
	*providerschema.Data
}

// NewAppEndpointActivationStatus creates the datasource instance.
func NewAppEndpointActivationStatus() datasource.DataSource {
	return &AppEndpointActivationStatusDS{}
}

// Metadata sets the datasource name.
func (d *AppEndpointActivationStatusDS) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_activation_status"
}

// Schema describes required identifiers and computed state.
func (d *AppEndpointActivationStatusDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the activation status of an App Endpoint by calling get App Endpoint and returning its state.",
		Attributes: map[string]schema.Attribute{
			"organization_id":   schema.StringAttribute{Required: true, MarkdownDescription: "The GUID4 ID of the Capella organization."},
			"project_id":        schema.StringAttribute{Required: true, MarkdownDescription: "The GUID4 ID of the Capella project."},
			"cluster_id":        schema.StringAttribute{Required: true, MarkdownDescription: "The GUID4 ID of the Capella Cluster."},
			"app_service_id":    schema.StringAttribute{Required: true, MarkdownDescription: "The GUID4 ID of the App Service."},
			"app_endpoint_name": schema.StringAttribute{Required: true, MarkdownDescription: "The name of the App Endpoint."},
			"state":             schema.StringAttribute{Computed: true, MarkdownDescription: "The current activation state of the App Endpoint as reported by Capella."},
		},
	}
}

// Read performs GET App Endpoint and extracts state.
func (d *AppEndpointActivationStatusDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpointActivationStatus
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
		config.AppEndpointName.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Endpoint Activation Status",
			fmt.Sprintf("Could not read activation status for the App Endpoint: %s on App Service: %s: %s", config.AppEndpointName.ValueString(), config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var getResp appservice.GetAppEndpointStateResp
	if err := json.Unmarshal(response.Body, &getResp); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing app endpoint response",
			"error during unmarshalling: "+err.Error(),
		)
		return
	}

	state := &providerschema.AppEndpointActivationStatus{
		OrganizationId:  config.OrganizationId,
		ProjectId:       config.ProjectId,
		ClusterId:       config.ClusterId,
		AppServiceId:    config.AppServiceId,
		AppEndpointName: config.AppEndpointName,
		State:           types.StringValue(getResp.State),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Configure wires provider data.
func (d *AppEndpointActivationStatusDS) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
