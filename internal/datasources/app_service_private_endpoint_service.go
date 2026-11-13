package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServicePrivateEndpointService)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServicePrivateEndpointService)(nil)
)

// AppServicePrivateEndpointService is the App Service private endpoint service data source implementation.
type AppServicePrivateEndpointService struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpointService is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpointService() datasource.DataSource {
	return &AppServicePrivateEndpointService{}
}

// Metadata returns the data source type name.
func (d *AppServicePrivateEndpointService) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoint_service"
}

// Schema defines the schema for the data source.
func (d *AppServicePrivateEndpointService) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointServiceSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppServicePrivateEndpointService) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Data = data
}

// Read refreshes the Terraform state with the private endpoint service state and, if requested, the CLI command.
func (d *AppServicePrivateEndpointService) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppServicePrivateEndpointServiceData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoint Service",
			fmt.Sprintf("Could not read private endpoint service state for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
		)
		return
	}

	var stateResp apigen.GetAppServicePrivateEndpointStateResponse
	if err := json.Unmarshal(response.Body, &stateResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Private Endpoint Service",
			"Could not unmarshal private endpoint service state response: "+err.Error(),
		)
		return
	}

	config.State = types.StringPointerValue((*string)(stateResp.State))
	config.TargetState = types.StringPointerValue((*string)(stateResp.TargetState))
	config.Command = types.StringNull()

	if config.Aws != nil {
		commandRequest := apigen.CreateVPCEndpointCommandRequest{
			VpcID:     config.Aws.VpcID.ValueString(),
			SubnetIDs: *convertSubnetIDs(config.Aws.SubnetIDs),
		}

		cfg := api.EndpointCfg{Url: url + "/privateEndpointCommand", Method: http.MethodPost, SuccessStatus: http.StatusOK}
		response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, commandRequest, d.Token, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading App Service Private Endpoint Command",
				fmt.Sprintf("Could not read private endpoint command for App Service %s: %s", config.AppServiceId.ValueString(), api.ParseError(err)),
			)
			return
		}

		var commandResp apigen.CreatePrivateEndpointCommandResponse
		if err := json.Unmarshal(response.Body, &commandResp); err != nil {
			resp.Diagnostics.AddError(
				"Error Reading App Service Private Endpoint Command",
				"Could not unmarshal private endpoint command response: "+err.Error(),
			)
			return
		}

		config.Command = types.StringValue(commandResp.Command)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
