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
	_ datasource.DataSource              = &AWSPrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &AWSPrivateEndpointCommand{}
)

// AWSPrivateEndpointCommand is the data source implementation.
type AWSPrivateEndpointCommand struct {
	*providerschema.Data
}

// NewAWSPrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewAWSPrivateEndpointCommand() datasource.DataSource {
	return &AWSPrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *AWSPrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_private_endpoint_command"
}

// Schema defines the schema for the private endpoint command data source.
func (a *AWSPrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AwsPrivateEndpointCommandSchema()
}

// Read refreshes the Terraform state with the latest data of private endpoint.
func (a *AWSPrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AWSCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating AWS private endpoint command request",
			"Could not validate AWS private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	AWSCommandRequest := api.CreateVPCEndpointCommandRequest{
		VpcID:     state.VpcID.ValueString(),
		SubnetIDs: convertSubnetIDs(state.SubnetIDs),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpointCommand", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		AWSCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AWS private endpoint command",
			"Could not read AWS private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var AWSCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &AWSCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling AWS private endpoint command response",
			"Could not unmarshall AWS private endpoint command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(AWSCommandResponse.Command)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the private endpoint command data source.
func (a *AWSPrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validate ensures organization id, project id, cluster id, and VPC id are valued.
func (a *AWSPrivateEndpointCommand) validate(config providerschema.AWSCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.VpcID.IsNull() {
		return errors.ErrVPCIDMissing
	}
	return nil
}

// convertSubnetIDs converts terraform string to go string.
func convertSubnetIDs(subnets []types.String) *[]string {
	var convertedSubnets []string
	for _, s := range subnets {
		convertedSubnets = append(convertedSubnets, s.ValueString())
	}

	return &convertedSubnets
}
