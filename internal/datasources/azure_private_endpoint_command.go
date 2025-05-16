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
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &AzurePrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &AzurePrivateEndpointCommand{}
)

// AzurePrivateEndpointCommand is the data source implementation.
type AzurePrivateEndpointCommand struct {
	*providerschema.Data
}

// NewAzurePrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewAzurePrivateEndpointCommand() datasource.DataSource {
	return &AzurePrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *AzurePrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_private_endpoint_command"
}

// Schema defines the schema for the private endpoint command data source.
func (a *AzurePrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source to generate an Azure CLI command for setting up a private endpoint connection to a Capella cluster. Retrieves the command or script to create the private endpoint, which establishes a private connection between the specified VPC and the designated Capella private endpoint service.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"resource_group_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of your Azure resource group.",
			},
			"virtual_network": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The virtual network and subnet name",
			},
			"command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The generated Azure CLI command that can be used to create the private endpoint connection within Azure.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of private endpoint command .
func (a *AzurePrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AzureCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateAzureCommand(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Azure private endpoint command request",
			"Could not validate Azure private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	AzureCommandRequest := api.CreateAzurePrivateEndpointCommandRequest{
		ResourceGroupName: state.ResourceGroupName.ValueString(),
		VirtualNetwork:    state.VirtualNetwork.ValueString(),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpointCommand", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		AzureCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Azure private endpoint command",
			"Could not read Azure private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var AzureCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &AzureCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling Azure private endpoint command response",
			"Could not unmarshall Azure private endpoint command response, unexpected error: "+err.Error(),
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

// Configure adds the provider configured client to the private endpoint command data source.
func (a *AzurePrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validateAzureCommand ensures organization id, project id, cluster id, virtual network and resource group are valued.
func validateAzureCommand(config providerschema.AzureCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.VirtualNetwork.IsNull() {
		return errors.ErrVirtualNetworkMissing
	}
	if config.ResourceGroupName.IsNull() {
		return errors.ErrResourceGroupName
	}

	return nil
}
