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
	_ datasource.DataSource              = &GCPPrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &GCPPrivateEndpointCommand{}
)

// GCPPrivateEndpointCommand is the data source implementation.
type GCPPrivateEndpointCommand struct {
	*providerschema.Data
}

// NewGCPPrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewGCPPrivateEndpointCommand() datasource.DataSource {
	return &GCPPrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *GCPPrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_private_endpoint_command"
}

// Schema defines the schema for the private endpoint command data source.
func (a *GCPPrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source to generate an GCP CLI command for setting up a private endpoint connection to an operational cluster. Retrieves the command or script to create the private endpoint, which establishes a private connection between the specified VPC and the Capella private endpoint service.",
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
			"vpc_network_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of your GCP VPC where the private endpoint will be created.",
			},
			"subnet_ids": schema.SetAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of subnet IDs in your VPC where the private endpoint interface will be created. These subnets must be in the same VPC.",
			},
			"command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The generated GCP CLI command. Use this command to create the private endpoint connection within GCP.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of private endpoint.
func (a *GCPPrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.GCPCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating GCP private endpoint command request",
			"Could not validate GCP private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	GCPCommandRequest := api.CreateGCPEndpointCommandRequest{
		VpcNetworkID: state.VpcNetworkID.ValueString(),
		SubnetIDs:    convertSubnetIDs(state.SubnetIDs),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpointCommand", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		GCPCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading GCP private endpoint command",
			"Could not read GCP private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var GCPCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &GCPCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling GCP private endpoint command response",
			"Could not unmarshall GCP private endpoint command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(GCPCommandResponse.Command)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the private endpoint command data source.
func (a *GCPPrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validate ensures organization id, project id, cluster id, and VPC Network id are valued.
func (a *GCPPrivateEndpointCommand) validate(config providerschema.GCPCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.VpcNetworkID.IsNull() {
		return errors.ErrVPCIDMissing
	}
	return nil
}
