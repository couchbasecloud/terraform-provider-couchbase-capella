package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &AppServiceAzurePrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &AppServiceAzurePrivateEndpointCommand{}
)

// AppServiceAzurePrivateEndpointCommand is the data source implementation.
//
// NOTE: as of this writing, Couchbase Capella's private endpoints for App Services are
// AWS-only; Azure/GCP support is not confirmed GA (the generated OpenAPI client models the
// connection-command request body as an unresolved oneOf/anyOf union, consistent with the
// schema being written generically/forward-compatible even if only AWS is wired up in the
// backend today). This data source is included for API/schema parity with the cluster-level
// azure_private_endpoint_command resource, but has not been validated against a live
// Azure-backed App Service.
type AppServiceAzurePrivateEndpointCommand struct {
	*providerschema.Data
}

// NewAppServiceAzurePrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewAppServiceAzurePrivateEndpointCommand() datasource.DataSource {
	return &AppServiceAzurePrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *AppServiceAzurePrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_azure_private_endpoint_command"
}

// Schema defines the schema for the App Service private endpoint command data source.
func (a *AppServiceAzurePrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAzurePrivateEndpointCommandSchema()
}

// Read refreshes the Terraform state with the latest data of the App Service private endpoint command.
func (a *AppServiceAzurePrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceAzureCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateAppServiceAzureCommand(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service Azure private endpoint command request",
			"Could not validate App Service Azure private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	AzureCommandRequest := api.CreateAzurePrivateEndpointCommandRequest{
		ResourceGroupName: state.ResourceGroupName.ValueString(),
		VirtualNetwork:    state.VirtualNetwork.ValueString(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/privateEndpointCommand", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		AzureCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Azure private endpoint command",
			"Could not read App Service Azure private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var AzureCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &AzureCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling App Service Azure private endpoint command response",
			"Could not unmarshall App Service Azure private endpoint command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(AzureCommandResponse.Command)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint command data source.
func (a *AppServiceAzurePrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}

// validateAppServiceAzureCommand ensures organization id, project id, cluster id, app service
// id, virtual network, and resource group are valued.
func validateAppServiceAzureCommand(config providerschema.AppServiceAzureCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}
	if config.VirtualNetwork.IsNull() {
		return errors.ErrVirtualNetworkMissing
	}
	if config.ResourceGroupName.IsNull() {
		return errors.ErrResourceGroupName
	}

	return nil
}
